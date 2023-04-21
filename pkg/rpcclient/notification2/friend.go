package notification2

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	pbFriend "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type FriendNotificationSender struct {
	*rpcclient.MsgClient
	getUsersInfoMap func(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.UserInfo, error)
	getFriendsInfo  func(ctx context.Context, ownerUserID, friendUserID string) (resp *sdkws.FriendInfo, err error)
	getUsersInfo    func(ctx context.Context, userIDs []string, complete bool) ([]*sdkws.UserInfo, error)
}

func NewFriendNotificationSender(client discoveryregistry.SvcDiscoveryRegistry, getUsersInfoMap func(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.UserInfo, error)) *FriendNotificationSender {
	return &FriendNotificationSender{
		MsgClient:       rpcclient.NewMsgClient(client),
		getUsersInfoMap: getUsersInfoMap,
	}
}

func (c *FriendNotificationSender) getFromToUserNickname(ctx context.Context, fromUserID, toUserID string) (string, string, error) {
	users, err := c.getUsersInfoMap(ctx, []string{fromUserID, toUserID}, true)
	if err != nil {
		return "", "", nil
	}
	return users[fromUserID].Nickname, users[toUserID].Nickname, nil
}

func (c *FriendNotificationSender) friendNotification(ctx context.Context, fromUserID, toUserID string, contentType int32, m proto.Message) {
	var err error
	var tips sdkws.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		return
	}

	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}

	tips.JsonDetail, _ = marshaler.MarshalToString(m)

	fromUserNickname, toUserNickname, err := c.getFromToUserNickname(ctx, fromUserID, toUserID)
	if err != nil {
		return
	}
	cn := config.Config.Notification
	switch contentType {
	case constant.FriendApplicationNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplication.DefaultTips.Tips
	case constant.FriendApplicationApprovedNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplicationApproved.DefaultTips.Tips
	case constant.FriendApplicationRejectedNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplicationRejected.DefaultTips.Tips
	case constant.FriendAddedNotification:
		tips.DefaultTips = cn.FriendAdded.DefaultTips.Tips
	case constant.FriendDeletedNotification:
		tips.DefaultTips = cn.FriendDeleted.DefaultTips.Tips + toUserNickname
	case constant.FriendRemarkSetNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendRemarkSet.DefaultTips.Tips
	case constant.BlackAddedNotification:
		tips.DefaultTips = cn.BlackAdded.DefaultTips.Tips
	case constant.BlackDeletedNotification:
		tips.DefaultTips = cn.BlackDeleted.DefaultTips.Tips + toUserNickname
	case constant.UserInfoUpdatedNotification:
		tips.DefaultTips = cn.UserInfoUpdated.DefaultTips.Tips
	case constant.FriendInfoUpdatedNotification:
		tips.DefaultTips = cn.FriendInfoUpdated.DefaultTips.Tips + toUserNickname
	default:
		return
	}

	var n rpcclient.NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		return
	}
	c.Notification(ctx, &n)
}

func (c *FriendNotificationSender) FriendApplicationAddNotification(ctx context.Context, req *pbFriend.ApplyToAddFriendReq) {
	FriendApplicationTips := sdkws.FriendApplicationTips{FromToUserID: &sdkws.FromToUserID{}}
	FriendApplicationTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationTips.FromToUserID.ToUserID = req.ToUserID
	c.friendNotification(ctx, req.FromUserID, req.ToUserID, constant.FriendApplicationNotification, &FriendApplicationTips)
}

func (c *FriendNotificationSender) FriendApplicationAgreedNotification(ctx context.Context, req *pbFriend.RespondFriendApplyReq) {
	FriendApplicationApprovedTips := sdkws.FriendApplicationApprovedTips{FromToUserID: &sdkws.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	c.friendNotification(ctx, req.ToUserID, req.FromUserID, constant.FriendApplicationApprovedNotification, &FriendApplicationApprovedTips)
}

func (c *FriendNotificationSender) FriendApplicationRefusedNotification(ctx context.Context, req *pbFriend.RespondFriendApplyReq) {
	FriendApplicationApprovedTips := sdkws.FriendApplicationApprovedTips{FromToUserID: &sdkws.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	c.friendNotification(ctx, req.ToUserID, req.FromUserID, constant.FriendApplicationRejectedNotification, &FriendApplicationApprovedTips)
}

func (c *FriendNotificationSender) FriendAddedNotification(ctx context.Context, operationID, opUserID, fromUserID, toUserID string) {
	friendAddedTips := sdkws.FriendAddedTips{Friend: &sdkws.FriendInfo{}, OpUser: &sdkws.PublicUserInfo{}}
	user, err := c.getUsersInfo(ctx, []string{opUserID}, true)
	if err != nil {
		return
	}
	friendAddedTips.OpUser.UserID = user[0].UserID
	friendAddedTips.OpUser.Ex = user[0].Ex
	friendAddedTips.OpUser.Nickname = user[0].Nickname
	friendAddedTips.OpUser.FaceURL = user[0].FaceURL

	friend, err := c.getFriendsInfo(ctx, fromUserID, toUserID)
	if err != nil {
		return
	}
	friendAddedTips.Friend = friend
	c.friendNotification(ctx, fromUserID, toUserID, constant.FriendAddedNotification, &friendAddedTips)
}

func (c *FriendNotificationSender) FriendDeletedNotification(ctx context.Context, req *pbFriend.DeleteFriendReq) {
	friendDeletedTips := sdkws.FriendDeletedTips{FromToUserID: &sdkws.FromToUserID{}}
	friendDeletedTips.FromToUserID.FromUserID = req.OwnerUserID
	friendDeletedTips.FromToUserID.ToUserID = req.FriendUserID
	c.friendNotification(ctx, req.OwnerUserID, req.FriendUserID, constant.FriendDeletedNotification, &friendDeletedTips)
}

func (c *FriendNotificationSender) FriendRemarkSetNotification(ctx context.Context, fromUserID, toUserID string) {
	friendInfoChangedTips := sdkws.FriendInfoChangedTips{FromToUserID: &sdkws.FromToUserID{}}
	friendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	friendInfoChangedTips.FromToUserID.ToUserID = toUserID
	c.friendNotification(ctx, fromUserID, toUserID, constant.FriendRemarkSetNotification, &friendInfoChangedTips)
}

func (c *FriendNotificationSender) BlackAddedNotification(ctx context.Context, req *pbFriend.AddBlackReq) {
	blackAddedTips := sdkws.BlackAddedTips{FromToUserID: &sdkws.FromToUserID{}}
	blackAddedTips.FromToUserID.FromUserID = req.OwnerUserID
	blackAddedTips.FromToUserID.ToUserID = req.BlackUserID
	c.friendNotification(ctx, req.OwnerUserID, req.BlackUserID, constant.BlackAddedNotification, &blackAddedTips)
}

func (c *FriendNotificationSender) BlackDeletedNotification(ctx context.Context, req *pbFriend.RemoveBlackReq) {
	blackDeletedTips := sdkws.BlackDeletedTips{FromToUserID: &sdkws.FromToUserID{}}
	blackDeletedTips.FromToUserID.FromUserID = req.OwnerUserID
	blackDeletedTips.FromToUserID.ToUserID = req.BlackUserID
	c.friendNotification(ctx, req.OwnerUserID, req.BlackUserID, constant.BlackDeletedNotification, &blackDeletedTips)
}

func (c *FriendNotificationSender) FriendInfoUpdatedNotification(ctx context.Context, changedUserID string, needNotifiedUserID string, opUserID string) {
	selfInfoUpdatedTips := sdkws.UserInfoUpdatedTips{UserID: changedUserID}
	c.friendNotification(ctx, opUserID, needNotifiedUserID, constant.FriendInfoUpdatedNotification, &selfInfoUpdatedTips)
}
