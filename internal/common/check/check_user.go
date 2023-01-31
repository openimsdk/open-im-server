package check

import (
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"context"
	"errors"
)

func GetUsersInfo(ctx context.Context, args ...interface{}) ([]*server_api_params.UserInfo, error) {
	return nil, errors.New("TODO:GetUserInfo")
}
