package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"strings"
)

var MaxPullMsgNum = 100

func (r *RPCServer) GenPullSeqList(currentSeq uint32, operationID string, userID string) ([]uint32, error) {
	maxSeq, err := db.DB.GetUserMaxSeq(userID)
	if err != nil {
		log.Error(operationID, "GetUserMaxSeq failed ", userID, err.Error())
		return nil, utils.Wrap(err, "")
	}

	var seqList []uint32
	num := 0
	for i := currentSeq + 1; i < uint32(maxSeq); i++ {
		seqList = append(seqList, i)
		num++
		if num == MaxPullMsgNum {
			break
		}
	}
	log.Info(operationID, "GenPullSeqList ", seqList, "current seq", currentSeq)
	return seqList, nil
}
func (r *RPCServer) GetSingleUserMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserID string, platformID int) []*sdk_ws.MsgData {
	userConn := ws.getUserConn(pushToUserID, platformID)
	if userConn == nil {
		return []*sdk_ws.MsgData{msgData}
	}

	if msgData.Seq <= userConn.PushedMaxSeq {
		return nil
	}

	msgList := r.GetSingleUserMsg(operationID, msgData.Seq, pushToUserID)
	if msgList == nil {
		userConn.PushedMaxSeq = msgData.Seq
		return []*sdk_ws.MsgData{msgData}
	}
	msgList = append(msgList, msgData)

	for _, v := range msgList {
		if v.Seq > userConn.PushedMaxSeq {
			userConn.PushedMaxSeq = v.Seq
		}
	}
	return msgList
}

func (r *RPCServer) GetSingleUserMsg(operationID string, currentMsgSeq uint32, userID string) []*sdk_ws.MsgData {
	seqList, err := r.GenPullSeqList(currentMsgSeq, operationID, userID)
	if err != nil {
		log.Error(operationID, "GenPullSeqList failed ", err.Error(), currentMsgSeq, userID)
		return nil
	}
	if len(seqList) == 0 {
		log.Error(operationID, "GenPullSeqList len == 0 ", currentMsgSeq, userID)
		return nil
	}
	rpcReq := sdk_ws.PullMessageBySeqListReq{}
	rpcReq.SeqList = seqList
	rpcReq.UserID = userID
	rpcReq.OperationID = operationID
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	msgClient := pbChat.NewChatClient(grpcConn)
	reply, err := msgClient.PullMessageBySeqList(context.Background(), &rpcReq)
	if err != nil {
		log.Error(operationID, "PullMessageBySeqList failed ", err.Error(), rpcReq.String())
		return nil
	}
	if len(reply.List) == 0 {
		return nil
	}
	return reply.List
}

func (r *RPCServer) GetBatchUserMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserIDList []string, platformID int) map[string][]*sdk_ws.MsgData {
	user2PushMsg := make(map[string][]*sdk_ws.MsgData, 0)
	for _, v := range pushToUserIDList {
		user2PushMsg[v] = r.GetSingleUserMsgForPush(operationID, msgData, v, platformID)
	}
	return user2PushMsg
}
