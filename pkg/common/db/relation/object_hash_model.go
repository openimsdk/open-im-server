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

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type ObjectHashGorm struct {
	*MetaDB
}

func NewObjectHash(db *gorm.DB) relation.ObjectHashModelInterface {
	return &ObjectHashGorm{
		NewMetaDB(db, &relation.ObjectHashModel{}),
	}
}

func (o *ObjectHashGorm) NewTx(tx any) relation.ObjectHashModelInterface {
	return &ObjectHashGorm{
		NewMetaDB(tx.(*gorm.DB), &relation.ObjectHashModel{}),
	}
}

func (o *ObjectHashGorm) Take(
	ctx context.Context,
	hash string,
	engine string,
) (oh *relation.ObjectHashModel, err error) {
	oh = &relation.ObjectHashModel{}
	return oh, utils.Wrap1(o.DB.Where("hash = ? and engine = ?", hash, engine).Take(oh).Error)
}

func (o *ObjectHashGorm) Create(ctx context.Context, h []*relation.ObjectHashModel) (err error) {
	return utils.Wrap1(o.DB.Create(h).Error)
}

func (o *ObjectHashGorm) DeleteNoCitation(
	ctx context.Context,
	engine string,
	num int,
) (list []*relation.ObjectHashModel, err error) {
	err = o.DB.Table(relation.ObjectHashModelTableName, "as h").Select("h.*").
		Joins("LEFT JOIN "+relation.ObjectInfoModelTableName+" as i ON h.hash = i.hash").
		Where("h.engine = ? AND i.hash IS NULL", engine).
		Limit(num).
		Find(&list).Error
	return list, utils.Wrap1(err)
}
