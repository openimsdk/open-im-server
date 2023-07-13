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

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type User struct {
	conn   grpc.ClientConnInterface
	Client user.UserClient
	Discov discoveryregistry.SvcDiscoveryRegistry
}

// NewUser 新建用户
func NewUser(discov discoveryregistry.SvcDiscoveryRegistry) *User {
	conn, err := discov.GetConn(context.Background(), config.Config.RPCRegisterName.OpenImUserName)
	if err != nil {
		panic(err)
	}
	client := user.NewUserClient(conn)

	return &User{Discov: discov, Client: client, conn: conn}
}

type UserRPCClient User

// NewUserRPCClient 新建用户 RPC 客户端
func NewUserRPCClient(client discoveryregistry.SvcDiscoveryRegistry) UserRPCClient {
	return UserRPCClient(*NewUser(client))
}

func (u *UserRPCClient) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	resp, err := u.Client.GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(userIDs, utils.Slice(resp.UsersInfo, func(e *sdkws.UserInfo) string {
		return e.UserID
	})); len(ids) > 0 {
		return nil, errs.ErrUserIDNotFound.Wrap(strings.Join(ids, ","))
	}
	return resp.UsersInfo, nil
}

// GetUserInfo 获取指定用户信息
func (u *UserRPCClient) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

// GetUsersInfoMap 获取用户信息集合
func (u *UserRPCClient) GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return utils.SliceToMap(users, func(e *sdkws.UserInfo) string {
		return e.UserID
	}), nil
}

// 获取公有用户信息
func (u *UserRPCClient) GetPublicUserInfos(
	ctx context.Context,
	userIDs []string,
	complete bool,
) ([]*sdkws.PublicUserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.Slice(users, func(e *sdkws.UserInfo) *sdkws.PublicUserInfo {
		return &sdkws.PublicUserInfo{
			UserID:   e.UserID,
			Nickname: e.Nickname,
			FaceURL:  e.FaceURL,
			Ex:       e.Ex,
		}
	}), nil
}

// GetPublicUserInfo 获取用户信息
func (u *UserRPCClient) GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, []string{userID}, true)
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

// GetPublicUserInfoMap 获取用户信息集合
func (u *UserRPCClient) GetPublicUserInfoMap(
	ctx context.Context,
	userIDs []string,
	complete bool,
) (map[string]*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}

	return utils.SliceToMap(users, func(e *sdkws.PublicUserInfo) string {
		return e.UserID
	}), nil
}

// GetUserGlobalMsgRecvOpt 获取用户消息接收选项
func (u *UserRPCClient) GetUserGlobalMsgRecvOpt(ctx context.Context, userID string) (int32, error) {
	resp, err := u.Client.GetGlobalRecvMessageOpt(ctx, &user.GetGlobalRecvMessageOptReq{
		UserID: userID,
	})
	if err != nil {
		return 0, err
	}

	return resp.GlobalRecvMsgOpt, err
}

// Access token 验签
func (u *UserRPCClient) Access(ctx context.Context, ownerUserID string) error {
	_, err := u.GetUserInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}

	return tokenverify.CheckAccessV3(ctx, ownerUserID)
}
