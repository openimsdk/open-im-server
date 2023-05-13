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

package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/msg"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsUserNewestSeq struct {
	ReqIdentifier int    `json:"reqIdentifier" binding:"required"`
	SendID        string `json:"sendID" binding:"required"`
	OperationID   string `json:"operationID" binding:"required"`
	MsgIncr       int    `json:"msgIncr" binding:"required"`
}

func GetSeq(c *gin.Context) {
	params := paramsUserNewestSeq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	token := c.Request.Header.Get("token")
	if ok, err := token_verify.VerifyToken(token, params.SendID); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err" + err.Error()})
		return
	}
	pbData := sdk_ws.GetMaxAndMinSeqReq{}
	pbData.UserID = params.SendID
	pbData.OperationID = params.OperationID
	grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, pbData.OperationID)
	if grpcConn == nil {
		errMsg := pbData.OperationID + " getcdv3.GetDefaultConn == nil"
		log.NewError(pbData.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	msgClient := pbChat.NewMsgClient(grpcConn)
	reply, err := msgClient.GetMaxAndMinSeq(context.Background(), &pbData)
	if err != nil {
		log.NewError(params.OperationID, "UserGetSeq rpc failed, ", params, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "UserGetSeq rpc failed, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errCode":       reply.ErrCode,
		"errMsg":        reply.ErrMsg,
		"msgIncr":       params.MsgIncr,
		"reqIdentifier": params.ReqIdentifier,
		"data": gin.H{
			"maxSeq": reply.MaxSeq,
			"minSeq": reply.MinSeq,
		},
	})

}
