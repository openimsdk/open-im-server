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
	"time"
)

const (
	ObjectInfoModelTableName = "object_info"
)

type ObjectInfoModel struct {
	Name        string     `gorm:"column:name;primary_key"`
	Hash        string     `gorm:"column:hash"`
	ContentType string     `gorm:"column:content_type"`
	ValidTime   *time.Time `gorm:"column:valid_time"`
	CreateTime  time.Time  `gorm:"column:create_time"`
}

func (ObjectInfoModel) TableName() string {
	return ObjectInfoModelTableName
}

type ObjectInfoModelInterface interface {
	NewTx(tx any) ObjectInfoModelInterface
	SetObject(ctx context.Context, obj *ObjectInfoModel) error
	Take(ctx context.Context, name string) (*ObjectInfoModel, error)
	DeleteExpiration(ctx context.Context, expiration time.Time) error
}
