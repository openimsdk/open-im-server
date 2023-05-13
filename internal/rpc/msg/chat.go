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
	"Open_IM/pkg/common/token_verify"
	pbChat "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"context"
)

func (rpc *rpcChat) ClearMsg(_ context.Context, req *pbChat.ClearMsgReq) (*pbChat.ClearMsgResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc req: ", req.String())
	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CleanUpOneUserAllMsgFromRedis args", req.UserID)
	err := db.DB.CleanUpOneUserAllMsgFromRedis(req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpOneUserAllMsgFromRedis failed " + err.Error() + req.OperationID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CleanUpUserMsgFromMongo args", req.UserID)
	err = db.DB.CleanUpUserMsgFromMongo(req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpUserMsgFromMongo failed " + err.Error() + req.OperationID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}

	resp := pbChat.ClearMsgResp{ErrCode: 0}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return &resp, nil
}

func (rpc *rpcChat) SetMsgMinSeq(_ context.Context, req *pbChat.SetMsgMinSeqReq) (*pbChat.SetMsgMinSeqResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc req: ", req.String())
	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}, nil
	}
	if req.GroupID == "" {
		err := db.DB.SetUserMinSeq(req.UserID, req.MinSeq)
		if err != nil {
			errMsg := "SetUserMinSeq failed " + err.Error() + req.OperationID + req.UserID + utils.Uint32ToString(req.MinSeq)
			log.Error(req.OperationID, errMsg)
			return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
		}
		return &pbChat.SetMsgMinSeqResp{}, nil
	}
	err := db.DB.SetGroupUserMinSeq(req.GroupID, req.UserID, uint64(req.MinSeq))
	if err != nil {
		errMsg := "SetGroupUserMinSeq failed " + err.Error() + req.OperationID + req.GroupID + req.UserID + utils.Uint32ToString(req.MinSeq)
		log.Error(req.OperationID, errMsg)
		return &pbChat.SetMsgMinSeqResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
	}
	return &pbChat.SetMsgMinSeqResp{}, nil
}
