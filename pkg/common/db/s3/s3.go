package s3

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type PartLimit struct {
	MinPartSize int64 `json:"minPartSize"`
	MaxPartSize int64 `json:"maxPartSize"`
	MaxNumSize  int   `json:"maxNumSize"`
}

type InitiateMultipartUploadResult struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	UploadID string `json:"uploadID"`
}

type MultipartUploadRequest struct {
	UploadID  string      `json:"uploadId"`
	Bucket    string      `json:"bucket"`
	Key       string      `json:"key"`
	Method    string      `json:"method"`
	URL       string      `json:"url"`
	Query     url.Values  `json:"query"`
	Header    http.Header `json:"header"`
	PartKey   string      `json:"partKey"`
	PartSize  int64       `json:"partSize"`
	FirstPart int         `json:"firstPart"`
}

type Part struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"etag"`
}

type CompleteMultipartUploadResult struct {
	Location string `json:"location"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	ETag     string `json:"etag"`
}

type SignResult struct {
	Parts []SignPart `json:"parts"`
}

type ObjectInfo struct {
	ETag         string    `json:"etag"`
	Key          string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

type CopyObjectInfo struct {
	Key  string `json:"name"`
	ETag string `json:"etag"`
}

type SignPart struct {
	PartNumber int         `json:"partNumber"`
	URL        string      `json:"url"`
	Query      url.Values  `json:"query"`
	Header     http.Header `json:"header"`
}

type AuthSignResult struct {
	URL    string      `json:"url"`
	Query  url.Values  `json:"query"`
	Header http.Header `json:"header"`
	Parts  []SignPart  `json:"parts"`
}

type InitiateUpload struct {
	UploadID  string      `json:"uploadId"`
	Bucket    string      `json:"bucket"`
	Key       string      `json:"key"`
	Method    string      `json:"method"`
	URL       string      `json:"url"`
	Query     url.Values  `json:"query"`
	Header    http.Header `json:"header"`
	PartKey   string      `json:"partKey"`
	PartSize  int64       `json:"partSize"`
	FirstPart int         `json:"firstPart"`
}

type UploadedPart struct {
	PartNumber   int       `json:"partNumber"`
	LastModified time.Time `json:"lastModified"`
	ETag         string    `json:"etag"`
	Size         int64     `json:"size"`
}

type ListUploadedPartsResult struct {
	Key                  string         `xml:"Key"`
	UploadID             string         `xml:"UploadId"`
	NextPartNumberMarker int            `xml:"NextPartNumberMarker"`
	MaxParts             int            `xml:"MaxParts"`
	UploadedParts        []UploadedPart `xml:"Part"`
}

type AccessURLOption struct {
	ContentType        string `json:"contentType"`
	ContentDisposition string `json:"contentDisposition"`
}

type Interface interface {
	Engine() string
	PartLimit() *PartLimit

	InitiateMultipartUpload(ctx context.Context, name string) (*InitiateMultipartUploadResult, error)
	CompleteMultipartUpload(ctx context.Context, uploadID string, name string, parts []Part) (*CompleteMultipartUploadResult, error)

	PartSize(ctx context.Context, size int64) (int64, error)
	AuthSign(ctx context.Context, uploadID string, name string, expire time.Duration, partNumbers []int) (*AuthSignResult, error)

	PresignedPutObject(ctx context.Context, name string, expire time.Duration) (string, error)

	DeleteObject(ctx context.Context, name string) error

	CopyObject(ctx context.Context, src string, dst string) (*CopyObjectInfo, error)

	StatObject(ctx context.Context, name string) (*ObjectInfo, error)

	IsNotFound(err error) bool

	AbortMultipartUpload(ctx context.Context, uploadID string, name string) error
	ListUploadedParts(ctx context.Context, uploadID string, name string, partNumberMarker int, maxParts int) (*ListUploadedPartsResult, error)

	AccessURL(ctx context.Context, name string, expire time.Duration, opt *AccessURLOption) (string, error)
}
