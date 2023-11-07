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
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"

	"github.com/OpenIMSDK/tools/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/signer"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
)

const (
	unsignedPayload = "UNSIGNED-PAYLOAD"
)

const (
	minPartSize = 1024 * 1024 * 5        // 1MB
	maxPartSize = 1024 * 1024 * 1024 * 5 // 5GB
	maxNumSize  = 10000
)

const (
	maxImageWidth      = 1024
	maxImageHeight     = 1024
	maxImageSize       = 1024 * 1024 * 50
	imageThumbnailPath = "openim/thumbnail"
)

func NewMinio(cache cache.MinioCache) (s3.Interface, error) {
	u, err := url.Parse(config.Config.Object.Minio.Endpoint)
	if err != nil {
		return nil, err
	}
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(config.Config.Object.Minio.AccessKeyID, config.Config.Object.Minio.SecretAccessKey, config.Config.Object.Minio.SessionToken),
		Secure: u.Scheme == "https",
	}
	client, err := minio.New(u.Host, opts)
	if err != nil {
		return nil, err
	}
	m := &Minio{
		bucket: config.Config.Object.Minio.Bucket,
		core:   &minio.Core{Client: client},
		lock:   &sync.Mutex{},
		init:   false,
		cache:  cache,
	}
	if config.Config.Object.Minio.SignEndpoint == "" || config.Config.Object.Minio.SignEndpoint == config.Config.Object.Minio.Endpoint {
		m.opts = opts
		m.sign = m.core.Client
		m.prefix = u.Path
		u.Path = ""
		config.Config.Object.Minio.Endpoint = u.String()
		m.signEndpoint = config.Config.Object.Minio.Endpoint
	} else {
		su, err := url.Parse(config.Config.Object.Minio.SignEndpoint)
		if err != nil {
			return nil, err
		}
		m.opts = &minio.Options{
			Creds:  credentials.NewStaticV4(config.Config.Object.Minio.AccessKeyID, config.Config.Object.Minio.SecretAccessKey, config.Config.Object.Minio.SessionToken),
			Secure: su.Scheme == "https",
		}
		m.sign, err = minio.New(su.Host, m.opts)
		if err != nil {
			return nil, err
		}
		m.prefix = su.Path
		su.Path = ""
		config.Config.Object.Minio.SignEndpoint = su.String()
		m.signEndpoint = config.Config.Object.Minio.SignEndpoint
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := m.initMinio(ctx); err != nil {
		fmt.Println("init minio error:", err)
	}
	return m, nil
}

type Minio struct {
	bucket       string
	signEndpoint string
	location     string
	opts         *minio.Options
	core         *minio.Core
	sign         *minio.Client
	lock         sync.Locker
	init         bool
	prefix       string
	cache        cache.MinioCache
}

func (m *Minio) initMinio(ctx context.Context) error {
	if m.init {
		return nil
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.init {
		return nil
	}
	conf := config.Config.Object.Minio
	exists, err := m.core.Client.BucketExists(ctx, conf.Bucket)
	if err != nil {
		return fmt.Errorf("check bucket exists error: %w", err)
	}
	if !exists {
		if err := m.core.Client.MakeBucket(ctx, conf.Bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("make bucket error: %w", err)
		}
	}
	if conf.PublicRead {
		policy := fmt.Sprintf(
			`{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject","s3:PutObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::%s/*"],"Sid": ""}]}`,
			conf.Bucket,
		)
		if err := m.core.Client.SetBucketPolicy(ctx, conf.Bucket, policy); err != nil {
			return err
		}
	}
	m.location, err = m.core.Client.GetBucketLocation(ctx, conf.Bucket)
	if err != nil {
		return err
	}
	func() {
		if conf.SignEndpoint == "" || conf.SignEndpoint == conf.Endpoint {
			return
		}
		defer func() {
			if r := recover(); r != nil {
				m.sign = m.core.Client
				log.ZWarn(
					context.Background(),
					"set sign bucket location cache panic",
					errors.New("failed to get private field value"),
					"recover",
					fmt.Sprintf("%+v", r),
					"development version",
					"github.com/minio/minio-go/v7 v7.0.61",
				)
			}
		}()
		blc := reflect.ValueOf(m.sign).Elem().FieldByName("bucketLocCache")
		vblc := reflect.New(reflect.PtrTo(blc.Type()))
		*(*unsafe.Pointer)(vblc.UnsafePointer()) = unsafe.Pointer(blc.UnsafeAddr())
		vblc.Elem().Elem().Interface().(interface{ Set(string, string) }).Set(conf.Bucket, m.location)
	}()
	m.init = true
	return nil
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
	if err := m.initMinio(ctx); err != nil {
		return nil, err
	}
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
	if err := m.initMinio(ctx); err != nil {
		return nil, err
	}
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
	m.delObjectImageInfoKey(ctx, name, upload.Size)
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
	if err := m.initMinio(ctx); err != nil {
		return nil, err
	}
	creds, err := m.opts.Creds.Get()
	if err != nil {
		return nil, err
	}
	result := s3.AuthSignResult{
		URL:   m.signEndpoint + "/" + m.bucket + "/" + name,
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
		request = signer.SignV4Trailer(*request, creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken, m.location, nil)
		result.Parts[i] = s3.SignPart{
			PartNumber: partNumber,
			Query:      url.Values{"partNumber": {strconv.Itoa(partNumber)}},
			Header:     request.Header,
		}
	}
	if m.prefix != "" {
		result.URL = m.signEndpoint + m.prefix + "/" + m.bucket + "/" + name
	}
	return &result, nil
}

func (m *Minio) PresignedPutObject(ctx context.Context, name string, expire time.Duration) (string, error) {
	if err := m.initMinio(ctx); err != nil {
		return "", err
	}
	rawURL, err := m.sign.PresignedPutObject(ctx, m.bucket, name, expire)
	if err != nil {
		return "", err
	}
	if m.prefix != "" {
		rawURL.Path = path.Join(m.prefix, rawURL.Path)
	}
	return rawURL.String(), nil
}

func (m *Minio) DeleteObject(ctx context.Context, name string) error {
	if err := m.initMinio(ctx); err != nil {
		return err
	}
	return m.core.Client.RemoveObject(ctx, m.bucket, name, minio.RemoveObjectOptions{})
}

func (m *Minio) StatObject(ctx context.Context, name string) (*s3.ObjectInfo, error) {
	if err := m.initMinio(ctx); err != nil {
		return nil, err
	}
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
	if err := m.initMinio(ctx); err != nil {
		return nil, err
	}
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
	if err := m.initMinio(ctx); err != nil {
		return err
	}
	return m.core.AbortMultipartUpload(ctx, m.bucket, name, uploadID)
}

func (m *Minio) ListUploadedParts(ctx context.Context, uploadID string, name string, partNumberMarker int, maxParts int) (*s3.ListUploadedPartsResult, error) {
	if err := m.initMinio(ctx); err != nil {
		return nil, err
	}
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

func (m *Minio) PresignedGetObject(ctx context.Context, name string, expire time.Duration, query url.Values) (string, error) {
	if expire <= 0 {
		expire = time.Hour * 24 * 365 * 99 // 99 years
	} else if expire < time.Second {
		expire = time.Second
	}
	var (
		rawURL *url.URL
		err    error
	)
	if config.Config.Object.Minio.PublicRead {
		rawURL, err = makeTargetURL(m.sign, m.bucket, name, m.location, false, query)
	} else {
		rawURL, err = m.sign.PresignedGetObject(ctx, m.bucket, name, expire, query)
	}
	if err != nil {
		return "", err
	}
	if m.prefix != "" {
		rawURL.Path = path.Join(m.prefix, rawURL.Path)
	}
	return rawURL.String(), nil
}

func (m *Minio) AccessURL(ctx context.Context, name string, expire time.Duration, opt *s3.AccessURLOption) (string, error) {
	if err := m.initMinio(ctx); err != nil {
		return "", err
	}
	reqParams := make(url.Values)
	if opt != nil {
		if opt.ContentType != "" {
			reqParams.Set("response-content-type", opt.ContentType)
		}
		if opt.Filename != "" {
			reqParams.Set("response-content-disposition", `attachment; filename=`+strconv.Quote(opt.Filename))
		}
	}
	if opt.Image == nil || (opt.Image.Width < 0 && opt.Image.Height < 0 && opt.Image.Format == "") || (opt.Image.Width > maxImageWidth || opt.Image.Height > maxImageHeight) {
		return m.PresignedGetObject(ctx, name, expire, reqParams)
	}
	return m.getImageThumbnailURL(ctx, name, expire, opt.Image)
}

func (m *Minio) getObjectData(ctx context.Context, name string, limit int64) ([]byte, error) {
	object, err := m.core.Client.GetObject(ctx, m.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	if limit < 0 {
		return io.ReadAll(object)
	}
	return io.ReadAll(io.LimitReader(object, limit))
}
