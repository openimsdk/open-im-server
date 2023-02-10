package check

import (
	sdkws "Open_IM/pkg/proto/sdkws"
	"context"
	"errors"
)

//func GetUsersInfo(ctx context.Context, args ...interface{}) ([]*sdkws.UserInfo, error) {
//	return nil, errors.New("TODO:GetUserInfo")
//}

func NewUserCheck() *UserCheck {
	return &UserCheck{}
}

type UserCheck struct{}

func (u *UserCheck) GetUsersInfos(ctx context.Context, userIDs []string, complete bool) ([]*sdkws.UserInfo, error) {
	return nil, errors.New("todo")
}

func (u *UserCheck) GetUsersInfoMap(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.UserInfo, error) {
	return nil, errors.New("todo")
}

func (u *UserCheck) GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	return nil, errors.New("todo")
}

func (u *UserCheck) GetPublicUserInfos(ctx context.Context, userIDs []string, complete bool) ([]*sdkws.PublicUserInfo, error) {
	return nil, errors.New("todo")
}

func (u *UserCheck) GetPublicUserInfoMap(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.PublicUserInfo, error) {
	return nil, errors.New("todo")
}
