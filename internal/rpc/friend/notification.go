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

package friend

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/notification"
	"github.com/openimsdk/protocol/constant"
	pbfriend "github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"
)

type FriendNotificationSender struct {
	*rpcclient.NotificationSender
	// Target not found err
	getUsersInfo func(ctx context.Context, userIDs []string) ([]notification.CommonUser, error)
	// db controller
	db controller.FriendDatabase
}

type friendNotificationSenderOptions func(*FriendNotificationSender)

func WithFriendDB(db controller.FriendDatabase) friendNotificationSenderOptions {
	return func(s *FriendNotificationSender) {
		s.db = db
	}
}

func WithDBFunc(
	fn func(ctx context.Context, userIDs []string) (users []*relationtb.UserModel, err error),
) friendNotificationSenderOptions {
	return func(s *FriendNotificationSender) {
		f := func(ctx context.Context, userIDs []string) (result []notification.CommonUser, err error) {
			users, err := fn(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			for _, user := range users {
				result = append(result, user)
			}
			return result, nil
		}
		s.getUsersInfo = f
	}
}

func WithRpcFunc(
	fn func(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error),
) friendNotificationSenderOptions {
	return func(s *FriendNotificationSender) {
		f := func(ctx context.Context, userIDs []string) (result []notification.CommonUser, err error) {
			users, err := fn(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			for _, user := range users {
				result = append(result, user)
			}
			return result, err
		}
		s.getUsersInfo = f
	}
}

func NewFriendNotificationSender(
	conf *config.Notification,
	msgRpcClient *rpcclient.MessageRpcClient,
	opts ...friendNotificationSenderOptions,
) *FriendNotificationSender {
	f := &FriendNotificationSender{
		NotificationSender: rpcclient.NewNotificationSender(conf, rpcclient.WithRpcClient(msgRpcClient)),
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *FriendNotificationSender) getUsersInfoMap(
	ctx context.Context,
	userIDs []string,
) (map[string]*sdkws.UserInfo, error) {
	users, err := f.getUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*sdkws.UserInfo)
	for _, user := range users {
		result[user.GetUserID()] = user.(*sdkws.UserInfo)
	}
	return result, nil
}

//nolint:unused
func (f *FriendNotificationSender) getFromToUserNickname(
	ctx context.Context,
	fromUserID, toUserID string,
) (string, string, error) {
	users, err := f.getUsersInfoMap(ctx, []string{fromUserID, toUserID})
	if err != nil {
		return "", "", nil
	}
	return users[fromUserID].Nickname, users[toUserID].Nickname, nil
}

func (f *FriendNotificationSender) UserInfoUpdatedNotification(ctx context.Context, changedUserID string) {
	tips := sdkws.UserInfoUpdatedTips{UserID: changedUserID}
	f.Notification(ctx, mcontext.GetOpUserID(ctx), changedUserID, constant.UserInfoUpdatedNotification, &tips)
}

func (f *FriendNotificationSender) FriendApplicationAddNotification(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) {
	tips := sdkws.FriendApplicationTips{FromToUserID: &sdkws.FromToUserID{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
	}}
	f.Notification(ctx, req.FromUserID, req.ToUserID, constant.FriendApplicationNotification, &tips)
}

func (f *FriendNotificationSender) FriendApplicationAgreedNotification(
	ctx context.Context,
	req *pbfriend.RespondFriendApplyReq,
) {
	tips := sdkws.FriendApplicationApprovedTips{FromToUserID: &sdkws.FromToUserID{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
	}, HandleMsg: req.HandleMsg}
	f.Notification(ctx, req.ToUserID, req.FromUserID, constant.FriendApplicationApprovedNotification, &tips)
}

func (f *FriendNotificationSender) FriendApplicationRefusedNotification(
	ctx context.Context,
	req *pbfriend.RespondFriendApplyReq,
) {
	tips := sdkws.FriendApplicationApprovedTips{FromToUserID: &sdkws.FromToUserID{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
	}, HandleMsg: req.HandleMsg}
	f.Notification(ctx, req.ToUserID, req.FromUserID, constant.FriendApplicationRejectedNotification, &tips)
}

func (f *FriendNotificationSender) FriendAddedNotification(
	ctx context.Context,
	operationID, opUserID, fromUserID, toUserID string,
) error {
	tips := sdkws.FriendAddedTips{Friend: &sdkws.FriendInfo{}, OpUser: &sdkws.PublicUserInfo{}}
	user, err := f.getUsersInfo(ctx, []string{opUserID})
	if err != nil {
		return err
	}
	tips.OpUser.UserID = user[0].GetUserID()
	tips.OpUser.Ex = user[0].GetEx()
	tips.OpUser.Nickname = user[0].GetNickname()
	tips.OpUser.FaceURL = user[0].GetFaceURL()
	friends, err := f.db.FindFriendsWithError(ctx, fromUserID, []string{toUserID})
	if err != nil {
		return err
	}
	tips.Friend, err = convert.FriendDB2Pb(ctx, friends[0], f.getUsersInfoMap)
	if err != nil {
		return err
	}
	f.Notification(ctx, fromUserID, toUserID, constant.FriendAddedNotification, &tips)
	return nil
}

func (f *FriendNotificationSender) FriendDeletedNotification(ctx context.Context, req *pbfriend.DeleteFriendReq) {
	tips := sdkws.FriendDeletedTips{FromToUserID: &sdkws.FromToUserID{
		FromUserID: req.OwnerUserID,
		ToUserID:   req.FriendUserID,
	}}
	f.Notification(ctx, req.OwnerUserID, req.FriendUserID, constant.FriendDeletedNotification, &tips)
}

func (f *FriendNotificationSender) FriendRemarkSetNotification(ctx context.Context, fromUserID, toUserID string) {
	tips := sdkws.FriendInfoChangedTips{FromToUserID: &sdkws.FromToUserID{}}
	tips.FromToUserID.FromUserID = fromUserID
	tips.FromToUserID.ToUserID = toUserID
	f.Notification(ctx, fromUserID, toUserID, constant.FriendRemarkSetNotification, &tips)
}

func (f *FriendNotificationSender) FriendsInfoUpdateNotification(ctx context.Context, toUserID string, friendIDs []string) {
	tips := sdkws.FriendsInfoUpdateTips{FromToUserID: &sdkws.FromToUserID{}}
	tips.FromToUserID.ToUserID = toUserID
	tips.FriendIDs = friendIDs
	f.Notification(ctx, toUserID, toUserID, constant.FriendsInfoUpdateNotification, &tips)
}

func (f *FriendNotificationSender) BlackAddedNotification(ctx context.Context, req *pbfriend.AddBlackReq) {
	tips := sdkws.BlackAddedTips{FromToUserID: &sdkws.FromToUserID{}}
	tips.FromToUserID.FromUserID = req.OwnerUserID
	tips.FromToUserID.ToUserID = req.BlackUserID
	f.Notification(ctx, req.OwnerUserID, req.BlackUserID, constant.BlackAddedNotification, &tips)
}

func (f *FriendNotificationSender) BlackDeletedNotification(ctx context.Context, req *pbfriend.RemoveBlackReq) {
	blackDeletedTips := sdkws.BlackDeletedTips{FromToUserID: &sdkws.FromToUserID{
		FromUserID: req.OwnerUserID,
		ToUserID:   req.BlackUserID,
	}}
	f.Notification(ctx, req.OwnerUserID, req.BlackUserID, constant.BlackDeletedNotification, &blackDeletedTips)
}

func (f *FriendNotificationSender) FriendInfoUpdatedNotification(ctx context.Context, changedUserID string, needNotifiedUserID string) {
	tips := sdkws.UserInfoUpdatedTips{UserID: changedUserID}
	f.Notification(ctx, mcontext.GetOpUserID(ctx), needNotifiedUserID, constant.FriendInfoUpdatedNotification, &tips)
}
