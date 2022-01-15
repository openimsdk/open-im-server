package msg

import (
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	utils2 "Open_IM/pkg/common/utils"
	pbFriend "Open_IM/pkg/proto/friend"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"encoding/json"
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

func FriendApplicationAddedNotification(req *pbFriend.AddFriendReq) {
	var FriendApplicationAddedTips open_im_sdk.FriendApplicationAddedTips
	FriendApplicationAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
		return
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendApplicationAddedTips)
	tips.DefaultTips = fromUserNickname + " FriendApplicationAddedNotification " + toUserNickname

	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendApplicationAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

func FriendApplicationProcessedNotification(req *pbFriend.AddFriendResponseReq) {
	var FriendApplicationProcessedTips open_im_sdk.FriendApplicationProcessedTips
	FriendApplicationProcessedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendApplicationProcessedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
		return
	}

	var tips open_im_sdk.TipsComm
	tips.Detail, err = proto.Marshal(&FriendApplicationProcessedTips)
	if err != nil {
		log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), FriendApplicationProcessedTips)
		return
	}
	tips.DefaultTips = fromUserNickname + " FriendApplicationProcessedNotification " + toUserNickname

	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendApplicationProcessedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID
	n.Content, _ = json.Marshal(tips)
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(req.CommID.OperationID, "Marshal failed ", err.Error(), tips)
		return
	}
	Notification(&n)
}

func FriendAddedNotification(operationID, opUserID, fromUserID, toUserID string) {
	var FriendAddedTips open_im_sdk.FriendAddedTips
	user, err := imdb.GetUserByUserID(opUserID)
	if err != nil {
		log.NewError(operationID, "GetUserByUserID failed ", err.Error(), opUserID)
		return
	}
	utils.CopyStructFields(FriendAddedTips.OpUser, user)

	friend, err := imdb.GetFriendRelationshipFromFriend(fromUserID, toUserID)
	if err != nil {
		log.NewError(operationID, "GetFriendRelationshipFromFriend failed ", err.Error(), fromUserID, toUserID)
		return
	}
	utils2.FriendDBCopyOpenIM(FriendAddedTips.Friend, friend)

	fromUserNickname, toUserNickname, err := getFromToUserNickname(fromUserID, toUserID)
	if err != nil {
		log.Error(operationID, "getFromToUserNickname failed ", err.Error(), fromUserID, toUserID)
		return
	}

	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendAddedTips)
	tips.DefaultTips = fromUserNickname + " FriendAddedNotification " + toUserNickname
	var n NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = constant.FriendAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

//message FriendDeletedTips{
//  FriendInfo Friend = 1;
//}
func FriendDeletedNotification(req *pbFriend.DeleteFriendReq) {

	var FriendDeletedTips open_im_sdk.FriendDeletedTips
	FriendDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	FriendDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", req.CommID.FromUserID, req.CommID.ToUserID)
		return
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendDeletedTips)
	tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.FriendDeletedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

//message FriendInfoChangedTips{
//  FriendInfo Friend = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func FriendInfoChangedNotification(operationID, opUserID, fromUserID, toUserID string) {

	var FriendInfoChangedTips open_im_sdk.FriendInfoChangedTips
	FriendInfoChangedTips.FromToUserID.FromUserID = fromUserID
	FriendInfoChangedTips.FromToUserID.ToUserID = toUserID
	fromUserNickname, toUserNickname, err := getFromToUserNickname(fromUserID, toUserID)
	if err != nil {
		log.Error(operationID, "getFromToUserNickname failed ", fromUserID, toUserID)
		return
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(FriendInfoChangedTips)
	tips.DefaultTips = fromUserNickname + " FriendDeletedNotification " + toUserNickname
	var n NotificationMsg
	n.SendID = fromUserID
	n.RecvID = toUserID
	n.ContentType = constant.FriendInfoChangedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

func BlackAddedNotification(req *pbFriend.AddBlacklistReq) {

	var BlackAddedTips open_im_sdk.BlackAddedTips
	BlackAddedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	BlackAddedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", req.CommID.FromUserID, req.CommID.ToUserID)
		return
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(BlackAddedTips)
	tips.DefaultTips = fromUserNickname + " BlackAddedNotification " + toUserNickname
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.BlackAddedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

//message BlackDeletedTips{
//  BlackInfo Black = 1;
//}
func BlackDeletedNotification(req *pbFriend.RemoveBlacklistReq) {
	var BlackDeletedTips open_im_sdk.BlackDeletedTips
	BlackDeletedTips.FromToUserID.FromUserID = req.CommID.FromUserID
	BlackDeletedTips.FromToUserID.ToUserID = req.CommID.ToUserID
	fromUserNickname, toUserNickname, err := getFromToUserNickname(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.Error(req.CommID.OperationID, "getFromToUserNickname failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
		return
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(BlackDeletedTips)
	tips.DefaultTips = fromUserNickname + " BlackDeletedNotification " + toUserNickname
	var n NotificationMsg
	n.SendID = req.CommID.FromUserID
	n.RecvID = req.CommID.ToUserID
	n.ContentType = constant.BlackDeletedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = req.CommID.OperationID
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}

//message SelfInfoUpdatedTips{
//  UserInfo SelfUserInfo = 1;
//  PublicUserInfo OpUser = 2;
//  uint64 OperationTime = 3;
//}
func SelfInfoUpdatedNotification(operationID, userID string) {
	var SelfInfoUpdatedTips open_im_sdk.SelfInfoUpdatedTips
	SelfInfoUpdatedTips.UserID = userID
	u, err := imdb.GetUserByUserID(userID)
	if err != nil {
		log.NewError(operationID, "FindUserByUID failed ", err.Error(), userID)
		return
	}
	var tips open_im_sdk.TipsComm
	tips.Detail, _ = json.Marshal(SelfInfoUpdatedTips)
	var n NotificationMsg
	n.SendID = userID
	n.RecvID = userID
	n.ContentType = constant.SelfInfoUpdatedNotification
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID

	tips.DefaultTips = u.Nickname + " SelfInfoUpdatedNotification "
	n.Content, _ = json.Marshal(tips)
	Notification(&n)
}
