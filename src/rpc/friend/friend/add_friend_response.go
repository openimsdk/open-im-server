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

func (s *friendServer) AddedFriend(ctx context.Context, req *pbFriend.AddedFriendReq) (*pbFriend.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc add friend response is server,args=%s", req.String())
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	err = im_mysql_model.UpdateFriendRelationshipToFriendReq(req.Uid, claims.UID, req.Flag)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,update friend request table failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrMysql.ErrCode, ErrorMsg: config.ErrMysql.ErrMsg}, nil
	}
	log.Info(req.Token, req.OperationID, "rpc add friend response success return,userid=%s,flag=%d", req.Uid, req.Flag)
	//Change the status of the friend request form
	if req.Flag == constant.FriendFlag {
		//Establish friendship after find friend relationship not exists
		_, err := im_mysql_model.FindFriendRelationshipFromFriend(claims.UID, req.Uid)
		if err == nil {
			return &pbFriend.CommonResp{ErrorCode: 0, ErrorMsg: "You are already friends"}, nil
		}
		//Establish two single friendship
		err = im_mysql_model.InsertToFriend(claims.UID, req.Uid, req.Flag)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s,create friendship failed", err.Error())
			return &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: config.ErrAddFriend.ErrMsg}, nil
		}
		err = im_mysql_model.InsertToFriend(req.Uid, claims.UID, req.Flag)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s,create friendship failed", err.Error())
			return &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: config.ErrAddFriend.ErrMsg}, nil
		}
		//Push message when establish friends successfully
		senderInfo, errSend := im_mysql_model.FindUserByUID(claims.UID)
		if errSend == nil {
			logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
				SendID:      claims.UID,
				RecvID:      req.Uid,
				Content:     content_struct.NewContentStructString(0, "", senderInfo.Name+" agreed to add you as a friend."),
				SendTime:    utils.GetCurrentTimestampBySecond(),
				MsgFrom:     constant.SysMsgType,        //Notification message identification
				ContentType: constant.AgreeAddFriendTip, //Add friend flag
				SessionType: constant.SingleChatType,
				OperationID: req.OperationID,
			})
		}
	}
	return &pbFriend.CommonResp{}, nil
}
