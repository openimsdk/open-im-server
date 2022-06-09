package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"context"
	"strings"
)

var MaxPullMsgNum = 100

func (r *RPCServer) GenPullSeqList(currentSeq uint32, operationID string, userID string) ([]uint32, error) {
	return nil, nil
}
func (r *RPCServer) GetMergeSingleMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserID string, platformID int) []*sdk_ws.MsgData {
	return nil
	//ws.getUserConn(pushToUserID, platformID)
	//msgData.Seq
	//msgList := r.GetSingleMsgForPush(operationID, msgData, pushToUserID, platformID)

}
func (r *RPCServer) GetSingleMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserID string, platformID string) []*sdk_ws.MsgData {
	seqList, err := r.GenPullSeqList(msgData.Seq, operationID, pushToUserID)
	if err != nil {
		log.Error(operationID, "GenPullSeqList failed ", err.Error(), msgData.Seq, pushToUserID)
		return nil
	}
	rpcReq := sdk_ws.PullMessageBySeqListReq{}
	rpcReq.SeqList = seqList
	rpcReq.UserID = pushToUserID
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

func (r *RPCServer) GetBatchMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserIDList []string, platformID string) map[string][]*sdk_ws.MsgData {
	return nil
}

func (r *RPCServer) GetMaxSeq(userID string) (uint32, error) {
	return 0, nil
}
