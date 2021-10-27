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

package creation

import (
	"context"
	"fmt"
	"strings"

	"github.com/dolthub/dolt/go/libraries/doltcore/doltdb"
	"github.com/dolthub/dolt/go/libraries/doltcore/schema"
	"github.com/dolthub/dolt/go/libraries/doltcore/table/editor"
)

type CreateIndexReturn struct {
	NewTable *doltdb.Table
	Sch      schema.Schema
	OldIndex schema.Index
	NewIndex schema.Index
}

// CreateIndex creates the given index on the given table with the given schema. Returns the updated table, updated schema, and created index.
func CreateIndex(
	ctx context.Context,
	table *doltdb.Table,
	indexName string,
	columns []string,
	isUnique bool,
	isUserDefined bool,
	comment string,
	opts editor.Options,
) (*CreateIndexReturn, error) {
	sch, err := table.GetSchema(ctx)
	if err != nil {
		return nil, err
	}

	// get the real column names as CREATE INDEX columns are case-insensitive
	var realColNames []string
	allTableCols := sch.GetAllCols()
	for _, indexCol := range columns {
		tableCol, ok := allTableCols.GetByNameCaseInsensitive(indexCol)
		if !ok {
			return nil, fmt.Errorf("column `%s` does not exist for the table", indexCol)
		}
		realColNames = append(realColNames, tableCol.Name)
	}

	if indexName == "" {
		indexName = strings.Join(realColNames, "")
		_, ok := sch.Indexes().GetByNameCaseInsensitive(indexName)
		var i int
		for ok {
			i++
			indexName = fmt.Sprintf("%s_%d", strings.Join(realColNames, ""), i)
			_, ok = sch.Indexes().GetByNameCaseInsensitive(indexName)
		}
	}
	if !doltdb.IsValidIndexName(indexName) {
		return nil, fmt.Errorf("invalid index name `%s` as they must match the regular expression %s", indexName, doltdb.IndexNameRegexStr)
	}

	// if an index was already created for the column set but was not generated by the user then we replace it
	replacingIndex := false
	existingIndex, ok := sch.Indexes().GetIndexByColumnNames(realColNames...)
	if ok && !existingIndex.IsUserDefined() {
		replacingIndex = true
		_, err = sch.Indexes().RemoveIndex(existingIndex.Name())
		if err != nil {
			return nil, err
		}
	}

	// create the index metadata, will error if index names are taken or an index with the same columns in the same order exists
	index, err := sch.Indexes().AddIndexByColNames(
		indexName,
		realColNames,
		schema.IndexProperties{
			IsUnique:      isUnique,
			IsUserDefined: isUserDefined,
			Comment:       comment,
		},
	)
	if err != nil {
		return nil, err
	}

	// update the table schema with the new index
	newTable, err := table.UpdateSchema(ctx, sch)
	if err != nil {
		return nil, err
	}

	if replacingIndex { // verify that the pre-existing index data is valid
		newTable, err = newTable.RenameIndexRowData(ctx, existingIndex.Name(), index.Name())
		if err != nil {
			return nil, err
		}
		err = newTable.VerifyIndexRowData(ctx, index.Name())
		if err != nil {
			return nil, err
		}
	} else { // set the index row data and get a new root with the updated table
		//indexRowData, err := editor.RebuildIndex(ctx, newTable, index.Name(), opts)
		//if err != nil {
		//	return nil, err
		//}
		//newTable, err = newTable.SetIndexRowData(ctx, index.Name(), indexRowData)
		//if err != nil {
		//	return nil, err
		//}
	}
	return &CreateIndexReturn{
		NewTable: newTable,
		Sch:      sch,
		OldIndex: existingIndex,
		NewIndex: index,
	}, nil
}
