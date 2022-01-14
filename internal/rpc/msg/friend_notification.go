package msg

import (
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbFriend "Open_IM/pkg/proto/friend"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"encoding/json"
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

func getFromToUserNickname(fromUserID, toUserID string) (string, string) {
	from, err1 := imdb.GetUserByUserID(fromUserID)
	to, err2 := imdb.GetUserByUserID(toUserID)
	if err1 != nil || err2 != nil {
		log.NewError("FindUserByUID failed ", err1, err2, fromUserID, toUserID)
	}
	fromNickname, toNickname := "", ""
	if from != nil {
		fromNickname = from.Nickname
	}
	if to != nil {
		toNickname = to.Nickname
	}
	return fromNickname, toNickname
}

func FriendApplicationAddedNotification(req *pbFriend.AddFriendReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendApplicationAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var FriendApplicationAddedTips open_im_sdk.FriendApplicationAddedTips
	FriendApplicationAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.OperationID, req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendApplicationAddedTips)
	tips.DefaultTips = fromUserNickname + " FriendApplicationAddedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n, true)
}

func FriendApplicationProcessedNotification(req *pbFriend.AddFriendResponseReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendApplicationProcessedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var FriendApplicationProcessedTips open_im_sdk.FriendApplicationProcessedTips
	FriendApplicationProcessedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationProcessedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.OperationID, req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendApplicationProcessedTips)
	tips.DefaultTips = fromUserNickname + " FriendApplicationProcessedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

func FriendAddedNotification(operationID, opUserID, fromUserID, toUserID string) {
	var n NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = constant.FriendAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var FriendAddedTips open_im_sdk.FriendAddedTips

	user, err := imdb.GetUserByUserID(opUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), opUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.OpUser, user)
	}

	friend, err := imdb.GetFriendRelationshipFromFriend(fromUserID, toUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), fromUserID, toUserID)
	} else {
		FriendAddedTips.Friend.Remark = friend.Remark
	}

	from, err := imdb.GetUserByUserID(fromUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), fromUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.Friend, from)
	}

	to, err := imdb.GetUserByUserID(toUserID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), toUserID)

	} else {
		utils.CopyStructFields(FriendAddedTips.Friend.FriendUser, to)
	}

	fromUserNickname, toUserNickname := from.Nickname, to.Nickname
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendAddedTips)
	tips.DefaultTips = fromUserNickname + " FriendAddedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

//message FriendDeletedTips{
//  FriendInfo Friend = 1;
//}
func FriendDeletedNotification(req *pbFriend.DeleteFriendReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendDeletedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var FriendDeletedTips open_im_sdk.FriendDeletedTips
	FriendDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendDeletedTips)
	tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

//message FriendInfoChangedTips{
//  FriendInfo Friend = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func FriendInfoChangedNotification(operationID, opUserID, fromUserID, toUserID string) {
	var n NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = constant.FriendInfoChangedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var FriendInfoChangedTips open_im_sdk.FriendInfoChangedTips
	FriendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	FriendInfoChangedTips.FromToUserID.ToUserID = toUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(fromUserID, toUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendInfoChangedTips)
	tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

func BlackAddedNotification(req *pbFriend.AddBlacklistReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.BlackAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var BlackAddedTips open_im_sdk.BlackAddedTips
	BlackAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	BlackAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(BlackAddedTips)
	tips.DefaultTips = fromUserNickname + " BlackAddedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

//message BlackDeletedTips{
//  BlackInfo Black = 1;
//}
func BlackDeletedNotification(req *pbFriend.RemoveBlacklistReq) {
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.BlackDeletedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID

	var BlackDeletedTips open_im_sdk.BlackDeletedTips
	BlackDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	BlackDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(BlackDeletedTips)
	tips.DefaultTips = fromUserNickname + " BlackDeletedNotification " + toUserNickname
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

//message SelfInfoUpdatedTips{
//  UserInfo SelfUserInfo = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func SelfInfoUpdatedNotification(operationID, userID string) {
	var n NotificationMsg
	n.SendID = userID
	n.RecvID = userID
	n.ContentType = constant.SelfInfoUpdatedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	var SelfInfoUpdatedTips open_im_sdk.SelfInfoUpdatedTips
	SelfInfoUpdatedTips.UserID = userID

	var tips open_im_sdk.TipsComm
	u, err := imdb.GetUserByUserID(userID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), userID)
	}

	tips.Detail, _ = json.Marshal(SelfInfoUpdatedTips)
	tips.DefaultTips = u.Nickname + " SelfInfoUpdatedNotification "
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}
