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

package convert

import (
	"reflect"
	"testing"

	"github.com/openimsdk/protocol/sdkws"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func TestUsersDB2Pb(t *testing.T) {
	type args struct {
		users []*relationtb.User
	}
	tests := []struct {
		name       string
		args       args
		wantResult []*sdkws.UserInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := UsersDB2Pb(tt.args.users); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("UsersDB2Pb() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestUserPb2DB(t *testing.T) {
	type args struct {
		user *sdkws.UserInfo
	}
	tests := []struct {
		name string
		args args
		want *relationtb.User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserPb2DB(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserPb2DB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserPb2DBMap(t *testing.T) {
	user := &sdkws.UserInfo{
		Nickname:         "TestUser",
		FaceURL:          "http://openim.io/logo.jpg",
		Ex:               "Extra Data",
		AppMangerLevel:   1,
		GlobalRecvMsgOpt: 2,
	}

	expected := map[string]any{
		"nickname":            "TestUser",
		"face_url":            "http://openim.io/logo.jpg",
		"ex":                  "Extra Data",
		"app_manager_level":   int32(1),
		"global_recv_msg_opt": int32(2),
	}

	result := UserPb2DBMap(user)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("UserPb2DBMap returned unexpected map. Got %v, want %v", result, expected)
	}
}
