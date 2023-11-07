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

package cmd

import (
	"testing"

	"github.com/OpenIMSDK/protocol/constant"
)

func TestApiCmd_AddApi(t *testing.T) {
	type fields struct {
		RootCmd *RootCmd
	}
	type args struct {
		f func(port int) error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		// No 'want' field here since we're testing side effects, not return values.
	}{
		{
			name: "adds API with valid function",
			fields: fields{
				RootCmd: &RootCmd{
					// setup RootCmd properties as needed for the test
				},
			},
			args: args{
				f: func(port int) error {
					// implement a mock function or check side effects
					return nil
				},
			},
		},
		// Add more test cases as needed.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ApiCmd{
				RootCmd: tt.fields.RootCmd,
			}
			a.AddApi(tt.args.f)
			// Test the side effects or behavior of AddApi here.
		})
	}
}

func TestApiCmd_GetPortFromConfig(t *testing.T) {
	type fields struct {
		RootCmd *RootCmd
	}
	type args struct {
		portType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "gets API port from config",
			fields: fields{
				RootCmd: &RootCmd{
					// setup RootCmd properties as needed for the test
				},
			},
			args: args{
				portType: constant.FlagPort,
			},
			want: 8080, // This should be the expected port number according to your configuration
		},
		// Add more test cases for different portTypes and expected results.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ApiCmd{
				RootCmd: tt.fields.RootCmd,
			}
			if got := a.GetPortFromConfig(tt.args.portType); got != tt.want {
				t.Errorf("ApiCmd.GetPortFromConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
