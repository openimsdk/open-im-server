package obj

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7/pkg/s3utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var conf = config.Config.Object.Tencent

func create_url(bucket, region, source string) string {
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", bucket, region, source)
}

func NewCosClient() (Interface, error) {
	u, err := url.Parse(create_url(conf.Bucket, conf.Region, ""))
	if err != nil {
		return nil, fmt.Errorf("tencent cos url parse %w", err)
	}
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: conf.SecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	return &cosImpl{
		client:     c,
		tempBucket: conf.Bucket,
	}, err
}

type cosImpl struct {
	tempBucket string // 上传桶
	//dataBucket string // 永久桶
	urlstr string // 访问地址
	client *cos.Client
}

func (c *cosImpl) Name() string {
	return "tx_oss cos"
}

func (c *cosImpl) MinFragmentSize() int64 {
	return 1024 * 1024 * 5 // 每个分片最小大小 tx_oss.absMinPartSize
}

func (c *cosImpl) MaxFragmentNum() int {
	return 1000 // 最大分片数量 tx_oss.maxPartsCount
}

func (c *cosImpl) MinExpirationTime() time.Duration {
	return time.Hour * 24
}

func (c *cosImpl) TempBucket() string {
	return c.tempBucket
}

func (c *cosImpl) DataBucket() string {
	//return c.dataBucket
	return ""
}

func (c *cosImpl) PresignedGetURL(ctx context.Context, bucket string, name string, expires time.Duration, opt *HeaderOption) (string, error) {
	// 参考文档：https://cloud.tencent.com/document/product/436/14116
	// 获取对象访问 URL，用于匿名下载和分发
	presignedGetURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodGet, name, conf.SecretID, conf.SecretKey, time.Hour, nil)
	if err != nil {
		return "", err
	}
	return presignedGetURL.String(), nil
}

func (c *cosImpl) PresignedPutURL(ctx context.Context, args *ApplyPutArgs) (string, error) {
	// 参考文档：https://cloud.tencent.com/document/product/436/14114

	if args.Effective <= 0 {
		return "", errors.New("EffectiveTime <= 0")
	}
	_, err := c.GetObjectInfo(ctx, &BucketObject{
		Bucket: c.tempBucket,
		Name:   args.Name,
	})
	if err == nil {
		return "", fmt.Errorf("minio bucket %s name %s already exists", args.Bucket, args.Name)
	} else if !c.IsNotFound(err) {
		return "", err
	}
	// 获取预签名 URL
	presignedPutURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodPut, args.Name, conf.SecretID, conf.SecretKey, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("minio apply error: %w", err)
	}
	return presignedPutURL.String(), nil
}

func (c *cosImpl) GetObjectInfo(ctx context.Context, args *BucketObject) (*ObjectInfo, error) {
	// https://cloud.tencent.com/document/product/436/7745
	// 新增参数 id 代表指定版本，如果不指定，默认查询对象最新版本
	head, err := c.client.Object.Head(ctx, args.Name, nil, "")
	if err != nil {
		return nil, err
	}
	size, _ := strconv.ParseInt(head.Header.Get("Content-Length"), 10, 64)
	return &ObjectInfo{
		Size: size,
		Hash: head.Header.Get("ETag"),
	}, nil
}

func (c *cosImpl) CopyObject(ctx context.Context, src *BucketObject, dst *BucketObject) error {
	srcURL := create_url(src.Bucket, conf.Region, src.Name)
	_, _, err := c.client.Object.Copy(ctx, dst.Name, srcURL, nil)
	if err == nil {
		_, err = c.client.Object.Delete(ctx, srcURL, nil)
		if err != nil {
			// Error
		}
	}
	return err
}

func (c *cosImpl) DeleteObject(ctx context.Context, info *BucketObject) error {
	_, err := c.client.Object.Delete(ctx, info.Name)
	return err
}

func (c *cosImpl) ComposeObject(ctx context.Context, src []BucketObject, dst *BucketObject) error {
	//TODO implement me
	panic("implement me")
}

func (c *cosImpl) IsNotFound(err error) bool {
	ok, err := c.client.Object.IsExist(context.Background(), c.tempBucket)
	if err == nil && ok {
		fmt.Printf("object exists\n")
		return true
	} else if err != nil {
		fmt.Printf("head object failed: %v\n", err)
		return false
	} else {
		fmt.Printf("object does not exist\n")
		return false
	}
}

func (c *cosImpl) CheckName(name string) error {
	return s3utils.CheckValidObjectName(name)
}

func (c *cosImpl) PutObject(ctx context.Context, info *BucketObject, reader io.Reader, size int64) (*ObjectInfo, error) {
	/*// 采用高级接口, 无需用户指定 size
	update, _, err := c.client.Object.Upload(
		ctx, info.Name, info.Bucket, nil,
	)
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{
		Hash: update.ETag,
	}, nil*/
	// Case1 使用 Put 上传对象
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			// 如果不是必要操作，建议上传文件时不要给单个文件设置权限，避免达到限制。若不设置默认继承桶的权限。
			XCosACL: "private",
		},
	}
	resp, err := c.client.Object.Put(ctx, info.Name, reader, opt)
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{
		Hash: resp.Header.Get("ETag"),
	}, nil
}

func (c *cosImpl) GetObject(ctx context.Context, info *BucketObject) (SizeReader, error) {
	opt := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
	}
	update, err := c.client.Object.Download(
		ctx, info.Name, info.Bucket, opt,
	)
	if err != nil {
		return nil, err
	}
	size, _ := strconv.ParseInt(update.Header.Get("Content-Length"), 10, 64)
	body := update.Body
	if body == nil {
		return nil, errors.New("response body is nil")
	}
	readCloser, ok := body.(io.ReadCloser)
	if !ok {
		return nil, errors.New("failed to convert response to ReadCloser")
	}
	return NewSizeReader(readCloser, size), nil
}
