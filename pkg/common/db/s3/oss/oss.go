// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oss

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	minPartSize = 1024 * 1024 * 1        // 1MB
	maxPartSize = 1024 * 1024 * 1024 * 5 // 5GB
	maxNumSize  = 10000
)

func NewOSS() (s3.Interface, error) {
	conf := config.Config.Object.Oss
	if conf.BucketURL == "" {
		return nil, errors.New("bucket url is empty")
	}
	client, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(conf.Bucket)
	if err != nil {
		return nil, err
	}
	if conf.BucketURL[len(conf.BucketURL)-1] != '/' {
		conf.BucketURL += "/"
	}
	return &OSS{
		bucketURL:   conf.BucketURL,
		bucket:      bucket,
		credentials: client.Config.GetCredentials(),
	}, nil
}

type OSS struct {
	bucketURL   string
	bucket      *oss.Bucket
	credentials oss.Credentials
}

func (o *OSS) Engine() string {
	return "ali-oss"
}

func (o *OSS) PartLimit() *s3.PartLimit {
	return &s3.PartLimit{
		MinPartSize: minPartSize,
		MaxPartSize: maxPartSize,
		MaxNumSize:  maxNumSize,
	}
}

func (o *OSS) InitiateMultipartUpload(ctx context.Context, name string) (*s3.InitiateMultipartUploadResult, error) {
	result, err := o.bucket.InitiateMultipartUpload(name)
	if err != nil {
		return nil, err
	}
	return &s3.InitiateMultipartUploadResult{
		UploadID: result.UploadID,
		Bucket:   result.Bucket,
		Key:      result.Key,
	}, nil
}

func (o *OSS) CompleteMultipartUpload(ctx context.Context, uploadID string, name string, parts []s3.Part) (*s3.CompleteMultipartUploadResult, error) {
	ossParts := make([]oss.UploadPart, len(parts))
	for i, part := range parts {
		ossParts[i] = oss.UploadPart{
			PartNumber: part.PartNumber,
			ETag:       strings.ToUpper(part.ETag),
		}
	}
	result, err := o.bucket.CompleteMultipartUpload(oss.InitiateMultipartUploadResult{
		UploadID: uploadID,
		Bucket:   o.bucket.BucketName,
		Key:      name,
	}, ossParts)
	if err != nil {
		return nil, err
	}
	return &s3.CompleteMultipartUploadResult{
		Location: result.Location,
		Bucket:   result.Bucket,
		Key:      result.Key,
		ETag:     strings.ToLower(strings.ReplaceAll(result.ETag, `"`, ``)),
	}, nil
}

func (o *OSS) PartSize(ctx context.Context, size int64) (int64, error) {
	if size <= 0 {
		return 0, errors.New("size must be greater than 0")
	}
	if size > maxPartSize*maxNumSize {
		return 0, fmt.Errorf("size must be less than %db", maxPartSize*maxNumSize)
	}
	if size <= minPartSize*maxNumSize {
		return minPartSize, nil
	}
	partSize := size / maxNumSize
	if size%maxNumSize != 0 {
		partSize++
	}
	return partSize, nil
}

func (o *OSS) AuthSign(ctx context.Context, uploadID string, name string, expire time.Duration, partNumbers []int) (*s3.AuthSignResult, error) {
	result := s3.AuthSignResult{
		URL:    o.bucketURL + name,
		Query:  url.Values{"uploadId": {uploadID}},
		Header: make(http.Header),
		Parts:  make([]s3.SignPart, len(partNumbers)),
	}
	for i, partNumber := range partNumbers {
		rawURL := fmt.Sprintf(`%s%s?partNumber=%d&uploadId=%s`, o.bucketURL, name, partNumber, uploadID)
		request, err := http.NewRequestWithContext(ctx, http.MethodPut, rawURL, nil)
		if err != nil {
			return nil, err
		}
		if o.credentials.GetSecurityToken() != "" {
			request.Header.Set(oss.HTTPHeaderOssSecurityToken, o.credentials.GetSecurityToken())
		}
		request.Header.Set(oss.HTTPHeaderHost, request.Host)
		request.Header.Set(oss.HTTPHeaderDate, time.Now().UTC().Format(http.TimeFormat))
		authorization := fmt.Sprintf(`OSS %s:%s`, o.credentials.GetAccessKeyID(), o.getSignedStr(request, fmt.Sprintf(`/%s/%s?partNumber=%d&uploadId=%s`, o.bucket.BucketName, name, partNumber, uploadID), o.credentials.GetAccessKeySecret()))
		request.Header.Set(oss.HTTPHeaderAuthorization, authorization)
		result.Parts[i] = s3.SignPart{
			PartNumber: partNumber,
			Query:      url.Values{"partNumber": {strconv.Itoa(partNumber)}},
			URL:        request.URL.String(),
			Header:     request.Header,
		}
	}
	return &result, nil
}

func (o *OSS) PresignedPutObject(ctx context.Context, name string, expire time.Duration) (string, error) {
	return o.bucket.SignURL(name, http.MethodPut, int64(expire/time.Second))
}

func (o *OSS) StatObject(ctx context.Context, name string) (*s3.ObjectInfo, error) {
	header, err := o.bucket.GetObjectMeta(name)
	if err != nil {
		return nil, err
	}
	res := &s3.ObjectInfo{Key: name}
	if res.ETag = strings.ToLower(strings.ReplaceAll(header.Get("ETag"), `"`, ``)); res.ETag == "" {
		return nil, errors.New("StatObject etag not found")
	}
	if contentLengthStr := header.Get("Content-Length"); contentLengthStr == "" {
		return nil, errors.New("StatObject content-length not found")
	} else {
		res.Size, err = strconv.ParseInt(contentLengthStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("StatObject content-length parse error: %w", err)
		}
		if res.Size < 0 {
			return nil, errors.New("StatObject content-length must be greater than 0")
		}
	}
	if lastModified := header.Get("Last-Modified"); lastModified == "" {
		return nil, errors.New("StatObject last-modified not found")
	} else {
		res.LastModified, err = time.Parse(http.TimeFormat, lastModified)
		if err != nil {
			return nil, fmt.Errorf("StatObject last-modified parse error: %w", err)
		}
	}
	return res, nil
}

func (o *OSS) DeleteObject(ctx context.Context, name string) error {
	return o.bucket.DeleteObject(name)
}

func (o *OSS) CopyObject(ctx context.Context, src string, dst string) (*s3.CopyObjectInfo, error) {
	result, err := o.bucket.CopyObject(src, dst)
	if err != nil {
		return nil, err
	}
	return &s3.CopyObjectInfo{
		Key:  dst,
		ETag: strings.ToLower(strings.ReplaceAll(result.ETag, `"`, ``)),
	}, nil
}

func (o *OSS) IsNotFound(err error) bool {
	switch e := err.(type) {
	case oss.ServiceError:
		return e.StatusCode == http.StatusNotFound || e.Code == "NoSuchKey"
	case *oss.ServiceError:
		return e.StatusCode == http.StatusNotFound || e.Code == "NoSuchKey"
	default:
		return false
	}
}

func (o *OSS) AbortMultipartUpload(ctx context.Context, uploadID string, name string) error {
	return o.bucket.AbortMultipartUpload(oss.InitiateMultipartUploadResult{
		UploadID: uploadID,
		Key:      name,
		Bucket:   o.bucket.BucketName,
	})
}

func (o *OSS) ListUploadedParts(ctx context.Context, uploadID string, name string, partNumberMarker int, maxParts int) (*s3.ListUploadedPartsResult, error) {
	result, err := o.bucket.ListUploadedParts(oss.InitiateMultipartUploadResult{
		UploadID: uploadID,
		Key:      name,
		Bucket:   o.bucket.BucketName,
	}, oss.MaxUploads(100), oss.MaxParts(maxParts), oss.PartNumberMarker(partNumberMarker))
	if err != nil {
		return nil, err
	}
	res := &s3.ListUploadedPartsResult{
		Key:           result.Key,
		UploadID:      result.UploadID,
		MaxParts:      result.MaxParts,
		UploadedParts: make([]s3.UploadedPart, len(result.UploadedParts)),
	}
	res.NextPartNumberMarker, _ = strconv.Atoi(result.NextPartNumberMarker)
	for i, part := range result.UploadedParts {
		res.UploadedParts[i] = s3.UploadedPart{
			PartNumber:   part.PartNumber,
			LastModified: part.LastModified,
			ETag:         part.ETag,
			Size:         int64(part.Size),
		}
	}
	return res, nil
}

func (o *OSS) AccessURL(ctx context.Context, name string, expire time.Duration, opt *s3.AccessURLOption) (string, error) {
	//var opts []oss.Option
	//if opt != nil {
	//	if opt.ContentType != "" {
	//		opts = append(opts, oss.ContentType(opt.ContentType))
	//	}
	//	if opt.ContentDisposition != "" {
	//		opts = append(opts, oss.ContentDisposition(opt.ContentDisposition))
	//	}
	//}
	if expire <= 0 {
		expire = time.Hour * 24 * 365 * 99 // 99 years
	} else if expire < time.Second {
		expire = time.Second
	}
	return o.bucket.SignURL(name, http.MethodGet, int64(expire/time.Second))
}
