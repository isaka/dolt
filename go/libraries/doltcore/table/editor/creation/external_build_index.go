// Copyright 2024 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package creation

import (
	"context"
	"errors"
	"io"

	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/dolt/go/libraries/doltcore/doltdb/durable"
	"github.com/dolthub/dolt/go/libraries/doltcore/schema"
	"github.com/dolthub/dolt/go/libraries/doltcore/sqle/index"
	"github.com/dolthub/dolt/go/store/prolly"
	"github.com/dolthub/dolt/go/store/prolly/sort"
	"github.com/dolthub/dolt/go/store/prolly/tree"
	"github.com/dolthub/dolt/go/store/types"
	"github.com/dolthub/dolt/go/store/util/tempfiles"
	"github.com/dolthub/dolt/go/store/val"
)

const (
	batchSize = 32 * 1024 * 1024 // 32MB
	fileMax   = 128
)

// BuildProllyIndexExternal builds unique and non-unique indexes with a
// single prolly tree materialization by presorting the index keys in an
// intermediate file format.
func BuildProllyIndexExternal(ctx *sql.Context, vrw types.ValueReadWriter, ns tree.NodeStore, sch schema.Schema, tableName string, idx schema.Index, primary prolly.Map, uniqCb DupEntryCb) (durable.Index, error) {
	iter, err := primary.IterAll(ctx)
	if err != nil {
		return nil, err
	}
	p := primary.Pool()

	keyDesc, _ := idx.Schema().GetMapDescriptors(ns)
	if schema.IsKeyless(sch) {
		keyDesc = prolly.AddHashToSchema(keyDesc)
	}

	prefixDesc := keyDesc.PrefixDesc(idx.Count())
	secondaryBld, err := index.NewSecondaryKeyBuilder(ctx, tableName, sch, idx, keyDesc, p, ns)
	if err != nil {
		return nil, err
	}

	if idx.IsVector() {
		return BuildProximityIndex(ctx, ns, idx, keyDesc, prefixDesc, iter, secondaryBld, uniqCb)
	}

	sorter := sort.NewTupleSorter(batchSize, fileMax, func(t1, t2 val.Tuple) bool {
		return keyDesc.Compare(ctx, t1, t2) < 0
	}, tempfiles.MovableTempFileProvider)
	defer sorter.Close()

	for {
		k, v, err := iter.Next(ctx)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		idxKey, err := secondaryBld.SecondaryKeyFromRow(ctx, k, v)
		if err != nil {
			return nil, err
		}

		if uniqCb != nil && prefixDesc.HasNulls(idxKey) {
			continue
		}

		if err := sorter.Insert(ctx, idxKey); err != nil {
			return nil, err
		}
	}

	sortedKeys, err := sorter.Flush(ctx)
	if err != nil {
		return nil, err
	}
	defer sortedKeys.Close()

	it, err := sortedKeys.IterAll(ctx)
	if err != nil {
		return nil, err
	}
	defer it.Close()

	empty, err := durable.NewEmptyIndexFromTableSchema(ctx, vrw, ns, idx, sch)
	secondary, err := durable.ProllyMapFromIndex(empty)
	if err != nil {
		return nil, err
	}

	tupIter := &tupleIterWithCb{iter: it, prefixDesc: prefixDesc, uniqCb: uniqCb}
	ret, err := prolly.MutateMapWithTupleIter(ctx, secondary, tupIter)
	if err != nil {
		return nil, err
	}
	if tupIter.err != nil {
		return nil, tupIter.err
	}

	return durable.IndexFromProllyMap(ret), nil
}

// func BuildProximityIndexExternal(ctx *sql.Context, vrw types.ValueReadWriter, ns tree.NodeStore, sch schema.Schema, tableName string, idx schema.Index, primary prolly.Map, uniqCb DupEntryCb) (durable.Index, error) {
func BuildProximityIndex(
	ctx *sql.Context,
	ns tree.NodeStore,
	idx schema.Index,
	keyDesc val.TupleDesc,
	prefixDesc val.TupleDesc,
	iter prolly.MapIter,
	secondaryBld index.SecondaryKeyBuilder,
	uniqCb DupEntryCb,
) (durable.Index, error) {
	// Secondary indexes have no non-key columns
	valDesc := val.NewTupleDescriptor()
	proximityMapBuilder, err := prolly.NewProximityMapBuilder(ctx, ns, idx.VectorProperties().DistanceType, keyDesc, valDesc, prolly.DefaultLogChunkSize)
	if err != nil {
		return nil, err
	}
	for {
		k, v, err := iter.Next(ctx)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		idxKey, err := secondaryBld.SecondaryKeyFromRow(ctx, k, v)
		if err != nil {
			return nil, err
		}

		if uniqCb != nil && prefixDesc.HasNulls(idxKey) {
			continue
		}

		if err := proximityMapBuilder.Insert(ctx, idxKey, val.EmptyTuple); err != nil {
			return nil, err
		}
	}
	proximityMap, err := proximityMapBuilder.Flush(ctx)
	return durable.IndexFromProximityMap(proximityMap), nil
}

type tupleIterWithCb struct {
	iter sort.KeyIter
	err  error

	prefixDesc val.TupleDesc
	uniqCb     DupEntryCb
	lastKey    val.Tuple
}

var _ prolly.TupleIter = (*tupleIterWithCb)(nil)

func (t *tupleIterWithCb) Next(ctx context.Context) (val.Tuple, val.Tuple) {
	for {
		curKey, err := t.iter.Next(ctx)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				t.err = err
			}
			return nil, nil
		}
		if t.lastKey != nil && t.prefixDesc.Compare(ctx, t.lastKey, curKey) == 0 && t.uniqCb != nil {
			// register a constraint violation if |key| collides with |lastKey|
			if err := t.uniqCb(ctx, t.lastKey, curKey); err != nil {
				t.err = err
				return nil, nil
			}

		}
		t.lastKey = curKey
		return curKey, val.EmptyTuple
	}
}
