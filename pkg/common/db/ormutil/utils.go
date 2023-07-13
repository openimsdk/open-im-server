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

package ormutil

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

func GormPage[E any](db *gorm.DB, pageNumber, showNumber int32) (uint32, []*E, error) {
	var count int64
	var model E
	if err := db.Model(&model).Count(&count).Error; err != nil {
		return 0, nil, errs.Wrap(err)
	}
	var es []*E
	if err := db.Limit(int(showNumber)).Offset(int((pageNumber - 1) * showNumber)).Find(&es).Error; err != nil {
		return 0, nil, errs.Wrap(err)
	}
	return uint32(count), es, nil
}

func GormSearch[E any](db *gorm.DB, fields []string, value string, pageNumber, showNumber int32) (uint32, []*E, error) {
	if len(fields) > 0 && value != "" {
		val := "%" + value + "%"
		arr := make([]string, 0, len(fields))
		vals := make([]interface{}, 0, len(fields))
		for _, field := range fields {
			arr = append(arr, fmt.Sprintf("`%s` like ?", field))
			vals = append(vals, val)
		}
		db = db.Where(strings.Join(arr, " or "), vals...)
	}
	return GormPage[E](db, pageNumber, showNumber)
}

func GormIn[E any](db **gorm.DB, field string, es []E) {
	if len(es) == 0 {
		return
	}
	*db = (*db).Where(field+" in (?)", es)
}

func MapCount(db *gorm.DB, field string) (map[string]uint32, error) {
	var items []struct {
		ID    string `gorm:"column:id"`
		Count uint32 `gorm:"column:count"`
	}
	if err := db.Select(field + " as id, count(1) as count").Group(field).Find(&items).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	m := make(map[string]uint32)
	for _, item := range items {
		m[item.ID] = item.Count
	}
	return m, nil
}
