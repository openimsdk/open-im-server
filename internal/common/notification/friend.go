package notification

import (
	"Open_IM/internal/common/check"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func getFromToUserNickname(userCheck *check.UserCheck, fromUserID, toUserID string) (string, string, error) {
	users, err := userCheck.GetUsersInfoMap(context.Background(), []string{fromUserID, toUserID}, true)
	if err != nil {
		return "", "", nil
	}
	return users[fromUserID].Nickname, users[toUserID].Nickname, nil
}

func friendNotification(operationID, fromUserID, toUserID string, contentType int32, m proto.Message) {
	log.Info(operationID, utils.GetSelfFuncName(), "args: ", contentType)
	var err error
	var tips sdkws.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), m.String())
		return
	}

	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}

	tips.JsonDetail, _ = marshaler.MarshalToString(m)

	fromUserNickname, toUserNickname, err := getFromToUserNickname(fromUserID, toUserID)
	if err != nil {
		log.Error(operationID, "getFromToUserNickname failed ", err.Error(), fromUserID, toUserID)
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
		log.Error(operationID, "contentType failed ", contentType)
		return
	}

	var n NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}

func (c *Check) FriendApplicationAddNotification(ctx context.Context, req *pbFriend.ApplyToAddFriendReq) {
	FriendApplicationTips := sdkws.FriendApplicationTips{FromToUserID: &sdkws.FromToUserID{}}
	FriendApplicationTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationTips.FromToUserID.ToUserID = req.ToUserID
	friendNotification(tracelog.GetOperationID(ctx), req.FromUserID, req.ToUserID, constant.FriendApplicationNotification, &FriendApplicationTips)
}

func (c *Check) FriendApplicationAgreedNotification(ctx context.Context, req *pbFriend.RespondFriendApplyReq) {
	FriendApplicationApprovedTips := sdkws.FriendApplicationApprovedTips{FromToUserID: &sdkws.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	friendNotification(tracelog.GetOperationID(ctx), req.ToUserID, req.FromUserID, constant.FriendApplicationApprovedNotification, &FriendApplicationApprovedTips)
}

func (c *Check) FriendApplicationRefusedNotification(ctx context.Context, req *pbFriend.RespondFriendApplyReq) {
	FriendApplicationApprovedTips := sdkws.FriendApplicationApprovedTips{FromToUserID: &sdkws.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	friendNotification(tracelog.GetOperationID(ctx), req.ToUserID, req.FromUserID, constant.FriendApplicationRejectedNotification, &FriendApplicationApprovedTips)
}

func (c *Check) FriendAddedNotification(ctx context.Context, operationID, opUserID, fromUserID, toUserID string) {
	friendAddedTips := sdkws.FriendAddedTips{Friend: &sdkws.FriendInfo{}, OpUser: &sdkws.PublicUserInfo{}}
	user, err := check.NewUserCheck().GetUsersInfos(context.Background(), []string{opUserID}, true)
	if err != nil {
		return
	}
	friendAddedTips.OpUser.UserID = user[0].UserID
	friendAddedTips.OpUser.Ex = user[0].Ex
	friendAddedTips.OpUser.Nickname = user[0].Nickname
	friendAddedTips.OpUser.FaceURL = user[0].FaceURL

	friend, err := check.GetFriendsInfo(ctx, fromUserID, toUserID)
	if err != nil {
		return
	}
	friendAddedTips.Friend = friend
	friendNotification(operationID, fromUserID, toUserID, constant.FriendAddedNotification, &friendAddedTips)
}

func (c *Check) FriendDeletedNotification(ctx context.Context, req *pbFriend.DeleteFriendReq) {
	friendDeletedTips := sdkws.FriendDeletedTips{FromToUserID: &sdkws.FromToUserID{}}
	friendDeletedTips.FromToUserID.FromUserID = req.OwnerUserID
	friendDeletedTips.FromToUserID.ToUserID = req.FriendUserID
	friendNotification(tracelog.GetOperationID(ctx), req.OwnerUserID, req.FriendUserID, constant.FriendDeletedNotification, &friendDeletedTips)
}

func (c *Check) FriendRemarkSetNotification(ctx context.Context, fromUserID, toUserID string) {
	friendInfoChangedTips := sdkws.FriendInfoChangedTips{FromToUserID: &sdkws.FromToUserID{}}
	friendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	friendInfoChangedTips.FromToUserID.ToUserID = toUserID
	friendNotification(tracelog.GetOperationID(ctx), fromUserID, toUserID, constant.FriendRemarkSetNotification, &friendInfoChangedTips)
}

func (c *Check) BlackAddedNotification(ctx context.Context, req *pbFriend.AddBlackReq) {
	blackAddedTips := sdkws.BlackAddedTips{FromToUserID: &sdkws.FromToUserID{}}
	blackAddedTips.FromToUserID.FromUserID = req.OwnerUserID
	blackAddedTips.FromToUserID.ToUserID = req.BlackUserID
	friendNotification(tracelog.GetOperationID(ctx), req.OwnerUserID, req.BlackUserID, constant.BlackAddedNotification, &blackAddedTips)
}

func (c *Check) BlackDeletedNotification(ctx context.Context, req *pbFriend.RemoveBlackReq) {
	blackDeletedTips := sdkws.BlackDeletedTips{FromToUserID: &sdkws.FromToUserID{}}
	blackDeletedTips.FromToUserID.FromUserID = req.OwnerUserID
	blackDeletedTips.FromToUserID.ToUserID = req.BlackUserID
	friendNotification(tracelog.GetOperationID(ctx), req.OwnerUserID, req.BlackUserID, constant.BlackDeletedNotification, &blackDeletedTips)
}

// send to myself
func (c *Check) UserInfoUpdatedNotification(ctx context.Context, opUserID string, changedUserID string) {
	selfInfoUpdatedTips := sdkws.UserInfoUpdatedTips{UserID: changedUserID}
	friendNotification(tracelog.GetOperationID(ctx), opUserID, changedUserID, constant.UserInfoUpdatedNotification, &selfInfoUpdatedTips)
}

func (c *Check) FriendInfoUpdatedNotification(ctx context.Context, changedUserID string, needNotifiedUserID string, opUserID string) {
	selfInfoUpdatedTips := sdkws.UserInfoUpdatedTips{UserID: changedUserID}
	friendNotification(tracelog.GetOperationID(ctx), opUserID, needNotifiedUserID, constant.FriendInfoUpdatedNotification, &selfInfoUpdatedTips)
}
