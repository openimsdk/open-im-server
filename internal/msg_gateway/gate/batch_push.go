package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/msg"
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

func (r *RPCServer) GetSingleUserMsgForPushPlatforms(operationID string, msgData *sdk_ws.MsgData, pushToUserID string, platformIDList []int) map[int]*sdk_ws.MsgDataList {
	user2PushMsg := make(map[int]*sdk_ws.MsgDataList, 0)
	for _, v := range platformIDList {
		user2PushMsg[v] = r.GetSingleUserMsgForPush(operationID, msgData, pushToUserID, v)
		//log.Info(operationID, "GetSingleUserMsgForPush", msgData.Seq, pushToUserID, v, "len:", len(user2PushMsg[v]))
	}
	return user2PushMsg
}

func (r *RPCServer) GetSingleUserMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserID string, platformID int) *sdk_ws.MsgDataList {
	//msgData.MsgDataList = nil
	return &sdk_ws.MsgDataList{MsgDataList: []*sdk_ws.MsgData{msgData}}

	//userConn := ws.getUserConn(pushToUserID, platformID)
	//if userConn == nil {
	//	log.Debug(operationID, "userConn == nil")
	//	return []*sdk_ws.MsgData{msgData}
	//}
	//
	//if msgData.Seq <= userConn.PushedMaxSeq {
	//	log.Debug(operationID, "msgData.Seq <= userConn.PushedMaxSeq", msgData.Seq, userConn.PushedMaxSeq)
	//	return nil
	//}
	//
	//msgList := r.GetSingleUserMsg(operationID, msgData.Seq, pushToUserID)
	//if msgList == nil {
	//	log.Debug(operationID, "GetSingleUserMsg msgList == nil", msgData.Seq, userConn.PushedMaxSeq)
	//	userConn.PushedMaxSeq = msgData.Seq
	//	return []*sdk_ws.MsgData{msgData}
	//}
	//msgList = append(msgList, msgData)
	//
	//for _, v := range msgList {
	//	if v.Seq > userConn.PushedMaxSeq {
	//		userConn.PushedMaxSeq = v.Seq
	//	}
	//}
	//log.Debug(operationID, "GetSingleUserMsg msgList len ", len(msgList), userConn.PushedMaxSeq)
	//return msgList
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
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, rpcReq.OperationID)
	if grpcConn == nil {
		errMsg := "getcdv3.GetConn == nil"
		log.NewError(rpcReq.OperationID, errMsg)
		return nil
	}

	msgClient := pbChat.NewMsgClient(grpcConn)
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

//func (r *RPCServer) GetBatchUserMsgForPush(operationID string, msgData *sdk_ws.MsgData, pushToUserIDList []string, platformID int) map[string][]*sdk_ws.MsgData {
//	user2PushMsg := make(map[string][]*sdk_ws.MsgData, 0)
//	for _, v := range pushToUserIDList {
//		user2PushMsg[v] = r.GetSingleUserMsgForPush(operationID, msgData, v, platformID)
//	}
//	return user2PushMsg
//}
