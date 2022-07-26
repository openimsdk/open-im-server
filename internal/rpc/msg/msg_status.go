package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"context"
	goRedis "github.com/go-redis/redis/v8"
)

func (rpc *rpcChat) SetSendMsgFailedFlag(_ context.Context, req *pbMsg.SetSendMsgFailedFlagReq) (resp *pbMsg.SetSendMsgFailedFlagResp, err error) {
	resp = &pbMsg.SetSendMsgFailedFlagResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.String())
	if err := db.DB.SetSendMsgFailedFlag(req.OperationID); err != nil {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}

func (rpc *rpcChat) GetSendMsgStatus(_ context.Context, req *pbMsg.GetSendMsgStatusReq) (resp *pbMsg.GetSendMsgStatusResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.String())
	resp = &pbMsg.GetSendMsgStatusResp{}
	if err := db.DB.GetSendMsgStatus(req.OperationID); err != nil {
		if err == goRedis.Nil {
			resp.Status = 0
			return resp, nil
		} else {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			resp.ErrMsg = err.Error()
			resp.ErrCode = constant.ErrDB.ErrCode
			return resp, nil
		}
	}
	resp.Status = 1
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}
