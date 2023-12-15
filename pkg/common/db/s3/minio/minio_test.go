package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	config.Config.Object.Minio.Bucket = "openim"
	config.Config.Object.Minio.AccessKeyID = "root"
	config.Config.Object.Minio.SecretAccessKey = "openIM123"
	config.Config.Object.Minio.Endpoint = "http://172.16.8.38:10005"
	tmp, err := NewMinio(nil)
	if err != nil {
		panic(err)
	}
	min := tmp.(*Minio)
	cli := min.core.Client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	policy := minio.NewPostPolicy()
	_ = policy.SetExpires(time.Now().Add(time.Hour))
	_ = policy.SetKey("test.txt")
	_ = policy.SetBucket(config.Config.Object.Minio.Bucket)
	policy.SetContentType("text/plain")
	u, fd, err := cli.PresignedPostPolicy(ctx, policy)
	if err != nil {
		panic(err)
	}
	t.Log(u)
	for k, v := range fd {
		t.Log(k, v)
	}
}
