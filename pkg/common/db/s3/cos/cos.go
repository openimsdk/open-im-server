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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
	"github.com/tencentyun/cos-go-sdk-v5"
)

const (
	minPartSize int64 = 1024 * 1024 * 1        // 1MB
	maxPartSize int64 = 1024 * 1024 * 1024 * 5 // 5GB
	maxNumSize  int64 = 1000
)

const (
	imagePng  = "png"
	imageJpg  = "jpg"
	imageJpeg = "jpeg"
	imageGif  = "gif"
	imageWebp = "webp"
)

const successCode = http.StatusOK

type Config struct {
	BucketURL    string
	SecretID     string
	SecretKey    string
	SessionToken string
	PublicRead   bool
}

func NewCos(conf Config) (s3.Interface, error) {
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
		publicRead: conf.PublicRead,
		copyURL:    u.Host + "/",
		client:     client,
		credential: client.GetCredential(),
	}, nil
}

type Cos struct {
	publicRead bool
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
		return 0, fmt.Errorf("COS size must be less than the maximum allowed limit")
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
	switch e := errs.Unwrap(err).(type) {
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
	var imageMogr string
	var option cos.PresignedURLOptions
	if opt != nil {
		query := make(url.Values)
		if opt.Image != nil {
			// https://cloud.tencent.com/document/product/436/44880
			style := make([]string, 0, 2)
			wh := make([]string, 2)
			if opt.Image.Width > 0 {
				wh[0] = strconv.Itoa(opt.Image.Width)
			}
			if opt.Image.Height > 0 {
				wh[1] = strconv.Itoa(opt.Image.Height)
			}
			if opt.Image.Width > 0 || opt.Image.Height > 0 {
				style = append(style, strings.Join(wh, "x"))
			}
			switch opt.Image.Format {
			case
				imagePng,
				imageJpg,
				imageJpeg,
				imageGif,
				imageWebp:
				style = append(style, "format/"+opt.Image.Format)
			}
			if len(style) > 0 {
				imageMogr = "imageMogr2/thumbnail/" + strings.Join(style, "/") + "/ignore-error/1"
			}
		}
		if opt.ContentType != "" {
			query.Set("response-content-type", opt.ContentType)
		}
		if opt.Filename != "" {
			query.Set("response-content-disposition", `attachment; filename=`+strconv.Quote(opt.Filename))
		}
		if len(query) > 0 {
			option.Query = &query
		}
	}
	if expire <= 0 {
		expire = time.Hour * 24 * 365 * 99 // 99 years
	} else if expire < time.Second {
		expire = time.Second
	}
	rawURL, err := c.getPresignedURL(ctx, name, expire, &option)
	if err != nil {
		return "", err
	}
	if imageMogr != "" {
		if rawURL.RawQuery == "" {
			rawURL.RawQuery = imageMogr
		} else {
			rawURL.RawQuery = rawURL.RawQuery + "&" + imageMogr
		}
	}
	return rawURL.String(), nil
}

func (c *Cos) getPresignedURL(ctx context.Context, name string, expire time.Duration, opt *cos.PresignedURLOptions) (*url.URL, error) {
	if !c.publicRead {
		return c.client.Object.GetPresignedURL(ctx, http.MethodGet, name, c.credential.SecretID, c.credential.SecretKey, expire, opt)
	}
	return c.client.Object.GetObjectURL(name), nil
}

func (c *Cos) FormData(ctx context.Context, name string, size int64, contentType string, duration time.Duration) (*s3.FormData, error) {
	// https://cloud.tencent.com/document/product/436/14690
	now := time.Now()
	expiration := now.Add(duration)
	keyTime := fmt.Sprintf("%d;%d", now.Unix(), expiration.Unix())
	conditions := []any{
		map[string]string{"q-sign-algorithm": "sha1"},
		map[string]string{"q-ak": c.credential.SecretID},
		map[string]string{"q-sign-time": keyTime},
		map[string]string{"key": name},
	}
	if contentType != "" {
		conditions = append(conditions, map[string]string{"Content-Type": contentType})
	}
	policy := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": conditions,
	}
	policyJson, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	signKey := hmacSha1val(c.credential.SecretKey, keyTime)
	strToSign := sha1val(string(policyJson))
	signature := hmacSha1val(signKey, strToSign)

	fd := &s3.FormData{
		URL:     c.client.BaseURL.BucketURL.String(),
		File:    "file",
		Expires: expiration,
		FormData: map[string]string{
			"policy":                base64.StdEncoding.EncodeToString(policyJson),
			"q-sign-algorithm":      "sha1",
			"q-ak":                  c.credential.SecretID,
			"q-key-time":            keyTime,
			"q-signature":           signature,
			"key":                   name,
			"success_action_status": strconv.Itoa(successCode),
		},
		SuccessCodes: []int{successCode},
	}
	if contentType != "" {
		fd.FormData["Content-Type"] = contentType
	}
	if c.credential.SessionToken != "" {
		fd.FormData["x-cos-security-token"] = c.credential.SessionToken
	}
	return fd, nil
}

func hmacSha1val(key, msg string) string {
	v := hmac.New(sha1.New, []byte(key))
	v.Write([]byte(msg))
	return hex.EncodeToString(v.Sum(nil))
}

func sha1val(msg string) string {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(msg))
	return hex.EncodeToString(sha1Hash.Sum(nil))
}
