package msggateway

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	pbChat "OpenIM/pkg/proto/msg"
	sdkws "OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
	"google.golang.org/grpc"
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

func (r *RPCServer) GetSingleUserMsgForPushPlatforms(operationID string, msgData *sdkws.MsgData, pushToUserID string, platformIDList []int) map[int]*sdkws.MsgDataList {
	user2PushMsg := make(map[int]*sdkws.MsgDataList, 0)
	for _, v := range platformIDList {
		user2PushMsg[v] = r.GetSingleUserMsgForPush(operationID, msgData, pushToUserID, v)
		//log.Info(operationID, "GetSingleUserMsgForPush", msgData.Seq, pushToUserID, v, "len:", len(user2PushMsg[v]))
	}
	return user2PushMsg
}

func (r *RPCServer) GetSingleUserMsgForPush(operationID string, msgData *sdkws.MsgData, pushToUserID string, platformID int) *sdkws.MsgDataList {
	//msgData.MsgDataList = nil
	return &sdkws.MsgDataList{MsgDataList: []*sdkws.MsgData{msgData}}

	//userConn := ws.getUserConn(pushToUserID, platformID)
	//if userConn == nil {
	//	log.Debug(operationID, "userConn == nil")
	//	return []*sdkws.MsgData{msgData}
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
	//	return []*sdkws.MsgData{msgData}
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

func (r *RPCServer) GetSingleUserMsg(operationID string, currentMsgSeq uint32, userID string) []*sdkws.MsgData {
	seqs, err := r.GenPullSeqList(currentMsgSeq, operationID, userID)
	if err != nil {
		log.Error(operationID, "GenPullSeqList failed ", err.Error(), currentMsgSeq, userID)
		return nil
	}
	if len(seqs) == 0 {
		log.Error(operationID, "GenPullSeqList len == 0 ", currentMsgSeq, userID)
		return nil
	}
	rpcReq := sdkws.PullMessageBySeqsReq{}
	//rpcReq.Seqs = seqs
	rpcReq.UserID = userID
	var grpcConn *grpc.ClientConn

	msgClient := pbChat.NewMsgClient(grpcConn)
	reply, err := msgClient.PullMessageBySeqs(context.Background(), &rpcReq)
	if err != nil {
		log.Error(operationID, "PullMessageBySeqList failed ", err.Error(), rpcReq.String())
		return nil
	}
	if len(reply.List) == 0 {
		return nil
	}
	return reply.List
}

//func (r *RPCServer) GetBatchUserMsgForPush(operationID string, msgData *sdkws.MsgData, pushToUserIDList []string, platformID int) map[string][]*sdkws.MsgData {
//	user2PushMsg := make(map[string][]*sdkws.MsgData, 0)
//	for _, v := range pushToUserIDList {
//		user2PushMsg[v] = r.GetSingleUserMsgForPush(operationID, msgData, v, platformID)
//	}
//	return user2PushMsg
//}
