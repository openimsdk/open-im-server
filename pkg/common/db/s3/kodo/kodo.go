package kodo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	awss3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

const (
	minPartSize = 1024 * 1024 * 1        // 1MB
	maxPartSize = 1024 * 1024 * 1024 * 5 // 5GB
	maxNumSize  = 10000
)

type Kodo struct {
	AccessKey     string
	SecretKey     string
	Region        string
	Token         string
	Endpoint      string
	BucketURL     string
	Auth          *auth.Credentials
	Client        *awss3.Client
	PresignClient *awss3.PresignClient
}

func NewKodo() (s3.Interface, error) {
	conf := config.Config.Object.Kodo
	//init client
	cfg, err := awss3config.LoadDefaultConfig(context.TODO(),
		awss3config.WithRegion(conf.Bucket),
		awss3config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: conf.Endpoint}, nil
			})),
		awss3config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			conf.AccessKeyID,
			conf.AccessKeySecret,
			conf.SessionToken),
		),
	)
	if err != nil {
		panic(err)
	}
	client := awss3.NewFromConfig(cfg)
	presignClient := awss3.NewPresignClient(client)

	return &Kodo{
		AccessKey:     conf.AccessKeyID,
		SecretKey:     conf.AccessKeySecret,
		Region:        conf.Bucket,
		BucketURL:     conf.BucketURL,
		Auth:          auth.New(conf.AccessKeyID, conf.AccessKeySecret),
		Client:        client,
		PresignClient: presignClient,
	}, nil
}

func (k Kodo) Engine() string {
	return "kodo"
}

func (k Kodo) PartLimit() *s3.PartLimit {
	return &s3.PartLimit{
		MinPartSize: minPartSize,
		MaxPartSize: maxPartSize,
		MaxNumSize:  maxNumSize,
	}
}

func (k Kodo) InitiateMultipartUpload(ctx context.Context, name string) (*s3.InitiateMultipartUploadResult, error) {
	result, err := k.Client.CreateMultipartUpload(ctx, &awss3.CreateMultipartUploadInput{
		Bucket: aws.String(k.Region),
		Key:    aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	return &s3.InitiateMultipartUploadResult{
		UploadID: aws.ToString(result.UploadId),
		Bucket:   aws.ToString(result.Bucket),
		Key:      aws.ToString(result.Key),
	}, nil
}

func (k Kodo) CompleteMultipartUpload(ctx context.Context, uploadID string, name string, parts []s3.Part) (*s3.CompleteMultipartUploadResult, error) {
	kodoParts := make([]awss3types.CompletedPart, len(parts))
	for i, part := range parts {
		kodoParts[i] = awss3types.CompletedPart{
			PartNumber: aws.Int32(int32(part.PartNumber)),
			ETag:       aws.String(part.ETag),
		}
	}
	result, err := k.Client.CompleteMultipartUpload(ctx, &awss3.CompleteMultipartUploadInput{
		Bucket:          aws.String(k.Region),
		Key:             aws.String(name),
		UploadId:        aws.String(uploadID),
		MultipartUpload: &awss3types.CompletedMultipartUpload{Parts: kodoParts},
	})
	if err != nil {
		return nil, err
	}
	return &s3.CompleteMultipartUploadResult{
		Location: aws.ToString(result.Location),
		Bucket:   aws.ToString(result.Bucket),
		Key:      aws.ToString(result.Key),
		ETag:     strings.ToLower(strings.ReplaceAll(aws.ToString(result.ETag), `"`, ``)),
	}, nil
}

func (k Kodo) PartSize(ctx context.Context, size int64) (int64, error) {
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

func (k Kodo) AuthSign(ctx context.Context, uploadID string, name string, expire time.Duration, partNumbers []int) (*s3.AuthSignResult, error) {
	result := s3.AuthSignResult{
		URL:    k.BucketURL + "/" + name,
		Query:  url.Values{"uploadId": {uploadID}},
		Header: make(http.Header),
		Parts:  make([]s3.SignPart, len(partNumbers)),
	}
	for i, partNumber := range partNumbers {
		part, _ := k.PresignClient.PresignUploadPart(ctx, &awss3.UploadPartInput{
			Bucket:     aws.String(k.Region),
			UploadId:   aws.String(uploadID),
			Key:        aws.String(name),
			PartNumber: aws.Int32(int32(partNumber)),
		})
		result.Parts[i] = s3.SignPart{
			PartNumber: partNumber,
			URL:        part.URL,
			Header:     part.SignedHeader,
		}
	}
	return &result, nil

}

func (k Kodo) PresignedPutObject(ctx context.Context, name string, expire time.Duration) (string, error) {
	object, err := k.PresignClient.PresignPutObject(ctx, &awss3.PutObjectInput{
		Bucket: aws.String(k.Region),
		Key:    aws.String(name),
	}, func(po *awss3.PresignOptions) {
		po.Expires = expire
	})
	return object.URL, err

}

func (k Kodo) DeleteObject(ctx context.Context, name string) error {
	_, err := k.Client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: aws.String(k.Region),
		Key:    aws.String(name),
	})
	return err
}

func (k Kodo) CopyObject(ctx context.Context, src string, dst string) (*s3.CopyObjectInfo, error) {
	result, err := k.Client.CopyObject(ctx, &awss3.CopyObjectInput{
		Bucket:     aws.String(k.Region),
		CopySource: aws.String(k.Region + "/" + src),
		Key:        aws.String(dst),
	})
	if err != nil {
		return nil, err
	}
	return &s3.CopyObjectInfo{
		Key:  dst,
		ETag: strings.ToLower(strings.ReplaceAll(aws.ToString(result.CopyObjectResult.ETag), `"`, ``)),
	}, nil
}

func (k Kodo) StatObject(ctx context.Context, name string) (*s3.ObjectInfo, error) {
	info, err := k.Client.HeadObject(ctx, &awss3.HeadObjectInput{
		Bucket: aws.String(k.Region),
		Key:    aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	res := &s3.ObjectInfo{Key: name}
	res.Size = aws.ToInt64(info.ContentLength)
	res.ETag = strings.ToLower(strings.ReplaceAll(aws.ToString(info.ETag), `"`, ``))
	return res, nil
}

func (k Kodo) IsNotFound(err error) bool {
	return true
}

func (k Kodo) AbortMultipartUpload(ctx context.Context, uploadID string, name string) error {
	_, err := k.Client.AbortMultipartUpload(ctx, &awss3.AbortMultipartUploadInput{
		UploadId: aws.String(uploadID),
		Bucket:   aws.String(k.Region),
		Key:      aws.String(name),
	})
	return err
}

func (k Kodo) ListUploadedParts(ctx context.Context, uploadID string, name string, partNumberMarker int, maxParts int) (*s3.ListUploadedPartsResult, error) {
	result, err := k.Client.ListParts(ctx, &awss3.ListPartsInput{
		Key:              aws.String(name),
		UploadId:         aws.String(uploadID),
		Bucket:           aws.String(k.Region),
		MaxParts:         aws.Int32(int32(maxParts)),
		PartNumberMarker: aws.String(strconv.Itoa(partNumberMarker)),
	})
	if err != nil {
		return nil, err
	}
	res := &s3.ListUploadedPartsResult{
		Key:           aws.ToString(result.Key),
		UploadID:      aws.ToString(result.UploadId),
		MaxParts:      int(aws.ToInt32(result.MaxParts)),
		UploadedParts: make([]s3.UploadedPart, len(result.Parts)),
	}
	// int to string
	NextPartNumberMarker, err := strconv.Atoi(aws.ToString(result.NextPartNumberMarker))
	if err != nil {
		return nil, err
	}
	res.NextPartNumberMarker = NextPartNumberMarker
	for i, part := range result.Parts {
		res.UploadedParts[i] = s3.UploadedPart{
			PartNumber:   int(aws.ToInt32(part.PartNumber)),
			LastModified: aws.ToTime(part.LastModified),
			ETag:         aws.ToString(part.ETag),
			Size:         aws.ToInt64(part.Size),
		}
	}
	return res, nil
}

func (k Kodo) AccessURL(ctx context.Context, name string, expire time.Duration, opt *s3.AccessURLOption) (string, error) {
	//get object head
	info, err := k.Client.HeadObject(ctx, &awss3.HeadObjectInput{
		Bucket: aws.String(k.Region),
		Key:    aws.String(name),
	})
	if err != nil {
		return "", errors.New("AccessURL object not found")
	}
	if opt != nil {
		if opt.ContentType != aws.ToString(info.ContentType) {
			//修改文件类型
			err := k.SetObjectContentType(ctx, name, opt.ContentType)
			if err != nil {
				return "", errors.New("AccessURL setContentType error")
			}
		}
	}
	imageMogr := ""
	//image dispose
	if opt != nil {
		if opt.Image != nil {
			//https://developer.qiniu.com/dora/8255/the-zoom
			process := ""
			if opt.Image.Width > 0 {
				process += strconv.Itoa(opt.Image.Width) + "x"
			}
			if opt.Image.Height > 0 {
				if opt.Image.Width > 0 {
					process += strconv.Itoa(opt.Image.Height)
				} else {
					process += "x" + strconv.Itoa(opt.Image.Height)
				}
			}
			imageMogr = "imageMogr2/thumbnail/" + process
		}
	}
	//expire
	deadline := time.Now().Add(time.Second * expire).Unix()
	domain := k.BucketURL
	query := url.Values{}
	if opt != nil && opt.Filename != "" {
		query.Add("attname", opt.Filename)
	}
	privateURL := storage.MakePrivateURLv2WithQuery(k.Auth, domain, name, query, deadline)
	if imageMogr != "" {
		privateURL += "&" + imageMogr
	}
	return privateURL, nil
}

func (k *Kodo) SetObjectContentType(ctx context.Context, name string, contentType string) error {
	//set object content-type
	_, err := k.Client.CopyObject(ctx, &awss3.CopyObjectInput{
		Bucket:            aws.String(k.Region),
		CopySource:        aws.String(k.Region + "/" + name),
		Key:               aws.String(name),
		ContentType:       aws.String(contentType),
		MetadataDirective: awss3types.MetadataDirectiveReplace,
	})
	return err
}
