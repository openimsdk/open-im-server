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

package controller

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3/cont"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"path/filepath"
	"time"
)

type S3Database interface {
	PartLimit() *s3.PartLimit
	PartSize(ctx context.Context, size int64) (int64, error)
	AuthSign(ctx context.Context, uploadID string, partNumbers []int) (*s3.AuthSignResult, error)
	InitiateMultipartUpload(ctx context.Context, hash string, size int64, expire time.Duration, maxParts int) (*cont.InitiateUploadResult, error)
	CompleteMultipartUpload(ctx context.Context, uploadID string, parts []string) (*cont.UploadResult, error)
	AccessURL(ctx context.Context, name string, expire time.Duration) (time.Time, string, error)
	SetObject(ctx context.Context, info *relation.ObjectModel) error
}

func NewS3Database(s3 s3.Interface, obj relation.ObjectInfoModelInterface) S3Database {
	return &s3Database{
		s3:  cont.New(s3),
		obj: obj,
	}
}

type s3Database struct {
	s3  *cont.Controller
	obj relation.ObjectInfoModelInterface
}

func (s *s3Database) PartSize(ctx context.Context, size int64) (int64, error) {
	return s.s3.PartSize(ctx, size)
}

func (s *s3Database) PartLimit() *s3.PartLimit {
	return s.s3.PartLimit()
}

func (s *s3Database) AuthSign(ctx context.Context, uploadID string, partNumbers []int) (*s3.AuthSignResult, error) {
	return s.s3.AuthSign(ctx, uploadID, partNumbers)
}

func (s *s3Database) InitiateMultipartUpload(ctx context.Context, hash string, size int64, expire time.Duration, maxParts int) (*cont.InitiateUploadResult, error) {
	return s.s3.InitiateUpload(ctx, hash, size, expire, maxParts)
}

func (s *s3Database) CompleteMultipartUpload(ctx context.Context, uploadID string, parts []string) (*cont.UploadResult, error) {
	return s.s3.CompleteUpload(ctx, uploadID, parts)
}

func (s *s3Database) SetObject(ctx context.Context, info *relation.ObjectModel) error {
	return s.obj.SetObject(ctx, info)
}

func (s *s3Database) AccessURL(ctx context.Context, name string, expire time.Duration) (time.Time, string, error) {
	obj, err := s.obj.Take(ctx, name)
	if err != nil {
		return time.Time{}, "", err
	}
	opt := &s3.AccessURLOption{
		ContentType: obj.ContentType,
	}
	if filename := filepath.Base(obj.Name); filename != "" {
		opt.ContentDisposition = `attachment; filename=` + filename
	}
	expireTime := time.Now().Add(expire)
	rawURL, err := s.s3.AccessURL(ctx, obj.Key, expire, opt)
	if err != nil {
		return time.Time{}, "", err
	}
	return expireTime, rawURL, nil
}
