package startrpc

import (
	"context"
	"time"

	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/stability/ratelimit"
	"github.com/openimsdk/tools/stability/ratelimit/bbr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimiter struct {
	Enable       bool
	Window       time.Duration
	Bucket       int
	CPUThreshold int64
}

func NewRateLimiter(config *RateLimiter) ratelimit.Limiter {
	if !config.Enable {
		return nil
	}

	return bbr.NewBBRLimiter(
		bbr.WithWindow(config.Window),
		bbr.WithBucket(config.Bucket),
		bbr.WithCPUThreshold(config.CPUThreshold),
	)
}

func UnaryRateLimitInterceptor(limiter ratelimit.Limiter) grpc.ServerOption {
	if limiter == nil {
		return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			return handler(ctx, req)
		})
	}

	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		done, err := limiter.Allow()
		if err != nil {
			log.ZWarn(ctx, "rpc rate limited", err, "method", info.FullMethod)
			return nil, status.Errorf(codes.ResourceExhausted, "rpc request rate limit exceeded: %v, please try again later", err)
		}

		defer done(ratelimit.DoneInfo{})
		return handler(ctx, req)
	})
}

func StreamRateLimitInterceptor(limiter ratelimit.Limiter) grpc.ServerOption {
	if limiter == nil {
		return grpc.ChainStreamInterceptor(func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, ss)
		})
	}

	return grpc.ChainStreamInterceptor(func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		done, err := limiter.Allow()
		if err != nil {
			log.ZWarn(ss.Context(), "rpc rate limited", err, "method", info.FullMethod)
			return status.Errorf(codes.ResourceExhausted, "rpc request rate limit exceeded: %v, please try again later", err)
		}
		defer done(ratelimit.DoneInfo{})

		return handler(srv, ss)
	})
}
