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

package rpcclient

import (
	"context"
	"strings"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/utils/datautil"
)

// GetUsersInfo retrieves information for multiple users based on their user IDs.
func GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	if len(userIDs) == 0 {
		return []*sdkws.UserInfo{}, nil
	}
	resp, err := user.GetDesignateUsersCaller.Invoke(ctx, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if ids := datautil.Single(userIDs, datautil.Slice(resp.UsersInfo, func(e *sdkws.UserInfo) string {
		return e.UserID
	})); len(ids) > 0 {
		return nil, servererrs.ErrUserIDNotFound.WrapMsg(strings.Join(ids, ","))
	}
	return resp.UsersInfo, nil
}

// GetUserInfo retrieves information for a single user based on the provided user ID.
func GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	users, err := GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

// GetUsersInfoMap retrieves a map of user information indexed by their user IDs.
func GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users, err := GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMap(users, func(e *sdkws.UserInfo) string {
		return e.UserID
	}), nil
}

// GetPublicUserInfos retrieves public information for multiple users based on their user IDs.
func GetPublicUserInfos(
	ctx context.Context,
	userIDs []string,
) ([]*sdkws.PublicUserInfo, error) {
	users, err := GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return datautil.Slice(users, func(e *sdkws.UserInfo) *sdkws.PublicUserInfo {
		return &sdkws.PublicUserInfo{
			UserID:   e.UserID,
			Nickname: e.Nickname,
			FaceURL:  e.FaceURL,
			Ex:       e.Ex,
		}
	}), nil
}

// GetPublicUserInfo retrieves public information for a single user based on the provided user ID.
func GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := GetPublicUserInfos(ctx, []string{userID})
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

// GetPublicUserInfoMap retrieves a map of public user information indexed by their user IDs.
func GetPublicUserInfoMap(
	ctx context.Context,
	userIDs []string,
) (map[string]*sdkws.PublicUserInfo, error) {
	users, err := GetPublicUserInfos(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return datautil.SliceToMap(users, func(e *sdkws.PublicUserInfo) string {
		return e.UserID
	}), nil
}
