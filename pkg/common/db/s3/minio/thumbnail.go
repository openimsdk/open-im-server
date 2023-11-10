package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/minio/minio-go/v7"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
)

func (m *Minio) getImageThumbnailURL(ctx context.Context, name string, expire time.Duration, opt *s3.Image) (string, error) {
	var img image.Image
	info, err := m.cache.GetImageObjectKeyInfo(ctx, name, func(ctx context.Context) (info *cache.MinioImageInfo, err error) {
		info, img, err = m.getObjectImageInfo(ctx, name)
		return
	})
	if err != nil {
		return "", err
	}
	if !info.IsImg {
		return "", errs.ErrData.Wrap("object not image")
	}
	if opt.Width > info.Width || opt.Width <= 0 {
		opt.Width = info.Width
	}
	if opt.Height > info.Height || opt.Height <= 0 {
		opt.Height = info.Height
	}
	opt.Format = strings.ToLower(opt.Format)
	if opt.Format == formatJpg {
		opt.Format = formatJpeg
	}
	switch opt.Format {
	case formatPng, formatJpeg, formatGif:
	default:
		opt.Format = ""
	}
	reqParams := make(url.Values)
	if opt.Width == info.Width && opt.Height == info.Height && (opt.Format == info.Format || opt.Format == "") {
		reqParams.Set("response-content-type", "image/"+info.Format)
		return m.PresignedGetObject(ctx, name, expire, reqParams)
	}
	if opt.Format == "" {
		switch opt.Format {
		case formatGif:
			opt.Format = formatGif
		case formatJpeg:
			opt.Format = formatJpeg
		case formatPng:
			opt.Format = formatPng
		default:
			opt.Format = formatPng
		}
	}
	key, err := m.cache.GetThumbnailKey(ctx, name, opt.Format, opt.Width, opt.Height, func(ctx context.Context) (string, error) {
		if img == nil {
			reader, err := m.core.Client.GetObject(ctx, m.bucket, name, minio.GetObjectOptions{})
			if err != nil {
				return "", err
			}
			defer reader.Close()
			img, _, err = ImageStat(reader)
			if err != nil {
				return "", err
			}
		}
		thumbnail := resizeImage(img, opt.Width, opt.Height)
		buf := bytes.NewBuffer(nil)
		switch opt.Format {
		case formatPng:
			err = png.Encode(buf, thumbnail)
		case formatJpeg:
			err = jpeg.Encode(buf, thumbnail, nil)
		case formatGif:
			err = gif.Encode(buf, thumbnail, nil)
		}
		cacheKey := filepath.Join(imageThumbnailPath, info.Etag, fmt.Sprintf("image_w%d_h%d.%s", opt.Width, opt.Height, opt.Format))
		if _, err := m.core.Client.PutObject(ctx, m.bucket, cacheKey, buf, int64(buf.Len()), minio.PutObjectOptions{}); err != nil {
			return "", err
		}
		return cacheKey, nil
	})
	if err != nil {
		return "", err
	}
	reqParams.Set("response-content-type", "image/"+opt.Format)
	return m.PresignedGetObject(ctx, key, expire, reqParams)
}

func (m *Minio) getObjectImageInfo(ctx context.Context, name string) (*cache.MinioImageInfo, image.Image, error) {
	fileInfo, err := m.StatObject(ctx, name)
	if err != nil {
		return nil, nil, err
	}
	if fileInfo.Size > maxImageSize {
		return nil, nil, errors.New("file size too large")
	}
	imageData, err := m.getObjectData(ctx, name, fileInfo.Size)
	if err != nil {
		return nil, nil, err
	}
	var info cache.MinioImageInfo
	imageInfo, format, err := ImageStat(bytes.NewReader(imageData))
	if err == nil {
		info.IsImg = true
		info.Format = format
		info.Width, info.Height = ImageWidthHeight(imageInfo)
	} else {
		info.IsImg = false
	}
	info.Etag = fileInfo.ETag
	return &info, imageInfo, nil
}

func (m *Minio) delObjectImageInfoKey(ctx context.Context, key string, size int64) {
	if size > 0 && size > maxImageSize {
		return
	}
	if err := m.cache.DelObjectImageInfoKey(key).ExecDel(ctx); err != nil {
		log.ZError(ctx, "DelObjectImageInfoKey failed", err, "key", key)
	}
}
