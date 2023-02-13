package check

import (
	sdkws "Open_IM/pkg/proto/sdkws"
	"context"
	"errors"
)

type GroupChecker struct{}

func NewGroupChecker() GroupChecker {
	return GroupChecker{}
}

func (GroupChecker) GetGroupInfo(ctx context.Context, groupID string) (*sdkws.GroupInfo, error) {
	return nil, errors.New("TODO:GetUserInfo")
}

func (GroupChecker) GetGroupMemberInfo(ctx context.Context, groupID string, userID string) (*sdkws.GroupMemberFullInfo, error) {
	return nil, errors.New("TODO:GetUserInfo")
}
