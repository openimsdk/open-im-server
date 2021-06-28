package user

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbChat "Open_IM/src/proto/chat"
	pbFriend "Open_IM/src/proto/friend"
	pbUser "Open_IM/src/proto/user"
	"Open_IM/src/push/logic"
	"Open_IM/src/utils"
	"context"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"strings"
)

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbUser.UpdateUserInfoReq) (*pbUser.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc modify user is server,args=%s", req.String())
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbUser.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: err.Error()}, nil
	}
	err = im_mysql_model.UpDateUserInfo(claims.UID, req.Name, req.Icon, req.Mobile, req.Birth, req.Email, req.Ex, req.Gender)
	if err != nil {
		log.Error(req.Token, req.OperationID, "update user some attribute failed,err=%s", err.Error())
		return &pbUser.CommonResp{ErrorCode: config.ErrModifyUserInfo.ErrCode, ErrorMsg: config.ErrModifyUserInfo.ErrMsg}, nil
	}
	go func() {
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
		client := pbFriend.NewFriendClient(etcdConn)
		defer etcdConn.Close()
		newReq := &pbFriend.GetFriendListReq{
			OperationID: req.OperationID,
			Token:       req.Token,
		}

		RpcResp, err := client.GetFriendList(context.Background(), newReq)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s,call get friend list rpc server failed", err)
			log.ErrorByKv("get friend list rpc server failed", req.OperationID, "err", err.Error(), "req", req.String())
		}
		if RpcResp.ErrorCode != 0 {
			log.ErrorByKv("get friend list rpc server failed", req.OperationID, "err", err.Error(), "req", req.String())
		}
		for _, v := range RpcResp.Data {
			logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
				SendID:      claims.UID,
				RecvID:      v.Uid,
				Content:     claims.UID + "'s info has changed",
				SendTime:    utils.GetCurrentTimestampBySecond(),
				MsgFrom:     constant.SysMsgType,
				ContentType: constant.SetSelfInfoTip,
				SessionType: constant.SingleChatType,
				OperationID: req.OperationID,
				Token:       req.Token,
			})

		}
	}()
	return &pbUser.CommonResp{}, nil
}
