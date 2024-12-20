package rpcli

import "github.com/openimsdk/protocol/auth"

func NewAuthClient(cli auth.AuthClient) *AuthClient {
	return &AuthClient{cli}
}

type AuthClient struct {
	auth.AuthClient
}
