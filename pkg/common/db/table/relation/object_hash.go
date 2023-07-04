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
	ObjectHashModelTableName = "object_hash"
)

type ObjectHashModel struct {
	Hash       string    `gorm:"column:hash;primary_key;size:32"`
	Engine     string    `gorm:"column:engine;primary_key;size:16"`
	Size       int64     `gorm:"column:size"`
	Bucket     string    `gorm:"column:bucket"`
	Name       string    `gorm:"column:name"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (ObjectHashModel) TableName() string {
	return ObjectHashModelTableName
}

type ObjectHashModelInterface interface {
	NewTx(tx any) ObjectHashModelInterface
	Take(ctx context.Context, hash string, engine string) (*ObjectHashModel, error)
	Create(ctx context.Context, h []*ObjectHashModel) error
	DeleteNoCitation(ctx context.Context, engine string, num int) (list []*ObjectHashModel, err error)
}
