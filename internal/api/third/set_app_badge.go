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

package apiThird

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetAppBadge(c *gin.Context) {
	var (
		req  api.SetAppBadgeReq
		resp api.SetAppBadgeResp
	)
	if err := c.Bind(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)

	var ok bool
	var errInfo, opUserID string
	ok, opUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	if !token_verify.CheckAccess(opUserID, req.FromUserID) {
		log.NewError(req.OperationID, "CheckAccess false ", opUserID, req.FromUserID)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "no permission"})
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req, opUserID)
	err := db.DB.SetUserBadgeUnreadCountSum(req.FromUserID, int(req.AppUnreadCount))
	if err != nil {
		errMsg := req.OperationID + " " + "SetUserBadgeUnreadCountSum failed " + err.Error() + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		resp.ErrCode = 500
		resp.ErrMsg = errMsg
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}
