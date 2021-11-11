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

type paramsGroupApplicationResponse struct {
	OperationID      string `json:"operationID" binding:"required"`
	GroupID          string `json:"groupID" binding:"required"`
	FromUserID       string `json:"fromUserID" binding:"required"`
	FromUserNickName string `json:"fromUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceUrl"`
	ToUserID         string `json:"toUserID" binding:"required"`
	ToUserNickName   string `json:"toUserNickName"`
	ToUserFaceUrl    string `json:"toUserFaceUrl"`
	AddTime          int64  `json:"addTime"`
	RequestMsg       string `json:"requestMsg"`
	HandledMsg       string `json:"handledMsg"`
	Type             int32  `json:"type"`
	HandleStatus     int32  `json:"handleStatus"`
	HandleResult     int32  `json:"handleResult"`
}

func newGroupApplicationResponse(params *paramsGroupApplicationResponse) *group.GroupApplicationResponseReq {
	pbData := group.GroupApplicationResponseReq{
		OperationID:      params.OperationID,
		GroupID:          params.GroupID,
		FromUserID:       params.FromUserID,
		FromUserNickName: params.FromUserNickName,
		FromUserFaceUrl:  params.FromUserFaceUrl,
		ToUserID:         params.ToUserID,
		ToUserNickName:   params.ToUserNickName,
		ToUserFaceUrl:    params.ToUserFaceUrl,
		AddTime:          params.AddTime,
		RequestMsg:       params.RequestMsg,
		HandledMsg:       params.HandledMsg,
		Type:             params.Type,
		HandleStatus:     params.HandleStatus,
		HandleResult:     params.HandleResult,
	}
	return &pbData
}

func ApplicationGroupResponse(c *gin.Context) {
	log.Info("", "", "api GroupApplicationResponse init ....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := group.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsGroupApplicationResponse{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	pbData := newGroupApplicationResponse(&params)

	token := c.Request.Header.Get("token")
	if claims, err := utils.ParseToken(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	} else {
		pbData.OwnerID = claims.UID
	}

	log.Info("", "", "api GroupApplicationResponse is server, [data: %s]", pbData.String())
	reply, err := client.GroupApplicationResponse(context.Background(), pbData)
	if err != nil {
		log.Error("", "", "api GroupApplicationResponse call rpc fail, [data: %s] [err: %s]", pbData.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info("", "", "api GroupApplicationResponse call rpc success, [data: %s] [reply: %s]", pbData.String(), reply.String())

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
	})

}
