// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
