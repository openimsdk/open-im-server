package msg

import (
	"Open_IM/internal/common/check"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	pbFriend "Open_IM/pkg/proto/friend"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func getFromToUserNickname(fromUserID, toUserID string) (string, string, error) {
	users, err := check.GetUsersInfo(context.Background(), fromUserID, toUserID)
	if err != nil {
		return "", "", nil
	}
	if users[0].UserID == fromUserID {
		return users[0].Nickname, users[1].Nickname, nil
	}
	return users[1].Nickname, users[0].Nickname, nil

}

func friendNotification(operationID, fromUserID, toUserID string, contentType int32, m proto.Message) {
	log.Info(operationID, utils.GetSelfFuncName(), "args: ", contentType)
	var err error
	var tips open_im_sdk.TipsComm
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

func FriendApplicationAddNotification(ctx context.Context, req *pbFriend.AddFriendReq) {
	FriendApplicationTips := open_im_sdk.FriendApplicationTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	FriendApplicationTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationTips.FromToUserID.ToUserID = req.ToUserID
	friendNotification(tracelog.GetOperationID(ctx), req.FromUserID, req.ToUserID, constant.FriendApplicationNotification, &FriendApplicationTips)
}

func FriendApplicationAgreedNotification(ctx context.Context, req *pbFriend.RespondFriendApplyReq) {
	FriendApplicationApprovedTips := open_im_sdk.FriendApplicationApprovedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	friendNotification(tracelog.GetOperationID(ctx), req.ToUserID, req.FromUserID, constant.FriendApplicationApprovedNotification, &FriendApplicationApprovedTips)
}

func FriendApplicationRefusedNotification(ctx context.Context, req *pbFriend.RespondFriendApplyReq) {
	FriendApplicationApprovedTips := open_im_sdk.FriendApplicationApprovedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	friendNotification(tracelog.GetOperationID(ctx), req.ToUserID, req.FromUserID, constant.FriendApplicationRejectedNotification, &FriendApplicationApprovedTips)
}

func FriendAddedNotification(ctx context.Context, operationID, opUserID, fromUserID, toUserID string) {
	friendAddedTips := open_im_sdk.FriendAddedTips{Friend: &open_im_sdk.FriendInfo{}, OpUser: &open_im_sdk.PublicUserInfo{}}
	user, err := check.GetUsersInfo(context.Background(), opUserID)
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

func FriendDeletedNotification(ctx context.Context, req *pbFriend.DeleteFriendReq) {
	friendDeletedTips := open_im_sdk.FriendDeletedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	friendDeletedTips.FromToUserID.FromUserID = req.OwnerUserID
	friendDeletedTips.FromToUserID.ToUserID = req.FriendUserID
	friendNotification(tracelog.GetOperationID(ctx), req.OwnerUserID, req.FriendUserID, constant.FriendDeletedNotification, &friendDeletedTips)
}

func FriendRemarkSetNotification(ctx context.Context, fromUserID, toUserID string) {
	friendInfoChangedTips := open_im_sdk.FriendInfoChangedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	friendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	friendInfoChangedTips.FromToUserID.ToUserID = toUserID
	friendNotification(tracelog.GetOperationID(ctx), fromUserID, toUserID, constant.FriendRemarkSetNotification, &friendInfoChangedTips)
}

func BlackAddedNotification(ctx context.Context, req *pbFriend.AddBlackReq) {
	blackAddedTips := open_im_sdk.BlackAddedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	blackAddedTips.FromToUserID.FromUserID = req.OwnerUserID
	blackAddedTips.FromToUserID.ToUserID = req.BlackUserID
	friendNotification(tracelog.GetOperationID(ctx), req.OwnerUserID, req.BlackUserID, constant.BlackAddedNotification, &blackAddedTips)
}

func BlackDeletedNotification(ctx context.Context, req *pbFriend.RemoveBlackReq) {
	blackDeletedTips := open_im_sdk.BlackDeletedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	blackDeletedTips.FromToUserID.FromUserID = req.OwnerUserID
	blackDeletedTips.FromToUserID.ToUserID = req.BlackUserID
	friendNotification(tracelog.GetOperationID(ctx), req.OwnerUserID, req.BlackUserID, constant.BlackDeletedNotification, &blackDeletedTips)
}

// send to myself
func UserInfoUpdatedNotification(ctx context.Context, opUserID string, changedUserID string) {
	selfInfoUpdatedTips := open_im_sdk.UserInfoUpdatedTips{UserID: changedUserID}
	friendNotification(tracelog.GetOperationID(ctx), opUserID, changedUserID, constant.UserInfoUpdatedNotification, &selfInfoUpdatedTips)
}

func FriendInfoUpdatedNotification(ctx context.Context, changedUserID string, needNotifiedUserID string, opUserID string) {
	selfInfoUpdatedTips := open_im_sdk.UserInfoUpdatedTips{UserID: changedUserID}
	friendNotification(tracelog.GetOperationID(ctx), opUserID, needNotifiedUserID, constant.FriendInfoUpdatedNotification, &selfInfoUpdatedTips)
}
