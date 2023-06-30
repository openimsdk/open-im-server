package obj

import (
	"context"
	"io"
	"net/http"
	"time"
)

type BucketObject struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

type ApplyPutArgs struct {
	Bucket        string
	Name          string
	Effective     time.Duration // 申请有效时间
	Header        http.Header   // header
	MaxObjectSize int64
}

type HeaderOption struct {
	ContentType string
	Filename    string
}

type ObjectInfo struct {
	Size int64
	Hash string
}

type SizeReader interface {
	io.ReadCloser
	Size() int64
}

func NewSizeReader(r io.ReadCloser, size int64) SizeReader {
	if r == nil {
		return nil
	}
	return &sizeReader{
		size:       size,
		ReadCloser: r,
	}
}

type sizeReader struct {
	size int64
	io.ReadCloser
}

func (r *sizeReader) Size() int64 {
	return r.size
}

type Interface interface {
	// Name 存储名字
	Name() string
	// MinFragmentSize 最小允许的分片大小
	MinFragmentSize() int64
	// MaxFragmentNum 最大允许的分片数量
	MaxFragmentNum() int
	// MinExpirationTime 最小过期时间
	MinExpirationTime() time.Duration
	// TempBucket 临时桶名，用于上传
	TempBucket() string
	// DataBucket 永久存储的桶名
	DataBucket() string
	// PresignedGetURL 通过桶名和对象名返回URL
	PresignedGetURL(ctx context.Context, bucket string, name string, expires time.Duration, opt *HeaderOption) (string, error)
	// PresignedPutURL 申请上传,返回PUT的上传地址
	PresignedPutURL(ctx context.Context, args *ApplyPutArgs) (string, error)
	// GetObjectInfo 获取对象信息
	GetObjectInfo(ctx context.Context, args *BucketObject) (*ObjectInfo, error)
	// CopyObject 复制对象
	CopyObject(ctx context.Context, src *BucketObject, dst *BucketObject) error
	// DeleteObject 删除对象(不存在返回nil)
	DeleteObject(ctx context.Context, info *BucketObject) error
	// ComposeObject 合并对象
	ComposeObject(ctx context.Context, src []BucketObject, dst *BucketObject) error
	// IsNotFound 判断是不是不存在导致的错误
	IsNotFound(err error) bool
	// CheckName 检查名字是否可用
	CheckName(name string) error
	// PutObject 上传文件
	PutObject(ctx context.Context, info *BucketObject, reader io.Reader, size int64) (*ObjectInfo, error)
	// GetObject 下载文件
	GetObject(ctx context.Context, info *BucketObject) (SizeReader, error)
}
