// Copyright 2021 Dolthub, Inc.
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

package prolly

import (
	"context"

	"github.com/dolthub/dolt/go/store/skip"
	"github.com/dolthub/dolt/go/store/val"
)

type memoryMap struct {
	list    *skip.List
	keyDesc val.TupleDesc
}

func newMemoryMap(keyDesc val.TupleDesc, tups ...val.Tuple) (tm memoryMap) {
	if len(tups)%2 != 0 {
		panic("tuples must be key-value pairs")
	}

	tm.keyDesc = keyDesc

	// todo(andy): fix allocation for |tm.compare|
	tm.list = skip.NewSkipList(tm.compare)
	for i := 0; i < len(tups); i += 2 {
		tm.list.Put(tups[i], tups[i+1])
	}

	return
}

func (mm memoryMap) compare(left, right []byte) int {
	return int(mm.keyDesc.Compare(left, right))
}

func (mm memoryMap) Count() uint64 {
	return uint64(mm.list.Count())
}

func (mm memoryMap) Put(key, val val.Tuple) (ok bool) {
	ok = !mm.list.Full()
	if ok {
		mm.list.Put(key, val)
	}
	return
}

func (mm memoryMap) Get(_ context.Context, key val.Tuple, cb KeyValueFn) error {
	var value val.Tuple
	v, ok := mm.list.Get(key)
	if ok {
		value = v
	} else {
		key = nil
	}

	return cb(key, value)
}

func (mm memoryMap) Has(_ context.Context, key val.Tuple) (ok bool, err error) {
	_, ok = mm.list.Get(key)
	return
}

// IterAll returns a MapIterator that iterates over the entire Map.
func (mm memoryMap) IterAll(ctx context.Context) (MapRangeIter, error) {
	rng := Range{
		Start:   RangeCut{Unbound: true},
		Stop:    RangeCut{Unbound: true},
		KeyDesc: mm.keyDesc,
		Reverse: false,
	}
	return mm.IterValueRange(ctx, rng)
}

// IterValueRange returns a MapIterator that iterates over an ValueRange.
func (mm memoryMap) IterValueRange(ctx context.Context, rng Range) (MapRangeIter, error) {
	var iter *skip.ListIter
	if rng.Start.Unbound {
		if rng.Reverse {
			iter = mm.list.IterAtEnd()
		} else {
			iter = mm.list.IterAtStart()
		}
	} else {
		iter = mm.list.IterAt(rng.Start.Key)
	}

	tc := memTupleCursor{iter: iter}
	if err := startInRange(ctx, rng, tc); err != nil {
		return MapRangeIter{}, err
	}

	return MapRangeIter{
		memCur: tc,
		rng:    rng,
	}, nil
}

func (mm memoryMap) mutations() mutationIter {
	return memTupleCursor{iter: mm.list.IterAtStart()}
}

type memTupleCursor struct {
	iter    *skip.ListIter
	reverse bool
}

var _ tupleCursor = mapTupleCursor{}
var _ mutationIter = memTupleCursor{}

func (it memTupleCursor) next() (key, value val.Tuple) {
	key, value = it.iter.Current()
	if key == nil {
		return
	} else if it.reverse {
		it.iter.Retreat()
	} else {
		it.iter.Advance()
	}
	return
}

func (it memTupleCursor) current() (key, value val.Tuple) {
	return it.iter.Current()
}

func (it memTupleCursor) advance(context.Context) (err error) {
	it.iter.Advance()
	return
}

func (it memTupleCursor) retreat(context.Context) (err error) {
	it.iter.Retreat()
	return
}

func (it memTupleCursor) count() int {
	return it.iter.Count()
}

func (it memTupleCursor) close() error {
	return nil
}
