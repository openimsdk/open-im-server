package objstorage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"path"
	"strconv"
	"time"
)

func NewController(i Interface, kv KV) (*Controller, error) {
	if err := i.Init(); err != nil {
		return nil, err
	}
	return &Controller{
		i:  i,
		kv: kv,
	}, nil
}

type Controller struct {
	i Interface
	//i  *minioImpl
	kv KV
}

func (c *Controller) key(v string) string {
	return "OBJECT_STORAGE:" + c.i.Name() + ":" + v
}

func (c *Controller) putKey(v string) string {
	return c.key("put:" + v)
}

func (c *Controller) pathKey(v string) string {
	return c.key("path:" + v)
}

func (c *Controller) ApplyPut(ctx context.Context, args *FragmentPutArgs) (*PutAddr, error) {
	if data, err := c.kv.Get(ctx, c.pathKey(args.Hash)); err == nil {
		// 服务器已存在
		var src BucketFile
		if err := json.Unmarshal([]byte(data), &src); err != nil {
			return nil, err
		}
		var bucket string
		if args.ClearTime <= 0 {
			bucket = c.i.PermanentBucket()
		} else {
			bucket = c.i.ClearBucket()
		}
		dst := &BucketFile{
			Bucket: bucket,
			Name:   args.Name,
		}
		// 直接拷贝一份
		err := c.i.CopyObjectInfo(ctx, &src, dst)
		if err == nil {
			info, err := c.i.GetObjectInfo(ctx, dst)
			if err != nil {
				return nil, err
			}
			return &PutAddr{
				ResourceURL: info.URL,
			}, nil
		} else if !c.i.IsNotFound(err) {
			return nil, err
		}
	} else if !c.kv.IsNotFound(err) {
		return nil, err
	}
	// 上传逻辑
	name := args.Name
	effective := time.Now().Add(args.EffectiveTime)
	prefix := c.Prefix(&args.PutArgs)
	if minSize := c.i.MinMultipartSize(); args.FragmentSize > 0 && args.FragmentSize < minSize {
		args.FragmentSize = minSize
	}
	var pack int64
	if args.FragmentSize <= 0 || args.Size <= args.FragmentSize {
		pack = 1
	} else {
		pack = args.Size / args.FragmentSize
		if args.Size%args.FragmentSize > 0 {
			pack++
		}
	}
	p := path.Join(path.Dir(args.Name), time.Now().Format("20060102"))
	info := putInfo{
		Bucket:       c.i.UploadBucket(),
		Fragments:    make([]string, 0, pack),
		FragmentSize: args.FragmentSize,
		Name:         name,
		Hash:         args.Hash,
		Size:         args.Size,
	}
	if args.ClearTime > 0 {
		t := time.Now().Add(args.ClearTime).UnixMilli()
		info.ClearTime = &t
	}
	putURLs := make([]string, 0, pack)
	for i := int64(1); i <= pack; i++ {
		name := prefix + "_" + strconv.FormatInt(i, 10) + path.Ext(args.Name)
		name = path.Join(p, name)
		info.Fragments = append(info.Fragments, name)
		args.Name = name
		put, err := c.i.ApplyPut(ctx, &ApplyPutArgs{
			Bucket:    info.Bucket,
			Name:      name,
			Effective: args.EffectiveTime,
			Header:    args.Header,
		})
		if err != nil {
			return nil, err
		}
		putURLs = append(putURLs, put.URL)
	}
	data, err := json.Marshal(&info)
	if err != nil {
		return nil, err
	}
	if err := c.kv.Set(ctx, c.putKey(prefix), string(data), args.EffectiveTime); err != nil {
		return nil, err
	}
	var fragmentSize int64
	if pack == 1 {
		fragmentSize = args.Size
	} else {
		fragmentSize = args.FragmentSize
	}
	return &PutAddr{
		PutURLs:       putURLs,
		FragmentSize:  fragmentSize,
		PutID:         prefix,
		EffectiveTime: effective,
	}, nil
}

func (c *Controller) ConfirmPut(ctx context.Context, putID string) (*ObjectInfo, error) {
	data, err := c.kv.Get(ctx, c.putKey(putID))
	if err != nil {
		return nil, err
	}
	var info putInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}
	var total int64
	src := make([]BucketFile, len(info.Fragments))
	for i, fragment := range info.Fragments {
		state, err := c.i.GetObjectInfo(ctx, &BucketFile{
			Bucket: info.Bucket,
			Name:   fragment,
		})
		if err != nil {
			return nil, err
		}
		total += state.Size
		src[i] = BucketFile{
			Bucket: info.Bucket,
			Name:   fragment,
		}
	}
	if total != info.Size {
		return nil, fmt.Errorf("incomplete upload %d/%d", total, info.Size)
	}
	var dst *BucketFile
	if info.ClearTime == nil {
		dst = &BucketFile{
			Bucket: c.i.PermanentBucket(),
			Name:   info.Name,
		}
	} else {
		dst = &BucketFile{
			Bucket: c.i.ClearBucket(),
			Name:   info.Name,
		}
	}
	if err := c.i.MergeObjectInfo(ctx, src, dst); err != nil { // SourceInfo 0 is too small (2) and it is not the last part
		return nil, err
	}
	obj, err := c.i.GetObjectInfo(ctx, dst)
	if err != nil {
		return nil, err
	}
	go func() {
		err := c.kv.Del(ctx, c.putKey(putID))
		if err != nil {
			log.Println("del key:", err)
		}
		for _, b := range src {
			err = c.i.DeleteObjectInfo(ctx, &b)
			if err != nil {
				log.Println("del obj:", err)
			}
		}
	}()
	return obj, nil
}

func (c *Controller) Prefix(args *PutArgs) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(args.Name)
	buf.WriteString("~~~@~@~~~")
	buf.WriteString(strconv.FormatInt(args.Size, 10))
	buf.WriteString(",")
	buf.WriteString(args.Hash)
	buf.WriteString(",")
	buf.WriteString(strconv.FormatInt(int64(args.ClearTime), 10))
	buf.WriteString(",")
	buf.WriteString(strconv.FormatInt(int64(args.EffectiveTime), 10))
	buf.WriteString(",")
	buf.WriteString(c.i.Name())
	r := make([]byte, 16)
	rand.Read(r)
	buf.Write(r)
	md5v := md5.Sum(buf.Bytes())
	return hex.EncodeToString(md5v[:])
}

type putInfo struct {
	Bucket       string
	Fragments    []string
	FragmentSize int64
	Size         int64
	Name         string
	Hash         string
	ClearTime    *int64
}
