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
	//Establish a latest relationship in the friend request table
	err = im_mysql_model.ReplaceIntoFriendReq(claims.UID, req.Uid, constant.NotFriendFlag, req.ReqMessage)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,create friend request ship failed", err.Error())
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
