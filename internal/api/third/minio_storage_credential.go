package apiThird

import (
	apiStruct "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/constant"
	http "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"

	"github.com/gin-gonic/gin"
	_ "github.com/minio/minio-go/v7/pkg/credentials"
)

func MinioStorageCredential(c *gin.Context) {
	var (
		req  apiStruct.MinioStorageCredentialReq
		resp apiStruct.MiniostorageCredentialResp
	)
	ok, _ := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError("", "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		http.RespHttp200(c, constant.ErrAccess, nil)
		return
	}
	//var stsOpts cr.STSAssumeRoleOptions
	//stsOpts.AccessKey = minioUsername
	//stsOpts.SecretKey = minioPassword
	log.NewInfo("0", req, resp)
	http.RespHttp200(c, constant.OK, nil)
}
