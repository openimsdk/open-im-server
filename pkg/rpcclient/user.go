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

	"github.com/OpenIMSDK/Open-IM-Server/pkg/authverify"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)

type User struct {
	conn   grpc.ClientConnInterface
	Client user.UserClient
	Discov discoveryregistry.SvcDiscoveryRegistry
}

func NewUser(discov discoveryregistry.SvcDiscoveryRegistry) *User {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		panic(err)
	}
	client := user.NewUserClient(conn)
	return &User{Discov: discov, Client: client, conn: conn}
}

type UserRpcClient User

func NewUserRpcClientByUser(user *User) *UserRpcClient {
	rpc := UserRpcClient(*user)
	return &rpc
}

func NewUserRpcClient(client discoveryregistry.SvcDiscoveryRegistry) UserRpcClient {
	return UserRpcClient(*NewUser(client))
}

func (u *UserRpcClient) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
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

func (u *UserRpcClient) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (u *UserRpcClient) GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdkws.UserInfo) string {
		return e.UserID
	}), nil
}

func (u *UserRpcClient) GetPublicUserInfos(
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

func (u *UserRpcClient) GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (u *UserRpcClient) GetPublicUserInfoMap(
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

func (u *UserRpcClient) GetUserGlobalMsgRecvOpt(ctx context.Context, userID string) (int32, error) {
	resp, err := u.Client.GetGlobalRecvMessageOpt(ctx, &user.GetGlobalRecvMessageOptReq{
		UserID: userID,
	})
	if err != nil {
		return 0, err
	}
	return resp.GlobalRecvMsgOpt, err
}

func (u *UserRpcClient) Access(ctx context.Context, ownerUserID string) error {
	_, err := u.GetUserInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}
	return authverify.CheckAccessV3(ctx, ownerUserID)
}

func (u *UserRpcClient) GetAllUserIDs(ctx context.Context, pageNumber, showNumber int32) ([]string, error) {
	resp, err := u.Client.GetAllUserID(ctx, &user.GetAllUserIDReq{Pagination: &sdkws.RequestPagination{PageNumber: pageNumber, ShowNumber: showNumber}})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (u *UserRpcClient) SetUserStatus(ctx context.Context, userID string, status int32, platformID int) error {
	_, err := u.Client.SetUserStatus(ctx, &user.SetUserStatusReq{StatusList: []*user.OnlineStatus{{UserID: userID, Status: status, PlatformID: int32(platformID)}}})
	return err
}
