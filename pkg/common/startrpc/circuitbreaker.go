package startrpc

import (
	"context"
	"time"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/openimsdk/tools/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreaker struct {
	Enable  bool          `yaml:"enable"`
	Success float64       `yaml:"success"` // success rate threshold (0.0-1.0)
	Request int64         `yaml:"request"` // request threshold
	Bucket  int           `yaml:"bucket"`  // number of buckets
	Window  time.Duration `yaml:"window"`  // time window for statistics
}

func NewCircuitBreaker(config *CircuitBreaker) circuitbreaker.CircuitBreaker {
	if !config.Enable {
		return nil
	}

	return sre.NewBreaker(
		sre.WithWindow(config.Window),
		sre.WithBucket(config.Bucket),
		sre.WithSuccess(config.Success),
		sre.WithRequest(config.Request),
	)
}

func UnaryCircuitBreakerInterceptor(breaker circuitbreaker.CircuitBreaker) grpc.ServerOption {
	if breaker == nil {
		return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			return handler(ctx, req)
		})
	}

	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if err := breaker.Allow(); err != nil {
			log.ZWarn(ctx, "rpc circuit breaker open", err, "method", info.FullMethod)
			return nil, status.Error(codes.Unavailable, "service unavailable due to circuit breaker")
		}

		resp, err = handler(ctx, req)

		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Internal {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		} else {
			breaker.MarkSuccess()
		}

		return resp, err

	})
}

func StreamCircuitBreakerInterceptor(breaker circuitbreaker.CircuitBreaker) grpc.ServerOption {
	if breaker == nil {
		return grpc.ChainStreamInterceptor(func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, ss)
		})
	}

	return grpc.ChainStreamInterceptor(func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := breaker.Allow(); err != nil {
			log.ZWarn(ss.Context(), "rpc circuit breaker open", err, "method", info.FullMethod)
			return status.Error(codes.Unavailable, "service unavailable due to circuit breaker")
		}

		err := handler(srv, ss)

		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Internal {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
		} else {
			breaker.MarkSuccess()
		}

		return err
	})
}
