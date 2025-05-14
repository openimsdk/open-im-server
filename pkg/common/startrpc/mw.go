package startrpc

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"google.golang.org/grpc"
)

func grpcServerIMAdminUserID(imAdminUserID []string) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = authverify.WithIMAdminUserIDs(ctx, imAdminUserID)
		return handler(ctx, req)
	})
}
