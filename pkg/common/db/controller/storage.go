package controller

import "C"
import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/obj"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/google/uuid"
	"io"
	"net/url"
	"path"
	"strconv"
	"time"
)

const (
	hashPrefix     = "hash"
	tempPrefix     = "temp"
	fragmentPrefix = "fragment_"
	urlsName       = "urls.json"
)

type S3Database interface {
	ApplyPut(ctx context.Context, req *third.ApplyPutReq) (*third.ApplyPutResp, error)
	GetPut(ctx context.Context, req *third.GetPutReq) (*third.GetPutResp, error)
	ConfirmPut(ctx context.Context, req *third.ConfirmPutReq) (*third.ConfirmPutResp, error)
	GetUrl(ctx context.Context, req *third.GetUrlReq) (*third.GetUrlResp, error)
	GetHashInfo(ctx context.Context, req *third.GetHashInfoReq) (*third.GetHashInfoResp, error)
	CleanExpirationObject(ctx context.Context, t time.Time)
}

func NewS3Database(obj obj.Interface, hash relation.ObjectHashModelInterface, info relation.ObjectInfoModelInterface, put relation.ObjectPutModelInterface, url *url.URL) S3Database {
	return &s3Database{
		url:  url,
		obj:  obj,
		hash: hash,
		info: info,
		put:  put,
	}
}

type s3Database struct {
	url  *url.URL
	obj  obj.Interface
	hash relation.ObjectHashModelInterface
	info relation.ObjectInfoModelInterface
	put  relation.ObjectPutModelInterface
}

// today 今天的日期
func (c *s3Database) today() string {
	return time.Now().Format("20060102")
}

// fragmentName 根据序号生成文件名
func (c *s3Database) fragmentName(index int) string {
	return fragmentPrefix + strconv.Itoa(index+1)
}

// getFragmentNum 获取分片大小和分片数量
func (c *s3Database) getFragmentNum(fragmentSize int64, objectSize int64) (int64, int) {
	if size := c.obj.MinFragmentSize(); fragmentSize < size {
		fragmentSize = size
	}
	if fragmentSize <= 0 || objectSize <= fragmentSize {
		return objectSize, 1
	} else {
		num := int(objectSize / fragmentSize)
		if objectSize%fragmentSize > 0 {
			num++
		}
		if n := c.obj.MaxFragmentNum(); num > n {
			num = n
		}
		return fragmentSize, num
	}
}

func (c *s3Database) CheckHash(hash string) error {
	val, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	if len(val) != md5.Size {
		return errs.ErrArgs.Wrap("invalid hash")
	}
	return nil
}

func (c *s3Database) urlName(name string) string {
	u := url.URL{
		Scheme:      c.url.Scheme,
		Opaque:      c.url.Opaque,
		User:        c.url.User,
		Host:        c.url.Host,
		Path:        c.url.Path,
		RawPath:     c.url.RawPath,
		OmitHost:    c.url.OmitHost,
		ForceQuery:  c.url.ForceQuery,
		RawQuery:    c.url.RawQuery,
		Fragment:    c.url.Fragment,
		RawFragment: c.url.RawFragment,
	}
	v := make(url.Values, 1)
	v.Set("name", name)
	u.RawQuery = v.Encode()
	return u.String()
}

func (c *s3Database) UUID() string {
	return uuid.New().String()
}

func (c *s3Database) HashName(hash string) string {
	return path.Join(hashPrefix, hash+"_"+c.today()+"_"+c.UUID())
}

func (c *s3Database) isNotFound(err error) bool {
	return relation.IsNotFound(err)
}

func (c *s3Database) ApplyPut(ctx context.Context, req *third.ApplyPutReq) (*third.ApplyPutResp, error) {
	if err := c.CheckHash(req.Hash); err != nil {
		return nil, err
	}
	if err := c.obj.CheckName(req.Name); err != nil {
		return nil, err
	}
	if req.ValidTime != 0 && req.ValidTime <= time.Now().UnixMilli() {
		return nil, errors.New("invalid ValidTime")
	}
	var expirationTime *time.Time
	if req.ValidTime != 0 {
		expirationTime = utils.ToPtr(time.UnixMilli(req.ValidTime))
	}
	if hash, err := c.hash.Take(ctx, req.Hash, c.obj.Name()); err == nil {
		o := relation.ObjectInfoModel{
			Name:        req.Name,
			Hash:        hash.Hash,
			ValidTime:   expirationTime,
			ContentType: req.ContentType,
			CreateTime:  time.Now(),
		}
		if err := c.info.SetObject(ctx, &o); err != nil {
			return nil, err
		}
		return &third.ApplyPutResp{Url: c.urlName(o.Name)}, nil // 服务器已存在
	} else if !c.isNotFound(err) {
		return nil, err
	}
	// 新上传
	var fragmentNum int
	const effective = time.Hour * 24 * 2
	req.FragmentSize, fragmentNum = c.getFragmentNum(req.FragmentSize, req.Size)
	put := relation.ObjectPutModel{
		PutID:         req.PutID,
		Hash:          req.Hash,
		Name:          req.Name,
		ObjectSize:    req.Size,
		ContentType:   req.ContentType,
		FragmentSize:  req.FragmentSize,
		ValidTime:     expirationTime,
		EffectiveTime: time.Now().Add(effective),
	}
	if put.PutID == "" {
		put.PutID = c.UUID()
	}
	if v, err := c.put.Take(ctx, put.PutID); err == nil {
		now := time.Now().UnixMilli()
		if v.EffectiveTime.UnixMilli() <= now {
			if err := c.put.DelPut(ctx, []string{v.PutID}); err != nil {
				return nil, err
			}
		} else {
			return nil, errs.ErrDuplicateKey.Wrap(fmt.Sprintf("duplicate put id %s", put.PutID))
		}
	} else if !c.isNotFound(err) {
		return nil, err
	}
	put.Path = path.Join(tempPrefix, c.today(), req.Hash, put.PutID)
	putURLs := make([]string, 0, fragmentNum)
	for i := 0; i < fragmentNum; i++ {
		url, err := c.obj.PresignedPutURL(ctx, &obj.ApplyPutArgs{
			Bucket:        c.obj.TempBucket(),
			Name:          path.Join(put.Path, c.fragmentName(i)),
			Effective:     effective,
			MaxObjectSize: req.FragmentSize,
		})
		if err != nil {
			return nil, err
		}
		putURLs = append(putURLs, url)
	}
	urlsJsonData, err := json.Marshal(putURLs)
	if err != nil {
		return nil, err
	}
	t := md5.Sum(urlsJsonData)
	put.PutURLsHash = hex.EncodeToString(t[:])
	_, err = c.obj.PutObject(ctx, &obj.BucketObject{Bucket: c.obj.TempBucket(), Name: path.Join(put.Path, urlsName)}, bytes.NewReader(urlsJsonData), int64(len(urlsJsonData)))
	if err != nil {
		return nil, err
	}
	put.CreateTime = time.Now()
	if err := c.put.Create(ctx, []*relation.ObjectPutModel{&put}); err != nil {
		return nil, err
	}
	return &third.ApplyPutResp{
		PutID:        put.PutID,
		FragmentSize: put.FragmentSize,
		PutURLs:      putURLs,
		ValidTime:    put.EffectiveTime.UnixMilli(),
	}, nil
}

func (c *s3Database) GetPut(ctx context.Context, req *third.GetPutReq) (*third.GetPutResp, error) {
	up, err := c.put.Take(ctx, req.PutID)
	if err != nil {
		return nil, err
	}
	reader, err := c.obj.GetObject(ctx, &obj.BucketObject{Bucket: c.obj.TempBucket(), Name: path.Join(up.Path, urlsName)})
	if err != nil {
		return nil, err
	}
	urlsData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	t := md5.Sum(urlsData)
	if h := hex.EncodeToString(t[:]); h != up.PutURLsHash {
		return nil, fmt.Errorf("invalid put urls hash %s %s", h, up.PutURLsHash)
	}
	var urls []string
	if err := json.Unmarshal(urlsData, &urls); err != nil {
		return nil, err
	}
	_, fragmentNum := c.getFragmentNum(up.FragmentSize, up.ObjectSize)
	if len(urls) != fragmentNum {
		return nil, fmt.Errorf("invalid urls length %d fragment %d", len(urls), fragmentNum)
	}
	fragments := make([]*third.GetPutFragment, fragmentNum)
	for i := 0; i < fragmentNum; i++ {
		name := path.Join(up.Path, c.fragmentName(i))
		o, err := c.obj.GetObjectInfo(ctx, &obj.BucketObject{
			Bucket: c.obj.TempBucket(),
			Name:   name,
		})
		if err != nil {
			if c.obj.IsNotFound(err) {
				fragments[i] = &third.GetPutFragment{Url: urls[i]}
				continue
			}
			return nil, err
		}
		fragments[i] = &third.GetPutFragment{Size: o.Size, Hash: o.Hash, Url: urls[i]}
	}
	var validTime int64
	if up.ValidTime != nil {
		validTime = up.ValidTime.UnixMilli()
	}
	return &third.GetPutResp{
		FragmentSize: up.FragmentSize,
		Size:         up.ObjectSize,
		Name:         up.Name,
		Hash:         up.Hash,
		Fragments:    fragments,
		PutURLsHash:  up.PutURLsHash,
		ContentType:  up.ContentType,
		ValidTime:    validTime,
	}, nil
}

func (c *s3Database) ConfirmPut(ctx context.Context, req *third.ConfirmPutReq) (_ *third.ConfirmPutResp, _err error) {
	put, err := c.put.Take(ctx, req.PutID)
	if err != nil {
		return nil, err
	}
	_, pack := c.getFragmentNum(put.FragmentSize, put.ObjectSize)
	defer func() {
		if _err == nil {
			// 清理上传的碎片
			err := c.obj.DeleteObject(ctx, &obj.BucketObject{Bucket: c.obj.TempBucket(), Name: put.Path})
			if err != nil {
				log.ZError(ctx, "deleteObject failed", err, "Bucket", c.obj.TempBucket(), "Path", put.Path)
			}
		}
	}()
	now := time.Now().UnixMilli()
	if put.EffectiveTime.UnixMilli() < now {
		return nil, errs.ErrFileUploadedExpired.Wrap("put expired")
	}
	if put.ValidTime != nil && put.ValidTime.UnixMilli() < now {
		return nil, errs.ErrFileUploadedExpired.Wrap("object expired")
	}
	if hash, err := c.hash.Take(ctx, put.Hash, c.obj.Name()); err == nil {
		o := relation.ObjectInfoModel{
			Name:        put.Name,
			Hash:        hash.Hash,
			ValidTime:   put.ValidTime,
			ContentType: put.ContentType,
			CreateTime:  time.Now(),
		}
		if err := c.info.SetObject(ctx, &o); err != nil {
			return nil, err
		}
		defer func() {
			err := c.obj.DeleteObject(ctx, &obj.BucketObject{
				Bucket: c.obj.TempBucket(),
				Name:   put.Path,
			})
			if err != nil {
				log.ZError(ctx, "DeleteObject", err, "Bucket", c.obj.TempBucket(), "Path", put.Path)
			}
		}()
		// 服务端已存在
		return &third.ConfirmPutResp{
			Url: c.urlName(o.Name),
		}, nil
	} else if !c.isNotFound(err) {
		return nil, err
	}
	src := make([]obj.BucketObject, pack)
	for i := 0; i < pack; i++ {
		name := path.Join(put.Path, c.fragmentName(i))
		o, err := c.obj.GetObjectInfo(ctx, &obj.BucketObject{
			Bucket: c.obj.TempBucket(),
			Name:   name,
		})
		if err != nil {
			return nil, err
		}
		if i+1 == pack { // 最后一个
			size := put.ObjectSize - put.FragmentSize*int64(i)
			if size != o.Size {
				return nil, fmt.Errorf("last fragment %d size %d not equal to %d hash %s", i, o.Size, size, o.Hash)
			}
		} else {
			if o.Size != put.FragmentSize {
				return nil, fmt.Errorf("fragment %d size %d not equal to %d hash %s", i, o.Size, put.FragmentSize, o.Hash)
			}
		}
		src[i] = obj.BucketObject{
			Bucket: c.obj.TempBucket(),
			Name:   name,
		}
	}
	dst := &obj.BucketObject{
		Bucket: c.obj.DataBucket(),
		Name:   c.HashName(put.Hash),
	}
	if len(src) == 1 { // 未分片直接触发copy
		// 检查数据完整性,避免脏数据
		o, err := c.obj.GetObjectInfo(ctx, &src[0])
		if err != nil {
			return nil, err
		}
		if put.ObjectSize != o.Size {
			return nil, fmt.Errorf("size mismatching should %d reality %d", put.ObjectSize, o.Size)
		}
		if put.Hash != o.Hash {
			return nil, fmt.Errorf("hash mismatching should %s reality %s", put.Hash, o.Hash)
		}
		if err := c.obj.CopyObject(ctx, &src[0], dst); err != nil {
			return nil, err
		}
	} else {
		tempBucket := &obj.BucketObject{
			Bucket: c.obj.TempBucket(),
			Name:   path.Join(put.Path, "merge_"+c.UUID()),
		}
		defer func() { // 清理合成的文件
			if err := c.obj.DeleteObject(ctx, tempBucket); err != nil {
				log.ZError(ctx, "DeleteObject", err, "Bucket", tempBucket.Bucket, "Path", tempBucket.Name)
			}
		}()
		err := c.obj.ComposeObject(ctx, src, tempBucket)
		if err != nil {
			return nil, err
		}
		info, err := c.obj.GetObjectInfo(ctx, tempBucket)
		if err != nil {
			return nil, err
		}
		if put.ObjectSize != info.Size {
			return nil, fmt.Errorf("size mismatch should %d reality %d", put.ObjectSize, info.Size)
		}
		if put.Hash != info.Hash {
			return nil, fmt.Errorf("hash mismatch should %s reality %s", put.Hash, info.Hash)
		}
		if err := c.obj.CopyObject(ctx, tempBucket, dst); err != nil {
			return nil, err
		}
	}
	h := &relation.ObjectHashModel{
		Hash:       put.Hash,
		Engine:     c.obj.Name(),
		Size:       put.ObjectSize,
		Bucket:     c.obj.DataBucket(),
		Name:       dst.Name,
		CreateTime: time.Now(),
	}
	if err := c.hash.Create(ctx, []*relation.ObjectHashModel{h}); err != nil {
		return nil, err
	}
	o := &relation.ObjectInfoModel{
		Name:        put.Name,
		Hash:        put.Hash,
		ContentType: put.ContentType,
		ValidTime:   put.ValidTime,
		CreateTime:  time.Now(),
	}
	if err := c.info.SetObject(ctx, o); err != nil {
		return nil, err
	}
	if err := c.put.DelPut(ctx, []string{put.PutID}); err != nil {
		log.ZError(ctx, "DelPut", err, "PutID", put.PutID)
	}
	return &third.ConfirmPutResp{
		Url: c.urlName(o.Name),
	}, nil
}

func (c *s3Database) GetUrl(ctx context.Context, req *third.GetUrlReq) (*third.GetUrlResp, error) {
	info, err := c.info.Take(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if info.ValidTime != nil && info.ValidTime.Before(time.Now()) {
		return nil, errs.ErrRecordNotFound.Wrap("object expired")
	}
	hash, err := c.hash.Take(ctx, info.Hash, c.obj.Name())
	if err != nil {
		return nil, err
	}
	opt := obj.HeaderOption{ContentType: info.ContentType}
	if req.Attachment {
		opt.Filename = info.Name
	}
	u, err := c.obj.PresignedGetURL(ctx, hash.Bucket, hash.Name, time.Duration(req.Expires)*time.Millisecond, &opt)
	if err != nil {
		return nil, err
	}
	return &third.GetUrlResp{
		Url:  u,
		Size: hash.Size,
		Hash: hash.Hash,
	}, nil
}

func (c *s3Database) CleanExpirationObject(ctx context.Context, t time.Time) {
	// 清理上传产生的临时文件
	c.cleanPutTemp(ctx, t, 10)
	// 清理hash引用全过期的文件
	c.cleanExpirationObject(ctx, t)
	// 清理没有引用的hash对象
	c.clearNoCitation(ctx, c.obj.Name(), 10)
}

func (c *s3Database) cleanPutTemp(ctx context.Context, t time.Time, num int) {
	for {
		puts, err := c.put.FindExpirationPut(ctx, t, num)
		if err != nil {
			log.ZError(ctx, "FindExpirationPut", err, "Time", t, "Num", num)
			return
		}
		if len(puts) == 0 {
			return
		}
		for _, put := range puts {
			err := c.obj.DeleteObject(ctx, &obj.BucketObject{Bucket: c.obj.TempBucket(), Name: put.Path})
			if err != nil {
				log.ZError(ctx, "DeleteObject", err, "Bucket", c.obj.TempBucket(), "Path", put.Path)
				return
			}
		}
		ids := utils.Slice(puts, func(e *relation.ObjectPutModel) string { return e.PutID })
		err = c.put.DelPut(ctx, ids)
		if err != nil {
			log.ZError(ctx, "DelPut", err, "PutID", ids)
			return
		}
	}
}

func (c *s3Database) cleanExpirationObject(ctx context.Context, t time.Time) {
	err := c.info.DeleteExpiration(ctx, t)
	if err != nil {
		log.ZError(ctx, "DeleteExpiration", err, "Time", t)
	}
}

func (c *s3Database) clearNoCitation(ctx context.Context, engine string, limit int) {
	for {
		list, err := c.hash.DeleteNoCitation(ctx, engine, limit)
		if err != nil {
			log.ZError(ctx, "DeleteNoCitation", err, "Engine", engine, "Limit", limit)
			return
		}
		if len(list) == 0 {
			return
		}
		var hasErr bool
		for _, h := range list {
			err := c.obj.DeleteObject(ctx, &obj.BucketObject{Bucket: h.Bucket, Name: h.Name})
			if err != nil {
				hasErr = true
				log.ZError(ctx, "DeleteObject", err, "Bucket", h.Bucket, "Path", h.Name)
				continue
			}
		}
		if hasErr {
			return
		}
	}
}

func (c *s3Database) GetHashInfo(ctx context.Context, req *third.GetHashInfoReq) (*third.GetHashInfoResp, error) {
	if err := c.CheckHash(req.Hash); err != nil {
		return nil, err
	}
	o, err := c.hash.Take(ctx, req.Hash, c.obj.Name())
	if err != nil {
		return nil, err
	}
	return &third.GetHashInfoResp{
		Hash: o.Hash,
		Size: o.Size,
	}, nil
}
