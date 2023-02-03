package check

import (
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"errors"
)

type GroupChecker struct {
}

func NewGroupChecker() *GroupChecker {
	return &GroupChecker{}
}

func (g *GroupChecker) GetGroupInfo(groupID string) (*server_api_params.GroupInfo, error) {
	return nil, errors.New("TODO:GetUserInfo")
}
