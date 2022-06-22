package register

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type ParamsSetPassword struct {
	Email            string `json:"email"`
	Nickname         string `json:"nickname"`
	PhoneNumber      string `json:"phoneNumber"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"verificationCode"`
	Platform         int32  `json:"platform" binding:"required,min=1,max=7"`
	Ex               string `json:"ex"`
	FaceURL          string `json:"faceURL"`
	OperationID      string `json:"operationID" binding:"required"`
}

func SetPassword(c *gin.Context) {
	params := ParamsSetPassword{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	var account string
	if params.Email != "" {
		account = params.Email
	} else {
		account = params.PhoneNumber
	}
	if params.Nickname == "" {
		params.Nickname = account
	}
	if params.VerificationCode != config.Config.Demo.SuperCode {
		accountKey := account + "_" + constant.VerificationCodeForRegisterSuffix
		v, err := db.DB.GetAccountCode(accountKey)
		if err != nil || v != params.VerificationCode {
			log.NewError(params.OperationID, "password Verification code error", account, params.VerificationCode)
			data := make(map[string]interface{})
			data["PhoneNumber"] = account
			c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code error!", "data": data})
			return
		}
	}
	//userID := utils.Base64Encode(account)

	userID := utils.Md5(params.OperationID + strconv.FormatInt(time.Now().UnixNano(), 10))
	bi := big.NewInt(0)
	bi.SetString(userID[0:8], 16)
	userID = bi.String()

	url := config.Config.Demo.ImAPIURL + "/auth/user_register"
	openIMRegisterReq := api.UserRegisterReq{}
	openIMRegisterReq.OperationID = params.OperationID
	openIMRegisterReq.Platform = params.Platform
	openIMRegisterReq.UserID = userID
	openIMRegisterReq.Nickname = params.Nickname
	openIMRegisterReq.Secret = config.Config.Secret
	openIMRegisterReq.FaceURL = params.FaceURL
	openIMRegisterResp := api.UserRegisterResp{}
	bMsg, err := http2.Post(url, openIMRegisterReq, 2)
	if err != nil {
		log.NewError(params.OperationID, "request openIM register error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": err.Error()})
		return
	}
	err = json.Unmarshal(bMsg, &openIMRegisterResp)
	if err != nil || openIMRegisterResp.ErrCode != 0 {
		log.NewError(params.OperationID, "request openIM register error", account, "err", "resp: ", openIMRegisterResp.ErrCode)
		if err != nil {
			log.NewError(params.OperationID, utils.GetSelfFuncName(), err.Error())
		}
		c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": "register failed: " + openIMRegisterResp.ErrMsg})
		return
	}
	log.Info(params.OperationID, "begin store mysql", account, params.Password, params.FaceURL, params.Nickname)
	err = im_mysql_model.SetPassword(account, params.Password, params.Ex, userID)
	if err != nil {
		log.NewError(params.OperationID, "set phone number password error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": err.Error()})
		return
	}
	log.Info(params.OperationID, "end setPassword", account, params.Password)
	// demo onboarding
	onboardingProcess(params.OperationID, userID, params.Nickname, params.FaceURL)
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": openIMRegisterResp.UserToken})
	return
}
