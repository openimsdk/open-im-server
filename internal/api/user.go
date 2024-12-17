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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msggateway"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

type UserApi struct {
	Discov                discovery.SvcDiscoveryRegistry
	MessageGateWayRpcName string
}

func NewUserApi(discov discovery.SvcDiscoveryRegistry, messageGateWayRpcName string) UserApi {
	return UserApi{
		Discov:                discov,
		MessageGateWayRpcName: messageGateWayRpcName,
	}
}

func (u *UserApi) UserRegister(c *gin.Context) {
	a2r.CallV2(c, user.UserRegisterCaller.Invoke)
}

// UpdateUserInfo is deprecated. Use UpdateUserInfoEx
func (u *UserApi) UpdateUserInfo(c *gin.Context) {
	a2r.CallV2(c, user.UpdateUserInfoCaller.Invoke)
}

func (u *UserApi) UpdateUserInfoEx(c *gin.Context) {
	a2r.CallV2(c, user.UpdateUserInfoExCaller.Invoke)
}
func (u *UserApi) SetGlobalRecvMessageOpt(c *gin.Context) {
	a2r.CallV2(c, user.SetGlobalRecvMessageOptCaller.Invoke)
}

func (u *UserApi) GetUsersPublicInfo(c *gin.Context) {
	a2r.CallV2(c, user.GetDesignateUsersCaller.Invoke)
}

func (u *UserApi) GetAllUsersID(c *gin.Context) {
	a2r.CallV2(c, user.GetAllUserIDCaller.Invoke)
}

func (u *UserApi) AccountCheck(c *gin.Context) {
	a2r.CallV2(c, user.AccountCheckCaller.Invoke)
}

func (u *UserApi) GetUsers(c *gin.Context) {
	a2r.CallV2(c, user.GetPaginationUsersCaller.Invoke)
}

// GetUsersOnlineStatus Get user online status.
func (u *UserApi) GetUsersOnlineStatus(c *gin.Context) {
	var req msggateway.GetUsersOnlineStatusReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	conns, err := u.Discov.GetConns(c, u.MessageGateWayRpcName)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	var wsResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	flag := false

	// Online push message
	for _, v := range conns {
		msgClient := msggateway.NewMsgGatewayClient(v)
		reply, err := msgClient.GetUsersOnlineStatus(c, &req)
		if err != nil {
			log.ZDebug(c, "GetUsersOnlineStatus rpc error", err)

			parseError := apiresp.ParseError(err)
			if parseError.ErrCode == errs.NoPermissionError {
				apiresp.GinError(c, err)
				return
			}
		} else {
			wsResult = append(wsResult, reply.SuccessResult...)
		}
	}
	// Traversing the userIDs in the api request body
	for _, v1 := range req.UserIDs {
		flag = false
		res := new(msggateway.GetUsersOnlineStatusResp_SuccessResult)
		// Iterate through the online results fetched from various gateways
		for _, v2 := range wsResult {
			// If matches the above description on the line, and vice versa
			if v2.UserID == v1 {
				flag = true
				res.UserID = v1
				res.Status = constant.Online
				res.DetailPlatformStatus = append(res.DetailPlatformStatus, v2.DetailPlatformStatus...)
				break
			}
		}
		if !flag {
			res.UserID = v1
			res.Status = constant.Offline
		}
		respResult = append(respResult, res)
	}
	apiresp.GinSuccess(c, respResult)
}

func (u *UserApi) UserRegisterCount(c *gin.Context) {
	a2r.CallV2(c, user.UserRegisterCountCaller.Invoke)
}

// GetUsersOnlineTokenDetail Get user online token details.
func (u *UserApi) GetUsersOnlineTokenDetail(c *gin.Context) {
	var wsResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*msggateway.SingleDetail
	flag := false
	var req msggateway.GetUsersOnlineStatusReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	conns, err := u.Discov.GetConns(c, u.MessageGateWayRpcName)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	// Online push message
	for _, v := range conns {
		msgClient := msggateway.NewMsgGatewayClient(v)
		reply, err := msgClient.GetUsersOnlineStatus(c, &req)
		if err != nil {
			log.ZWarn(c, "GetUsersOnlineStatus rpc err", err)
			continue
		} else {
			wsResult = append(wsResult, reply.SuccessResult...)
		}
	}

	for _, v1 := range req.UserIDs {
		m := make(map[int32][]string, 10)
		flag = false
		temp := new(msggateway.SingleDetail)
		for _, v2 := range wsResult {
			if v2.UserID == v1 {
				flag = true
				temp.UserID = v1
				temp.Status = constant.Online
				for _, status := range v2.DetailPlatformStatus {
					if v, ok := m[status.PlatformID]; ok {
						m[status.PlatformID] = append(v, status.Token)
					} else {
						m[status.PlatformID] = []string{status.Token}
					}
				}
			}
		}
		for p, tokens := range m {
			t := new(msggateway.SinglePlatformToken)
			t.PlatformID = p
			t.Token = tokens
			t.Total = int32(len(tokens))
			temp.SinglePlatformToken = append(temp.SinglePlatformToken, t)
		}

		if flag {
			respResult = append(respResult, temp)
		}
	}

	apiresp.GinSuccess(c, respResult)
}

// SubscriberStatus Presence status of subscribed users.
func (u *UserApi) SubscriberStatus(c *gin.Context) {
	a2r.CallV2(c, user.SubscribeOrCancelUsersStatusCaller.Invoke)
}

// GetUserStatus Get the online status of the user.
func (u *UserApi) GetUserStatus(c *gin.Context) {
	a2r.CallV2(c, user.GetUserStatusCaller.Invoke)
}

// GetSubscribeUsersStatus Get the online status of subscribers.
func (u *UserApi) GetSubscribeUsersStatus(c *gin.Context) {
	a2r.CallV2(c, user.GetSubscribeUsersStatusCaller.Invoke)
}

// ProcessUserCommandAdd user general function add.
func (u *UserApi) ProcessUserCommandAdd(c *gin.Context) {
	a2r.CallV2(c, user.ProcessUserCommandAddCaller.Invoke)
}

// ProcessUserCommandDelete user general function delete.
func (u *UserApi) ProcessUserCommandDelete(c *gin.Context) {
	a2r.CallV2(c, user.ProcessUserCommandDeleteCaller.Invoke)
}

// ProcessUserCommandUpdate  user general function update.
func (u *UserApi) ProcessUserCommandUpdate(c *gin.Context) {
	a2r.CallV2(c, user.ProcessUserCommandUpdateCaller.Invoke)
}

// ProcessUserCommandGet user general function get.
func (u *UserApi) ProcessUserCommandGet(c *gin.Context) {
	a2r.CallV2(c, user.ProcessUserCommandGetCaller.Invoke)
}

// ProcessUserCommandGet user general function get all.
func (u *UserApi) ProcessUserCommandGetAll(c *gin.Context) {
	a2r.CallV2(c, user.ProcessUserCommandGetAllCaller.Invoke)
}

func (u *UserApi) AddNotificationAccount(c *gin.Context) {
	a2r.CallV2(c, user.AddNotificationAccountCaller.Invoke)
}

func (u *UserApi) UpdateNotificationAccountInfo(c *gin.Context) {
	a2r.CallV2(c, user.UpdateNotificationAccountInfoCaller.Invoke)
}

func (u *UserApi) SearchNotificationAccount(c *gin.Context) {
	a2r.CallV2(c, user.SearchNotificationAccountCaller.Invoke)
}
