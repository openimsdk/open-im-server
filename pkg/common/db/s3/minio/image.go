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

	// Calculating scaling
	scaleWidth := float64(maxWidth) / float64(imgWidth)
	scaleHeight := float64(maxHeight) / float64(imgHeight)

	// If both are 0, then no scaling is done and the original image is returned
	if maxWidth == 0 && maxHeight == 0 {
		return img
	}

	// If both width and height are greater than 0, select a smaller zoom ratio to maintain the aspect ratio
	if maxWidth > 0 && maxHeight > 0 {
		scale := scaleWidth
		if scaleHeight < scaleWidth {
			scale = scaleHeight
		}

		// Calculate Thumbnail Size
		thumbnailWidth := int(float64(imgWidth) * scale)
		thumbnailHeight := int(float64(imgHeight) * scale)

		// Thumbnails are generated using the Resample method of the "image" library.
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

	// If only width or height is specified, thumbnails are generated based on the maximum not to exceed rule
	if maxWidth > 0 {
		thumbnailWidth := maxWidth
		thumbnailHeight := int(float64(imgHeight) * scaleWidth)

		// Thumbnails are generated using the Resample method of the "image" library.
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

		// Thumbnails are generated using the Resample method of the "image" library.
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

	// By default, the original image is returned
	return img
}
