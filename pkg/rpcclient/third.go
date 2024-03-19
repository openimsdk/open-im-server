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

	util "github.com/openimsdk/open-im-server/v3/pkg/util/genutil"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/discoveryregistry"
	"google.golang.org/grpc"
)

type Third struct {
	conn       grpc.ClientConnInterface
	Client     third.ThirdClient
	discov     discoveryregistry.SvcDiscoveryRegistry
	GrafanaUrl string
}

func NewThird(discov discoveryregistry.SvcDiscoveryRegistry, rpcRegisterName, grafanaUrl string) *Third {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		util.ExitWithError(err)
	}
	client := third.NewThirdClient(conn)
	if err != nil {
		util.ExitWithError(err)
	}
	return &Third{discov: discov, Client: client, conn: conn, GrafanaUrl: grafanaUrl}
}
