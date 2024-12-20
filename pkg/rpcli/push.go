package rpcli

import (
	"github.com/openimsdk/protocol/push"
)

func NewPushMsgServiceClient(cli push.PushMsgServiceClient) *PushMsgServiceClient {
	return &PushMsgServiceClient{cli}
}

type PushMsgServiceClient struct {
	push.PushMsgServiceClient
}
