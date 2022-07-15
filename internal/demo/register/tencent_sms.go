package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"encoding/json"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	v20210111 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentSMS struct {
	client *v20210111.Client
}

func (t TencentSMS) SendSms(code int, phoneNumber string) (resp interface{}, err error) {
	request := v20210111.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(config.Config.Demo.TencentSMS.AppID)
	request.SignName = common.StringPtr(config.Config.Demo.TencentSMS.SignName)
	request.TemplateId = common.StringPtr(config.Config.Demo.TencentSMS.VerificationCodeTemplateCode)
	request.TemplateParamSet = common.StringPtrs([]string{utils.IntToString(code)})
	request.PhoneNumberSet = common.StringPtrs([]string{phoneNumber})
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := t.client.SendSms(request)
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		return response, err
	}
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return response, err
	}
	b, _ := json.Marshal(response.Response)
	log.Debug("tencent send message is ", code, phoneNumber, string(b))
	return response, nil
}

func NewTencentSMS() (*TencentSMS, error) {
	var a TencentSMS
	credential := common.NewCredential(
		config.Config.Demo.TencentSMS.SecretID,
		config.Config.Demo.TencentSMS.SecretKey,
	)
	cpf := profile.NewClientProfile()
	client, err := v20210111.NewClient(credential, config.Config.Demo.TencentSMS.Region, cpf)
	if err != nil {
		return &a, err
	}
	a.client = client
	return &a, nil

}
