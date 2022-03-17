package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
)

func (rpc *rpcChat) DelMsgList(_ context.Context, req *commonPb.DelMsgListReq) (*commonPb.DelMsgListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &commonPb.DelMsgListResp{}
	if err := db.DB.DelMsgLogic(req.UserID, req.SeqList, req.OperationID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DelMsg failed", err.Error())
		resp.ErrMsg = constant.ErrDB.ErrMsg
		resp.ErrCode = constant.ErrDB.ErrCode
		return resp, err
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}
