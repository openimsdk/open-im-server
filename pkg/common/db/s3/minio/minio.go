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

package minio

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/signer"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	unsignedPayload = "UNSIGNED-PAYLOAD"
)

const (
	minPartSize = 1024 * 1024 * 5        // 1MB
	maxPartSize = 1024 * 1024 * 1024 * 5 // 5GB
	maxNumSize  = 10000
)

func NewMinio() (s3.Interface, error) {
	conf := config.Config.Object.Minio
	u, err := url.Parse(conf.Endpoint)
	if err != nil {
		return nil, err
	}
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, conf.SessionToken),
		Secure: u.Scheme == "https",
	}
	client, err := minio.New(u.Host, opts)
	if err != nil {
		return nil, err
	}
	exists, err := client.BucketExists(context.Background(), conf.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), conf.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	return &Minio{
		bucket:    conf.Bucket,
		bucketURL: conf.Endpoint + "/" + conf.Bucket + "/",
		opts:      opts,
		core:      &minio.Core{Client: client},
	}, nil
}

type Minio struct {
	bucket    string
	bucketURL string
	opts      *minio.Options
	core      *minio.Core
}

func (m *Minio) Engine() string {
	return "minio"
}

func (m *Minio) PartLimit() *s3.PartLimit {
	return &s3.PartLimit{
		MinPartSize: minPartSize,
		MaxPartSize: maxPartSize,
		MaxNumSize:  maxNumSize,
	}
}

func (m *Minio) InitiateMultipartUpload(ctx context.Context, name string) (*s3.InitiateMultipartUploadResult, error) {
	uploadID, err := m.core.NewMultipartUpload(ctx, m.bucket, name, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &s3.InitiateMultipartUploadResult{
		Bucket:   m.bucket,
		Key:      name,
		UploadID: uploadID,
	}, nil
}

func (m *Minio) CompleteMultipartUpload(ctx context.Context, uploadID string, name string, parts []s3.Part) (*s3.CompleteMultipartUploadResult, error) {
	minioParts := make([]minio.CompletePart, len(parts))
	for i, part := range parts {
		minioParts[i] = minio.CompletePart{
			PartNumber: part.PartNumber,
			ETag:       strings.ToLower(part.ETag),
		}
	}
	upload, err := m.core.CompleteMultipartUpload(ctx, m.bucket, name, uploadID, minioParts, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &s3.CompleteMultipartUploadResult{
		Location: upload.Location,
		Bucket:   upload.Bucket,
		Key:      upload.Key,
		ETag:     strings.ToLower(upload.ETag),
	}, nil
}

func (m *Minio) PartSize(ctx context.Context, size int64) (int64, error) {
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

func (m *Minio) AuthSign(ctx context.Context, uploadID string, name string, expire time.Duration, partNumbers []int) (*s3.AuthSignResult, error) {
	creds, err := m.opts.Creds.Get()
	if err != nil {
		return nil, err
	}
	result := s3.AuthSignResult{
		URL:   m.bucketURL + name,
		Query: url.Values{"uploadId": {uploadID}},
		Parts: make([]s3.SignPart, len(partNumbers)),
	}
	for i, partNumber := range partNumbers {
		rawURL := result.URL + "?partNumber=" + strconv.Itoa(partNumber) + "&uploadId=" + uploadID
		request, err := http.NewRequestWithContext(ctx, http.MethodPut, rawURL, nil)
		if err != nil {
			return nil, err
		}
		request.Header.Set("X-Amz-Content-Sha256", unsignedPayload)
		request = signer.SignV4Trailer(*request, creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken, "us-east-1", nil)
		result.Parts[i] = s3.SignPart{
			PartNumber: partNumber,
			URL:        request.URL.String(),
			Query:      url.Values{"partNumber": {strconv.Itoa(partNumber)}},
			Header:     request.Header,
		}
	}
	return &result, nil
}

func (m *Minio) PresignedPutObject(ctx context.Context, name string, expire time.Duration) (string, error) {
	rawURL, err := m.core.Client.PresignedPutObject(ctx, m.bucket, name, expire)
	if err != nil {
		return "", err
	}
	return rawURL.String(), nil
}

func (m *Minio) DeleteObject(ctx context.Context, name string) error {
	return m.core.Client.RemoveObject(ctx, m.bucket, name, minio.RemoveObjectOptions{})
}

func (m *Minio) StatObject(ctx context.Context, name string) (*s3.ObjectInfo, error) {
	info, err := m.core.Client.StatObject(ctx, m.bucket, name, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &s3.ObjectInfo{
		ETag:         strings.ToLower(info.ETag),
		Key:          info.Key,
		Size:         info.Size,
		LastModified: info.LastModified,
	}, nil
}

func (m *Minio) CopyObject(ctx context.Context, src string, dst string) (*s3.CopyObjectInfo, error) {
	result, err := m.core.Client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: m.bucket,
		Object: dst,
	}, minio.CopySrcOptions{
		Bucket: m.bucket,
		Object: src,
	})
	if err != nil {
		return nil, err
	}
	return &s3.CopyObjectInfo{
		Key:  dst,
		ETag: strings.ToLower(result.ETag),
	}, nil
}

func (m *Minio) IsNotFound(err error) bool {
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

func (m *Minio) AbortMultipartUpload(ctx context.Context, uploadID string, name string) error {
	return m.core.AbortMultipartUpload(ctx, m.bucket, name, uploadID)
}

func (m *Minio) ListUploadedParts(ctx context.Context, uploadID string, name string, partNumberMarker int, maxParts int) (*s3.ListUploadedPartsResult, error) {
	result, err := m.core.ListObjectParts(ctx, m.bucket, name, uploadID, partNumberMarker, maxParts)
	if err != nil {
		return nil, err
	}
	res := &s3.ListUploadedPartsResult{
		Key:                  result.Key,
		UploadID:             result.UploadID,
		MaxParts:             result.MaxParts,
		NextPartNumberMarker: result.NextPartNumberMarker,
		UploadedParts:        make([]s3.UploadedPart, len(result.ObjectParts)),
	}
	for i, part := range result.ObjectParts {
		res.UploadedParts[i] = s3.UploadedPart{
			PartNumber:   part.PartNumber,
			LastModified: part.LastModified,
			ETag:         part.ETag,
			Size:         part.Size,
		}
	}
	return res, nil
}

func (m *Minio) AccessURL(ctx context.Context, name string, expire time.Duration, opt *s3.AccessURLOption) (string, error) {
	//reqParams := make(url.Values)
	//if opt != nil {
	//	if opt.ContentType != "" {
	//		reqParams.Set("Content-Type", opt.ContentType)
	//	}
	//	if opt.ContentDisposition != "" {
	//		reqParams.Set("Content-Disposition", opt.ContentDisposition)
	//	}
	//}
	if expire <= 0 {
		expire = time.Hour * 24 * 365 * 99 // 99 years
	} else if expire < time.Second {
		expire = time.Second
	}
	u, err := m.core.Client.PresignedGetObject(ctx, m.bucket, name, expire, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
