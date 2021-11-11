package user

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbUser "Open_IM/pkg/proto/user"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsStruct struct {
	OperationID string   `json:"operationID" binding:"required"`
	UIDList     []string `json:"uidList"`
	Platform    int32    `json:"platform"`
	Name        string   `json:"name"`
	Icon        string   `json:"icon"`
	Gender      int32    `json:"gender"`
	Mobile      string   `json:"mobile"`
	Birth       string   `json:"birth"`
	Email       string   `json:"email"`
	Ex          string   `json:"ex"`
	Uid         string   `json:"uid"`
}

func UpdateUserInfo(c *gin.Context) {
	log.InfoByKv("api update userinfo init...", "")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsStruct{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbUser.UpdateUserInfoReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
		Name:        params.Name,
		Icon:        params.Icon,
		Gender:      params.Gender,
		Mobile:      params.Mobile,
		Birth:       params.Birth,
		Email:       params.Email,
		Ex:          params.Ex,
		Uid:         params.Uid,
	}
	log.InfoByKv("api update user info is server", req.OperationID, req.Token)
	RpcResp, err := client.UpdateUserInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call get user info rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByKv("call update user info rpc server success", params.OperationID)
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	log.InfoByKv("api update user info return success", params.OperationID, "args=%s", RpcResp.String())
}
