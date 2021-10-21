package friend

import (
	"Open_IM/internal/push/content_struct"
	"Open_IM/internal/push/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/chat"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
)

func (s *friendServer) AddFriendResponse(ctx context.Context, req *pbFriend.AddFriendResponseReq) (*pbFriend.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc add friend response is server,args=%s", req.String())
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	//Check there application before agreeing or refuse to a friend's application
	if _, err = im_mysql_model.FindFriendApplyFromFriendReqByUid(req.Uid, claims.UID); err != nil {
		log.Error(req.Token, req.OperationID, "No such application record")
		return &pbFriend.CommonResp{ErrorCode: config.ErrAgreeToAddFriend.ErrCode, ErrorMsg: config.ErrAgreeToAddFriend.ErrMsg}, nil
	}
	//Change friend request status flag
	err = im_mysql_model.UpdateFriendRelationshipToFriendReq(req.Uid, claims.UID, req.Flag)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,update friend request table failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrMysql.ErrCode, ErrorMsg: config.ErrAgreeToAddFriend.ErrMsg}, nil
	}
	log.Info(req.Token, req.OperationID, "rpc add friend response success return,userid=%s,flag=%d", req.Uid, req.Flag)
	//Change the status of the friend request form
	if req.Flag == constant.FriendFlag {
		//Establish friendship after find friend relationship not exists
		_, err := im_mysql_model.FindFriendRelationshipFromFriend(claims.UID, req.Uid)
		//fixme If there is an error, it means that there is no friend record or database err, if no friend record should be inserted,Continue down execution
		if err != nil {
			log.Error("", req.OperationID, err.Error())
		}
		//Establish two single friendship
		err = im_mysql_model.InsertToFriend(claims.UID, req.Uid, req.Flag)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s,create friendship failed", err.Error())
		}
		err = im_mysql_model.InsertToFriend(req.Uid, claims.UID, req.Flag)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s,create friendship failed", err.Error())
		}
		//Push message when establish friends successfully
		//senderInfo, errSend := im_mysql_model.FindUserByUID(claims.UID)
		//if errSend == nil {
		//	logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
		//		SendID:      claims.UID,
		//		RecvID:      req.Uid,
		//		Content:     content_struct.NewContentStructString(1, "", senderInfo.Name+" agreed to add you as a friend."),
		//		SendTime:    utils.GetCurrentTimestampBySecond(),
		//		MsgFrom:     constant.UserMsgType,                //Notification message identification
		//		ContentType: constant.AcceptFriendApplicationTip, //Add friend flag
		//		SessionType: constant.SingleChatType,
		//		OperationID: req.OperationID,
		//	})
		//}
	}
	if req.Flag == constant.RefuseFriendFlag {
		senderInfo, errSend := im_mysql_model.FindUserByUID(claims.UID)
		if errSend == nil {
			logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
				SendID:      claims.UID,
				RecvID:      req.Uid,
				Content:     content_struct.NewContentStructString(0, "", senderInfo.Name+" refuse to add you as a friend."),
				SendTime:    utils.GetCurrentTimestampBySecond(),
				MsgFrom:     constant.UserMsgType,                //Notification message identification
				ContentType: constant.RefuseFriendApplicationTip, //Add friend flag
				SessionType: constant.SingleChatType,
				OperationID: req.OperationID,
			})
		}
	}
	return &pbFriend.CommonResp{}, nil
}
