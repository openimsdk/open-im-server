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

// docURL: https://docs.aws.amazon.com/AmazonS3/latest/API/Welcome.html

package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
)

const (
	minPartSize int64 = 1024 * 1024 * 1        // 1MB
	maxPartSize int64 = 1024 * 1024 * 1024 * 5 // 5GB
	maxNumSize  int64 = 10000
)

// const (
// 	imagePng  = "png"
// 	imageJpg  = "jpg"
// 	imageJpeg = "jpeg"
// 	imageGif  = "gif"
// 	imageWebp = "webp"
// )

// const successCode = http.StatusOK

// const (
// 	videoSnapshotImagePng = "png"
// 	videoSnapshotImageJpg = "jpg"
// )

func NewAWS() (s3.Interface, error) {
	conf := config.Config.Object.Aws
	credential := credentials.NewStaticCredentials(
		conf.AccessKeyID,     // accessKey
		conf.AccessKeySecret, // secretKey
		"")                   // stoken

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region), // The area where the bucket is located
		Credentials: credential,
	})

	if err != nil {
		return nil, err
	}
	return &Aws{
		bucket:     conf.Bucket,
		client:     sdk.New(sess),
		credential: credential,
	}, nil
}

type Aws struct {
	bucket     string
	client     *sdk.S3
	credential *credentials.Credentials
}

func (a *Aws) Engine() string {
	return "aws"
}

func (a *Aws) InitiateMultipartUpload(ctx context.Context, name string) (*s3.InitiateMultipartUploadResult, error) {
	input := &sdk.CreateMultipartUploadInput{
		Bucket: aws.String(a.bucket), // TODO: To be verified whether it is required
		Key:    aws.String(name),
	}
	result, err := a.client.CreateMultipartUploadWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return &s3.InitiateMultipartUploadResult{
		Bucket:   *result.Bucket,
		Key:      *result.Key,
		UploadID: *result.UploadId,
	}, nil
}

func (a *Aws) CompleteMultipartUpload(ctx context.Context, uploadID string, name string, parts []s3.Part) (*s3.CompleteMultipartUploadResult, error) {
	sdkParts := make([]*sdk.CompletedPart, len(parts))
	for i, part := range parts {
		sdkParts[i] = &sdk.CompletedPart{
			ETag:       aws.String(part.ETag),
			PartNumber: aws.Int64(int64(part.PartNumber)),
		}
	}
	input := &sdk.CompleteMultipartUploadInput{
		Bucket:   aws.String(a.bucket), // TODO: To be verified whether it is required
		Key:      aws.String(name),
		UploadId: aws.String(uploadID),
		MultipartUpload: &sdk.CompletedMultipartUpload{
			Parts: sdkParts,
		},
	}
	result, err := a.client.CompleteMultipartUploadWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return &s3.CompleteMultipartUploadResult{
		Location: *result.Location,
		Bucket:   *result.Bucket,
		Key:      *result.Key,
		ETag:     *result.ETag,
	}, nil
}

func (a *Aws) PartSize(ctx context.Context, size int64) (int64, error) {
	if size <= 0 {
		return 0, errors.New("size must be greater than 0")
	}
	if size > maxPartSize*maxNumSize {
		return 0, fmt.Errorf("AWS size must be less than the maximum allowed limit")
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

func (a *Aws) DeleteObject(ctx context.Context, name string) error {
	_, err := a.client.DeleteObjectWithContext(ctx, &sdk.DeleteObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(name),
	})
	return err
}

func (a *Aws) CopyObject(ctx context.Context, src string, dst string) (*s3.CopyObjectInfo, error) {
	result, err := a.client.CopyObjectWithContext(ctx, &sdk.CopyObjectInput{
		Bucket:     aws.String(a.bucket),
		Key:        aws.String(dst),
		CopySource: aws.String(src),
	})
	if err != nil {
		return nil, err
	}
	return &s3.CopyObjectInfo{
		ETag: *result.CopyObjectResult.ETag,
		Key:  dst,
	}, nil
}

func (a *Aws) IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case sdk.ErrCodeNoSuchKey:
			return true
		default:
			return false
		}
	}
	return false
}

func (a *Aws) AbortMultipartUpload(ctx context.Context, uploadID string, name string) error {
	_, err := a.client.AbortMultipartUploadWithContext(ctx, &sdk.AbortMultipartUploadInput{
		Bucket:   aws.String(a.bucket),
		Key:      aws.String(name),
		UploadId: aws.String(uploadID),
	})
	return err
}

func (a *Aws) ListUploadedParts(ctx context.Context, uploadID string, name string, partNumberMarker int, maxParts int) (*s3.ListUploadedPartsResult, error) {
	result, err := a.client.ListPartsWithContext(ctx, &sdk.ListPartsInput{
		Bucket:           aws.String(a.bucket),
		Key:              aws.String(name),
		UploadId:         aws.String(uploadID),
		MaxParts:         aws.Int64(int64(maxParts)),
		PartNumberMarker: aws.Int64(int64(partNumberMarker)),
	})
	if err != nil {
		return nil, err
	}
	parts := make([]s3.UploadedPart, len(result.Parts))
	for i, part := range result.Parts {
		parts[i] = s3.UploadedPart{
			PartNumber:   int(*part.PartNumber),
			LastModified: *part.LastModified,
			Size:         *part.Size,
			ETag:         *part.ETag,
		}
	}
	return &s3.ListUploadedPartsResult{
		Key:                  *result.Key,
		UploadID:             *result.UploadId,
		NextPartNumberMarker: int(*result.NextPartNumberMarker),
		MaxParts:             int(*result.MaxParts),
		UploadedParts:        parts,
	}, nil
}

func (a *Aws) PartLimit() *s3.PartLimit {
	return &s3.PartLimit{
		MinPartSize: minPartSize,
		MaxPartSize: maxPartSize,
		MaxNumSize:  maxNumSize,
	}
}

func (a *Aws) PresignedPutObject(ctx context.Context, name string, expire time.Duration) (string, error) {
	req, _ := a.client.PutObjectRequest(&sdk.PutObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(name),
	})
	url, err := req.Presign(expire)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (a *Aws) StatObject(ctx context.Context, name string) (*s3.ObjectInfo, error) {
	result, err := a.client.GetObjectWithContext(ctx, &sdk.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	res := &s3.ObjectInfo{
		Key:          name,
		ETag:         *result.ETag,
		Size:         *result.ContentLength,
		LastModified: *result.LastModified,
	}
	return res, nil
}

// AccessURL todo.
func (a *Aws) AccessURL(ctx context.Context, name string, expire time.Duration, opt *s3.AccessURLOption) (string, error) {
	// todo
	return "", nil
}

func (a *Aws) FormData(ctx context.Context, name string, size int64, contentType string, duration time.Duration) (*s3.FormData, error) {
	// todo
	return nil, nil
}

func (a *Aws) AuthSign(ctx context.Context, uploadID string, name string, expire time.Duration, partNumbers []int) (*s3.AuthSignResult, error) {
	// todo
	return nil, nil
}
