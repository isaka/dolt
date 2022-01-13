// Copyright 2019 Dolthub, Inc.
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

package dtestutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dolthub/dolt/go/libraries/doltcore/doltdb"
	"github.com/dolthub/dolt/go/libraries/doltcore/env"
	"github.com/dolthub/dolt/go/libraries/doltcore/table/editor"
	"github.com/dolthub/dolt/go/libraries/utils/filesys"
	"github.com/dolthub/dolt/go/store/types"
)

const (
	TestHomeDir = "/user/bheni"
	WorkingDir  = "/user/bheni/datasets/states"
)

func testHomeDirFunc() (string, error) {
	return TestHomeDir, nil
}

func CreateTestEnv() *env.DoltEnv {
	const name = "billy bob"
	const email = "bigbillieb@fake.horse"
	initialDirs := []string{TestHomeDir, WorkingDir}
	fs := filesys.NewInMemFS(initialDirs, nil, WorkingDir)
	dEnv := env.Load(context.Background(), testHomeDirFunc, fs, doltdb.InMemDoltDB, "test")
	cfg, _ := dEnv.Config.GetConfig(env.GlobalConfig)
	cfg.SetStrings(map[string]string{
		env.UserNameKey:  name,
		env.UserEmailKey: email,
	})
	err := dEnv.InitRepo(context.Background(), types.Format_Default, name, email, env.DefaultInitBranch)

	if err != nil {
		panic("Failed to initialize environment:" + err.Error())
	}

	return dEnv
}

func CreateEnvWithSeedData(t *testing.T) *env.DoltEnv {
	dEnv := CreateTestEnv()
	imt, sch := CreateTestDataTable(true)

	ctx := context.Background()
	vrw := dEnv.DoltDB.ValueReadWriter()

	rowMap, err := types.NewMap(ctx, vrw)
	require.NoError(t, err)
	me := rowMap.Edit()
	for i := 0; i < imt.NumRows(); i++ {
		r, err := imt.GetRow(i)
		require.NoError(t, err)
		k, v := r.NomsMapKey(sch), r.NomsMapValue(sch)
		me.Set(k, v)
	}
	rowMap, err = me.Map(ctx)
	require.NoError(t, err)

	ai := sch.Indexes().AllIndexes()
	sch.Indexes().Merge(ai...)

	tbl, err := doltdb.NewTable(ctx, vrw, sch, rowMap, nil, nil)
	require.NoError(t, err)
	tbl, err = editor.RebuildAllIndexes(ctx, tbl, editor.TestEditorOptions(vrw))
	require.NoError(t, err)

	sch, err = tbl.GetSchema(ctx)
	require.NoError(t, err)
	rows, err := tbl.GetNomsRowData(ctx)
	require.NoError(t, err)
	indexes, err := tbl.GetIndexData(ctx)
	require.NoError(t, err)
	err = putTableToWorking(ctx, dEnv, sch, rows, indexes, TableName, nil)
	require.NoError(t, err)

	return dEnv
}
