package apiThird

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	url2 "net/url"
)

func MinioInit() {
	minioUrl, err := url2.Parse(config.Config.Credential.Minio.Endpoint)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "parse failed, please check config/config.yaml", err.Error())
		return
	}
	minioClient, err := minio.New(minioUrl.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Config.Credential.Minio.AccessKeyID, config.Config.Credential.Minio.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "init minio client failed", err.Error())
		return
	}
	opt := minio.MakeBucketOptions{
		Region:        config.Config.Credential.Minio.Location,
		ObjectLocking: false,
	}
	err = minioClient.MakeBucket(context.Background(), config.Config.Credential.Minio.Bucket, opt)
	if err != nil {
		log.NewInfo("", utils.GetSelfFuncName(), err.Error())
		exists, err := minioClient.BucketExists(context.Background(), config.Config.Credential.Minio.Bucket)
		if err == nil && exists {
			log.NewInfo("", utils.GetSelfFuncName(), "We already own %s\n", config.Config.Credential.Minio.Bucket)
		} else {
			if err != nil {
				log.NewError("", utils.GetSelfFuncName(), err.Error())
			}
			log.NewError("", utils.GetSelfFuncName(), "create bucket failed and bucket not exists")
			return
		}
	}
	// 自动化桶public的代码
	//err = minioClient.SetBucketPolicy(context.Background(), config.Config.Credential.Minio.Bucket, policy.BucketPolicyReadWrite)
	//if err != nil {
	//	log.NewError("", utils.GetSelfFuncName(), "SetBucketPolicy failed please set in 	", err.Error())
	//	return`z
	//}
	log.NewInfo("", utils.GetSelfFuncName(), "minio create and set policy success")
}
