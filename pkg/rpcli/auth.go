package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/auth"
	"google.golang.org/grpc"
)

func NewAuthClient(cc grpc.ClientConnInterface) *AuthClient {
	return &AuthClient{auth.NewAuthClient(cc)}
}

type AuthClient struct {
	auth.AuthClient
}

func (x *AuthClient) KickTokens(ctx context.Context, tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}
	return ignoreResp(x.AuthClient.KickTokens(ctx, &auth.KickTokensReq{Tokens: tokens}))
}

func (x *AuthClient) InvalidateToken(ctx context.Context, req *auth.InvalidateTokenReq) error {
	return ignoreResp(x.AuthClient.InvalidateToken(ctx, req))
}

func (x *AuthClient) ParseToken(ctx context.Context, token string) (*auth.ParseTokenResp, error) {
	return x.AuthClient.ParseToken(ctx, &auth.ParseTokenReq{Token: token})
}
