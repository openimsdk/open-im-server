package rpcli

import "github.com/openimsdk/protocol/third"

func NewThirdClient(cli third.ThirdClient) *ThirdClient {
	return &ThirdClient{cli}
}

type ThirdClient struct {
	third.ThirdClient
}
