package checks

import (
	"context"
	"net"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"

	s3minio "github.com/openimsdk/open-im-server/v3/pkg/common/db/s3/minio"
)

const (
	minioHealthCheckDuration = 1 * time.Second
)

type MinioCheck struct {
	s3minio.Config
	UseSSL bool   `yaml:"useSSL"`
	ApiURL string `yaml:"apiURL"`
}

func CheckMinio(ctx context.Context, config MinioCheck) error {

	if config.Endpoint == "" || config.AccessKeyID == "" || config.SecretAccessKey == "" {
		logMsg := "Missing configuration for MinIO: endpoint, accessKeyID, or secretAccessKey"
		log.CInfo(ctx, logMsg, "Config", config)
		return errs.New(logMsg)
	}

	endpointURL, err := url.Parse(config.Endpoint)
	if err != nil {
		return errs.WrapMsg(err, "Failed to parse MinIO endpoint URL")
	}
	secure := endpointURL.Scheme == "https" || config.UseSSL

	minioClient, err := minio.New(endpointURL.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return errs.WrapMsg(err, "Failed to initialize MinIO client", "Endpoint", config.Endpoint)
	}

	cancel, err := minioClient.HealthCheck(minioHealthCheckDuration)
	if err != nil {
		return errs.WrapMsg(err, "MinIO client health check failed")
	}
	defer cancel()

	if minioClient.IsOffline() {
		return errs.New("minio client is offline").Wrap()
	}

	apiURLHost, _ := exactIP(config.ApiURL)
	signEndpointHost, _ := exactIP(config.SignEndpoint)
	if apiURLHost == "127.0.0.1" || signEndpointHost == "127.0.0.1" {
		logMsg := "Warning: MinIO ApiURL or SignEndpoint contains localhost"
		log.CInfo(ctx, logMsg, "ApiURL", config.ApiURL, "SignEndpoint", config.SignEndpoint)
	}

	return nil
}

func exactIP(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", errs.WrapMsg(err, "URL parse error")
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}
	return host, nil
}
