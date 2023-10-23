// Copyright © 2023 OpenIM. All rights reserved.
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

package relation

import (
	"context"

	"gorm.io/gorm"
)

type MetaDB struct {
	DB    *gorm.DB
	table any
}

func NewMetaDB(db *gorm.DB, table any) *MetaDB {
	return &MetaDB{
		DB:    db,
		table: table,
	}
}

func (g *MetaDB) db(ctx context.Context) *gorm.DB {
	db := g.DB.WithContext(ctx).Model(g.table)

	return db
}
