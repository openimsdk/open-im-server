package prometheus

import (
	"context"
	"encoding/json"
	"time"

	"Open_IM/pkg/common/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func UnaryServerInterceptorProme(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	remote, _ := peer.FromContext(ctx)
	remoteAddr := remote.Addr.String()

	in, _ := json.Marshal(req)
	inStr := string(in)
	log.NewInfo("ip", remoteAddr, "access_start", info.FullMethod, "in", inStr)

	start := time.Now()
	defer func() {
		out, _ := json.Marshal(resp)
		outStr := string(out)
		duration := int64(time.Since(start) / time.Millisecond)
		if duration >= 500 {
			log.NewInfo("ip", remoteAddr, "access_end", info.FullMethod, "in", inStr, "out", outStr, "err", err, "duration/ms", duration)
		} else {
			log.NewInfo("ip", remoteAddr, "access_end", info.FullMethod, "in", inStr, "out", outStr, "err", err, "duration/ms", duration)
		}
	}()
	resp, err = handler(ctx, req)
	return
}
