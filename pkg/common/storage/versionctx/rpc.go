package versionctx

import (
	"context"

	"google.golang.org/grpc"
)

func EnableVersionCtx() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(enableVersionCtxInterceptor)
}

func enableVersionCtxInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return handler(WithVersionLog(ctx), req)
}
