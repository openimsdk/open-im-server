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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
)

const (
	minPartSize = 1024 * 1024 * 1        // 1MB
	maxPartSize = 1024 * 1024 * 1024 * 5 // 5GB
	maxNumSize  = 10000
)

const (
	imagePng  = "png"
	imageJpg  = "jpg"
	imageJpeg = "jpeg"
	imageGif  = "gif"
	imageWebp = "webp"
)

const successCode = http.StatusOK

const (
	videoSnapshotImagePng = "png"
	videoSnapshotImageJpg = "jpg"
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
		um:          *(*urlMaker)(reflect.ValueOf(bucket.Client.Conn).Elem().FieldByName("url").UnsafePointer()),
	}, nil
}

type OSS struct {
	bucketURL   string
	bucket      *oss.Bucket
	credentials oss.Credentials
	um          urlMaker
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
		request, err := http.NewRequest(http.MethodPut, rawURL, nil)
		if err != nil {
			return nil, err
		}
		if o.credentials.GetSecurityToken() != "" {
			request.Header.Set(oss.HTTPHeaderOssSecurityToken, o.credentials.GetSecurityToken())
		}
		now := time.Now().UTC().Format(http.TimeFormat)
		request.Header.Set(oss.HTTPHeaderHost, request.Host)
		request.Header.Set(oss.HTTPHeaderDate, now)
		request.Header.Set(oss.HttpHeaderOssDate, now)
		signHeader(*o.bucket.Client.Conn, request, fmt.Sprintf(`/%s/%s?partNumber=%d&uploadId=%s`, o.bucket.BucketName, name, partNumber, uploadID))
		delete(request.Header, oss.HTTPHeaderDate)
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
	publicRead := config.Config.Object.Oss.PublicRead
	var opts []oss.Option
	if opt != nil {
		if opt.Image != nil {
			// 文档地址: https://help.aliyun.com/zh/oss/user-guide/resize-images-4?spm=a2c4g.11186623.0.0.4b3b1e4fWW6yji
			var format string
			switch opt.Image.Format {
			case
				imagePng,
				imageJpg,
				imageJpeg,
				imageGif,
				imageWebp:
				format = opt.Image.Format
			default:
				opt.Image.Format = imageJpg
			}
			// https://oss-console-img-demo-cn-hangzhou.oss-cn-hangzhou.aliyuncs.com/example.jpg?x-oss-process=image/resize,h_100,m_lfit
			process := "image/resize,m_lfit"
			if opt.Image.Width > 0 {
				process += ",w_" + strconv.Itoa(opt.Image.Width)
			}
			if opt.Image.Height > 0 {
				process += ",h_" + strconv.Itoa(opt.Image.Height)
			}
			process += ",format," + format
			opts = append(opts, oss.Process(process))
		}
		if !publicRead {
			if opt.ContentType != "" {
				opts = append(opts, oss.ResponseContentType(opt.ContentType))
			}
			if opt.Filename != "" {
				opts = append(opts, oss.ResponseContentDisposition(`attachment; filename=`+strconv.Quote(opt.Filename)))
			}
		}
	}
	if expire <= 0 {
		expire = time.Hour * 24 * 365 * 99 // 99 years
	} else if expire < time.Second {
		expire = time.Second
	}
	if !publicRead {
		return o.bucket.SignURL(name, http.MethodGet, int64(expire/time.Second), opts...)
	}
	rawParams, err := oss.GetRawParams(opts)
	if err != nil {
		return "", err
	}
	params := getURLParams(*o.bucket.Client.Conn, rawParams)
	return getURL(o.um, o.bucket.BucketName, name, params).String(), nil
}

func (o *OSS) FormData(ctx context.Context, name string, size int64, contentType string, duration time.Duration) (*s3.FormData, error) {
	// https://help.aliyun.com/zh/oss/developer-reference/postobject?spm=a2c4g.11186623.0.0.1cb83cebkP55nn
	expires := time.Now().Add(duration)
	conditions := []any{
		map[string]string{"bucket": o.bucket.BucketName},
		map[string]string{"key": name},
	}
	if size > 0 {
		conditions = append(conditions, []any{"content-length-range", 0, size})
	}
	policy := map[string]any{
		"expiration": expires.Format("2006-01-02T15:04:05.000Z"),
		"conditions": conditions,
	}
	policyJson, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	policyStr := base64.StdEncoding.EncodeToString(policyJson)
	h := hmac.New(sha1.New, []byte(o.credentials.GetAccessKeySecret()))
	if _, err := io.WriteString(h, policyStr); err != nil {
		return nil, err
	}
	fd := &s3.FormData{
		URL:     o.bucketURL,
		File:    "file",
		Expires: expires,
		FormData: map[string]string{
			"key":                   name,
			"policy":                policyStr,
			"OSSAccessKeyId":        o.credentials.GetAccessKeyID(),
			"success_action_status": strconv.Itoa(successCode),
			"signature":             base64.StdEncoding.EncodeToString(h.Sum(nil)),
		},
		SuccessCodes: []int{successCode},
	}
	if contentType != "" {
		fd.FormData["x-oss-content-type"] = contentType
	}
	return fd, nil
}
