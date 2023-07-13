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
	ObjectInfoModelTableName = "object"
)

type ObjectModel struct {
	Name        string    `gorm:"column:name;primary_key"`
	UserID      string    `gorm:"column:user_id"`
	Hash        string    `gorm:"column:hash"`
	Key         string    `gorm:"column:key"`
	Size        int64     `gorm:"column:size"`
	ContentType string    `gorm:"column:content_type"`
	Cause       string    `gorm:"column:cause"`
	CreateTime  time.Time `gorm:"column:create_time"`
}

func (ObjectModel) TableName() string {
	return ObjectInfoModelTableName
}

type ObjectInfoModelInterface interface {
	NewTx(tx any) ObjectInfoModelInterface
	SetObject(ctx context.Context, obj *ObjectModel) error
	Take(ctx context.Context, name string) (*ObjectModel, error)
}
