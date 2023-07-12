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

const FriendRequestModelTableName = "friend_requests"

type FriendRequestModel struct {
	FromUserID    string    `gorm:"column:from_user_id;primary_key;size:64"`
	ToUserID      string    `gorm:"column:to_user_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:255"`
	CreateTime    time.Time `gorm:"column:create_time; autoCreateTime"`
	HandlerUserID string    `gorm:"column:handler_user_id;size:64"`
	HandleMsg     string    `gorm:"column:handle_msg;size:255"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

func (FriendRequestModel) TableName() string {
	return FriendRequestModelTableName
}

type FriendRequestModelInterface interface {
	// 插入多条记录
	Create(ctx context.Context, friendRequests []*FriendRequestModel) (err error)
	// 删除记录
	Delete(ctx context.Context, fromUserID, toUserID string) (err error)
	// 更新零值
	UpdateByMap(ctx context.Context, formUserID string, toUserID string, args map[string]interface{}) (err error)
	// 更新多条记录 （非零值）
	Update(ctx context.Context, friendRequest *FriendRequestModel) (err error)
	// 获取来指定用户的好友申请  未找到 不返回错误
	Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *FriendRequestModel, err error)
	Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *FriendRequestModel, err error)
	// 获取toUserID收到的好友申请列表
	FindToUserID(
		ctx context.Context,
		toUserID string,
		pageNumber, showNumber int32,
	) (friendRequests []*FriendRequestModel, total int64, err error)
	// 获取fromUserID发出去的好友申请列表
	FindFromUserID(
		ctx context.Context,
		fromUserID string,
		pageNumber, showNumber int32,
	) (friendRequests []*FriendRequestModel, total int64, err error)

	NewTx(tx any) FriendRequestModelInterface
}
