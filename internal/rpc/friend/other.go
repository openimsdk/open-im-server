package friend

import (
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"context"
	"errors"
)

func GetUserInfo(ctx context.Context, userID string) (*server_api_params.PublicUserInfo, error) {
	return nil, errors.New("TODO:GetUserInfo")
}

func GetPublicUserInfoBatch(ctx context.Context, userIDs []string) ([]*server_api_params.PublicUserInfo, error) {
	if len(userIDs) == 0 {
		return []*server_api_params.PublicUserInfo{}, nil
	}
	return nil, errors.New("TODO:GetUserInfo")
}

func GetUserInfoList(ctx context.Context, userIDs []string) ([]*server_api_params.UserInfo, error) {
	if len(userIDs) == 0 {
		return []*server_api_params.UserInfo{}, nil
	}
	return nil, errors.New("TODO:GetUserInfo")
}
