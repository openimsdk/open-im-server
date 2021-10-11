package user

import (
	"Open_IM/internal/push/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	pbFriend "Open_IM/pkg/proto/friend"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"strings"
)

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbUser.UpdateUserInfoReq) (*pbUser.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc modify user is server,args=%s", req.String())
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbUser.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: err.Error()}, nil
	}

	ownerUid := ""
	//if claims.UID == config.Config.AppManagerUid {
	if utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		ownerUid = req.Uid
	} else {
		ownerUid = claims.UID
	}

	err = im_mysql_model.UpDateUserInfo(ownerUid, req.Name, req.Icon, req.Mobile, req.Birth, req.Email, req.Ex, req.Gender)
	if err != nil {
		log.Error(req.Token, req.OperationID, "update user some attribute failed,err=%s", err.Error())
		return &pbUser.CommonResp{ErrorCode: config.ErrModifyUserInfo.ErrCode, ErrorMsg: config.ErrModifyUserInfo.ErrMsg}, nil
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	newReq := &pbFriend.GetFriendListReq{
		OperationID: req.OperationID,
		Token:       req.Token,
	}

	RpcResp, err := client.GetFriendList(context.Background(), newReq)
	if err != nil {
		log.ErrorByKv("get friend list rpc server failed", req.OperationID, "err", err.Error(), "req", req.String())
		return &pbUser.CommonResp{}, nil
	}
	if RpcResp.ErrorCode != 0 {
		log.ErrorByKv("get friend list rpc server failed", req.OperationID, "err", err.Error(), "req", req.String())
		return &pbUser.CommonResp{}, nil
	}
	self, err := im_mysql_model.FindUserByUID(ownerUid)
	if err != nil {
		log.ErrorByKv("get self info failed", req.OperationID, "err", err.Error(), "req", req.String())
		return &pbUser.CommonResp{}, nil
	}
	var name, faceUrl string
	if self != nil {
		name, faceUrl = self.Name, self.Icon
	}
	for _, v := range RpcResp.Data {
		logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
			SendID:         ownerUid,
			RecvID:         v.Uid,
			SenderNickName: name,
			SenderFaceURL:  faceUrl,
			Content:        ownerUid + "'s info has changed",
			SendTime:       utils.GetCurrentTimestampByNano(),
			MsgFrom:        constant.SysMsgType,
			ContentType:    constant.SetSelfInfoTip,
			SessionType:    constant.SingleChatType,
			OperationID:    req.OperationID,
			Token:          req.Token,
		})

	}
	return &pbUser.CommonResp{}, nil
}
