package msg

import (
	commonDB "Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"context"
	go_redis "github.com/go-redis/redis/v8"
)

func (rpc *rpcChat) GetSuperGroupMsg(context context.Context, req *msg.GetSuperGroupMsgReq) (*msg.GetSuperGroupMsgResp, error) {
	log.Debug(req.OperationID, utils.GetSelfFuncName(), req.String())
	resp := new(msg.GetSuperGroupMsgResp)
	redisMsgList, failedSeqList, err := commonDB.DB.GetMessageListBySeq(req.GroupID, []uint32{req.Seq}, req.OperationID)
	if err != nil {
		if err != go_redis.Nil {
			promePkg.PromeAdd(promePkg.MsgPullFromRedisFailedCounter, len(failedSeqList))
			log.Error(req.OperationID, "get message from redis exception", err.Error(), failedSeqList)
		} else {
			log.Debug(req.OperationID, "get message from redis is nil", failedSeqList)
		}
		msgList, _, err1 := commonDB.DB.GetSuperGroupMsgBySeqListMongo(req.GroupID, failedSeqList, req.OperationID)
		if err1 != nil {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoFailedCounter, len(failedSeqList))
			log.Error(req.OperationID, "GetSuperGroupMsg data error", req.String(), err.Error())
			resp.ErrCode = 201
			resp.ErrMsg = err.Error()
			return resp, nil
		} else {
			promePkg.PromeAdd(promePkg.MsgPullFromMongoSuccessCounter, len(msgList))
			redisMsgList = append(redisMsgList, msgList...)
			for _, m := range msgList {
				resp.MsgData = m
			}

		}
	} else {
		promePkg.PromeAdd(promePkg.MsgPullFromRedisSuccessCounter, len(redisMsgList))
		for _, m := range redisMsgList {
			resp.MsgData = m
		}
	}
	log.Debug(req.OperationID, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}

func (rpc *rpcChat) GetWriteDiffMsg(context context.Context, req *msg.GetWriteDiffMsgReq) (*msg.GetWriteDiffMsgResp, error) {
	panic("implement me")
}
