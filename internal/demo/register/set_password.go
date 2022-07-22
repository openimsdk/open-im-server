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
	UserID           string `json:"userID"`
	Email            string `json:"email"`
	Nickname         string `json:"nickname"`
	PhoneNumber      string `json:"phoneNumber"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"verificationCode"`
	Platform         int32  `json:"platform" binding:"required,min=1,max=7"`
	Ex               string `json:"ex"`
	FaceURL          string `json:"faceURL"`
	OperationID      string `json:"operationID" binding:"required"`
	AreaCode         string `json:"areaCode"`
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
	} else if params.PhoneNumber != "" {
		account = params.PhoneNumber
	} else {
		account = params.UserID
	}
	if params.Nickname == "" {
		params.Nickname = account
	}
	if params.UserID == "" {
		if (config.Config.Demo.UseSuperCode && params.VerificationCode != config.Config.Demo.SuperCode) || !config.Config.Demo.UseSuperCode {
			accountKey := params.AreaCode + account + "_" + constant.VerificationCodeForRegisterSuffix
			v, err := db.DB.GetAccountCode(accountKey)
			if err != nil || v != params.VerificationCode {
				log.NewError(params.OperationID, "password Verification code error", account, params.VerificationCode)
				data := make(map[string]interface{})
				data["PhoneNumber"] = account
				c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code error!", "data": data})
				return
			}
		}
	}
	//userID := utils.Base64Encode(account)
	var userID string
	if params.UserID == "" {
		userID = utils.Md5(params.OperationID + strconv.FormatInt(time.Now().UnixNano(), 10))
		bi := big.NewInt(0)
		bi.SetString(userID[0:8], 16)
		userID = bi.String()
	} else {
		userID = params.UserID
	}

	url := config.Config.Demo.ImAPIURL + "/auth/user_register"
	openIMRegisterReq := api.UserRegisterReq{}
	openIMRegisterReq.OperationID = params.OperationID
	openIMRegisterReq.Platform = params.Platform
	openIMRegisterReq.UserID = userID
	openIMRegisterReq.Nickname = params.Nickname
	openIMRegisterReq.Secret = config.Config.Secret
	openIMRegisterReq.FaceURL = params.FaceURL
	openIMRegisterResp := api.UserRegisterResp{}
	log.NewDebug(params.OperationID, utils.GetSelfFuncName(), "register req:", openIMRegisterReq)
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
	log.Info(params.OperationID, "begin store mysql", account, params.Password, "info", params.FaceURL, params.Nickname)
	err = im_mysql_model.SetPassword(account, params.Password, params.Ex, userID, params.AreaCode)
	if err != nil {
		log.NewError(params.OperationID, "set phone number password error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.RegisterFailed, "errMsg": err.Error()})
		return
	}
	log.Info(params.OperationID, "end setPassword", account, params.Password)
	// demo onboarding
	if params.UserID == "" {
		select {
		case Ch <- OnboardingProcessReq{
			OperationID: params.OperationID,
			UserID:      userID,
			NickName:    params.Nickname,
			FaceURL:     params.FaceURL,
			PhoneNumber: params.AreaCode + params.PhoneNumber,
			Email:       params.Email,
		}:
		case <-time.After(time.Second * 2):
			log.NewWarn(params.OperationID, utils.GetSelfFuncName(), "to ch timeOut")
		}
	}
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": openIMRegisterResp.UserToken})
	return
}
