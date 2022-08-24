package register

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ParamsLogin struct {
	UserID      string `json:"userID"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Platform    int32  `json:"platform"`
	OperationID string `json:"operationID" binding:"required"`
	AreaCode    string `json:"areaCode"`
}

func Login(c *gin.Context) {
	params := ParamsLogin{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	var account string
	if params.Email != "" {
		account = params.Email
	} else if params.PhoneNumber != "" {
		account = params.PhoneNumber
	} else {
		account = params.UserID
	}

	r, err := im_mysql_model.GetRegister(account, params.AreaCode, params.UserID)
	if err != nil {
		log.NewError(params.OperationID, "user have not register", params.Password, account, err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NotRegistered, "errMsg": "Mobile phone number is not registered"})
		return
	}
	if r.Password != params.Password {
		log.NewError(params.OperationID, "password  err", params.Password, account, r.Password, r.Account)
		c.JSON(http.StatusOK, gin.H{"errCode": constant.PasswordErr, "errMsg": "password err"})
		return
	}
	var userID string
	if r.UserID != "" {
		userID = r.UserID
	} else {
		userID = r.Account
	}
	ip := c.Request.Header.Get("X-Forward-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	url := fmt.Sprintf("http://%s:%d/auth/user_token", utils.ServerIP, config.Config.Api.GinPort[0])
	openIMGetUserToken := api.UserTokenReq{}
	openIMGetUserToken.OperationID = params.OperationID
	openIMGetUserToken.Platform = params.Platform
	openIMGetUserToken.Secret = config.Config.Secret
	openIMGetUserToken.UserID = userID
	openIMGetUserToken.LoginIp = ip
	loginIp := c.Request.Header.Get("X-Forward-For")
	if loginIp == "" {
		loginIp = c.ClientIP()
	}
	openIMGetUserToken.LoginIp = loginIp
	openIMGetUserTokenResp := api.UserTokenResp{}
	bMsg, err := http2.Post(url, openIMGetUserToken, 2)
	if err != nil {
		log.NewError(params.OperationID, "request openIM get user token error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.GetIMTokenErr, "errMsg": err.Error()})
		return
	}
	err = json.Unmarshal(bMsg, &openIMGetUserTokenResp)
	if err != nil || openIMGetUserTokenResp.ErrCode != 0 {
		log.NewError(params.OperationID, "request get user token", account, "err", "")
		if openIMGetUserTokenResp.ErrCode == constant.LoginLimit {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.LoginLimit, "errMsg": "用户登录被限制"})
		} else {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.GetIMTokenErr, "errMsg": ""})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": openIMGetUserTokenResp.UserToken})

}
