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
	ObjectPutModelTableName = "object_put"
)

type ObjectPutModel struct {
	PutID         string     `gorm:"column:put_id;primary_key"`
	Hash          string     `gorm:"column:hash"`
	Path          string     `gorm:"column:path"`
	Name          string     `gorm:"column:name"`
	ContentType   string     `gorm:"column:content_type"`
	ObjectSize    int64      `gorm:"column:object_size"`
	FragmentSize  int64      `gorm:"column:fragment_size"`
	PutURLsHash   string     `gorm:"column:put_urls_hash"`
	ValidTime     *time.Time `gorm:"column:valid_time"`
	EffectiveTime time.Time  `gorm:"column:effective_time"`
	CreateTime    time.Time  `gorm:"column:create_time"`
}

func (ObjectPutModel) TableName() string {
	return ObjectPutModelTableName
}

type ObjectPutModelInterface interface {
	NewTx(tx any) ObjectPutModelInterface
	Create(ctx context.Context, m []*ObjectPutModel) error
	Take(ctx context.Context, putID string) (*ObjectPutModel, error)
	SetCompleted(ctx context.Context, putID string) error
	FindExpirationPut(ctx context.Context, expirationTime time.Time, num int) ([]*ObjectPutModel, error)
	DelPut(ctx context.Context, ids []string) error
}
