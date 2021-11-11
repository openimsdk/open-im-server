package apiAuth

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAuth "Open_IM/pkg/proto/auth"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsUserRegister struct {
	Secret   string `json:"secret" binding:"required,max=32"`
	Platform int32  `json:"platform" binding:"required,min=1,max=7"`
	UID      string `json:"uid" binding:"required,min=1,max=64"`
	Name     string `json:"name" binding:"required,min=1,max=64"`
	Icon     string `json:"icon" binding:"omitempty,max=1024"`
	Gender   int32  `json:"gender" binding:"omitempty,oneof=0 1 2"`
	Mobile   string `json:"mobile" binding:"omitempty,max=32"`
	Birth    string `json:"birth" binding:"omitempty,max=16"`
	Email    string `json:"email" binding:"omitempty,max=64"`
	Ex       string `json:"ex" binding:"omitempty,max=1024"`
}

func newUserRegisterReq(params *paramsUserRegister) *pbAuth.UserRegisterReq {
	pbData := pbAuth.UserRegisterReq{
		UID:    params.UID,
		Name:   params.Name,
		Icon:   params.Icon,
		Gender: params.Gender,
		Mobile: params.Mobile,
		Birth:  params.Birth,
		Email:  params.Email,
		Ex:     params.Ex,
	}
	return &pbData
}

func UserRegister(c *gin.Context) {
	log.Info("", "", "api user_register init ....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName)
	client := pbAuth.NewAuthClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsUserRegister{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if params.Secret != config.Config.Secret {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "not authorized"})
		return
	}
	pbData := newUserRegisterReq(&params)

	log.Info("", "", "api user_register is server, [data: %s]", pbData.String())
	reply, err := client.UserRegister(context.Background(), pbData)
	if err != nil || !reply.Success {
		log.Error("", "", "api user_register call rpc fail, [data: %s] [err: %s]", pbData.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info("", "", "api user_register call rpc success, [data: %s] [reply: %s]", pbData.String(), reply.String())

	pbDataToken := &pbAuth.UserTokenReq{
		Platform: params.Platform,
		UID:      params.UID,
	}
	replyToken, err := client.UserToken(context.Background(), pbDataToken)
	if err != nil {
		log.Error("", "", "api user_register call rpc fail, [data: %s] [err: %s]", pbData.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info("", "", "api user_register call success, [data: %s] [reply: %s]", pbData.String(), reply.String())

	if replyToken.ErrCode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"errCode": replyToken.ErrCode,
			"errMsg":  replyToken.ErrMsg,
			"data": gin.H{
				"uid":         pbData.UID,
				"token":       replyToken.Token,
				"expiredTime": replyToken.ExpiredTime,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"errCode": replyToken.ErrCode,
			"errMsg":  replyToken.ErrMsg,
		})
	}
}
