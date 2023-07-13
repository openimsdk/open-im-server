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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

type ObjectInfoGorm struct {
	*MetaDB
}

func NewObjectInfo(db *gorm.DB) relation.ObjectInfoModelInterface {
	return &ObjectInfoGorm{
		NewMetaDB(db, &relation.ObjectModel{}),
	}
}

func (o *ObjectInfoGorm) NewTx(tx any) relation.ObjectInfoModelInterface {
	return &ObjectInfoGorm{
		NewMetaDB(tx.(*gorm.DB), &relation.ObjectModel{}),
	}
}

func (o *ObjectInfoGorm) SetObject(ctx context.Context, obj *relation.ObjectModel) (err error) {
	if err := o.DB.WithContext(ctx).Where("name = ?", obj.Name).FirstOrCreate(obj).Error; err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (o *ObjectInfoGorm) Take(ctx context.Context, name string) (info *relation.ObjectModel, err error) {
	info = &relation.ObjectModel{}
	return info, errs.Wrap(o.DB.WithContext(ctx).Where("name = ?", name).Take(info).Error)
}
