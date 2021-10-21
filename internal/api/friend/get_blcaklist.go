package friend

import (
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsGetBlackList struct {
	OperationID string `json:"operationID" binding:"required"`
}

type blackListUserInfo struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Gender int32  `json:"gender"`
	Mobile string `json:"mobile"`
	Birth  string `json:"birth"`
	Email  string `json:"email"`
	Ex     string `json:"ex"`
}

func GetBlacklist(c *gin.Context) {
	log.Info("", "", "api get blacklist init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsGetBlackList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetBlacklistReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, fmt.Sprintf("api get blacklist is server"))
	RpcResp, err := client.GetBlacklist(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call get_friend_list rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get blacklist rpc server failed"})
		return
	}
	log.InfoByArgs("call get blacklist rpc server success,args=%s", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		userBlackList := make([]blackListUserInfo, 0)
		for _, friend := range RpcResp.Data {
			var fi blackListUserInfo
			fi.UID = friend.Uid
			fi.Name = friend.Name
			fi.Icon = friend.Icon
			fi.Gender = friend.Gender
			fi.Mobile = friend.Mobile
			fi.Birth = friend.Birth
			fi.Email = friend.Email
			fi.Ex = friend.Ex
			userBlackList = append(userBlackList, fi)
		}
		resp := gin.H{
			"errCode": RpcResp.ErrorCode,
			"errMsg":  RpcResp.ErrorMsg,
			"data":    userBlackList,
		}
		c.JSON(http.StatusOK, resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
		c.JSON(http.StatusOK, resp)
	}
	log.InfoByArgs("api get black list success return,get args=%s,return=%s", req.String(), RpcResp.String())
}
