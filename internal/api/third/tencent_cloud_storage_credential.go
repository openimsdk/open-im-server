package apiThird

import (
	"Open_IM/pkg/common/config"
	log2 "Open_IM/pkg/common/log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// paramsTencentCloudStorageCredential struct
type paramsTencentCloudStorageCredential struct {
	Token       string `json:"token"`
	OperationID string `json:"operationID"`
}

// resultTencentCredential struct
type resultTencentCredential struct {
	ErrCode int         `json:"errCode`
	ErrMsg  string      `json:"errMsg"`
	Region  string      `json:"region"`
	Bucket  string      `json:"bucket"`
	Data    interface{} `json:"data"`
}

var lastTime int64
var lastRes *sts.CredentialResult

// @Summary
// @Schemes
// @Description get Tencent cloud storage credential
// @Tags third
// @Accept json
// @Produce json
// @Param body body apiThird.paramsTencentCloudStorageCredential true "get Tencent cloud storage credential params"
// @Param token header string true "token"
// @Success 200 {object} apiThird.resultTencentCredential
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /third/user_register [post]
func TencentCloudStorageCredential(c *gin.Context) {
	params := paramsTencentCloudStorageCredential{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "Parameter parsing error，please check the parameters and request service again"})
		return
	}

	log2.Info(params.Token, params.OperationID, "api TencentUpLoadCredential call start...")

	if time.Now().Unix()-lastTime < 10 && lastRes != nil {
		c.JSON(http.StatusOK, gin.H{
			"errCode": 0,
			"errMsg":  "",
			"region":  config.Config.Credential.Tencent.Region,
			"bucket":  config.Config.Credential.Tencent.Bucket,
			"data":    lastRes,
		})
		return
	}

	lastTime = time.Now().Unix()

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
			"bucket":  "",
			"region":  "",
			"data":    res,
		})
		return
	}
	log2.Info(c.Request.Header.Get("token"), c.PostForm("optionID"), "api TencentUpLoadCredential cli.GetCredential success res = %v, res.Credentials = %v", res, res.Credentials)

	lastRes = res

	c.JSON(http.StatusOK, gin.H{
		"errCode": 0,
		"errMsg":  "",
		"region":  config.Config.Credential.Tencent.Region,
		"bucket":  config.Config.Credential.Tencent.Bucket,
		"data":    res,
	})
}
