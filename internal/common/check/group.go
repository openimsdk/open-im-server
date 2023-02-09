package check

import (
	sdkws "Open_IM/pkg/proto/sdkws"
	"errors"
)

type GroupChecker struct {
}

func NewGroupChecker() *GroupChecker {
	return &GroupChecker{}
}

func (g *GroupChecker) GetGroupInfo(groupID string) (*sdkws.GroupInfo, error) {
	return nil, errors.New("TODO:GetUserInfo")
}
