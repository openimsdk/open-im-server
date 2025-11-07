package convert

import (
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"reflect"
	"testing"

	"github.com/openimsdk/protocol/sdkws"
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
