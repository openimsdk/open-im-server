package group

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/log"
	"Open_IM/src/proto/group"
	"Open_IM/src/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"net/http"
	"strings"
)

type paramsGroupApplicationList struct {
	OperationID string `json:"operationID" binding:"required"`
}

func newUserRegisterReq(params *paramsGroupApplicationList) *group.GetGroupApplicationListReq {
	pbData := group.GetGroupApplicationListReq{
		OperationID: params.OperationID,
	}
	return &pbData
}

type paramsGroupApplicationListRet struct {
	GroupID          string `json:"groupID"`
	FromUserID       string `json:"fromUserID"`
	FromUserNickName string `json:"fromUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceUrl"`
	ToUserID         string `json:"toUserID"`
	AddTime          int64  `json:"addTime"`
	RequestMsg       string `json:"requestMsg"`
	HandledMsg       string `json:"handledMsg"`
	Type             int32  `json:"type"`
	HandleStatus     int32  `json:"handleStatus"`
	HandleResult     int32  `json:"handleResult"`
}

func GetGroupApplicationList(c *gin.Context) {
	log.Info("", "", "api GetGroupApplicationList init ....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := group.NewGroupClient(etcdConn)
	defer etcdConn.Close()

	params := paramsGroupApplicationList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	pbData := newUserRegisterReq(&params)

	token := c.Request.Header.Get("token")
	if claims, err := utils.ParseToken(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	} else {
		pbData.UID = claims.UID
	}

	log.Info("", "", "api GetGroupApplicationList is server, [data: %s]", pbData.String())
	reply, err := client.GetGroupApplicationList(context.Background(), pbData)
	if err != nil {
		log.Error("", "", "api GetGroupApplicationList call rpc fail, [data: %s] [err: %s]", pbData.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info("", "", "api GetGroupApplicationList call rpc success, [data: %s] [reply: %s]", pbData.String(), reply.String())

	var userReq []paramsGroupApplicationListRet
	for i := 0; i < len(reply.Data.User); i++ {
		req := paramsGroupApplicationListRet{}
		req.GroupID = reply.Data.User[i].GroupID
		req.FromUserID = reply.Data.User[i].FromUserID
		req.FromUserNickName = reply.Data.User[i].FromUserNickName
		req.FromUserFaceUrl = reply.Data.User[i].FromUserFaceUrl
		req.ToUserID = reply.Data.User[i].ToUserID
		req.RequestMsg = reply.Data.User[i].RequestMsg
		req.HandledMsg = reply.Data.User[i].HandledMsg
		req.Type = reply.Data.User[i].Type
		req.HandleStatus = reply.Data.User[i].HandleStatus
		req.HandleResult = reply.Data.User[i].HandleResult
		userReq = append(userReq, req)
	}

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
		"data": gin.H{
			"count": reply.Data.Count,
			"user":  userReq,
		},
	})

}
