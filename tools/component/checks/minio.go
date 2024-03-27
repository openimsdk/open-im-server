package checks

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/jsonutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/openimsdk/tools/log"

	s3minio "github.com/openimsdk/open-im-server/v3/pkg/common/db/s3/minio"
)

const (
	minioHealthCheckDuration = 1
	mongoConnTimeout         = 5 * time.Second
	MaxRetry                 = 300
)

type MinioConfig struct {
	s3minio.Config
	UseSSL string
	ApiURL string
}

// CheckMinio checks the MinIO connection.
func CheckMinio(minioStu MinioConfig) error {
	if minioStu.Endpoint == "" || minioStu.AccessKeyID == "" || minioStu.SecretAccessKey == "" {
		log.CInfo(nil, "Missing configuration for MinIO", "endpoint", minioStu.Endpoint, "accessKeyID", minioStu.AccessKeyID, "secretAccessKey", minioStu.SecretAccessKey)
		return errs.New("missing configuration for endpoint, accessKeyID, or secretAccessKey").Wrap()
	}

	minioInfo, err := jsonutil.JsonMarshal(minioStu)
	if err != nil {
		log.CInfo(nil, "MinioStu Marshal failed", "error", err)
		return errs.WrapMsg(err, "minioStu Marshal failed")
	}
	logJsonInfo := string(minioInfo)

	u, err := url.Parse(minioStu.Endpoint)
	if err != nil {
		log.CInfo(nil, "URL parse failed", "error", err, "minioInfo", logJsonInfo)
		return errs.WrapMsg(err, "url parse failed")
	}

	secure := u.Scheme == "https" || minioStu.UseSSL == "true"

	minioClient, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(minioStu.AccessKeyID, minioStu.SecretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		log.CInfo(nil, "Initialize MinIO client failed", "error", err, "minioInfo", logJsonInfo)
		return errs.WrapMsg(err, "initialize minio client failed")
	}

	cancel, err := minioClient.HealthCheck(time.Duration(minioHealthCheckDuration) * time.Second)
	if err != nil {
		log.CInfo(nil, "MinIO client health check failed", "error", err, "minioInfo", logJsonInfo)
		return errs.WrapMsg(err, "minio client health check failed")
	}
	defer cancel()

	if minioClient.IsOffline() {
		log.CInfo(nil, "MinIO client is offline", "minioInfo", logJsonInfo)
		return errors.New("minio client is offline")
	}

	apiURL, err := exactIP(minioStu.ApiURL)
	if err != nil {
		return err
	}
	signEndPoint, err := exactIP(minioStu.SignEndpoint)
	if err != nil {
		return err
	}

	if apiURL == "127.0.0.1" {
		log.CInfo(nil, "Warning, MinIOStu.apiURL contains localhost", "apiURL", minioStu.ApiURL)
	}
	if signEndPoint == "127.0.0.1" {
		log.CInfo(nil, "Warning, MinIOStu.signEndPoint contains localhost", "signEndPoint", minioStu.SignEndpoint)
	}
	return nil
}

func exactIP(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.CInfo(nil, "URL parse error", "error", err, "url", urlStr)
		return "", errs.WrapMsg(err, "url parse error")
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host // Assume the entire host part is the host name if split fails
	}
	if strings.HasSuffix(host, ":") {
		host = host[:len(host)-1]
	}
	return host, nil
}
