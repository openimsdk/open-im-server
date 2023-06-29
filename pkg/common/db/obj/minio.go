package obj

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"io"
	"net/http"
	"net/url"
	"time"
)

func NewMinioInterface() (Interface, error) {
	conf := config.Config.Object.Minio
	u, err := url.Parse(conf.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("minio endpoint parse %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("invalid minio endpoint scheme %s", u.Scheme)
	}
	client, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
		Secure: u.Scheme == "https",
	})
	if err != nil {
		return nil, fmt.Errorf("minio new client %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	for _, bucket := range utils.Distinct([]string{conf.TempBucket, conf.DataBucket}) {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return nil, fmt.Errorf("minio bucket %s exists %w", bucket, err)
		}
		if exists {
			continue
		}
		opt := minio.MakeBucketOptions{
			Region:        conf.Location,
			ObjectLocking: conf.IsDistributedMod,
		}
		if err := client.MakeBucket(ctx, bucket, opt); err != nil {
			return nil, fmt.Errorf("minio make bucket %s %w", bucket, err)
		}
	}
	return &minioImpl{
		client:     client,
		tempBucket: conf.TempBucket,
		dataBucket: conf.DataBucket,
	}, nil
}

type minioImpl struct {
	tempBucket string // 上传桶
	dataBucket string // 永久桶
	urlstr     string // 访问地址
	client     *minio.Client
}

func (m *minioImpl) Name() string {
	return "minio"
}

func (m *minioImpl) MinFragmentSize() int64 {
	return 1024 * 1024 * 5 // 每个分片最小大小 minio.absMinPartSize
}

func (m *minioImpl) MaxFragmentNum() int {
	return 1000 // 最大分片数量 minio.maxPartsCount
}

func (m *minioImpl) MinExpirationTime() time.Duration {
	return time.Hour * 24
}

func (m *minioImpl) TempBucket() string {
	return m.tempBucket
}

func (m *minioImpl) DataBucket() string {
	return m.dataBucket
}

func (m *minioImpl) PresignedGetURL(ctx context.Context, bucket string, name string, expires time.Duration, opt *HeaderOption) (string, error) {
	var reqParams url.Values
	if opt != nil {
		reqParams = make(url.Values)
		if opt.ContentType != "" {
			reqParams.Set("response-content-type", opt.ContentType)
		}
		if opt.Filename != "" {
			reqParams.Set("response-content-disposition", "attachment;filename="+opt.Filename)
		}
	}
	u, err := m.client.PresignedGetObject(ctx, bucket, name, expires, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (m *minioImpl) PresignedPutURL(ctx context.Context, args *ApplyPutArgs) (string, error) {
	if args.Effective <= 0 {
		return "", errors.New("EffectiveTime <= 0")
	}
	_, err := m.GetObjectInfo(ctx, &BucketObject{
		Bucket: m.tempBucket,
		Name:   args.Name,
	})
	if err == nil {
		return "", fmt.Errorf("minio bucket %s name %s already exists", args.Bucket, args.Name)
	} else if !m.IsNotFound(err) {
		return "", err
	}
	u, err := m.client.PresignedPutObject(ctx, m.tempBucket, args.Name, args.Effective)
	if err != nil {
		return "", fmt.Errorf("minio apply error: %w", err)
	}
	return u.String(), nil
}

func (m *minioImpl) GetObjectInfo(ctx context.Context, args *BucketObject) (*ObjectInfo, error) {
	info, err := m.client.StatObject(ctx, args.Bucket, args.Name, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{
		Size: info.Size,
		Hash: info.ETag,
	}, nil
}

func (m *minioImpl) CopyObject(ctx context.Context, src *BucketObject, dst *BucketObject) error {
	_, err := m.client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: dst.Bucket,
		Object: dst.Name,
	}, minio.CopySrcOptions{
		Bucket: src.Bucket,
		Object: src.Name,
	})
	return err
}

func (m *minioImpl) DeleteObject(ctx context.Context, info *BucketObject) error {
	return m.client.RemoveObject(ctx, info.Bucket, info.Name, minio.RemoveObjectOptions{})
}

func (m *minioImpl) MoveObjectInfo(ctx context.Context, src *BucketObject, dst *BucketObject) error {
	if err := m.CopyObject(ctx, src, dst); err != nil {
		return err
	}
	return m.DeleteObject(ctx, src)
}

func (m *minioImpl) ComposeObject(ctx context.Context, src []BucketObject, dst *BucketObject) error {
	destOptions := minio.CopyDestOptions{
		Bucket: dst.Bucket,
		Object: dst.Name + ".temp",
	}
	sources := make([]minio.CopySrcOptions, len(src))
	for i, s := range src {
		sources[i] = minio.CopySrcOptions{
			Bucket: s.Bucket,
			Object: s.Name,
		}
	}
	_, err := m.client.ComposeObject(ctx, destOptions, sources...)
	if err != nil {
		return err
	}
	return m.MoveObjectInfo(ctx, &BucketObject{
		Bucket: destOptions.Bucket,
		Name:   destOptions.Object,
	}, &BucketObject{
		Bucket: dst.Bucket,
		Name:   dst.Name,
	})
}

func (m *minioImpl) IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	switch e := err.(type) {
	case minio.ErrorResponse:
		return e.StatusCode == http.StatusNotFound || e.Code == "NoSuchKey"
	case *minio.ErrorResponse:
		return e.StatusCode == http.StatusNotFound || e.Code == "NoSuchKey"
	default:
		return false
	}
}

func (m *minioImpl) PutObject(ctx context.Context, info *BucketObject, reader io.Reader, size int64) (*ObjectInfo, error) {
	update, err := m.client.PutObject(ctx, info.Bucket, info.Name, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{
		Size: update.Size,
		Hash: update.ETag,
	}, nil
}

func (m *minioImpl) GetObject(ctx context.Context, info *BucketObject) (SizeReader, error) {
	object, err := m.client.GetObject(ctx, info.Bucket, info.Name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	stat, err := object.Stat()
	if err != nil {
		return nil, err
	}
	return NewSizeReader(object, stat.Size), nil
}

func (m *minioImpl) CheckName(name string) error {
	return s3utils.CheckValidObjectName(name)
}
