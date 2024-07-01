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

package rpcclient

import (
	"context"

	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
)

type Third struct {
	conn       grpc.ClientConnInterface
	Client     third.ThirdClient
	discov     discovery.SvcDiscoveryRegistry
	GrafanaUrl string
}

func NewThird(discov discovery.SvcDiscoveryRegistry, rpcRegisterName, grafanaUrl string) *Third {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := third.NewThirdClient(conn)
	if err != nil {
		program.ExitWithError(err)
	}
	return &Third{discov: discov, Client: client, conn: conn, GrafanaUrl: grafanaUrl}
}
func (t *Third) DeleteOutdatedData(ctx context.Context, expires int64) error {
	_, err := t.Client.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{ExpireTime: expires})
	return err
}
