package objstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go"
	"net/url"
	"time"
)

func NewMinio() Interface {
	return &minioImpl{}
}

type minioImpl struct {
	uploadBucket    string // 上传桶
	permanentBucket string // 永久桶
	clearBucket     string // 自动清理桶
	client          *minio.Client
}

func (m *minioImpl) Init() error {
	client, err := minio.New("127.0.0.1:9000", "minioadmin", "minioadmin", false)
	if err != nil {
		return fmt.Errorf("minio client error: %w", err)
	}
	m.client = client
	m.uploadBucket = "upload"
	m.permanentBucket = "permanent"
	m.clearBucket = "clear"
	return nil
}

func (m *minioImpl) Name() string {
	return "minio"
}

func (m *minioImpl) MinMultipartSize() int64 {
	return 1024 * 1024 * 5 // minio.absMinPartSize
}

func (m *minioImpl) UploadBucket() string {
	return m.uploadBucket
}

func (m *minioImpl) PermanentBucket() string {
	return m.permanentBucket
}

func (m *minioImpl) ClearBucket() string {
	return m.clearBucket
}

func (m *minioImpl) urlReplace(u *url.URL) {

}

func (m *minioImpl) ApplyPut(ctx context.Context, args *ApplyPutArgs) (*PutRes, error) {
	if args.Effective <= 0 {
		return nil, errors.New("EffectiveTime <= 0")
	}
	_, err := m.GetObjectInfo(ctx, &BucketFile{
		Bucket: m.uploadBucket,
		Name:   args.Name,
	})
	if err == nil {
		return nil, fmt.Errorf("minio bucket %s name %s already exists", args.Bucket, args.Name)
	} else if !m.IsNotFound(err) {
		return nil, err
	}
	effective := time.Now().Add(args.Effective)
	u, err := m.client.PresignedPutObject(m.uploadBucket, args.Name, args.Effective)
	if err != nil {
		return nil, fmt.Errorf("minio apply error: %w", err)
	}
	m.urlReplace(u)
	return &PutRes{
		URL:           u.String(),
		Bucket:        m.uploadBucket,
		Name:          args.Name,
		EffectiveTime: effective,
	}, nil
}

func (m *minioImpl) GetObjectInfo(ctx context.Context, args *BucketFile) (*ObjectInfo, error) {
	info, err := m.client.StatObject(args.Bucket, args.Name, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{
		URL:  "", // todo
		Size: info.Size,
		Hash: info.ETag,
	}, nil
}

func (m *minioImpl) CopyObjectInfo(ctx context.Context, src *BucketFile, dst *BucketFile) error {
	destination, err := minio.NewDestinationInfo(dst.Bucket, dst.Name, nil, nil)
	if err != nil {
		return err
	}
	return m.client.CopyObject(destination, minio.NewSourceInfo(src.Bucket, src.Name, nil))
}

func (m *minioImpl) DeleteObjectInfo(ctx context.Context, info *BucketFile) error {
	return m.client.RemoveObject(info.Bucket, info.Name)
}

func (m *minioImpl) MoveObjectInfo(ctx context.Context, src *BucketFile, dst *BucketFile) error {
	if err := m.CopyObjectInfo(ctx, src, dst); err != nil {
		return err
	}
	return m.DeleteObjectInfo(ctx, src)
}

func (m *minioImpl) MergeObjectInfo(ctx context.Context, src []BucketFile, dst *BucketFile) error {
	switch len(src) {
	case 0:
		return errors.New("src empty")
	case 1:
		return m.CopyObjectInfo(ctx, &src[0], dst)
	}
	destination, err := minio.NewDestinationInfo(dst.Bucket, dst.Name, nil, nil)
	if err != nil {
		return err
	}
	sources := make([]minio.SourceInfo, len(src))
	for i, s := range src {
		sources[i] = minio.NewSourceInfo(s.Bucket, s.Name, nil)
	}
	return m.client.ComposeObject(destination, sources) // todo
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
