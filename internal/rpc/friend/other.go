package friend

import (
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"context"
	"errors"
)

func GetPublicUserInfoBatch(ctx context.Context, userIDs []string) ([]*server_api_params.PublicUserInfo, error) {
	if len(userIDs) == 0 {
		return []*server_api_params.PublicUserInfo{}, nil
	}
	return nil, errors.New("TODO:GetUserInfo")
}
