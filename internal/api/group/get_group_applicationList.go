package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
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
	ID               string `json:"id"`
	GroupID          string `json:"groupID"`
	FromUserID       string `json:"fromUserID"`
	ToUserID         string `json:"toUserID"`
	Flag             int32  `json:"flag"`
	RequestMsg       string `json:"reqMsg"`
	HandledMsg       string `json:"handledMsg"`
	AddTime          int64  `json:"createTime"`
	FromUserNickname string `json:"fromUserNickName"`
	ToUserNickname   string `json:"toUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceURL"`
	ToUserFaceUrl    string `json:"toUserFaceURL"`
	HandledUser      string `json:"handledUser"`
	Type             int32  `json:"type"`
	HandleStatus     int32  `json:"handleStatus"`
	HandleResult     int32  `json:"handleResult"`
}

func GetGroupApplicationList(c *gin.Context) {
	log.Info("", "", "api GetGroupApplicationList init ....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := group.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

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

	unProcessCount := 0
	userReq := make([]paramsGroupApplicationListRet, 0)
	if reply != nil && reply.Data != nil && reply.Data.User != nil {
		for i := 0; i < len(reply.Data.User); i++ {
			req := paramsGroupApplicationListRet{}
			req.ID = reply.Data.User[i].ID
			req.GroupID = reply.Data.User[i].GroupID
			req.FromUserID = reply.Data.User[i].FromUserID
			req.ToUserID = reply.Data.User[i].ToUserID
			req.Flag = reply.Data.User[i].Flag
			req.RequestMsg = reply.Data.User[i].RequestMsg
			req.HandledMsg = reply.Data.User[i].HandledMsg
			req.AddTime = reply.Data.User[i].AddTime
			req.FromUserNickname = reply.Data.User[i].FromUserNickname
			req.ToUserNickname = reply.Data.User[i].ToUserNickname
			req.FromUserFaceUrl = reply.Data.User[i].FromUserFaceUrl
			req.ToUserFaceUrl = reply.Data.User[i].ToUserFaceUrl
			req.HandledUser = reply.Data.User[i].HandledUser
			req.Type = reply.Data.User[i].Type
			req.HandleStatus = reply.Data.User[i].HandleStatus
			req.HandleResult = reply.Data.User[i].HandleResult
			userReq = append(userReq, req)

			if req.Flag == 0 {
				unProcessCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
		"data": gin.H{
			"count": unProcessCount,
			"user":  userReq,
		},
	})

}
