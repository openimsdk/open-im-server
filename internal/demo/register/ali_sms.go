package register

import (
	"Open_IM/pkg/common/config"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

type AliSMS struct {
	client *dysmsapi20170525.Client
}

func (a AliSMS) SendSms(code int, phoneNumber string) (resp interface{}, err error) {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumber),
		SignName:      tea.String(config.Config.Demo.AliSMSVerify.SignName),
		TemplateCode:  tea.String(config.Config.Demo.AliSMSVerify.VerificationCodeTemplateCode),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%d\"}", code)),
	}
	response, err := a.client.SendSms(sendSmsRequest)
	if err != nil {
		//log.NewError(params.OperationID, "sendSms error", account, "err", err.Error())
		//c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		return resp, err
	}
	if *response.Body.Code != "OK" {
		//log.NewError(params.OperationID, "alibabacloud sendSms error", account, "err", response.Body.Code, response.Body.Message)
		//c.JSON(http.StatusOK, gin.H{"errCode": constant.SmsSendCodeErr, "errMsg": "Enter the superCode directly in the verification code box, SuperCode can be configured in config.xml"})
		//return
		return resp, errors.New("alibabacloud sendSms error")
	}
	return resp, nil
}

func NewAliSMS() (*AliSMS, error) {
	var a AliSMS
	client, err := createClient(tea.String(config.Config.Demo.AliSMSVerify.AccessKeyID), tea.String(config.Config.Demo.AliSMSVerify.AccessKeySecret))
	if err != nil {
		return &a, err
	}
	a.client = client
	return &a, nil
}
func createClient(accessKeyId *string, accessKeySecret *string) (result *dysmsapi20170525.Client, err error) {
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
