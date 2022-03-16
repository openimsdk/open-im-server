package msg

import (
	"Open_IM/pkg/common/log"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
)

func (rpc *rpcChat) DelMsgList(_ context.Context, req *commonPb.DelMsgListReq) (*commonPb.DelMsgListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &commonPb.DelMsgListResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	return resp, nil
}
