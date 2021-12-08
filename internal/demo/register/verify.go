package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

type paramsCertification struct {
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	VerificationCode string `json:"verificationCode"`
}

func Verify(c *gin.Context) {
	log.InfoByKv("Verify api is statrting...", "")
	params := paramsCertification{}

	if err := c.BindJSON(&params); err != nil {
		log.ErrorByKv("request params json parsing failed", "", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	var account string
	if params.Email != "" {
		account = params.Email
	} else {
		account = params.PhoneNumber
	}

	if params.VerificationCode == config.Config.Demo.SuperCode {
		log.InfoByKv("Super Code Verified successfully", account)
		data := make(map[string]interface{})
		data["account"] = account
		data["verificationCode"] = params.VerificationCode
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verified successfully!", "data": data})
		return
	}
	log.NewInfo("0", "params.VerificationCode != config.Config.Demo.SuperCode", params.VerificationCode, config.Config.Demo)
	log.InfoByKv("begin get form redis", account)
	v, err := redis.String(db.DB.Exec("GET", account))
	log.InfoByKv("redis phone number and verificating Code", account, v)
	if err != nil {
		log.ErrorByKv("Verification code expired", account, "err", err.Error())
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "Verification code expired!", "data": data})
		return
	}
	if params.VerificationCode == v {
		log.InfoByKv("Verified successfully", account)
		data := make(map[string]interface{})
		data["account"] = account
		data["verificationCode"] = params.VerificationCode
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verified successfully!", "data": data})
		return
	} else {
		log.InfoByKv("Verification code error", account, params.VerificationCode)
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "Verification code error!", "data": data})
	}
}
