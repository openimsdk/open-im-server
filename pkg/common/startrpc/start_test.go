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

package startrpc

import (
	"testing"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"google.golang.org/grpc"
)

func TestStart(t *testing.T) {
	type args struct {
		rpcPort         int
		rpcRegisterName string
		prometheusPort  int
		rpcFn           func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error
		options         []grpc.ServerOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Start(tt.args.rpcPort, tt.args.rpcRegisterName, tt.args.prometheusPort, tt.args.rpcFn, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
