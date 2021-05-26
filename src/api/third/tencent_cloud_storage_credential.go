package apiThird

import (
	"Open_IM/src/common/config"
	log2 "Open_IM/src/common/log"
	"github.com/gin-gonic/gin"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"net/http"
	"time"
)

type paramsTencentCloudStorageCredential struct {
	Token       string `json:"token"`
	OperationID string `json:"operationID"`
}

func TencentCloudStorageCredential(c *gin.Context) {
	params := paramsTencentCloudStorageCredential{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "Parameter parsing error，please check the parameters and request service again"})
		return
	}

	log2.Info(params.Token, params.OperationID, "api TencentUpLoadCredential call start...")

	cli := sts.NewClient(
		config.Config.Credential.Tencent.SecretID,
		config.Config.Credential.Tencent.SecretKey,
		nil,
	)
	log2.Info(c.Request.Header.Get("token"), c.PostForm("optionID"), "api TencentUpLoadCredential sts.NewClient cli = %v", cli)

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          config.Config.Credential.Tencent.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:" + config.Config.Credential.Tencent.Region + ":uid/" + config.Config.Credential.Tencent.AppID + ":" + config.Config.Credential.Tencent.Bucket + "/*",
					},
				},
			},
		},
	}
	log2.Info(c.Request.Header.Get("token"), c.PostForm("optionID"), "api TencentUpLoadCredential sts.CredentialOptions opt = %v", opt)

	res, err := cli.GetCredential(opt)
	if err != nil {
		log2.Error(c.Request.Header.Get("token"), c.PostForm("optionID"), "api TencentUpLoadCredential cli.GetCredential err = %s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"errCode": config.ErrTencentCredential.ErrCode,
			"errMsg":  err.Error(),
			"data":    res,
		})
		return
	}
	log2.Info(c.Request.Header.Get("token"), c.PostForm("optionID"), "api TencentUpLoadCredential cli.GetCredential success res = %v, res.Credentials = %v", res, res.Credentials)

	c.JSON(http.StatusOK, gin.H{
		"errCode": 0,
		"errMsg":  "",
		"data":    res,
	})
}
