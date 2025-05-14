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
	"path/filepath"
	"time"

	redisCache "github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/tools/s3"
	"github.com/openimsdk/tools/s3/cont"
	"github.com/redis/go-redis/v9"
)

type S3Database interface {
	PartLimit() (*s3.PartLimit, error)
	PartSize(ctx context.Context, size int64) (int64, error)
	AuthSign(ctx context.Context, uploadID string, partNumbers []int) (*s3.AuthSignResult, error)
	InitiateMultipartUpload(ctx context.Context, hash string, size int64, expire time.Duration, maxParts int, contentType string) (*cont.InitiateUploadResult, error)
	CompleteMultipartUpload(ctx context.Context, uploadID string, parts []string) (*cont.UploadResult, error)
	AccessURL(ctx context.Context, name string, expire time.Duration, opt *s3.AccessURLOption) (time.Time, string, error)
	SetObject(ctx context.Context, info *model.Object) error
	StatObject(ctx context.Context, name string) (*s3.ObjectInfo, error)
	FormData(ctx context.Context, name string, size int64, contentType string, duration time.Duration) (*s3.FormData, error)
	FindExpirationObject(ctx context.Context, engine string, expiration time.Time, needDelType []string, count int64) ([]*model.Object, error)
	DeleteSpecifiedData(ctx context.Context, engine string, name []string) error
	DelS3Key(ctx context.Context, engine string, keys ...string) error
	GetKeyCount(ctx context.Context, engine string, key string) (int64, error)
}

func NewS3Database(rdb redis.UniversalClient, s3 s3.Interface, obj database.ObjectInfo) S3Database {
	return &s3Database{
		s3:      cont.New(redisCache.NewS3Cache(rdb, s3), s3),
		cache:   redisCache.NewObjectCacheRedis(rdb, obj),
		s3cache: redisCache.NewS3Cache(rdb, s3),
		db:      obj,
	}
}

type s3Database struct {
	s3      *cont.Controller
	cache   cache.ObjectCache
	s3cache cont.S3Cache
	db      database.ObjectInfo
}

func (s *s3Database) PartSize(ctx context.Context, size int64) (int64, error) {
	return s.s3.PartSize(ctx, size)
}

func (s *s3Database) PartLimit() (*s3.PartLimit, error) {
	return s.s3.PartLimit()
}

func (s *s3Database) AuthSign(ctx context.Context, uploadID string, partNumbers []int) (*s3.AuthSignResult, error) {
	return s.s3.AuthSign(ctx, uploadID, partNumbers)
}

func (s *s3Database) InitiateMultipartUpload(ctx context.Context, hash string, size int64, expire time.Duration, maxParts int, contentType string) (*cont.InitiateUploadResult, error) {
	return s.s3.InitiateUploadContentType(ctx, hash, size, expire, maxParts, contentType)
}

func (s *s3Database) CompleteMultipartUpload(ctx context.Context, uploadID string, parts []string) (*cont.UploadResult, error) {
	return s.s3.CompleteUpload(ctx, uploadID, parts)
}

func (s *s3Database) SetObject(ctx context.Context, info *model.Object) error {
	info.Engine = s.s3.Engine()
	if err := s.db.SetObject(ctx, info); err != nil {
		return err
	}
	return s.cache.DelObjectName(info.Engine, info.Name).ChainExecDel(ctx)
}

func (s *s3Database) AccessURL(ctx context.Context, name string, expire time.Duration, opt *s3.AccessURLOption) (time.Time, string, error) {
	obj, err := s.cache.GetName(ctx, s.s3.Engine(), name)
	if err != nil {
		return time.Time{}, "", err
	}
	if opt == nil {
		opt = &s3.AccessURLOption{}
	}
	if opt.ContentType == "" {
		opt.ContentType = obj.ContentType
	}
	if opt.Filename == "" {
		opt.Filename = filepath.Base(obj.Name)
	}
	expireTime := time.Now().Add(expire)
	rawURL, err := s.s3.AccessURL(ctx, obj.Key, expire, opt)
	if err != nil {
		return time.Time{}, "", err
	}
	return expireTime, rawURL, nil
}

func (s *s3Database) StatObject(ctx context.Context, name string) (*s3.ObjectInfo, error) {
	return s.s3.StatObject(ctx, name)
}

func (s *s3Database) FormData(ctx context.Context, name string, size int64, contentType string, duration time.Duration) (*s3.FormData, error) {
	return s.s3.FormData(ctx, name, size, contentType, duration)
}

func (s *s3Database) FindExpirationObject(ctx context.Context, engine string, expiration time.Time, needDelType []string, count int64) ([]*model.Object, error) {
	return s.db.FindExpirationObject(ctx, engine, expiration, needDelType, count)
}

func (s *s3Database) GetKeyCount(ctx context.Context, engine string, key string) (int64, error) {
	return s.db.GetKeyCount(ctx, engine, key)
}

func (s *s3Database) DeleteSpecifiedData(ctx context.Context, engine string, name []string) error {
	return s.db.Delete(ctx, engine, name)
}

func (s *s3Database) DelS3Key(ctx context.Context, engine string, keys ...string) error {
	return s.s3cache.DelS3Key(ctx, engine, keys...)
}
