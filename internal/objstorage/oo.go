package objstorage

import "context"

type Interface interface {
	Init() error
	Name() string
	UploadBucket() string
	PermanentBucket() string
	ClearBucket() string
	ApplyPut(ctx context.Context, args *ApplyPutArgs) (*PutRes, error)
	GetObjectInfo(ctx context.Context, args *BucketFile) (*ObjectInfo, error)
	CopyObjetInfo(ctx context.Context, src *BucketFile, dst *BucketFile) error
	DeleteObjetInfo(ctx context.Context, info *BucketFile) error
	MoveObjetInfo(ctx context.Context, src *BucketFile, dst *BucketFile) error
	MergeObjectInfo(ctx context.Context, src []BucketFile, dst *BucketFile) error
	IsNotFound(err error) bool
}
