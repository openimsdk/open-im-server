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

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
)

// User represents a structure holding connection details for the User RPC client.
type User struct {
	conn                  grpc.ClientConnInterface
	Client                user.UserClient
	Discov                discovery.SvcDiscoveryRegistry
	MessageGateWayRpcName string
	imAdminUserID         []string
}

// NewUser initializes and returns a User instance based on the provided service discovery registry.
func NewUser(discov discovery.SvcDiscoveryRegistry, rpcRegisterName, messageGateWayRpcName string,
	imAdminUserID []string) *User {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := user.NewUserClient(conn)
	return &User{Discov: discov, Client: client,
		conn:                  conn,
		MessageGateWayRpcName: messageGateWayRpcName,
		imAdminUserID:         imAdminUserID}
}

// UserRpcClient represents the structure for a User RPC client.
type UserRpcClient User

// NewUserRpcClientByUser initializes a UserRpcClient based on the provided User instance.
func NewUserRpcClientByUser(user *User) *UserRpcClient {
	rpc := UserRpcClient(*user)
	return &rpc
}

// NewUserRpcClient initializes a UserRpcClient based on the provided service discovery registry.
func NewUserRpcClient(client discovery.SvcDiscoveryRegistry, rpcRegisterName string,
	imAdminUserID []string) UserRpcClient {
	return UserRpcClient(*NewUser(client, rpcRegisterName, "", imAdminUserID))
}

// GetUsersInfo retrieves information for multiple users based on their user IDs.
func (u *UserRpcClient) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	if len(userIDs) == 0 {
		return []*sdkws.UserInfo{}, nil
	}
	resp, err := u.Client.GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
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
func (u *UserRpcClient) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

// GetUsersInfoMap retrieves a map of user information indexed by their user IDs.
func (u *UserRpcClient) GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMap(users, func(e *sdkws.UserInfo) string {
		return e.UserID
	}), nil
}

// GetPublicUserInfos retrieves public information for multiple users based on their user IDs.
func (u *UserRpcClient) GetPublicUserInfos(
	ctx context.Context,
	userIDs []string,
) ([]*sdkws.PublicUserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
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
func (u *UserRpcClient) GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, []string{userID})
	if err != nil {
		return nil, err
	}

	return users[0], nil
}

// GetPublicUserInfoMap retrieves a map of public user information indexed by their user IDs.
func (u *UserRpcClient) GetPublicUserInfoMap(
	ctx context.Context,
	userIDs []string,
) (map[string]*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return datautil.SliceToMap(users, func(e *sdkws.PublicUserInfo) string {
		return e.UserID
	}), nil
}

// GetUserGlobalMsgRecvOpt retrieves the global message receive option for a user based on the provided user ID.
func (u *UserRpcClient) GetUserGlobalMsgRecvOpt(ctx context.Context, userID string) (int32, error) {
	resp, err := u.Client.GetGlobalRecvMessageOpt(ctx, &user.GetGlobalRecvMessageOptReq{
		UserID: userID,
	})
	if err != nil {
		return 0, err
	}
	return resp.GlobalRecvMsgOpt, nil
}

// Access verifies the access rights for the provided user ID.
func (u *UserRpcClient) Access(ctx context.Context, ownerUserID string) error {
	_, err := u.GetUserInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}
	return authverify.CheckAccessV3(ctx, ownerUserID, u.imAdminUserID)
}

// GetAllUserID retrieves all user IDs with pagination options.
func (u *UserRpcClient) GetAllUserID(ctx context.Context, pageNumber, showNumber int32) (*user.GetAllUserIDResp, error) {
	resp, err := u.Client.GetAllUserID(ctx, &user.GetAllUserIDReq{Pagination: &sdkws.RequestPagination{PageNumber: pageNumber, ShowNumber: showNumber}})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAllUserIDs retrieves all user IDs with pagination options.
func (u *UserRpcClient) GetAllUserIDs(ctx context.Context, pageNumber, showNumber int32) ([]string, error) {
	resp, err := u.Client.GetAllUserID(ctx, &user.GetAllUserIDReq{Pagination: &sdkws.RequestPagination{PageNumber: pageNumber, ShowNumber: showNumber}})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

// SetUserStatus sets the status for a user based on the provided user ID, status, and platform ID.
func (u *UserRpcClient) SetUserStatus(ctx context.Context, userID string, status int32, platformID int) error {
	_, err := u.Client.SetUserStatus(ctx, &user.SetUserStatusReq{
		UserID: userID,
		Status: status, PlatformID: int32(platformID),
	})
	return err
}

func (u *UserRpcClient) GetNotificationByID(ctx context.Context, userID string) error {
	_, err := u.Client.GetNotificationAccount(ctx, &user.GetNotificationAccountReq{
		UserID: userID,
	})
	return err
}

func (u *UserRpcClient) GetUsersOnlinePlatform(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	resp, err := u.Client.GetUserStatus(ctx, &user.GetUserStatusReq{UserIDs: userIDs, UserID: u.imAdminUserID[0]})
	if err != nil {
		return nil, err
	}
	return resp.StatusList, nil
}

func (u *UserRpcClient) GetUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	resp, err := u.GetUsersOnlinePlatform(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return resp[0].PlatformIDs, nil
}

func (u *UserRpcClient) GetAllOnlineUsers(ctx context.Context, cursor uint64) (*user.GetAllOnlineUsersResp, error) {
	return u.Client.GetAllOnlineUsers(ctx, &user.GetAllOnlineUsersReq{Cursor: cursor})
}
