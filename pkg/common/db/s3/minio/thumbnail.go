package minio

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"
)

//func (m *Minio) getHashImageInfo1(ctx context.Context, key string) (*cache.MinioImageInfo, image.Image, error) {
//
//	return nil, nil, nil
//}

func (m *Minio) get1(ctx context.Context, key string, format string, width int, height int, info *cache.MinioImageInfo, img image.Image) (string, error) {

	return "", nil
}

func (m *Minio) getHashImageInfo(ctx context.Context, name string, expire time.Duration, opt *s3.Image) (string, error) {
	info, img, err := m.cache.GetImageObjectKeyInfo(ctx, name, m.getObjectImageInfo)
	if err != nil {
		return "", err
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
		//if info.Format == formatGif {
		//	opt.Format = formatGif
		//} else {
		//	opt.Format = formatJpeg
		//}
	}
	if opt.Width == info.Width && opt.Height == info.Height && (opt.Format == info.Format || opt.Format == "") {
		return "", nil
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
		cacheKey := filepath.Join(pathInfo, info.Etag, fmt.Sprintf("image_w%d_h%d.%s", opt.Width, opt.Height, opt.Format))
		if _, err := m.core.Client.PutObject(ctx, m.bucket, cacheKey, buf, int64(buf.Len()), minio.PutObjectOptions{}); err != nil {
			return "", err
		}
		return cacheKey, nil
	})
	if err != nil {
		return "", err
	}
	return m.presignedGetObject(ctx, key, expire, reqParams)
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

func (m *Minio) GetImageThumbnail(ctx context.Context, name string, expire time.Duration, opt *s3.Image) (string, error) {
	fileInfo, err := m.StatObject(ctx, name)
	if err != nil {
		return "", err
	}
	if fileInfo.Size > maxImageSize {
		return "", errors.New("file size too large")
	}
	objectInfoPath := path.Join(pathInfo, fileInfo.ETag, "image.json")
	var (
		img  image.Image
		info minioImageInfo
	)
	data, err := m.getObjectData(ctx, objectInfoPath, maxImageInfoSize)
	if err == nil {
		if err := json.Unmarshal(data, &info); err != nil {
			return "", fmt.Errorf("unmarshal minio image info.json error: %w", err)
		}
		if info.NotImage {
			return "", errors.New("not image")
		}
	} else if m.IsNotFound(err) {
		reader, err := m.core.Client.GetObject(ctx, m.bucket, name, minio.GetObjectOptions{})
		if err != nil {
			return "", err
		}
		defer reader.Close()
		imageInfo, format, err := ImageStat(reader)
		if err == nil {
			info.NotImage = false
			info.Format = format
			info.Width, info.Height = ImageWidthHeight(imageInfo)
			img = imageInfo
		} else {
			info.NotImage = true
		}
		data, err := json.Marshal(&info)
		if err != nil {
			return "", err
		}
		if _, err := m.core.Client.PutObject(ctx, m.bucket, objectInfoPath, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{}); err != nil {
			return "", err
		}
	} else {
		return "", err
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
	case formatPng:
	case formatJpeg:
	case formatGif:
	default:
		if info.Format == formatGif {
			opt.Format = formatGif
		} else {
			opt.Format = formatJpeg
		}
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content-type", "image/"+opt.Format)
	if opt.Width == info.Width && opt.Height == info.Height && opt.Format == info.Format {
		return m.presignedGetObject(ctx, name, expire, reqParams)
	}
	cacheKey := filepath.Join(pathInfo, fileInfo.ETag, fmt.Sprintf("image_w%d_h%d.%s", opt.Width, opt.Height, opt.Format))
	if _, err := m.core.Client.StatObject(ctx, m.bucket, cacheKey, minio.StatObjectOptions{}); err == nil {
		return m.presignedGetObject(ctx, cacheKey, expire, reqParams)
	} else if !m.IsNotFound(err) {
		return "", err
	}
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
	if _, err := m.core.Client.PutObject(ctx, m.bucket, cacheKey, buf, int64(buf.Len()), minio.PutObjectOptions{}); err != nil {
		return "", err
	}
	return m.presignedGetObject(ctx, cacheKey, expire, reqParams)
}
