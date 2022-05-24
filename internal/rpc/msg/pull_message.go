package msg

import (
	"context"
	"github.com/garyburd/redigo/redis"

	commonDB "Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/chat"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
)

func (rpc *rpcChat) GetMaxAndMinSeq(_ context.Context, in *pbMsg.GetMaxAndMinSeqReq) (*pbMsg.GetMaxAndMinSeqResp, error) {
	log.NewInfo(in.OperationID, "rpc getMaxAndMinSeq is arriving", in.String())
	//seq, err := model.GetBiggestSeqFromReceive(in.UserID)
	maxSeq, err1 := commonDB.DB.GetUserMaxSeq(in.UserID)
	minSeq, err2 := commonDB.DB.GetUserMinSeq(in.UserID)
	resp := new(pbMsg.GetMaxAndMinSeqResp)
	if err1 == nil {
		resp.MaxSeq = uint32(maxSeq)
	} else if err1 == redis.ErrNil {
		resp.MaxSeq = 0
	} else {
		log.NewError(in.OperationID, "getMaxSeq from redis error", in.String(), err1.Error())
		resp.ErrCode = 200
		resp.ErrMsg = "redis get err"
	}
	if err2 == nil {
		resp.MinSeq = uint32(minSeq)
	} else if err2 == redis.ErrNil {
		resp.MinSeq = 0
	} else {
		log.NewError(in.OperationID, "getMaxSeq from redis error", in.String(), err2.Error())
		resp.ErrCode = 201
		resp.ErrMsg = "redis get err"
	}
	return resp, nil
}
func (rpc *rpcChat) PullMessageBySeqList(_ context.Context, in *open_im_sdk.PullMessageBySeqListReq) (*open_im_sdk.PullMessageBySeqListResp, error) {
	log.NewInfo(in.OperationID, "rpc PullMessageBySeqList is arriving", in.String())
	resp := new(open_im_sdk.PullMessageBySeqListResp)
	//msgList, err := commonDB.DB.GetMsgBySeqList(in.UserID, in.SeqList, in.OperationID)
	redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(in.UserID, in.SeqList, in.OperationID)
	if err != nil {
		if err != redis.ErrNil {
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
