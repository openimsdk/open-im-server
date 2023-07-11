package rpcclient

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"google.golang.org/grpc"
)

type Third struct {
	conn        grpc.ClientConnInterface
	Client      third.ThirdClient
	discov      discoveryregistry.SvcDiscoveryRegistry
	MinioClient *minio.Client
}

func NewThird(discov discoveryregistry.SvcDiscoveryRegistry) *Third {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImThirdName)
	if err != nil {
		panic(err)
	}
	client := third.NewThirdClient(conn)
	minioClient, err := minioInit()
	return &Third{discov: discov, Client: client, conn: conn, MinioClient: minioClient}
}

func minioInit() (*minio.Client, error) {
	minioClient := &minio.Client{}
	var initUrl string
	initUrl = config.Config.Object.Minio.Endpoint
	minioUrl, err := url.Parse(initUrl)
	if err != nil {
		return nil, err
	}
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(config.Config.Object.Minio.AccessKeyID, config.Config.Object.Minio.SecretAccessKey, ""),
		//Region: config.Config.Credential.Minio.Location,
	}
	if minioUrl.Scheme == "http" {
		opts.Secure = false
	} else if minioUrl.Scheme == "https" {
		opts.Secure = true
	}
	minioClient, err = minio.New(minioUrl.Host, opts)
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}
