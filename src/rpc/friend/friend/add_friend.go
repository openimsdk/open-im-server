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

func (s *friendServer) ImportFriend(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.ImportFriendResp, error) {
	log.Info(req.Token, req.OperationID, "ImportFriend come here,args=%s", req.String())
	var resp pbFriend.ImportFriendResp
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.ImportFriendResp{CommonResp: &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, FailedUidList: req.UidList}, nil
	}

	if !utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		log.Error(req.Token, req.OperationID, "not magager uid", claims.UID)
		return &pbFriend.ImportFriendResp{CommonResp: &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: "not authorized"}, FailedUidList: req.UidList}, nil
	}
	if _, err = im_mysql_model.FindUserByUID(req.OwnerUid); err != nil {
		log.Error(req.Token, req.OperationID, "this user not exists,cant not add friend", req.OwnerUid)
		return &pbFriend.ImportFriendResp{CommonResp: &pbFriend.CommonResp{ErrorCode: config.ErrAddFriend.ErrCode, ErrorMsg: "this user not exists,cant not add friend"}, FailedUidList: req.UidList}, nil
	}
	for _, v := range req.UidList {
		if _, err = im_mysql_model.FindUserByUID(v); err != nil {
			resp.CommonResp.ErrorMsg = "some uid establish failed"
			resp.CommonResp.ErrorCode = 408
			resp.FailedUidList = append(resp.FailedUidList, v)
		} else {
			if _, err = im_mysql_model.FindFriendRelationshipFromFriend(req.OwnerUid, v); err != nil {
				//Establish two single friendship
				err1 := im_mysql_model.InsertToFriend(req.OwnerUid, v, 1)
				if err1 != nil {
					resp.FailedUidList = append(resp.FailedUidList, v)
					log.Error(req.Token, req.OperationID, "err=%s,create friendship failed", err.Error())
				}
				err2 := im_mysql_model.InsertToFriend(v, req.OwnerUid, 1)
				if err2 != nil {
					log.Error(req.Token, req.OperationID, "err=%s,create friendship failed", err.Error())
				}
				if err1 == nil && err2 == nil {
					var name, faceUrl string
					n := content_struct.NotificationContent{1, constant.FriendAcceptTip, ""}
					r, err := im_mysql_model.FindUserByUID(v)
					if err != nil {
						log.ErrorByKv("get  info failed", req.OperationID, "err", err.Error(), "req", req.String())
					}
					if r != nil {
						name, faceUrl = r.Name, r.Icon
					}

					logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
						SendID:         v,
						RecvID:         req.OwnerUid,
						SenderFaceURL:  faceUrl,
						SenderNickName: name,
						Content:        n.ContentToString(),
						SendTime:       utils.GetCurrentTimestampByNano(),
						MsgFrom:        constant.UserMsgType,                //Notification message identification
						ContentType:    constant.AcceptFriendApplicationTip, //Add friend flag
						SessionType:    constant.SingleChatType,
						OperationID:    req.OperationID,
					})
				}
			}
		}

	}
	return &resp, nil

}
