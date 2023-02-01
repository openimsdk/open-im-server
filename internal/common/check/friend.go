package check

import (
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"context"
	"errors"
)

func GetFriendsInfo(ctx context.Context, ownerUserID, friendUserID string) (*server_api_params.FriendInfo, error) {
	return nil, errors.New("TODO:GetUserInfo")
}
