package rpcli

import "github.com/openimsdk/protocol/user"

func NewUserClient(cli user.UserClient) *UserClient {
	return &UserClient{cli}
}

type UserClient struct {
	user.UserClient
}
