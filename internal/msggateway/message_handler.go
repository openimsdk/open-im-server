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

package msggateway

import (
	"context"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/push"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/jsonutil"
	"google.golang.org/protobuf/proto"
)

type Req struct {
	ReqIdentifier int32  `json:"reqIdentifier" validate:"required"`
	Token         string `json:"token"`
	SendID        string `json:"sendID"        validate:"required"`
	OperationID   string `json:"operationID"   validate:"required"`
	MsgIncr       string `json:"msgIncr"       validate:"required"`
	Data          []byte `json:"data"`
}

func (r *Req) String() string {
	var tReq Req
	tReq.ReqIdentifier = r.ReqIdentifier
	tReq.Token = r.Token
	tReq.SendID = r.SendID
	tReq.OperationID = r.OperationID
	tReq.MsgIncr = r.MsgIncr
	return jsonutil.StructToJsonString(tReq)
}

var reqPool = sync.Pool{
	New: func() any {
		return new(Req)
	},
}

func getReq() *Req {
	req := reqPool.Get().(*Req)
	req.Data = nil
	req.MsgIncr = ""
	req.OperationID = ""
	req.ReqIdentifier = 0
	req.SendID = ""
	req.Token = ""
	return req
}

func freeReq(req *Req) {
	reqPool.Put(req)
}

type Resp struct {
	ReqIdentifier int32  `json:"reqIdentifier"`
	MsgIncr       string `json:"msgIncr"`
	OperationID   string `json:"operationID"`
	ErrCode       int    `json:"errCode"`
	ErrMsg        string `json:"errMsg"`
	Data          []byte `json:"data"`
}

func (r *Resp) String() string {
	var tResp Resp
	tResp.ReqIdentifier = r.ReqIdentifier
	tResp.MsgIncr = r.MsgIncr
	tResp.OperationID = r.OperationID
	tResp.ErrCode = r.ErrCode
	tResp.ErrMsg = r.ErrMsg
	return jsonutil.StructToJsonString(tResp)
}

type MessageHandler interface {
	GetSeq(context context.Context, data *Req) ([]byte, error)
	SendMessage(context context.Context, data *Req) ([]byte, error)
	SendSignalMessage(context context.Context, data *Req) ([]byte, error)
	PullMessageBySeqList(context context.Context, data *Req) ([]byte, error)
	UserLogout(context context.Context, data *Req) ([]byte, error)
	SetUserDeviceBackground(context context.Context, data *Req) ([]byte, bool, error)
}

var _ MessageHandler = (*GrpcHandler)(nil)

type GrpcHandler struct {
	msgRpcClient *rpcclient.MessageRpcClient
	pushClient   *rpcclient.PushRpcClient
	validate     *validator.Validate
}

func NewGrpcHandler(validate *validator.Validate, client discovery.SvcDiscoveryRegistry, rpcRegisterName *config.RpcRegisterName) *GrpcHandler {
	msgRpcClient := rpcclient.NewMessageRpcClient(client, rpcRegisterName.Msg)
	pushRpcClient := rpcclient.NewPushRpcClient(client, rpcRegisterName.Push)
	return &GrpcHandler{
		msgRpcClient: &msgRpcClient,
		pushClient:   &pushRpcClient, validate: validate,
	}
}

func (g GrpcHandler) GetSeq(ctx context.Context, data *Req) ([]byte, error) {
	req := sdkws.GetMaxSeqReq{}
	if err := proto.Unmarshal(data.Data, &req); err != nil {
		return nil, errs.WrapMsg(err, "GetSeq: error unmarshaling request", "action", "unmarshal", "dataType", "GetMaxSeqReq")
	}
	if err := g.validate.Struct(&req); err != nil {
		return nil, errs.WrapMsg(err, "GetSeq: validation failed", "action", "validate", "dataType", "GetMaxSeqReq")
	}
	resp, err := g.msgRpcClient.GetMaxSeq(ctx, &req)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, errs.WrapMsg(err, "GetSeq: error marshaling response", "action", "marshal", "dataType", "GetMaxSeqResp")
	}
	return c, nil
}

// SendMessage handles the sending of messages through gRPC. It unmarshals the request data,
// validates the message, and then sends it using the message RPC client.
func (g GrpcHandler) SendMessage(ctx context.Context, data *Req) ([]byte, error) {
	var msgData sdkws.MsgData
	if err := proto.Unmarshal(data.Data, &msgData); err != nil {
		return nil, errs.WrapMsg(err, "SendMessage: error unmarshaling message data", "action", "unmarshal", "dataType", "MsgData")
	}

	if err := g.validate.Struct(&msgData); err != nil {
		return nil, errs.WrapMsg(err, "SendMessage: message data validation failed", "action", "validate", "dataType", "MsgData")
	}

	req := msg.SendMsgReq{MsgData: &msgData}
	resp, err := g.msgRpcClient.SendMsg(ctx, &req)
	if err != nil {
		return nil, err
	}

	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, errs.WrapMsg(err, "SendMessage: error marshaling response", "action", "marshal", "dataType", "SendMsgResp")
	}

	return c, nil
}

func (g GrpcHandler) SendSignalMessage(context context.Context, data *Req) ([]byte, error) {
	resp, err := g.msgRpcClient.SendMsg(context, nil)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, errs.WrapMsg(err, "error marshaling response", "action", "marshal", "dataType", "SendMsgResp")
	}
	return c, nil
}

func (g GrpcHandler) PullMessageBySeqList(context context.Context, data *Req) ([]byte, error) {
	req := sdkws.PullMessageBySeqsReq{}
	if err := proto.Unmarshal(data.Data, &req); err != nil {
		return nil, errs.WrapMsg(err, "error unmarshaling request", "action", "unmarshal", "dataType", "PullMessageBySeqsReq")
	}
	if err := g.validate.Struct(data); err != nil {
		return nil, errs.WrapMsg(err, "validation failed", "action", "validate", "dataType", "PullMessageBySeqsReq")
	}
	resp, err := g.msgRpcClient.PullMessageBySeqList(context, &req)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, errs.WrapMsg(err, "error marshaling response", "action", "marshal", "dataType", "PullMessageBySeqsResp")
	}
	return c, nil
}

func (g GrpcHandler) UserLogout(context context.Context, data *Req) ([]byte, error) {
	req := push.DelUserPushTokenReq{}
	if err := proto.Unmarshal(data.Data, &req); err != nil {
		return nil, errs.WrapMsg(err, "error unmarshaling request", "action", "unmarshal", "dataType", "DelUserPushTokenReq")
	}
	resp, err := g.pushClient.DelUserPushToken(context, &req)
	if err != nil {
		return nil, err
	}
	c, err := proto.Marshal(resp)
	if err != nil {
		return nil, errs.WrapMsg(err, "error marshaling response", "action", "marshal", "dataType", "DelUserPushTokenResp")
	}
	return c, nil
}

func (g GrpcHandler) SetUserDeviceBackground(_ context.Context, data *Req) ([]byte, bool, error) {
	req := sdkws.SetAppBackgroundStatusReq{}
	if err := proto.Unmarshal(data.Data, &req); err != nil {
		return nil, false, errs.WrapMsg(err, "error unmarshaling request", "action", "unmarshal", "dataType", "SetAppBackgroundStatusReq")
	}
	if err := g.validate.Struct(data); err != nil {
		return nil, false, errs.WrapMsg(err, "validation failed", "action", "validate", "dataType", "SetAppBackgroundStatusReq")
	}
	return nil, req.IsBackground, nil
}
