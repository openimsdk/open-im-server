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
	"context"

	"github.com/openimsdk/protocol/constant"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/mcontext"
)

func (m *msgServer) SetSendMsgStatus(ctx context.Context, req *pbmsg.SetSendMsgStatusReq) (*pbmsg.SetSendMsgStatusResp, error) {
	resp := &pbmsg.SetSendMsgStatusResp{}
	if err := m.MsgDatabase.SetSendMsgStatus(ctx, mcontext.GetOperationID(ctx), req.Status); err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *msgServer) GetSendMsgStatus(ctx context.Context, req *pbmsg.GetSendMsgStatusReq) (*pbmsg.GetSendMsgStatusResp, error) {
	resp := &pbmsg.GetSendMsgStatusResp{}
	status, err := m.MsgDatabase.GetSendMsgStatus(ctx, mcontext.GetOperationID(ctx))
	if IsNotFound(err) {
		resp.Status = constant.MsgStatusNotExist
		return resp, nil
	} else if err != nil {
		return nil, err
	}
	resp.Status = status
	return resp, nil
}
