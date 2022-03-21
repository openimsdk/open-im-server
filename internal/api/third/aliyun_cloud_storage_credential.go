package apiThird

import (
	"Open_IM/pkg/common/config"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/gin-gonic/gin"
	"net/http"
)

type paramsAliyunCloudStorageCredential struct {
	Token       string `form:"token"`
	OperationID string `form:"operationID"`
}

func AliyunCloudStorageCredential(c *gin.Context) {
	params := paramsAliyunCloudStorageCredential{}
	if err := c.BindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "Parameter parsing error，please check the parameters and request service again"})
		return
	}

	credential := config.Config.Credential.Aliyun
	if credential.AccessKeyId == "" || credential.AccessKeySecret == "" || credential.Bucket == "" || credential.Region == "" || credential.RegionId == "" || credential.RoleArn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "Aliyun OSS config error"})
		return
	}
	//构建一个阿里云客户端, 用于发起请求。
	//构建阿里云客户端时，需要设置AccessKey ID和AccessKey Secret。
	client, err := sts.NewClientWithAccessKey(credential.RegionId, credential.AccessKeyId, credential.AccessKeySecret)

	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	type CredentialPolicyStatement struct {
		Action    []string                          `json:"action,omitempty"`
		Effect    string                            `json:"effect,omitempty"`
		Resource  []string                          `json:"resource,omitempty"`
		Condition map[string]map[string]interface{} `json:"condition,omitempty"`
	}

	type CredentialPolicy struct {
		Version   string                      `json:"version,omitempty"`
		Statement []CredentialPolicyStatement `json:"statement,omitempty"`
	}

	//设置参数。关于参数含义和设置方法，请参见《API参考》。
	request.RoleArn = credential.RoleArn
	request.RoleSessionName = params.OperationID
	request.Policy = "{\n    \"Version\": \"1\",\n    \"Statement\": [\n        {\n            \"Effect\": \"Allow\",\n            \"Action\": [\n                \"oss:PutObject\"\n            ],\n            \"Resource\": \"acs:oss:*:*:*\"\n        }\n    ]\n}"
	//request.DurationSeconds = "900"

	//发起请求，并得到响应。
	res, err := client.AssumeRole(request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errCode": config.ErrTencentCredential.ErrCode,
			"errMsg":  err.Error(),
			"bucket":  "",
			"region":  "",
			"data":    res,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"errCode": 0,
		"errMsg":  "",
		"region":  config.Config.Credential.Aliyun.Region,
		"bucket":  config.Config.Credential.Aliyun.Bucket,
		"data":    res,
	})
}
