// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func (rpc *rpcChat) SetSendMsgStatus(_ context.Context, req *pbMsg.SetSendMsgStatusReq) (resp *pbMsg.SetSendMsgStatusResp, err error) {
	resp = &pbMsg.SetSendMsgStatusResp{}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.String())
	if err := db.DB.SetSendMsgStatus(req.Status, req.OperationID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
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
	status, err := db.DB.GetSendMsgStatus(req.OperationID)
	if err != nil {
		resp.Status = constant.MsgStatusNotExist
		if err == goRedis.Nil {
			log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.OperationID, "not exist")
			return resp, nil
		} else {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			resp.ErrMsg = err.Error()
			resp.ErrCode = constant.ErrDB.ErrCode
			return resp, nil
		}
	}
	resp.Status = int32(status)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}
