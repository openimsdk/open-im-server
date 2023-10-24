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

	"github.com/OpenIMSDK/tools/ormutil"

	"gorm.io/gorm"

	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type GroupRequestGorm struct {
	*MetaDB
}

func NewGroupRequest(db *gorm.DB) relation.GroupRequestModelInterface {
	return &GroupRequestGorm{
		NewMetaDB(db, &relation.GroupRequestModel{}),
	}
}

func (g *GroupRequestGorm) NewTx(tx any) relation.GroupRequestModelInterface {
	return &GroupRequestGorm{NewMetaDB(tx.(*gorm.DB), &relation.GroupRequestModel{})}
}

func (g *GroupRequestGorm) Create(ctx context.Context, groupRequests []*relation.GroupRequestModel) (err error) {
	return utils.Wrap(g.DB.WithContext(ctx).Create(&groupRequests).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Delete(ctx context.Context, groupID string, userID string) (err error) {
	return utils.Wrap(
		g.DB.WithContext(ctx).
			Where("group_id = ? and user_id = ? ", groupID, userID).
			Delete(&relation.GroupRequestModel{}).
			Error,
		utils.GetSelfFuncName(),
	)
}

func (g *GroupRequestGorm) UpdateHandler(
	ctx context.Context,
	groupID string,
	userID string,
	handledMsg string,
	handleResult int32,
) (err error) {
	return utils.Wrap(
		g.DB.WithContext(ctx).
			Model(&relation.GroupRequestModel{}).
			Where("group_id = ? and user_id = ? ", groupID, userID).
			Updates(map[string]any{
				"handle_msg":    handledMsg,
				"handle_result": handleResult,
			}).
			Error,
		utils.GetSelfFuncName(),
	)
}

func (g *GroupRequestGorm) Take(
	ctx context.Context,
	groupID string,
	userID string,
) (groupRequest *relation.GroupRequestModel, err error) {
	groupRequest = &relation.GroupRequestModel{}

	return groupRequest, utils.Wrap(
		g.DB.WithContext(ctx).Where("group_id = ? and user_id = ? ", groupID, userID).Take(groupRequest).Error,
		utils.GetSelfFuncName(),
	)
}

func (g *GroupRequestGorm) Page(
	ctx context.Context,
	userID string,
	pageNumber, showNumber int32,
) (total uint32, groups []*relation.GroupRequestModel, err error) {
	return ormutil.GormSearch[relation.GroupRequestModel](
		g.DB.WithContext(ctx).Where("user_id = ?", userID),
		nil,
		"",
		pageNumber,
		showNumber,
	)
}

func (g *GroupRequestGorm) PageGroup(
	ctx context.Context,
	groupIDs []string,
	pageNumber, showNumber int32,
) (total uint32, groups []*relation.GroupRequestModel, err error) {
	return ormutil.GormPage[relation.GroupRequestModel](
		g.DB.WithContext(ctx).Where("group_id in ?", groupIDs),
		pageNumber,
		showNumber,
	)
}

func (g *GroupRequestGorm) FindGroupRequests(ctx context.Context, groupID string, userIDs []string) (total int64, groupRequests []*relation.GroupRequestModel, err error) {
	err = g.DB.WithContext(ctx).Where("group_id = ? and user_id in ?", groupID, userIDs).Find(&groupRequests).Error

	return int64(len(groupRequests)), groupRequests, utils.Wrap(err, utils.GetSelfFuncName())
}
