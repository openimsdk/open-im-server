package rpcli

import (
	"github.com/openimsdk/protocol/auth"
	"google.golang.org/grpc"
)

func NewAuthClient(cc grpc.ClientConnInterface) *AuthClient {
	return &AuthClient{auth.NewAuthClient(cc)}
}

type AuthClient struct {
	auth.AuthClient
}
