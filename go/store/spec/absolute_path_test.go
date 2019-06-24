// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package spec

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/liquidata-inc/ld/dolt/go/store/chunks"
	"github.com/liquidata-inc/ld/dolt/go/store/datas"
	"github.com/liquidata-inc/ld/dolt/go/store/hash"
	"github.com/liquidata-inc/ld/dolt/go/store/types"
	"github.com/stretchr/testify/assert"
)

func TestAbsolutePathToAndFromString(t *testing.T) {
	assert := assert.New(t)

	test := func(str string) {
		p, err := NewAbsolutePath(str)
		assert.NoError(err)
		assert.Equal(str, p.String())
	}

	h := types.Float(42).Hash() // arbitrary hash
	test(fmt.Sprintf("foo.bar[#%s]", h.String()))
	test(fmt.Sprintf("#%s.bar[42]", h.String()))
}

func TestAbsolutePaths(t *testing.T) {
	assert := assert.New(t)
	storage := &chunks.MemoryStorage{}
	db := datas.NewDatabase(storage.NewView())

	s0, s1 := types.String("foo"), types.String("bar")
	list := types.NewList(context.Background(), db, s0, s1)
	emptySet := types.NewSet(context.Background(), db)

	db.WriteValue(context.Background(), s0)
	db.WriteValue(context.Background(), s1)
	db.WriteValue(context.Background(), list)
	db.WriteValue(context.Background(), emptySet)

	var err error
	ds := db.GetDataset(context.Background(), "ds")
	ds, err = db.CommitValue(context.Background(), ds, list)
	assert.NoError(err)
	head := ds.Head()

	resolvesTo := func(exp types.Value, str string) {
		p, err := NewAbsolutePath(str)
		assert.NoError(err)
		act := p.Resolve(context.Background(), db)
		if exp == nil {
			assert.Nil(act)
		} else {
			assert.True(exp.Equals(act), "%s Expected %s Actual %s", str, types.EncodedValue(context.Background(), exp), types.EncodedValue(context.Background(), act))
		}
	}

	resolvesTo(head, "ds")
	resolvesTo(emptySet, "ds.parents")
	resolvesTo(list, "ds.value")
	resolvesTo(s0, "ds.value[0]")
	resolvesTo(s1, "ds.value[1]")
	resolvesTo(head, "#"+head.Hash().String())
	resolvesTo(list, "#"+list.Hash().String())
	resolvesTo(s0, "#"+s0.Hash().String())
	resolvesTo(s1, "#"+s1.Hash().String())
	resolvesTo(s0, "#"+list.Hash().String()+"[0]")
	resolvesTo(s1, "#"+list.Hash().String()+"[1]")

	resolvesTo(nil, "foo")
	resolvesTo(nil, "foo.parents")
	resolvesTo(nil, "foo.value")
	resolvesTo(nil, "foo.value[0]")
	resolvesTo(nil, "#"+types.String("baz").Hash().String())
	resolvesTo(nil, "#"+types.String("baz").Hash().String()+"[0]")
}

func TestReadAbsolutePaths(t *testing.T) {
	assert := assert.New(t)
	storage := &chunks.MemoryStorage{}
	db := datas.NewDatabase(storage.NewView())

	s0, s1 := types.String("foo"), types.String("bar")
	list := types.NewList(context.Background(), db, s0, s1)

	ds := db.GetDataset(context.Background(), "ds")
	_, err := db.CommitValue(context.Background(), ds, list)
	assert.NoError(err)

	vals, err := ReadAbsolutePaths(context.Background(), db, "ds.value[0]", "ds.value[1]")
	assert.NoError(err)

	assert.Equal(2, len(vals))
	assert.Equal("foo", string(vals[0].(types.String)))
	assert.Equal("bar", string(vals[1].(types.String)))

	vals, err = ReadAbsolutePaths(context.Background(), db, "!!#")
	assert.Nil(vals)
	assert.Equal("invalid input path '!!#'", err.Error())

	vals, err = ReadAbsolutePaths(context.Background(), db, "invalid.monkey")
	assert.Nil(vals)
	assert.Equal("input path 'invalid.monkey' does not exist in database", err.Error())
}

func TestAbsolutePathParseErrors(t *testing.T) {
	assert := assert.New(t)

	test := func(path, errMsg string) {
		p, err := NewAbsolutePath(path)
		assert.Equal(AbsolutePath{}, p)
		assert.Error(err)
		assert.Equal(errMsg, err.Error())
	}

	test("", "empty path")
	test(".foo", "invalid dataset name: .foo")
	test(".foo.bar.baz", "invalid dataset name: .foo.bar.baz")
	test("#", "invalid hash: ")
	test("#abc", "invalid hash: abc")
	invHash := strings.Repeat("z", hash.StringLen)
	test("#"+invHash, "invalid hash: "+invHash)
}
