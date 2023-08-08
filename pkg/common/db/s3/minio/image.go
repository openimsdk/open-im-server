package minio

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

const (
	formatPng  = "png"
	formatJpeg = "jpeg"
	formatJpg  = "jpg"
	formatGif  = "gif"
)

func ImageStat(reader io.Reader) (image.Image, string, error) {
	return image.Decode(reader)
}

func ImageWidthHeight(img image.Image) (int, int) {
	bounds := img.Bounds().Max
	return bounds.X, bounds.Y
}

func resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	imgWidth := bounds.Max.X
	imgHeight := bounds.Max.Y

	// 计算缩放比例
	scaleWidth := float64(maxWidth) / float64(imgWidth)
	scaleHeight := float64(maxHeight) / float64(imgHeight)

	// 如果都为0，则不缩放，返回原始图片
	if maxWidth == 0 && maxHeight == 0 {
		return img
	}

	// 如果宽度和高度都大于0，则选择较小的缩放比例，以保持宽高比
	if maxWidth > 0 && maxHeight > 0 {
		scale := scaleWidth
		if scaleHeight < scaleWidth {
			scale = scaleHeight
		}

		// 计算缩略图尺寸
		thumbnailWidth := int(float64(imgWidth) * scale)
		thumbnailHeight := int(float64(imgHeight) * scale)

		// 使用"image"库的Resample方法生成缩略图
		thumbnail := image.NewRGBA(image.Rect(0, 0, thumbnailWidth, thumbnailHeight))
		for y := 0; y < thumbnailHeight; y++ {
			for x := 0; x < thumbnailWidth; x++ {
				srcX := int(float64(x) / scale)
				srcY := int(float64(y) / scale)
				thumbnail.Set(x, y, img.At(srcX, srcY))
			}
		}

		return thumbnail
	}

	// 如果只指定了宽度或高度，则根据最大不超过的规则生成缩略图
	if maxWidth > 0 {
		thumbnailWidth := maxWidth
		thumbnailHeight := int(float64(imgHeight) * scaleWidth)

		// 使用"image"库的Resample方法生成缩略图
		thumbnail := image.NewRGBA(image.Rect(0, 0, thumbnailWidth, thumbnailHeight))
		for y := 0; y < thumbnailHeight; y++ {
			for x := 0; x < thumbnailWidth; x++ {
				srcX := int(float64(x) / scaleWidth)
				srcY := int(float64(y) / scaleWidth)
				thumbnail.Set(x, y, img.At(srcX, srcY))
			}
		}

		return thumbnail
	}

	if maxHeight > 0 {
		thumbnailWidth := int(float64(imgWidth) * scaleHeight)
		thumbnailHeight := maxHeight

		// 使用"image"库的Resample方法生成缩略图
		thumbnail := image.NewRGBA(image.Rect(0, 0, thumbnailWidth, thumbnailHeight))
		for y := 0; y < thumbnailHeight; y++ {
			for x := 0; x < thumbnailWidth; x++ {
				srcX := int(float64(x) / scaleHeight)
				srcY := int(float64(y) / scaleHeight)
				thumbnail.Set(x, y, img.At(srcX, srcY))
			}
		}

		return thumbnail
	}

	// 默认情况下，返回原始图片
	return img
}
