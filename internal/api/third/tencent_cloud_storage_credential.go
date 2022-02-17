package apiThird

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"github.com/fatih/structs"

	//"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"net/http"
	"time"
)

func TencentCloudStorageCredential(c *gin.Context) {
	req := api.TencentCloudStorageCredentialReq{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(req.OperationID, "TencentCloudStorageCredential args ", userID)

	cli := sts.NewClient(
		config.Config.Credential.Tencent.SecretID,
		config.Config.Credential.Tencent.SecretKey,
		nil,
	)

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
	res, err := cli.GetCredential(opt)
	resp := api.TencentCloudStorageCredentialResp{}
	if err != nil {
		resp.ErrCode = constant.ErrTencentCredential.ErrCode
		resp.ErrMsg = err.Error()
	} else {
		resp.CosData.Bucket = config.Config.Credential.Tencent.Bucket
		resp.CosData.Region = config.Config.Credential.Tencent.Region
		resp.CosData.CredentialResult = res
	}

	resp.Data = structs.Map(&resp.CosData)
	log.NewInfo(req.OperationID, "TencentCloudStorageCredential return ", resp)

	c.JSON(http.StatusOK, resp)
}
