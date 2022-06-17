package msg

import (
	"context"
	go_redis "github.com/go-redis/redis/v8"

	commonDB "Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
)

func (rpc *rpcChat) GetMaxAndMinSeq(_ context.Context, in *open_im_sdk.GetMaxAndMinSeqReq) (*open_im_sdk.GetMaxAndMinSeqResp, error) {
	log.NewInfo(in.OperationID, "rpc getMaxAndMinSeq is arriving", in.String())
	resp := new(open_im_sdk.GetMaxAndMinSeqResp)
	m := make(map[string]*open_im_sdk.MaxAndMinSeq)
	//seq, err := model.GetBiggestSeqFromReceive(in.UserID)
	maxSeq, err1 := commonDB.DB.GetUserMaxSeq(in.UserID)
	//minSeq, err2 := commonDB.DB.GetUserMinSeq(in.UserID)
	if err1 == nil {
		resp.MaxSeq = uint32(maxSeq)
		for _, v := range in.GroupIDList {
			x := new(open_im_sdk.MaxAndMinSeq)
			maxSeq, _ := commonDB.DB.GetUserMaxSeq(v)
			x.MaxSeq = uint32(maxSeq)
			m[v] = x
		}
		resp.GroupMaxAndMinSeq = m
	} else if err1 == go_redis.Nil {
		resp.MaxSeq = 0
	} else {
		log.NewError(in.OperationID, "getMaxSeq from redis error", in.String(), err1.Error())
		resp.ErrCode = 200
		resp.ErrMsg = "redis get err"
	}
	//if err2 == nil {
	//	resp.MinSeq = uint32(minSeq)
	//} else if err2 == redis.ErrNil {
	//	resp.MinSeq = 0
	//} else {
	//	log.NewError(in.OperationID, "getMaxSeq from redis error", in.String(), err2.Error())
	//	resp.ErrCode = 201
	//	resp.ErrMsg = "redis get err"
	//}
	return resp, nil
}
func (rpc *rpcChat) PullMessageBySeqList(_ context.Context, in *open_im_sdk.PullMessageBySeqListReq) (*open_im_sdk.PullMessageBySeqListResp, error) {
	log.NewInfo(in.OperationID, "rpc PullMessageBySeqList is arriving", in.String())
	resp := new(open_im_sdk.PullMessageBySeqListResp)
	m := make(map[string]*open_im_sdk.MsgDataList)
	//msgList, err := commonDB.DB.GetMsgBySeqList(in.UserID, in.SeqList, in.OperationID)
	redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(in.UserID, in.SeqList, in.OperationID)
	if err != nil {
		if err != go_redis.Nil {
			log.Error(in.OperationID, "get message from redis exception", err.Error(), failedSeqList)
		} else {
			log.Debug(in.OperationID, "get message from redis is nil", failedSeqList)
		}
		msgList, err1 := commonDB.DB.GetMsgBySeqListMongo2(in.UserID, failedSeqList, in.OperationID)
		if err1 != nil {
			log.Error(in.OperationID, "PullMessageBySeqList data error", in.String(), err.Error())
			resp.ErrCode = 201
			resp.ErrMsg = err.Error()
			return resp, nil
		} else {
			redisMsgList = append(redisMsgList, msgList...)
			resp.List = redisMsgList
		}
	} else {
		resp.List = redisMsgList
	}
	for k, v := range in.GroupSeqList {
		x := new(open_im_sdk.MsgDataList)
		redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(k, v.SeqList, in.OperationID)
		if err != nil {
			if err != go_redis.Nil {
				log.Error(in.OperationID, "get message from redis exception", err.Error(), failedSeqList)
			} else {
				log.Debug(in.OperationID, "get message from redis is nil", failedSeqList)
			}
			msgList, err1 := commonDB.DB.GetSuperGroupMsgBySeqListMongo(k, failedSeqList, in.OperationID)
			if err1 != nil {
				log.Error(in.OperationID, "PullMessageBySeqList data error", in.String(), err.Error())
				resp.ErrCode = 201
				resp.ErrMsg = err.Error()
				return resp, nil
			} else {
				redisMsgList = append(redisMsgList, msgList...)
				x.MsgDataList = redisMsgList
				m[k] = x
			}
		} else {
			x.MsgDataList = redisMsgList
			m[k] = x
		}
	}
	resp.GroupMsgDataList = m
	//respSingleMsgFormat = singleMsgHandleByUser(SingleMsgFormat, in.UserID)
	//respGroupMsgFormat = groupMsgHandleByUser(GroupMsgFormat)
	return resp, nil

}

type MsgFormats []*open_im_sdk.MsgData

// Implement the sort.Interface interface to get the number of elements method
func (s MsgFormats) Len() int {
	return len(s)
}

//Implement the sort.Interface interface comparison element method
func (s MsgFormats) Less(i, j int) bool {
	return s[i].SendTime < s[j].SendTime
}

//Implement the sort.Interface interface exchange element method
func (s MsgFormats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
