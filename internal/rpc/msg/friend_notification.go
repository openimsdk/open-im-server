package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	utils2 "Open_IM/pkg/common/utils"
	pbFriend "Open_IM/pkg/proto/friend"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"github.com/golang/protobuf/proto"
)

//message MemberInfoChangedTips{
//  int32 ChangeType = 1; //1:info changed; 2:mute
//  GroupMemberFullInfo OpUser = 2; //who do this
//  GroupMemberFullInfo FinalInfo = 3; //
//  uint64 MuteTime = 4;
//  GroupInfo Group = 5;
//}
//func MemberInfoChangedNotification(operationID string, group *immysql.Group, opUser *immysql.GroupMember, userFinalInfo *immysql.GroupMember) {

//}

//message FriendApplicationAddedTips{
//  PublicUserInfo OpUser = 1; //user1
//  FriendApplication Application = 2;
//  PublicUserInfo  OpedUser = 3; //user2
//}

func getFromToUserNickname(fromUserID, toUserID string) (string, string, error) {
	from, err := imdb.GetUserByUserID(fromUserID)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	to, err := imdb.GetUserByUserID(toUserID)
	if err != nil {
		return "", "", utils.Wrap(err, "")
	}
	return from.Nickname, to.Nickname, nil
}

func friendNotification(commID *pbFriend.CommID, contentType int32, m proto.Message) {
	log.Info(commID.OperationID, utils.GetSelfFuncName(), "args: ", commID, contentType)
	var err error
	var tips open_im_sdk.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(commID.OperationID, "Marshal failed ", err.Error(), m.String())
		return
	}

	fromUserNickname, toUserNickname, err := getFromToUserNickname(commID.FromUserID, commID.ToUserID)
	if err != nil {
		log.Error(commID.OperationID, "getFromToUserNickname failed ", err.Error(), commID.FromUserID, commID.ToUserID)
		return
	}
	cn := config.Config.Notification
	switch contentType {
	case constant.FriendApplicationNotification:

		tips.DefaultTips = fromUserNickname + cn.FriendApplicationAdded.DefaultTips.Tips
	case constant.FriendApplicationApprovedNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplicationApproved.DefaultTips.Tips
	case constant.FriendApplicationRejectedNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendApplicationRejected.DefaultTips.Tips
	case constant.FriendAddedNotification:
		tips.DefaultTips = cn.FriendAdded.DefaultTips.Tips
	case constant.FriendDeletedNotification:
		tips.DefaultTips = cn.FriendDeleted.DefaultTips.Tips + toUserNickname
	case constant.FriendRemarkSetNotification:
		tips.DefaultTips = fromUserNickname + cn.FriendInfoChanged.DefaultTips.Tips
	case constant.BlackAddedNotification:
		tips.DefaultTips = cn.BlackAdded.DefaultTips.Tips + toUserNickname
	case constant.BlackDeletedNotification:
		tips.DefaultTips = cn.BlackDeleted.DefaultTips.Tips + toUserNickname
	default:
		log.Error(commID.OperationID, "contentType failed ", contentType)
		return
	}

	var n NotificationMsg
	n.SendID = commID.FromUserID
	n.RecvID = commID.ToUserID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = commID.OperationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(commID.OperationID, "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}

func FriendApplicationNotification(req *pbFriend.AddFriendReq) {
	log.Info(req.CommID.OperationID, utils.GetSelfFuncName(), "args: ", req.String())
	var friendApplicationAddedTips open_im_sdk.FriendApplicationAddedTips
	friendApplicationAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	friendApplicationAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	friendNotification(req.CommID, constant.FriendApplicationNotification, &friendApplicationAddedTips)
	//fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
	//	return
	//}
	//var tips open_im_sdk.TipsComm
	//tips.Detail, err = proto.Marshal(&friendApplicationAddedTips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), friendApplicationAddedTips.String())
	//	return
	//}
	//tips.DefaultTips = fromUserNickname + " FriendApplicationAddedNotification " + toUserNickname
	//
	//var n NotificationMsg
	//n.SendID = req.CommID.FromUserID
	//n.RecvID = req.CommID.ToUserID
	//n.ContentType = constant.FriendApplicationAddedNotification
	//n.SessionType = constant.SingleChatType
	//n.MsgFrom = constant.SysMsgType
	//n.OperationID = req.CommID.OperationID
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), tips.String())
	//	return
	//}
	//Notification(&n)
}

func FriendApplicationApprovedNotification(req *pbFriend.AddFriendResponseReq) {
	FriendApplicationApprovedTips := open_im_sdk.FriendApplicationApprovedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg
	//fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
	//	return
	//}

	friendNotification(req.CommID, constant.FriendApplicationApprovedNotification, &FriendApplicationApprovedTips)

	//var tips open_im_sdk.TipsComm
	//tips.DefaultTips = fromUserNickname + " FriendApplicationProcessedNotification " + toUserNickname
	//tips.Detail, err = proto.Marshal(&friendApplicationProcessedTips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), friendApplicationProcessedTips.String())
	//	return
	//}
	//
	//var n NotificationMsg
	//n.SendID = req.CommID.FromUserID
	//n.RecvID = req.CommID.ToUserID
	//n.ContentType = constant.FriendApplicationProcessedNotification
	//n.SessionType = constant.SingleChatType
	//n.MsgFrom = constant.SysMsgType
	//n.OperationID = req.CommID.OperationID
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), tips.String())
	//	return
	//}
	//Notification(&n)
}

func FriendApplicationRejectedNotification(req *pbFriend.AddFriendResponseReq) {
	FriendApplicationApprovedTips := open_im_sdk.FriendApplicationApprovedTips{FromToUserID: &open_im_sdk.FromToUserID{}}
	FriendApplicationApprovedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationApprovedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	FriendApplicationApprovedTips.HandleMsg = req.HandleMsg

	friendNotification(req.CommID, constant.FriendApplicationRejectedNotification, &FriendApplicationApprovedTips)
}

//
//
//func FriendApplicationProcessedNotification(req *pbFriend.AddFriendResponseReq) {
//	var friendApplicationProcessedTips open_im_sdk.FriendApplicationProcessedTips
//	friendApplicationProcessedTips.FromToUserID.FromUserID = req.CommID.FromUserID
//	friendApplicationProcessedTips.FromToUserID.ToUserID = req.CommID.ToUserID
//	friendApplicationProcessedTips.HandleResult = req.HandleResult
//	//fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
//	//if err != nil {
//	//	log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
//	//	return
//	//}
//	if friendApplicationProcessedTips.HandleResult == 1 {
//		friendNotification(req.CommID, constant.FriendApplicationApprovedNotification, &friendApplicationProcessedTips)
//	} else if friendApplicationProcessedTips.HandleResult == -1 {
//		friendNotification(req.CommID, constant.FriendApplicationRejectedNotification, &friendApplicationProcessedTips)
//	} else {
//		log.Error(req.CommID.OperationID, "HandleResult failed ", friendApplicationProcessedTips.HandleResult)
//	}
//
//	//var tips open_im_sdk.TipsComm
//	//tips.DefaultTips = fromUserNickname + " FriendApplicationProcessedNotification " + toUserNickname
//	//tips.Detail, err = proto.Marshal(&friendApplicationProcessedTips)
//	//if err != nil {
//	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), friendApplicationProcessedTips.String())
//	//	return
//	//}
//	//
//	//var n NotificationMsg
//	//n.SendID = req.CommID.FromUserID
//	//n.RecvID = req.CommID.ToUserID
//	//n.ContentType = constant.FriendApplicationProcessedNotification
//	//n.SessionType = constant.SingleChatType
//	//n.MsgFrom = constant.SysMsgType
//	//n.OperationID = req.CommID.OperationID
//	//n.Content, err = proto.Marshal(&tips)
//	//if err != nil {
//	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), tips.String())
//	//	return
//	//}
//	//Notification(&n)
//}

func FriendAddedNotification(operationID, opUserID, fromUserID, toUserID string) {
	var friendAddedTips open_im_sdk.FriendAddedTips
	user, err := imdb.GetUserByUserID(opUserID)
	if err != nil {
		log.NewError(operationID, "GetUserByUserID failed ", err.Error(), opUserID)
		return
	}
	utils2.UserDBCopyOpenIMPublicUser(friendAddedTips.OpUser, user)
	friend, err := imdb.GetFriendRelationshipFromFriend(fromUserID, toUserID)
	if err != nil {
		log.NewError(operationID, "GetFriendRelationshipFromFriend failed ", err.Error(), fromUserID, toUserID)
		return
	}
	utils2.FriendDBCopyOpenIM(friendAddedTips.Friend, friend)
	commID := pbFriend.CommID{FromUserID: fromUserID, ToUserID: toUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(&commID, constant.FriendAddedNotification, &friendAddedTips)
	//fromUserNickname, toUserNickname, err := getFromToUserNickname(fromUserID, toUserID)
	//if err != nil {
	//	log.Error(operationID, "getFromToUserNickname failed ", err.Error(), fromUserID, toUserID)
	//	return
	//}
	//var tips open_im_sdk.TipsComm
	//tips.DefaultTips = fromUserNickname + " FriendAddedNotification " + toUserNickname
	//tips.Detail, err = proto.Marshal(&friendAddedTips)
	//if err != nil {
	//	log.Error(operationID, "Marshal failed ", err.Error(), friendAddedTips)
	//	return
	//}
	//
	//var n NotificationMsg
	//n.SendID = fromUserID
	//n.RecvID = toUserID
	//n.ContentType = constant.FriendAddedNotification
	//n.SessionType = constant.SingleChatType
	//n.MsgFrom = constant.SysMsgType
	//n.OperationID = operationID
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.Error(operationID, "Marshal failed ", err.Error(), tips.String())
	//	return
	//}
	//Notification(&n)
}

//message FriendDeletedTips{
//  FriendInfo Friend = 1;
//}
func FriendDeletedNotification(req *pbFriend.DeleteFriendReq) {
	var friendDeletedTips open_im_sdk.FriendDeletedTips
	friendDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	friendDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	//fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", req.CommID.FromUserID, req.CommID.ToUserID)
	//	return
	//}
	friendNotification(req.CommID, constant.FriendDeletedNotification, &friendDeletedTips)
	//var tips open_im_sdk.TipsComm
	//tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	//tips.Detail, err = proto.Marshal(&friendDeletedTips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), friendDeletedTips.String())
	//	return
	//}
	//
	//var n NotificationMsg
	//n.SendID = req.CommID.FromUserID
	//n.RecvID = req.CommID.ToUserID
	//n.ContentType = constant.FriendDeletedNotification
	//n.SessionType = constant.SingleChatType
	//n.MsgFrom = constant.SysMsgType
	//n.OperationID = req.CommID.OperationID
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), tips.String())
	//	return
	//}
	//Notification(&n)
}

//message FriendInfoChangedTips{
//  FriendInfo Friend = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func FriendRemarkSetNotification(operationID, opUserID, fromUserID, toUserID string) {
	var friendInfoChangedTips open_im_sdk.FriendInfoChangedTips
	friendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	friendInfoChangedTips.FromToUserID.ToUserID = toUserID
	//fromUserNickname, toUserNickname, err := getFromToUserNickname(fromUserID, toUserID)
	//if err != nil {
	//	log.Error(operationID, "getFromToUserNickname failed ", fromUserID, toUserID)
	//	return
	//}
	commID := pbFriend.CommID{FromUserID: fromUserID, ToUserID: toUserID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(&commID, constant.FriendRemarkSetNotification, &friendInfoChangedTips)
	//var tips open_im_sdk.TipsComm
	//tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	//tips.Detail, err = proto.Marshal(&friendInfoChangedTips)
	//if err != nil {
	//	log.Error(operationID, "Marshal failed ", err.Error(), friendInfoChangedTips.String())
	//	return
	//}
	//
	//var n NotificationMsg
	//n.SendID = fromUserID
	//n.RecvID = toUserID
	//n.ContentType = constant.FriendInfoChangedNotification
	//n.SessionType = constant.SingleChatType
	//n.MsgFrom = constant.SysMsgType
	//n.OperationID = operationID
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.Error(operationID, "Marshal failed ", err.Error(), tips.String())
	//	return
	//}
	//Notification(&n)
}

func BlackAddedNotification(req *pbFriend.AddBlacklistReq) {
	var blackAddedTips open_im_sdk.BlackAddedTips
	blackAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	blackAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	//fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", req.CommID.FromUserID, req.CommID.ToUserID)
	//	return
	//}
	//var tips open_im_sdk.TipsComm
	//tips.DefaultTips = fromUserNickname + " BlackAddedNotification " + toUserNickname
	//tips.Detail, err = proto.Marshal(&blackAddedTips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), blackAddedTips.String())
	//	return
	//}
	friendNotification(req.CommID, constant.BlackAddedNotification, &blackAddedTips)
	//var n NotificationMsg
	//n.SendID = req.CommID.FromUserID
	//n.RecvID = req.CommID.ToUserID
	//n.ContentType = constant.BlackAddedNotification
	//n.SessionType = constant.SingleChatType
	//n.MsgFrom = constant.SysMsgType
	//n.OperationID = req.CommID.OperationID
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), tips.String())
	//	return
	//}
	//Notification(&n)
}

//message BlackDeletedTips{
//  BlackInfo Black = 1;
//}
func BlackDeletedNotification(req *pbFriend.RemoveBlacklistReq) {
	var blackDeletedTips open_im_sdk.BlackDeletedTips
	blackDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	blackDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	//fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
	//	return
	//}
	//var tips open_im_sdk.TipsComm
	//tips.Detail, err = proto.Marshal(&blackDeletedTips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), blackDeletedTips.String())
	//	return
	//}
	friendNotification(req.CommID, constant.BlackDeletedNotification, &blackDeletedTips)
	//tips.DefaultTips = fromUserNickname + " BlackDeletedNotification " + toUserNickname
	//var n NotificationMsg
	//n.SendID = req.CommID.FromUserID
	//n.RecvID = req.CommID.ToUserID
	//n.ContentType = constant.BlackDeletedNotification
	//n.SessionType = constant.SingleChatType
	//n.MsgFrom = constant.SysMsgType
	//n.OperationID = req.CommID.OperationID
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), tips.String())
	//	return
	//}
	//Notification(&n)
}

//message SelfInfoUpdatedTips{
//  UserInfo SelfUserInfo = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func SelfInfoUpdatedNotification(operationID, userID string, opUserID string) {
	var selfInfoUpdatedTips open_im_sdk.SelfInfoUpdatedTips
	selfInfoUpdatedTips.UserID = userID
	commID := pbFriend.CommID{FromUserID: userID, ToUserID: userID, OpUserID: opUserID, OperationID: operationID}
	friendNotification(&commID, constant.SelfInfoUpdatedNotification, &selfInfoUpdatedTips)
	//u, err := imdb.GetUserByUserID(userID)
	//if err != nil {
	//	log.NewError(operationID, "FindUserByUID failed ", err.Error(), userID)
	//	return
	//}
	//var tips open_im_sdk.TipsComm
	//tips.Detail, err = proto.Marshal(&selfInfoUpdatedTips)
	//if err != nil {
	//	log.Error(operationID, "Marshal failed ", err.Error(), selfInfoUpdatedTips.String())
	//	return
	//}
	//var n NotificationMsg
	//n.SendID = userID
	//n.RecvID = userID
	//n.ContentType = constant.SelfInfoUpdatedNotification
	//n.SessionType = constant.SingleChatType
	//n.MsgFrom = constant.SysMsgType
	//n.OperationID = operationID
	//
	//tips.DefaultTips = u.Nickname + " SelfInfoUpdatedNotification "
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.Error(operationID, "Marshal failed ", err.Error(), tips.String())
	//	return
	//}
	//Notification(&n)
}
