package obj

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"net/http"
	"time"
)

func NewMinioClient() {

}

func NewMinioInterface() (Interface, error) {
	//client, err := minio.New("127.0.0.1:9000", &minio.Options{
	//	Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
	//	Secure: false,
	//})
	// todo 初始化连接和桶
	return &minioImpl{}, nil
}

type minioImpl struct {
	tempBucket      string // 上传桶
	permanentBucket string // 永久桶
	clearBucket     string // 自动清理桶
	urlstr          string // 访问地址
	client          *minio.Client
}

//func (m *minioImpl) Init() error {
//	client, err := minio.New("127.0.0.1:9000", &minio.Options{
//		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
//		Secure: false,
//	})
//	if err != nil {
//		return fmt.Errorf("minio client error: %w", err)
//	}
//	m.urlstr = "http://127.0.0.1:9000"
//	m.client = client
//	m.tempBucket = "temp"
//	m.permanentBucket = "permanent"
//	m.clearBucket = "clear"
//	return nil
//}

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

func (m *minioImpl) AppendHeader() http.Header {
	return map[string][]string{
		"x-amz-object-append": {"true"},
	}
}

func (m *minioImpl) TempBucket() string {
	return m.tempBucket
}

func (m *minioImpl) DataBucket() string {
	return m.permanentBucket
}

func (m *minioImpl) ClearBucket() string {
	return m.clearBucket
}

func (m *minioImpl) GetURL(bucket string, name string) string {
	return fmt.Sprintf("%s/%s/%s", m.urlstr, bucket, name)
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
		URL:  m.GetURL(args.Bucket, args.Name),
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
		return e.StatusCode == 404 && e.Code == "NoSuchKey"
	case *minio.ErrorResponse:
		return e.StatusCode == 404 && e.Code == "NoSuchKey"
	default:
		return false
	}
}

func (m *minioImpl) CheckName(name string) error {
	return s3utils.CheckValidObjectName(name)
}
