package rpcli

import (
	"github.com/openimsdk/protocol/msggateway"
)

func NewMsgGatewayClient(cli msggateway.MsgGatewayClient) *MsgGatewayClient {
	return &MsgGatewayClient{cli}
}

type MsgGatewayClient struct {
	msggateway.MsgGatewayClient
}
