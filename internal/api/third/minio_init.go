package apiThird

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	url2 "net/url"
)

var (
	MinioClient *minio.Client
)

func MinioInit() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "minio config: ", config.Config.Credential.Minio)
	var initUrl string
	if config.Config.Credential.Minio.EndpointInnerEnable {
		initUrl = config.Config.Credential.Minio.EndpointInner
	} else {
		initUrl = config.Config.Credential.Minio.Endpoint
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "use initUrl: ", initUrl)
	minioUrl, err := url2.Parse(initUrl)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "parse failed, please check config/config.yaml", err.Error())
		return
	}
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(config.Config.Credential.Minio.AccessKeyID, config.Config.Credential.Minio.SecretAccessKey, ""),
	}
	if minioUrl.Scheme == "http" {
		opts.Secure = false
	} else if minioUrl.Scheme == "https" {
		opts.Secure = true
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "Parse ok ", config.Config.Credential.Minio)
	MinioClient, err = minio.New(minioUrl.Host, opts)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "new ok ", config.Config.Credential.Minio)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "init minio client failed", err.Error())
		return
	}
	opt := minio.MakeBucketOptions{
		Region:        config.Config.Credential.Minio.Location,
		ObjectLocking: false,
	}
	err = MinioClient.MakeBucket(context.Background(), config.Config.Credential.Minio.Bucket, opt)
	if err != nil {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "MakeBucket failed ", err.Error())
		exists, err := MinioClient.BucketExists(context.Background(), config.Config.Credential.Minio.Bucket)
		if err == nil && exists {
			log.NewInfo(operationID, utils.GetSelfFuncName(), "We already own ", config.Config.Credential.Minio.Bucket)
		} else {
			if err != nil {
				log.NewInfo(operationID, utils.GetSelfFuncName(), err.Error())
			}
			log.NewInfo(operationID, utils.GetSelfFuncName(), "create bucket failed and bucket not exists")
			return
		}
	}
	// make app bucket
	err = MinioClient.MakeBucket(context.Background(), config.Config.Credential.Minio.AppBucket, opt)
	if err != nil {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "MakeBucket failed ", err.Error())
		exists, err := MinioClient.BucketExists(context.Background(), config.Config.Credential.Minio.Bucket)
		if err == nil && exists {
			log.NewInfo(operationID, utils.GetSelfFuncName(), "We already own ", config.Config.Credential.Minio.Bucket)
		} else {
			if err != nil {
				log.NewInfo(operationID, utils.GetSelfFuncName(), err.Error())
			}
			log.NewInfo(operationID, utils.GetSelfFuncName(), "create bucket failed and bucket not exists")
			return
		}
	}
	// 自动化桶public的代码
	policyJsonString := fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject","s3:PutObject"],
		"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"],"Sid": ""}]}`, config.Config.Credential.Minio.Bucket)
	err = MinioClient.SetBucketPolicy(context.Background(), config.Config.Credential.Minio.Bucket, policyJsonString)
	if err != nil {
		log.NewInfo("", utils.GetSelfFuncName(), "SetBucketPolicy failed please set in web", err.Error())
	}
	policyJsonString = fmt.Sprintf(`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject","s3:PutObject"],
		"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"],"Sid": ""}]}`, config.Config.Credential.Minio.AppBucket)
	err = MinioClient.SetBucketPolicy(context.Background(), config.Config.Credential.Minio.AppBucket, policyJsonString)
	if err != nil {
		log.NewInfo("", utils.GetSelfFuncName(), "SetBucketPolicy failed please set in web", err.Error())
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "minio create and set policy success")
}
