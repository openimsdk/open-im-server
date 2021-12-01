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
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"math/rand"
	"net/http"
	"time"
)

type paramsVerificationCode struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

func SendVerificationCode(c *gin.Context) {
	log.InfoByKv("sendCode api is statrting...", "")
	params := paramsVerificationCode{}

	if err := c.BindJSON(&params); err != nil {
		log.ErrorByKv("request params json parsing failed", params.PhoneNumber, params.Email, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	var account string
	if params.Email != "" {
		account = params.Email
	} else {
		account = params.PhoneNumber
	}

	queryParams := im_mysql_model.GetRegisterParams{
		Account: account,
	}
	_, err, rowsAffected := im_mysql_model.GetRegister(&queryParams)

	if err == nil && rowsAffected != 0 {
		log.ErrorByKv("The phone number has been registered", queryParams.Account, "err")
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "The phone number has been registered"})
		return
	}

	log.InfoByKv("begin sendSms", account)
	rand.Seed(time.Now().UnixNano())
	code := 100000 + rand.Intn(900000)
	log.NewDebug("", config.Config.Demo)
	if params.Email != "" {
		m := gomail.NewMessage()
		m.SetHeader(`From`, config.Config.Demo.Mail.SenderMail)
		m.SetHeader(`To`, []string{account}...)
		m.SetHeader(`Subject`, config.Config.Demo.Mail.Title)
		m.SetBody(`text/html`, fmt.Sprintf("%d", code))
		if err := gomail.NewDialer(config.Config.Demo.Mail.SmtpAddr, config.Config.Demo.Mail.SmtpPort, config.Config.Demo.Mail.SenderMail, config.Config.Demo.Mail.SenderAuthorizationCode).DialAndSend(m); err != nil {
			log.ErrorByKv("send mail error", account, "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.IntentionalError, "errMsg": err.Error()})
			return
		}
	} else {
		client, err := CreateClient(tea.String(config.Config.Demo.AliSMSVerify.AccessKeyID), tea.String(config.Config.Demo.AliSMSVerify.AccessKeySecret))
		if err != nil {
			log.ErrorByKv("create sendSms client err", "", "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.IntentionalError, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
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
			log.ErrorByKv("sendSms error", account, "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.IntentionalError, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
			return
		}
		if *response.Body.Code != "OK" {
			log.ErrorByKv("alibabacloud sendSms error", account, "err", response.Body.Code, response.Body.Message)
			c.JSON(http.StatusOK, gin.H{"errCode": constant.IntentionalError, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
			return
		}
	}

	log.InfoByKv("begin store redis", account)
	v, err := redis.Int(db.DB.Exec("TTL", account))
	if err != nil {
		log.ErrorByKv("get account from redis error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.IntentionalError, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		return
	}
	switch {
	case v == -2:
		_, err = db.DB.Exec("SET", account, code, "EX", 600)
		if err != nil {
			log.ErrorByKv("set redis error", account, "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": constant.IntentionalError, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
			return
		}
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code sent successfully!", "data": data})
		log.InfoByKv("send new verification code", account)
		return
	case v > 540:
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "Frequent operation!", "data": data})
		log.InfoByKv("frequent operation", account)
		return
	case v < 540:
		_, err = db.DB.Exec("SET", account, code, "EX", 600)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"errCode": constant.IntentionalError, "errMsg": "Enterthe superCode directly in the verification code box, SuperCode can be configured in config.xml"})
			return
		}
		data := make(map[string]interface{})
		data["account"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "Verification code has been reset!", "data": data})
		log.InfoByKv("Reset verification code", account)
		return
	}

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
