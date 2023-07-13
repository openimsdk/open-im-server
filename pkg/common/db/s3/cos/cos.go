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

package cos

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	minPartSize = 1024 * 1024 * 1        // 1MB
	maxPartSize = 1024 * 1024 * 1024 * 5 // 5GB
	maxNumSize  = 1000
)

func NewCos() (s3.Interface, error) {
	conf := config.Config.Object.Cos
	u, err := url.Parse(conf.BucketURL)
	if err != nil {
		panic(err)
	}
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     conf.SecretID,
			SecretKey:    conf.SecretKey,
			SessionToken: conf.SessionToken,
		},
	})
	return &Cos{
		copyURL:    u.Host + "/",
		client:     client,
		credential: client.GetCredential(),
	}, nil
}

type Cos struct {
	copyURL    string
	client     *cos.Client
	credential *cos.Credential
}

func (c *Cos) Engine() string {
	return "tencent-cos"
}

func (c *Cos) PartLimit() *s3.PartLimit {
	return &s3.PartLimit{
		MinPartSize: minPartSize,
		MaxPartSize: maxPartSize,
		MaxNumSize:  maxNumSize,
	}
}

func (c *Cos) InitiateMultipartUpload(ctx context.Context, name string) (*s3.InitiateMultipartUploadResult, error) {
	result, _, err := c.client.Object.InitiateMultipartUpload(ctx, name, nil)
	if err != nil {
		return nil, err
	}
	return &s3.InitiateMultipartUploadResult{
		UploadID: result.UploadID,
		Bucket:   result.Bucket,
		Key:      result.Key,
	}, nil
}

func (c *Cos) CompleteMultipartUpload(ctx context.Context, uploadID string, name string, parts []s3.Part) (*s3.CompleteMultipartUploadResult, error) {
	opts := &cos.CompleteMultipartUploadOptions{
		Parts: make([]cos.Object, len(parts)),
	}
	for i, part := range parts {
		opts.Parts[i] = cos.Object{
			PartNumber: part.PartNumber,
			ETag:       strings.ReplaceAll(part.ETag, `"`, ``),
		}
	}
	result, _, err := c.client.Object.CompleteMultipartUpload(ctx, name, uploadID, opts)
	if err != nil {
		return nil, err
	}
	return &s3.CompleteMultipartUploadResult{
		Location: result.Location,
		Bucket:   result.Bucket,
		Key:      result.Key,
		ETag:     result.ETag,
	}, nil
}

func (c *Cos) PartSize(ctx context.Context, size int64) (int64, error) {
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

func (c *Cos) AuthSign(ctx context.Context, uploadID string, name string, expire time.Duration, partNumbers []int) (*s3.AuthSignResult, error) {
	result := s3.AuthSignResult{
		URL:    c.client.BaseURL.BucketURL.String() + "/" + cos.EncodeURIComponent(name),
		Query:  url.Values{"uploadId": {uploadID}},
		Header: make(http.Header),
		Parts:  make([]s3.SignPart, len(partNumbers)),
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, result.URL, nil)
	if err != nil {
		return nil, err
	}
	cos.AddAuthorizationHeader(c.credential.SecretID, c.credential.SecretKey, c.credential.SessionToken, req, cos.NewAuthTime(expire))
	result.Header = req.Header
	for i, partNumber := range partNumbers {
		result.Parts[i] = s3.SignPart{
			PartNumber: partNumber,
			Query:      url.Values{"partNumber": {strconv.Itoa(partNumber)}},
		}
	}
	return &result, nil
}

func (c *Cos) PresignedPutObject(ctx context.Context, name string, expire time.Duration) (string, error) {
	rawURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodPut, name, c.credential.SecretID, c.credential.SecretKey, expire, nil)
	if err != nil {
		return "", err
	}
	return rawURL.String(), nil
}

func (c *Cos) DeleteObject(ctx context.Context, name string) error {
	_, err := c.client.Object.Delete(ctx, name)
	return err
}

func (c *Cos) StatObject(ctx context.Context, name string) (*s3.ObjectInfo, error) {
	if name != "" && name[0] == '/' {
		name = name[1:]
	}
	info, err := c.client.Object.Head(ctx, name, nil)
	if err != nil {
		return nil, err
	}
	res := &s3.ObjectInfo{Key: name}
	if res.ETag = strings.ToLower(strings.ReplaceAll(info.Header.Get("ETag"), `"`, "")); res.ETag == "" {
		return nil, errors.New("StatObject etag not found")
	}
	if contentLengthStr := info.Header.Get("Content-Length"); contentLengthStr == "" {
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
	if lastModified := info.Header.Get("Last-Modified"); lastModified == "" {
		return nil, errors.New("StatObject last-modified not found")
	} else {
		res.LastModified, err = time.Parse(http.TimeFormat, lastModified)
		if err != nil {
			return nil, fmt.Errorf("StatObject last-modified parse error: %w", err)
		}
	}
	return res, nil
}

func (c *Cos) CopyObject(ctx context.Context, src string, dst string) (*s3.CopyObjectInfo, error) {
	sourceURL := c.copyURL + src
	result, _, err := c.client.Object.Copy(ctx, dst, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	return &s3.CopyObjectInfo{
		Key:  dst,
		ETag: strings.ReplaceAll(result.ETag, `"`, ``),
	}, nil
}

func (c *Cos) IsNotFound(err error) bool {
	switch e := err.(type) {
	case *cos.ErrorResponse:
		return e.Response.StatusCode == http.StatusNotFound || e.Code == "NoSuchKey"
	default:
		return false
	}
}

func (c *Cos) AbortMultipartUpload(ctx context.Context, uploadID string, name string) error {
	_, err := c.client.Object.AbortMultipartUpload(ctx, name, uploadID)
	return err
}

func (c *Cos) ListUploadedParts(ctx context.Context, uploadID string, name string, partNumberMarker int, maxParts int) (*s3.ListUploadedPartsResult, error) {
	result, _, err := c.client.Object.ListParts(ctx, name, uploadID, &cos.ObjectListPartsOptions{
		MaxParts:         strconv.Itoa(maxParts),
		PartNumberMarker: strconv.Itoa(partNumberMarker),
	})
	if err != nil {
		return nil, err
	}
	res := &s3.ListUploadedPartsResult{
		Key:           result.Key,
		UploadID:      result.UploadID,
		UploadedParts: make([]s3.UploadedPart, len(result.Parts)),
	}
	res.MaxParts, _ = strconv.Atoi(result.MaxParts)
	res.NextPartNumberMarker, _ = strconv.Atoi(result.NextPartNumberMarker)
	for i, part := range result.Parts {
		lastModified, _ := time.Parse(http.TimeFormat, part.LastModified)
		res.UploadedParts[i] = s3.UploadedPart{
			PartNumber:   part.PartNumber,
			LastModified: lastModified,
			ETag:         part.ETag,
			Size:         part.Size,
		}
	}
	return res, nil
}

func (c *Cos) AccessURL(ctx context.Context, name string, expire time.Duration, opt *s3.AccessURLOption) (string, error) {
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
	rawURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodGet, name, c.credential.SecretID, c.credential.SecretKey, expire, nil)
	if err != nil {
		return "", err
	}
	return rawURL.String(), nil
}
