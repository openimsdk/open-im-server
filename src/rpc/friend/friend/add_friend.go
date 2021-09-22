package friend

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbChat "Open_IM/src/proto/chat"
	pbFriend "Open_IM/src/proto/friend"
	"Open_IM/src/push/content_struct"
	"Open_IM/src/push/logic"
	"Open_IM/src/utils"
	"context"
)

func (s *friendServer) AddFriend(ctx context.Context, req *pbFriend.AddFriendReq) (*pbFriend.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc add friend is server,userid=%s", req.Uid)
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	//Cannot add non-existent users
	if _, err = im_mysql_model.FindUserByUID(req.Uid); err != nil {
		log.Error(req.Token, req.OperationID, "this user not exists,cant not add friend")
		return &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: config.ErrSearchUserInfo.ErrMsg}, nil
	}

	//Establish a latest relationship in the friend request table
	err = im_mysql_model.ReplaceIntoFriendReq(claims.UID, req.Uid, constant.ApplicationFriendFlag, req.ReqMessage)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,create friend request record failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: config.ErrAddFriend.ErrMsg}, nil
	}
	log.Info(req.Token, req.OperationID, "rpc add friend  is success return,uid=%s", req.Uid)
	//Push message when add friend successfully
	senderInfo, errSend := im_mysql_model.FindUserByUID(claims.UID)
	receiverInfo, errReceive := im_mysql_model.FindUserByUID(req.Uid)
	if errSend == nil && errReceive == nil {
		logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
			SendID:      senderInfo.UID,
			RecvID:      receiverInfo.UID,
			Content:     content_struct.NewContentStructString(0, "", senderInfo.Name+" asked to add you as a friend"),
			SendTime:    utils.GetCurrentTimestampBySecond(),
			MsgFrom:     constant.SysMsgType,
			ContentType: constant.AddFriendTip,
			SessionType: constant.SingleChatType,
			OperationID: req.OperationID,
		})
	}
	return &pbFriend.CommonResp{}, nil
}

func (s *friendServer) ImportFriend(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "ImportFriendis server,userid=%s", req.OwnerUid)
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	if claims.UID != config.Config.AppManagerUid {
		log.Error(req.Token, req.OperationID, "not magager uid", claims.UID, config.Config.AppManagerUid)
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}

	if _, err = im_mysql_model.FindUserByUID(req.Uid); err != nil {
		log.Error(req.Token, req.OperationID, "this user not exists,cant not add friend", req.Uid)
		return &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: config.ErrSearchUserInfo.ErrMsg}, nil
	}

	if _, err = im_mysql_model.FindUserByUID(req.OwnerUid); err != nil {
		log.Error(req.Token, req.OperationID, "this user not exists,cant not add friend", req.OwnerUid)
		return &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: config.ErrSearchUserInfo.ErrMsg}, nil
	}

	_, err = im_mysql_model.FindFriendRelationshipFromFriend(req.OwnerUid, req.Uid)

	if err != nil {
		log.Error("", req.OperationID, err.Error())
	}
	//Establish two single friendship
	err = im_mysql_model.InsertToFriend(req.OwnerUid, req.Uid, 1)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,create friendship failed", err.Error())
	}
	err = im_mysql_model.InsertToFriend(req.Uid, req.OwnerUid, 1)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,create friendship failed", err.Error())
	}

	logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
		SendID:      req.OwnerUid,
		RecvID:      req.Uid,
		Content:     content_struct.NewContentStructString(0, "", " add you as a friend."),
		SendTime:    utils.GetCurrentTimestampBySecond(),
		MsgFrom:     constant.UserMsgType,                //Notification message identification
		ContentType: constant.AcceptFriendApplicationTip, //Add friend flag
		SessionType: constant.SingleChatType,
		OperationID: req.OperationID,
	})

	logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
		SendID:      req.Uid,
		RecvID:      req.OwnerUid,
		Content:     content_struct.NewContentStructString(0, "", " add you as a friend."),
		SendTime:    utils.GetCurrentTimestampBySecond(),
		MsgFrom:     constant.UserMsgType,                //Notification message identification
		ContentType: constant.AcceptFriendApplicationTip, //Add friend flag
		SessionType: constant.SingleChatType,
		OperationID: req.OperationID,
	})

	return &pbFriend.CommonResp{}, nil
}
