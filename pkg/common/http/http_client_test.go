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

package http

import (
	"context"
	"reflect"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func TestGet(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse []byte
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResponse, err := Get(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Get() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		ctx     context.Context
		url     string
		header  map[string]string
		data    interface{}
		timeout int
	}
	tests := []struct {
		name        string
		args        args
		wantContent []byte
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := Post(tt.args.ctx, tt.args.url, tt.args.header, tt.args.data, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("Post() = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}

func TestPostReturn(t *testing.T) {
	type args struct {
		ctx           context.Context
		url           string
		header        map[string]string
		input         interface{}
		output        interface{}
		timeOutSecond int
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
			if err := PostReturn(tt.args.ctx, tt.args.url, tt.args.header, tt.args.input, tt.args.output, tt.args.timeOutSecond); (err != nil) != tt.wantErr {
				t.Errorf("PostReturn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_callBackPostReturn(t *testing.T) {
	type args struct {
		ctx            context.Context
		url            string
		command        string
		input          interface{}
		output         callbackstruct.CallbackResp
		callbackConfig config.CallBackConfig
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
			if err := callBackPostReturn(tt.args.ctx, tt.args.url, tt.args.command, tt.args.input, tt.args.output, tt.args.callbackConfig); (err != nil) != tt.wantErr {
				t.Errorf("callBackPostReturn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCallBackPostReturn(t *testing.T) {
	type args struct {
		ctx            context.Context
		url            string
		req            callbackstruct.CallbackReq
		resp           callbackstruct.CallbackResp
		callbackConfig config.CallBackConfig
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
			if err := CallBackPostReturn(tt.args.ctx, tt.args.url, tt.args.req, tt.args.resp, tt.args.callbackConfig); (err != nil) != tt.wantErr {
				t.Errorf("CallBackPostReturn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
