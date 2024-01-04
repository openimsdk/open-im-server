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

package controller

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/pagination"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/tx"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type FriendDatabase interface {
	// CheckIn checks if user2 is in user1's friend list (inUser1Friends==true) and if user1 is in user2's friend list (inUser2Friends==true)
	CheckIn(ctx context.Context, user1, user2 string) (inUser1Friends bool, inUser2Friends bool, err error)

	// AddFriendRequest adds or updates a friend request
	AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error)

	// BecomeFriends first checks if the users are already in the friends table; if not, it inserts them as friends
	BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32) (err error)

	// RefuseFriendRequest refuses a friend request
	RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error)

	// AgreeFriendRequest accepts a friend request
	AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error)

	// Delete removes a friend or friends from the owner's friend list
	Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error)

	// UpdateRemark updates the remark for a friend
	UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error)

	// PageOwnerFriends retrieves the friend list of ownerUserID with pagination
	PageOwnerFriends(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendModel, err error)

	// PageInWhoseFriends finds the users who have friendUserID in their friend list with pagination
	PageInWhoseFriends(ctx context.Context, friendUserID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendModel, err error)

	// PageFriendRequestFromMe retrieves the friend requests sent by the user with pagination
	PageFriendRequestFromMe(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendRequestModel, err error)

	// PageFriendRequestToMe retrieves the friend requests received by the user with pagination
	PageFriendRequestToMe(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendRequestModel, err error)

	// FindFriendsWithError fetches specified friends of a user and returns an error if any do not exist
	FindFriendsWithError(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error)

	// FindFriendUserIDs retrieves the friend IDs of a user
	FindFriendUserIDs(ctx context.Context, ownerUserID string) (friendUserIDs []string, err error)

	// FindBothFriendRequests finds friend requests sent and received
	FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) (friends []*relation.FriendRequestModel, err error)

	// UpdateFriendPinStatus updates the pinned status of a friend
	UpdateFriendPinStatus(ctx context.Context, ownerUserID string, friendUserID string, isPinned bool) (err error)

	// UpdateFriendRemark updates the remark for a friend
	UpdateFriendRemark(ctx context.Context, ownerUserID string, friendUserID string, remark string) (err error)

	// UpdateFriendEx updates the 'ex' field for a friend
	UpdateFriendEx(ctx context.Context, ownerUserID string, friendUserID string, ex string) (err error)

}

type friendDatabase struct {
	friend        relation.FriendModelInterface
	friendRequest relation.FriendRequestModelInterface
	tx            tx.CtxTx
	cache         cache.FriendCache
}

func NewFriendDatabase(friend relation.FriendModelInterface, friendRequest relation.FriendRequestModelInterface, cache cache.FriendCache, tx tx.CtxTx) FriendDatabase {
	return &friendDatabase{friend: friend, friendRequest: friendRequest, cache: cache, tx: tx}
}

// ok 检查user2是否在user1的好友列表中(inUser1Friends==true) 检查user1是否在user2的好友列表中(inUser2Friends==true).
func (f *friendDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Friends bool, inUser2Friends bool, err error) {
	userID1FriendIDs, err := f.cache.GetFriendIDs(ctx, userID1)
	if err != nil {
		return
	}
	userID2FriendIDs, err := f.cache.GetFriendIDs(ctx, userID2)
	if err != nil {
		return
	}
	return utils.IsContain(userID2, userID1FriendIDs), utils.IsContain(userID1, userID2FriendIDs), nil
}

// 增加或者更新好友申请 如果之前有记录则更新，没有记录则新增.
func (f *friendDatabase) AddFriendRequest(ctx context.Context, fromUserID, toUserID string, reqMsg string, ex string) (err error) {
	return f.tx.Transaction(ctx, func(ctx context.Context) error {
		_, err := f.friendRequest.Take(ctx, fromUserID, toUserID)
		switch {
		case err == nil:
			m := make(map[string]any, 1)
			m["handle_result"] = 0
			m["handle_msg"] = ""
			m["req_msg"] = reqMsg
			m["ex"] = ex
			m["create_time"] = time.Now()
			return f.friendRequest.UpdateByMap(ctx, fromUserID, toUserID, m)
		case relation.IsNotFound(err):
			return f.friendRequest.Create(
				ctx,
				[]*relation.FriendRequestModel{{FromUserID: fromUserID, ToUserID: toUserID, ReqMsg: reqMsg, Ex: ex, CreateTime: time.Now(), HandleTime: time.Unix(0, 0)}},
			)
		default:
			return err
		}
	})
}

// (1)先判断是否在好友表 （在不在都不返回错误） (2)对于不在好友列表的 插入即可.
func (f *friendDatabase) BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32) (err error) {
	return f.tx.Transaction(ctx, func(ctx context.Context) error {
		cache := f.cache.NewCache()
		// 先find 找出重复的 去掉重复的
		fs1, err := f.friend.FindFriends(ctx, ownerUserID, friendUserIDs)
		if err != nil {
			return err
		}
		opUserID := mcontext.GetOperationID(ctx)
		for _, v := range friendUserIDs {
			fs1 = append(fs1, &relation.FriendModel{OwnerUserID: ownerUserID, FriendUserID: v, AddSource: addSource, OperatorUserID: opUserID})
		}
		fs11 := utils.DistinctAny(fs1, func(e *relation.FriendModel) string {
			return e.FriendUserID
		})

		err = f.friend.Create(ctx, fs11)
		if err != nil {
			return err
		}
		fs2, err := f.friend.FindReversalFriends(ctx, ownerUserID, friendUserIDs)
		if err != nil {
			return err
		}
		var newFriendIDs []string
		for _, v := range friendUserIDs {
			fs2 = append(fs2, &relation.FriendModel{OwnerUserID: v, FriendUserID: ownerUserID, AddSource: addSource, OperatorUserID: opUserID})
			newFriendIDs = append(newFriendIDs, v)
		}
		fs22 := utils.DistinctAny(fs2, func(e *relation.FriendModel) string {
			return e.OwnerUserID
		})
		err = f.friend.Create(ctx, fs22)
		if err != nil {
			return err
		}
		newFriendIDs = append(newFriendIDs, ownerUserID)
		cache = cache.DelFriendIDs(newFriendIDs...)
		return cache.ExecDel(ctx)

	})
}

// 拒绝好友申请 (1)检查是否有申请记录且为未处理状态 （没有记录返回错误） (2)修改申请记录 已拒绝.
func (f *friendDatabase) RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	fr, err := f.friendRequest.Take(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
	if err != nil {
		return err
	}
	if fr.HandleResult != 0 {
		return errs.ErrArgs.Wrap("the friend request has been processed")
	}
	log.ZDebug(ctx, "refuse friend request", "friendRequest db", fr, "friendRequest arg", friendRequest)
	friendRequest.HandleResult = constant.FriendResponseRefuse
	friendRequest.HandleTime = time.Now()
	err = f.friendRequest.Update(ctx, friendRequest)
	if err != nil {
		return err
	}
	return nil
}

// AgreeFriendRequest 同意好友申请  (1)检查是否有申请记录且为未处理状态 （没有记录返回错误） (2)检查是否好友（不返回错误）   (3) 建立双向好友关系（存在的忽略）.
func (f *friendDatabase) AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	return f.tx.Transaction(ctx, func(ctx context.Context) error {
		defer log.ZDebug(ctx, "return line")
		now := time.Now()
		fr, err := f.friendRequest.Take(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
		if err != nil {
			return err
		}
		if fr.HandleResult != 0 {
			return errs.ErrArgs.Wrap("the friend request has been processed")
		}
		friendRequest.HandlerUserID = mcontext.GetOpUserID(ctx)
		friendRequest.HandleResult = constant.FriendResponseAgree
		friendRequest.HandleTime = now
		err = f.friendRequest.Update(ctx, friendRequest)
		if err != nil {
			return err
		}

		fr2, err := f.friendRequest.Take(ctx, friendRequest.ToUserID, friendRequest.FromUserID)
		if err == nil && fr2.HandleResult == constant.FriendResponseNotHandle {
			fr2.HandlerUserID = mcontext.GetOpUserID(ctx)
			fr2.HandleResult = constant.FriendResponseAgree
			fr2.HandleTime = now
			err = f.friendRequest.Update(ctx, fr2)
			if err != nil {
				return err
			}
		} else if err != nil && (!relation.IsNotFound(err)) {
			return err
		}

		exists, err := f.friend.FindUserState(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
		if err != nil {
			return err
		}
		existsMap := utils.SliceSet(utils.Slice(exists, func(friend *relation.FriendModel) [2]string {
			return [...]string{friend.OwnerUserID, friend.FriendUserID} // 自己 - 好友
		}))
		var adds []*relation.FriendModel
		if _, ok := existsMap[[...]string{friendRequest.ToUserID, friendRequest.FromUserID}]; !ok { // 自己 - 好友
			adds = append(
				adds,
				&relation.FriendModel{
					OwnerUserID:    friendRequest.ToUserID,
					FriendUserID:   friendRequest.FromUserID,
					AddSource:      int32(constant.BecomeFriendByApply),
					OperatorUserID: friendRequest.FromUserID,
				},
			)
		}
		if _, ok := existsMap[[...]string{friendRequest.FromUserID, friendRequest.ToUserID}]; !ok { // 好友 - 自己
			adds = append(
				adds,
				&relation.FriendModel{
					OwnerUserID:    friendRequest.FromUserID,
					FriendUserID:   friendRequest.ToUserID,
					AddSource:      int32(constant.BecomeFriendByApply),
					OperatorUserID: friendRequest.FromUserID,
				},
			)
		}
		if len(adds) > 0 {
			if err := f.friend.Create(ctx, adds); err != nil {
				return err
			}
		}
		return f.cache.DelFriendIDs(friendRequest.ToUserID, friendRequest.FromUserID).ExecDel(ctx)
	})
}

// 删除好友  外部判断是否好友关系.
func (f *friendDatabase) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error) {
	if err := f.friend.Delete(ctx, ownerUserID, friendUserIDs); err != nil {
		return err
	}
	return f.cache.DelFriendIDs(append(friendUserIDs, ownerUserID)...).ExecDel(ctx)
}

// 更新好友备注 零值也支持.
func (f *friendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	if err := f.friend.UpdateRemark(ctx, ownerUserID, friendUserID, remark); err != nil {
		return err
	}
	return f.cache.DelFriend(ownerUserID, friendUserID).ExecDel(ctx)
}

// 获取ownerUserID的好友列表 无结果不返回错误.
func (f *friendDatabase) PageOwnerFriends(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendModel, err error) {
	return f.friend.FindOwnerFriends(ctx, ownerUserID, pagination)
}

// friendUserID在哪些人的好友列表中.
func (f *friendDatabase) PageInWhoseFriends(ctx context.Context, friendUserID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendModel, err error) {
	return f.friend.FindInWhoseFriends(ctx, friendUserID, pagination)
}

// 获取我发出去的好友申请  无结果不返回错误.
func (f *friendDatabase) PageFriendRequestFromMe(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendRequestModel, err error) {
	return f.friendRequest.FindFromUserID(ctx, userID, pagination)
}

// 获取我收到的的好友申请 无结果不返回错误.
func (f *friendDatabase) PageFriendRequestToMe(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendRequestModel, err error) {
	return f.friendRequest.FindToUserID(ctx, userID, pagination)
}

// 获取某人指定好友的信息 如果有好友不存在，也返回错误.
func (f *friendDatabase) FindFriendsWithError(ctx context.Context, ownerUserID string, friendUserIDs []string) (friends []*relation.FriendModel, err error) {
	friends, err = f.friend.FindFriends(ctx, ownerUserID, friendUserIDs)
	if err != nil {
		return
	}
	if len(friends) != len(friendUserIDs) {
		err = errs.ErrRecordNotFound.Wrap()
	}
	return
}

func (f *friendDatabase) FindFriendUserIDs(ctx context.Context, ownerUserID string) (friendUserIDs []string, err error) {
	return f.cache.GetFriendIDs(ctx, ownerUserID)
}

func (f *friendDatabase) FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) (friends []*relation.FriendRequestModel, err error) {
	return f.friendRequest.FindBothFriendRequests(ctx, fromUserID, toUserID)
}
func (f *friendDatabase) UpdateFriendPinStatus(ctx context.Context, ownerUserID string, friendUserID string, isPinned bool) (err error) {
	if err := f.friend.UpdatePinStatus(ctx, ownerUserID, friendUserID, isPinned); err != nil {
		return err
	}
	return f.cache.DelFriend(ownerUserID, friendUserID).ExecDel(ctx)
}
func (f *friendDatabase) UpdateFriendRemark(ctx context.Context, ownerUserID string, friendUserID string, remark string) (err error) {
	if err := f.friend.UpdateFriendRemark(ctx, ownerUserID, friendUserID, remark); err != nil {
		return err
	}
	return f.cache.DelFriend(ownerUserID, friendUserID).ExecDel(ctx)
}
func (f *friendDatabase) UpdateFriendEx(ctx context.Context, ownerUserID string, friendUserID string, ex string) (err error) {
	if err := f.friend.UpdateFriendEx(ctx, ownerUserID, friendUserID, ex); err != nil {
		return err
	}
	return f.cache.DelFriend(ownerUserID, friendUserID).ExecDel(ctx)
}
