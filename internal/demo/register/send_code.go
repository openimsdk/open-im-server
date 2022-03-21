package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"math/rand"
	"net/http"
	"time"
)

type paramsVerificationCode struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	OperationID string `json:"operationID" binding:"required"`
	UsedFor     int `json:"usedFor"`
}

func SendVerificationCode(c *gin.Context) {
	params := paramsVerificationCode{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("", "BindJSON failed", "err:", err.Error(), "phoneNumber", params.PhoneNumber, "email", params.Email)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	var account string
	if params.Email != "" {
		account = params.Email
	} else {
		account = params.PhoneNumber
	}
	var accountKey string
	if params.UsedFor == 0 {
		params.UsedFor = constant.VerificationCodeForRegister
	}
	switch params.UsedFor {
	case constant.VerificationCodeForRegister:
		_, err := im_mysql_model.GetRegister(account)
		if err == nil {
			log.NewError(params.OperationID, "The phone number has been registered", params)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.HasRegistered, "errMsg": "The phone number has been registered"})
			return
		}
		ok, err := db.DB.JudgeAccountEXISTS(account)
		if ok || err != nil {
			log.NewError(params.OperationID, "The phone number has been registered", params)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.RepeatSendCode, "errMsg": "The phone number has been registered"})
			return
		}
		accountKey = account + "_" + constant.VerificationCodeForRegisterSuffix

	case constant.VerificationCodeForReset:
		accountKey = account + "_" + constant.VerificationCodeForResetSuffix
	}
	rand.Seed(time.Now().UnixNano())
	code := 100000 + rand.Intn(900000)
	log.NewInfo(params.OperationID, params.UsedFor,"begin store redis", accountKey, code)
	err := db.DB.SetAccountCode(accountKey, code, config.Config.Demo.CodeTTL)
	if err != nil {
		log.NewError(params.OperationID, "set redis error", accountKey, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		return
	}
	log.NewDebug("", config.Config.Demo)
	if params.Email != "" {
		m := gomail.NewMessage()
		m.SetHeader(`From`, config.Config.Demo.Mail.SenderMail)
		m.SetHeader(`To`, []string{account}...)
		m.SetHeader(`Subject`, config.Config.Demo.Mail.Title)
		m.SetBody(`text/html`, fmt.Sprintf("%d", code))
		if err := gomail.NewDialer(config.Config.Demo.Mail.SmtpAddr, config.Config.Demo.Mail.SmtpPort, config.Config.Demo.Mail.SenderMail, config.Config.Demo.Mail.SenderAuthorizationCode).DialAndSend(m); err != nil {
			log.ErrorByKv("send mail error", account, "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.MailSendCodeErr, "errMsg": ""})
			return
		}
	} else {
		client, err := CreateClient(tea.String(config.Config.Demo.AliSMSVerify.AccessKeyID), tea.String(config.Config.Demo.AliSMSVerify.AccessKeySecret))
		if err != nil {
			log.NewError(params.OperationID, "create sendSms client err", "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
			return
		}

		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			PhoneNumbers:  tea.String(account),
			SignName:      tea.String(config.Config.Demo.AliSMSVerify.SignName),
			TemplateCode:  tea.String(config.Config.Demo.AliSMSVerify.VerificationCodeTemplateCode),
			TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%d\"}", code)),
		}

		response, err := client.SendSms(sendSmsRequest)
		if err != nil {
			log.NewError(params.OperationID, "sendSms error", account, "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
			return
		}
		if *response.Body.Code != "OK" {
			log.NewError(params.OperationID, "alibabacloud sendSms error", account, "err", response.Body.Code, response.Body.Message)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
			return
		}
	}

	data := make(map[string]interface{})
	data["account"] = account
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code has been set!", "data": data})
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (result *dysmsapi20170525.Client, err error) {
	c := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}

	// 访问的域名
	c.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result = &dysmsapi20170525.Client{}
	result, err = dysmsapi20170525.NewClient(c)
	return result, err
}
