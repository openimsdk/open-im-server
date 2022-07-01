package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"

	"github.com/gin-gonic/gin"
	"net/http"
)

type paramsCertification struct {
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	VerificationCode string `json:"verificationCode"`
	OperationID      string `json:"operationID" binding:"required"`
	UsedFor          int    `json:"usedFor"`
	AreaCode         string `json:"areaCode"`
}

func Verify(c *gin.Context) {
	params := paramsCertification{}
	operationID := params.OperationID
	if operationID == "" {
		operationID = utils.OperationIDGenerator()
	}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(operationID, "request params json parsing failed", "", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(operationID, "recv req: ", params)

	var account string
	if params.Email != "" {
		account = params.Email
	} else {
		account = params.AreaCode + params.PhoneNumber
	}

	if params.VerificationCode == config.Config.Demo.SuperCode {
		log.NewInfo(params.OperationID, "Super Code Verified successfully", account)
		data := make(map[string]interface{})
		data["account"] = account
		data["verificationCode"] = params.VerificationCode
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verified successfully!", "data": data})
		return
	}
	log.NewInfo(params.OperationID, " params.VerificationCode != config.Config.Demo.SuperCode", params.VerificationCode, config.Config.Demo)
	log.NewInfo(params.OperationID, "begin get form redis", account)
	if params.UsedFor == 0 {
		params.UsedFor = constant.VerificationCodeForRegister
	}
	var accountKey string
	switch params.UsedFor {
	case constant.VerificationCodeForRegister:
		accountKey = params.AreaCode + account + "_" + constant.VerificationCodeForRegisterSuffix
	case constant.VerificationCodeForReset:
		accountKey = params.AreaCode + account + "_" + constant.VerificationCodeForResetSuffix
	}

	code, err := db.DB.GetAccountCode(accountKey)
	log.NewInfo(params.OperationID, "redis phone number and verificating Code", "key: ", accountKey, "code: ", code, "params: ", params)
	if err != nil {
		log.NewError(params.OperationID, "Verification code expired", accountKey, "err", err.Error())
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code expired!", "data": data})
		return
	}
	if params.VerificationCode == code {
		log.Info(params.OperationID, "Verified successfully", account)
		data := make(map[string]interface{})
		data["account"] = account
		data["verificationCode"] = params.VerificationCode
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verified successfully!", "data": data})
		return
	} else {
		log.Info(params.OperationID, "Verification code error", account, params.VerificationCode)
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.CodeInvalidOrExpired, "errMsg": "Verification code error!", "data": data})
	}
}
