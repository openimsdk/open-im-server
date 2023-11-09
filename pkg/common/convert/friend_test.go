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
	"context"
	"reflect"
	"testing"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

func TestFriendPb2DB(t *testing.T) {
	type args struct {
		friend *sdkws.FriendInfo
	}
	tests := []struct {
		name string
		args args
		want *relation.FriendModel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FriendPb2DB(tt.args.friend); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FriendPb2DB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFriendDB2Pb(t *testing.T) {
	type args struct {
		ctx      context.Context
		friendDB *relation.FriendModel
		getUsers func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error)
	}
	tests := []struct {
		name    string
		args    args
		want    *sdkws.FriendInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FriendDB2Pb(tt.args.ctx, tt.args.friendDB, tt.args.getUsers)
			if (err != nil) != tt.wantErr {
				t.Errorf("FriendDB2Pb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FriendDB2Pb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFriendsDB2Pb(t *testing.T) {
	type args struct {
		ctx       context.Context
		friendsDB []*relation.FriendModel
		getUsers  func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error)
	}
	tests := []struct {
		name          string
		args          args
		wantFriendsPb []*sdkws.FriendInfo
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFriendsPb, err := FriendsDB2Pb(tt.args.ctx, tt.args.friendsDB, tt.args.getUsers)
			if (err != nil) != tt.wantErr {
				t.Errorf("FriendsDB2Pb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFriendsPb, tt.wantFriendsPb) {
				t.Errorf("FriendsDB2Pb() = %v, want %v", gotFriendsPb, tt.wantFriendsPb)
			}
		})
	}
}

func TestFriendRequestDB2Pb(t *testing.T) {
	type args struct {
		ctx            context.Context
		friendRequests []*relation.FriendRequestModel
		getUsers       func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error)
	}
	tests := []struct {
		name    string
		args    args
		want    []*sdkws.FriendRequest
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FriendRequestDB2Pb(tt.args.ctx, tt.args.friendRequests, tt.args.getUsers)
			if (err != nil) != tt.wantErr {
				t.Errorf("FriendRequestDB2Pb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FriendRequestDB2Pb() = %v, want %v", got, tt.want)
			}
		})
	}
}
