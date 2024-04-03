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
	"fmt"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
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

	// UpdateFriends updates fields for friends
	UpdateFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, val map[string]any) (err error)
}

type friendDatabase struct {
	friend        relation.FriendModelInterface
	friendRequest relation.FriendRequestModelInterface
	tx            tx.Tx
	cache         cache.FriendCache
}

func NewFriendDatabase(friend relation.FriendModelInterface, friendRequest relation.FriendRequestModelInterface, cache cache.FriendCache, tx tx.Tx) FriendDatabase {
	return &friendDatabase{friend: friend, friendRequest: friendRequest, cache: cache, tx: tx}
}

// CheckIn verifies if user2 is in user1's friend list (inUser1Friends returns true) and
// if user1 is in user2's friend list (inUser2Friends returns true).
func (f *friendDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Friends bool, inUser2Friends bool, err error) {
	// Retrieve friend IDs of userID1 from the cache
	userID1FriendIDs, err := f.cache.GetFriendIDs(ctx, userID1)
	if err != nil {
		err = fmt.Errorf("error retrieving friend IDs for user %s: %w", userID1, err)
		return
	}

	// Retrieve friend IDs of userID2 from the cache
	userID2FriendIDs, err := f.cache.GetFriendIDs(ctx, userID2)
	if err != nil {
		err = fmt.Errorf("error retrieving friend IDs for user %s: %w", userID2, err)
		return
	}

	// Check if userID2 is in userID1's friend list and vice versa
	inUser1Friends = datautil.Contain(userID2, userID1FriendIDs...)
	inUser2Friends = datautil.Contain(userID1, userID2FriendIDs...)
	return inUser1Friends, inUser2Friends, nil
}

// AddFriendRequest adds or updates a friend request.
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

// (1) First determine whether it is in the friends list (in or out does not return an error) (2) for not in the friends list can be inserted.
func (f *friendDatabase) BecomeFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, addSource int32) (err error) {
	return f.tx.Transaction(ctx, func(ctx context.Context) error {
		cache := f.cache.NewCache()
		// user find friends
		fs1, err := f.friend.FindFriends(ctx, ownerUserID, friendUserIDs)
		if err != nil {
			return err
		}
		opUserID := mcontext.GetOperationID(ctx)
		for _, v := range friendUserIDs {
			fs1 = append(fs1, &relation.FriendModel{OwnerUserID: ownerUserID, FriendUserID: v, AddSource: addSource, OperatorUserID: opUserID})
		}
		fs11 := datautil.DistinctAny(fs1, func(e *relation.FriendModel) string {
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
		fs22 := datautil.DistinctAny(fs2, func(e *relation.FriendModel) string {
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

// RefuseFriendRequest rejects a friend request. It first checks for an existing, unprocessed request.
// If no such request exists, it returns an error. Otherwise, it marks the request as refused.
func (f *friendDatabase) RefuseFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) error {
	// Attempt to retrieve the friend request from the database.
	fr, err := f.friendRequest.Take(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
	if err != nil {
		return fmt.Errorf("failed to retrieve friend request from %s to %s: %w", friendRequest.FromUserID, friendRequest.ToUserID, err)
	}

	// Check if the friend request has already been handled.
	if fr.HandleResult != 0 {
		return fmt.Errorf("friend request from %s to %s has already been processed", friendRequest.FromUserID, friendRequest.ToUserID)
	}

	// Log the action of refusing the friend request for debugging and auditing purposes.
	log.ZDebug(ctx, "Refusing friend request", map[string]interface{}{
		"DB_FriendRequest":  fr,
		"Arg_FriendRequest": friendRequest,
	})

	// Mark the friend request as refused and update the handle time.
	friendRequest.HandleResult = constant.FriendResponseRefuse
	friendRequest.HandleTime = time.Now()
	if err := f.friendRequest.Update(ctx, friendRequest); err != nil {
		return fmt.Errorf("failed to update friend request from %s to %s as refused: %w", friendRequest.FromUserID, friendRequest.ToUserID, err)
	}

	return nil
}

// AgreeFriendRequest accepts a friend request. It first checks for an existing, unprocessed request.
func (f *friendDatabase) AgreeFriendRequest(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	return f.tx.Transaction(ctx, func(ctx context.Context) error {
		now := time.Now()
		fr, err := f.friendRequest.Take(ctx, friendRequest.FromUserID, friendRequest.ToUserID)
		if err != nil {
			return err
		}
		if fr.HandleResult != 0 {
			return errs.ErrArgs.WrapMsg("the friend request has been processed")
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
		existsMap := datautil.SliceSet(datautil.Slice(exists, func(friend *relation.FriendModel) [2]string {
			return [...]string{friend.OwnerUserID, friend.FriendUserID} // My - Friend
		}))
		var adds []*relation.FriendModel
		if _, ok := existsMap[[...]string{friendRequest.ToUserID, friendRequest.FromUserID}]; !ok { // My - Friend
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
		if _, ok := existsMap[[...]string{friendRequest.FromUserID, friendRequest.ToUserID}]; !ok { // My - Friend
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

// Delete removes a friend relationship. It is assumed that the external caller has verified the friendship status.
func (f *friendDatabase) Delete(ctx context.Context, ownerUserID string, friendUserIDs []string) (err error) {
	if err := f.friend.Delete(ctx, ownerUserID, friendUserIDs); err != nil {
		return err
	}
	return f.cache.DelFriendIDs(append(friendUserIDs, ownerUserID)...).ExecDel(ctx)
}

// UpdateRemark updates the remark for a friend. Zero value for remark is also supported.
func (f *friendDatabase) UpdateRemark(ctx context.Context, ownerUserID, friendUserID, remark string) (err error) {
	if err := f.friend.UpdateRemark(ctx, ownerUserID, friendUserID, remark); err != nil {
		return err
	}
	return f.cache.DelFriend(ownerUserID, friendUserID).ExecDel(ctx)
}

// PageOwnerFriends retrieves the list of friends for the ownerUserID. It does not return an error if the result is empty.
func (f *friendDatabase) PageOwnerFriends(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendModel, err error) {
	return f.friend.FindOwnerFriends(ctx, ownerUserID, pagination)
}

// PageInWhoseFriends identifies in whose friend lists the friendUserID appears.
func (f *friendDatabase) PageInWhoseFriends(ctx context.Context, friendUserID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendModel, err error) {
	return f.friend.FindInWhoseFriends(ctx, friendUserID, pagination)
}

// PageFriendRequestFromMe retrieves friend requests sent by me. It does not return an error if the result is empty.
func (f *friendDatabase) PageFriendRequestFromMe(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendRequestModel, err error) {
	return f.friendRequest.FindFromUserID(ctx, userID, pagination)
}

// PageFriendRequestToMe retrieves friend requests received by me. It does not return an error if the result is empty.
func (f *friendDatabase) PageFriendRequestToMe(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, friends []*relation.FriendRequestModel, err error) {
	return f.friendRequest.FindToUserID(ctx, userID, pagination)
}

// FindFriendsWithError retrieves specified friends' information for ownerUserID. Returns an error if any friend does not exist.
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
func (f *friendDatabase) UpdateFriends(ctx context.Context, ownerUserID string, friendUserIDs []string, val map[string]any) (err error) {
	if len(val) == 0 {
		return nil
	}
	if err := f.friend.UpdateFriends(ctx, ownerUserID, friendUserIDs, val); err != nil {
		return err
	}
	return f.cache.DelFriends(ownerUserID, friendUserIDs).ExecDel(ctx)
}
