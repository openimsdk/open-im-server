package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsSearchFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
	OwnerUid    string `json:"ownerUid"`
}

func GetFriendsInfo(c *gin.Context) {
	log.Info("", "", fmt.Sprintf("api search friend init ...."))
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsSearchFriend{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetFriendsInfoReq{
		Uid:         params.UID,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api search_friend is server")
	RpcResp, err := client.GetFriendsInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call search friend rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call search friend rpc server failed"})
		return
	}
	log.InfoByArgs("call search friend rpc server success,args=%s", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		resp := gin.H{
			"errCode": RpcResp.ErrorCode,
			"errMsg":  RpcResp.ErrorMsg,
			"data": gin.H{
				"uid":     RpcResp.Data.Uid,
				"icon":    RpcResp.Data.Icon,
				"name":    RpcResp.Data.Name,
				"gender":  RpcResp.Data.Gender,
				"mobile":  RpcResp.Data.Mobile,
				"birth":   RpcResp.Data.Birth,
				"email":   RpcResp.Data.Email,
				"ex":      RpcResp.Data.Ex,
				"comment": RpcResp.Data.Comment,
			},
		}
		c.JSON(http.StatusOK, resp)
	} else {
		resp := gin.H{
			"errCode": RpcResp.ErrorCode,
			"errMsg":  RpcResp.ErrorMsg,
		}
		c.JSON(http.StatusOK, resp)
	}
	log.InfoByArgs("api search_friend success return,get args=%s,return=%s", req.String(), RpcResp.String())
}
