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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"time"

	"gorm.io/gorm"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

var _ relation.GroupModelInterface = (*GroupGorm)(nil)

type GroupGorm struct {
	*MetaDB
}

func NewGroupDB(db *gorm.DB) relation.GroupModelInterface {
	return &GroupGorm{NewMetaDB(db, &relation.GroupModel{})}
}

func (g *GroupGorm) NewTx(tx any) relation.GroupModelInterface {
	return &GroupGorm{NewMetaDB(tx.(*gorm.DB), &relation.GroupModel{})}
}

func (g *GroupGorm) Create(ctx context.Context, groups []*relation.GroupModel) (err error) {
	return utils.Wrap(g.DB.Create(&groups).Error, "")
}

func (g *GroupGorm) UpdateMap(ctx context.Context, groupID string, args map[string]interface{}) (err error) {
	return utils.Wrap(g.DB.Where("group_id = ?", groupID).Model(&relation.GroupModel{}).Updates(args).Error, "")
}

func (g *GroupGorm) UpdateStatus(ctx context.Context, groupID string, status int32) (err error) {
	return utils.Wrap(g.DB.Where("group_id = ?", groupID).Model(&relation.GroupModel{}).Updates(map[string]any{"status": status}).Error, "")
}

func (g *GroupGorm) Find(ctx context.Context, groupIDs []string) (groups []*relation.GroupModel, err error) {
	return groups, utils.Wrap(g.DB.Where("group_id in (?)", groupIDs).Find(&groups).Error, "")
}

func (g *GroupGorm) Take(ctx context.Context, groupID string) (group *relation.GroupModel, err error) {
	group = &relation.GroupModel{}
	return group, utils.Wrap(g.DB.Where("group_id = ?", groupID).Take(group).Error, "")
}

func (g *GroupGorm) Search(ctx context.Context, keyword string, pageNumber, showNumber int32) (total uint32, groups []*relation.GroupModel, err error) {
	db := g.DB
	db = db.WithContext(ctx).Where("status!=?", constant.GroupStatusDismissed)
	return ormutil.GormSearch[relation.GroupModel](db, []string{"name"}, keyword, pageNumber, showNumber)
}
func (g *GroupGorm) GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error) {
	return groupIDs, utils.Wrap(g.DB.Model(&relation.GroupModel{}).Where("group_type = ? ", groupType).Pluck("group_id", &groupIDs).Error, "")
}

func (g *GroupGorm) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	db := g.db(ctx).Model(&relation.GroupModel{})
	if before != nil {
		db = db.Where("create_time < ?", before)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (g *GroupGorm) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	var res []struct {
		Date  time.Time `gorm:"column:date"`
		Count int64     `gorm:"column:count"`
	}
	err := g.db(ctx).Model(&relation.GroupModel{}).Select("DATE(create_time) AS date, count(1) AS count").Where("create_time >= ? and create_time < ?", start, end).Group("date").Find(&res).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	v := make(map[string]int64)
	for _, r := range res {
		v[r.Date.Format("2006-01-02")] = r.Count
	}
	return v, nil
}
