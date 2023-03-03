package controller

import "C"
import (
	"OpenIM/pkg/common/db/obj"
	"OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/proto/third"
	"OpenIM/pkg/utils"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"path"
	"strconv"
	"time"
)

type S3Database interface {
	ApplyPut(ctx context.Context, req *third.ApplyPutReq) (*third.ApplyPutResp, error)
	GetPut(ctx context.Context, req *third.GetPutReq) (*third.GetPutResp, error)
	ConfirmPut(ctx context.Context, req *third.ConfirmPutReq) (*third.ConfirmPutResp, error)
}

func NewS3Database(obj obj.Interface, hash relation.ObjectHashModelInterface, info relation.ObjectInfoModelInterface, put relation.ObjectPutModelInterface) S3Database {
	return &s3Database{
		obj:  obj,
		hash: hash,
		info: info,
		put:  put,
	}
}

type s3Database struct {
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
	return "fragment_" + strconv.Itoa(index+1)
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
		return errors.New("hash value error")
	}
	return nil
}

func (c *s3Database) urlName(name string) string {
	if name[0] != '/' {
		name = "/" + name
	}
	return "http://127.0.0.1:8080" + name
}

func (c *s3Database) UUID() string {
	return uuid.New().String()
}

func (c *s3Database) HashName(hash string) string {
	return path.Join("hash", hash)
}

func (c *s3Database) isNotFound(err error) bool {
	return false
}

func (c *s3Database) ApplyPut(ctx context.Context, req *third.ApplyPutReq) (*third.ApplyPutResp, error) {
	if err := c.CheckHash(req.Hash); err != nil {
		return nil, err
	}
	if err := c.obj.CheckName(req.Name); err != nil {
		return nil, err
	}
	if req.CleanTime != 0 && req.CleanTime <= time.Now().UnixMilli() {
		return nil, errors.New("invalid CleanTime")
	}
	var expirationTime *time.Time
	if req.CleanTime != 0 {
		expirationTime = utils.ToPtr(time.UnixMilli(req.CleanTime))
	}
	if hash, err := c.hash.Take(ctx, req.Hash, c.obj.Name()); err == nil {
		o := relation.ObjectInfoModel{
			Name:           req.Name,
			Hash:           hash.Hash,
			ExpirationTime: expirationTime,
			CreateTime:     time.Now(),
		}
		if err := c.info.SetObject(ctx, &o); err != nil {
			return nil, err
		}
		return &third.ApplyPutResp{Url: c.urlName(o.Name)}, nil // 服务器已存在
	} else if !c.isNotFound(err) {
		return nil, err
	}
	// 新上传
	var pack int
	const effective = time.Hour * 24 * 2
	req.FragmentSize, pack = c.getFragmentNum(req.FragmentSize, req.Size)
	put := relation.ObjectPutModel{
		PutID:          c.UUID(),
		Hash:           req.Hash,
		Name:           req.Name,
		ObjectSize:     req.Size,
		FragmentSize:   req.FragmentSize,
		ExpirationTime: expirationTime,
		EffectiveTime:  time.Now().Add(effective),
	}
	put.Path = path.Join("upload", c.today(), req.Hash, put.PutID)
	putURLs := make([]string, 0, pack)
	for i := 0; i < pack; i++ {
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
	put.CreateTime = time.Now()
	if err := c.put.Create(ctx, []*relation.ObjectPutModel{&put}); err != nil {
		return nil, err
	}
	return &third.ApplyPutResp{
		PutID:        put.PutID,
		FragmentSize: put.FragmentSize,
		PutURLs:      putURLs,
	}, nil
}

func (c *s3Database) GetPut(ctx context.Context, req *third.GetPutReq) (*third.GetPutResp, error) {
	up, err := c.put.Take(ctx, req.PutID)
	if err != nil {
		return nil, err
	}
	if up.Complete {
		return nil, errors.New("up completed")
	}
	_, pack := c.getFragmentNum(up.FragmentSize, up.ObjectSize)
	fragments := make([]*third.GetPutFragment, pack)
	for i := 0; i < pack; i++ {
		name := path.Join(up.Path, c.fragmentName(i))
		o, err := c.obj.GetObjectInfo(ctx, &obj.BucketObject{
			Bucket: c.obj.TempBucket(),
			Name:   name,
		})
		if err != nil {
			if c.obj.IsNotFound(err) {
				fragments[i] = &third.GetPutFragment{}
				continue
			}
			return nil, err
		}
		fragments[i] = &third.GetPutFragment{Size: o.Size, Hash: o.Hash}
	}
	var cleanTime int64
	if up.ExpirationTime != nil {
		cleanTime = up.ExpirationTime.UnixMilli()
	}
	return &third.GetPutResp{
		FragmentSize: up.FragmentSize,
		Size:         up.ObjectSize,
		Name:         up.Name,
		Hash:         up.Hash,
		Fragments:    fragments,
		CleanTime:    cleanTime,
	}, nil
}

func (c *s3Database) ConfirmPut(ctx context.Context, req *third.ConfirmPutReq) (_ *third.ConfirmPutResp, _err error) {
	up, err := c.put.Take(ctx, req.PutID)
	if err != nil {
		return nil, err
	}
	_, pack := c.getFragmentNum(up.FragmentSize, up.ObjectSize)
	defer func() {
		if _err == nil {
			// 清理上传的碎片
			for i := 0; i < pack; i++ {
				name := path.Join(up.Path, c.fragmentName(i))
				err := c.obj.DeleteObjet(ctx, &obj.BucketObject{
					Bucket: c.obj.TempBucket(),
					Name:   name,
				})
				if err != nil {
					log.Printf("delete fragment %d %s %s failed %s\n", i, c.obj.TempBucket(), name, err)
				}
			}
		}
	}()
	if up.Complete {
		return nil, errors.New("put completed")
	}
	now := time.Now().UnixMilli()
	if up.EffectiveTime.UnixMilli() < now {
		return nil, errors.New("upload expired")
	}
	if up.ExpirationTime != nil && up.ExpirationTime.UnixMilli() < now {
		return nil, errors.New("object expired")
	}
	if hash, err := c.hash.Take(ctx, up.Hash, c.obj.Name()); err == nil {
		o := relation.ObjectInfoModel{
			Name:           up.Name,
			Hash:           hash.Hash,
			ExpirationTime: up.ExpirationTime,
			CreateTime:     time.Now(),
		}
		if err := c.info.SetObject(ctx, &o); err != nil {
			return nil, err
		}
		// 服务端已存在
		return &third.ConfirmPutResp{
			Url: c.urlName(o.Name),
		}, nil
	} else if c.isNotFound(err) {
		return nil, err
	}
	src := make([]obj.BucketObject, pack)
	for i := 0; i < pack; i++ {
		name := path.Join(up.Path, c.fragmentName(i))
		o, err := c.obj.GetObjectInfo(ctx, &obj.BucketObject{
			Bucket: c.obj.TempBucket(),
			Name:   name,
		})
		if err != nil {
			return nil, err
		}
		if i+1 == pack { // 最后一个
			size := up.ObjectSize - up.FragmentSize*int64(i)
			if size != o.Size {
				return nil, fmt.Errorf("last fragment %d size %d not equal to %d hash %s", i, o.Size, size, o.Hash)
			}
		} else {
			if o.Size != up.FragmentSize {
				return nil, fmt.Errorf("fragment %d size %d not equal to %d hash %s", i, o.Size, up.FragmentSize, o.Hash)
			}
		}
		src[i] = obj.BucketObject{
			Bucket: c.obj.TempBucket(),
			Name:   name,
		}
	}
	dst := &obj.BucketObject{
		Bucket: c.obj.DataBucket(),
		Name:   c.HashName(up.Hash),
	}
	if len(src) == 1 { // 未分片直接触发copy
		// 检查数据完整性,避免脏数据
		o, err := c.obj.GetObjectInfo(ctx, &src[0])
		if err != nil {
			return nil, err
		}
		if up.ObjectSize != o.Size {
			return nil, fmt.Errorf("size mismatching should %d reality %d", up.ObjectSize, o.Size)
		}
		if up.Hash != o.Hash {
			return nil, fmt.Errorf("hash mismatching should %s reality %s", up.Hash, o.Hash)
		}
		if err := c.obj.CopyObjet(ctx, &src[0], dst); err != nil {
			return nil, err
		}
	} else {
		tempBucket := &obj.BucketObject{
			Bucket: c.obj.TempBucket(),
			Name:   path.Join("merge", c.today(), req.PutID, c.UUID()),
		}
		defer func() { // 清理合成的文件
			if err := c.obj.DeleteObjet(ctx, tempBucket); err != nil {
				log.Printf("delete %s %s %s failed %s\n", c.obj.Name(), tempBucket.Bucket, tempBucket.Name, err)
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
		if up.ObjectSize != info.Size {
			return nil, fmt.Errorf("size mismatch should %d reality %d", up.ObjectSize, info.Size)
		}
		if up.Hash != info.Hash {
			return nil, fmt.Errorf("hash mismatch should %s reality %s", up.Hash, info.Hash)
		}
		if err := c.obj.CopyObjet(ctx, tempBucket, dst); err != nil {
			return nil, err
		}
	}
	o := &relation.ObjectInfoModel{
		Name:           up.Name,
		Hash:           up.Hash,
		ExpirationTime: up.ExpirationTime,
		CreateTime:     time.Now(),
	}
	h := &relation.ObjectHashModel{
		Hash:       up.Hash,
		Size:       up.ObjectSize,
		Engine:     c.obj.Name(),
		Bucket:     c.obj.DataBucket(),
		Name:       c.HashName(up.Hash),
		CreateTime: time.Now(),
	}
	if err := c.hash.Create(ctx, []*relation.ObjectHashModel{h}); err != nil {
		return nil, err
	}
	if err := c.info.SetObject(ctx, o); err != nil {
		return nil, err
	}
	if err := c.put.SetCompleted(ctx, up.PutID); err != nil {
		log.Printf("set uploaded %s failed %s\n", up.PutID, err)
	}
	return &third.ConfirmPutResp{
		Url: c.urlName(o.Name),
	}, nil
}
