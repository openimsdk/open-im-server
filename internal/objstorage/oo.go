package objstorage

import "context"

type Interface interface {
	Init() error
	Name() string
	MinMultipartSize() int64
	UploadBucket() string
	PermanentBucket() string
	ClearBucket() string
	ApplyPut(ctx context.Context, args *ApplyPutArgs) (*PutRes, error)
	GetObjectInfo(ctx context.Context, args *BucketFile) (*ObjectInfo, error)
	CopyObjectInfo(ctx context.Context, src *BucketFile, dst *BucketFile) error
	DeleteObjectInfo(ctx context.Context, info *BucketFile) error
	MoveObjectInfo(ctx context.Context, src *BucketFile, dst *BucketFile) error
	MergeObjectInfo(ctx context.Context, src []BucketFile, dst *BucketFile) error
	IsNotFound(err error) bool
}
