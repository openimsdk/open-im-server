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

// resizeImage resizes an image to a specified maximum width and height, maintaining the aspect ratio.
// If both maxWidth and maxHeight are set to 0, the original image is returned.
// If both are non-zero, the image is scaled to fit within the constraints while maintaining aspect ratio.
// If only one of maxWidth or maxHeight is non-zero, the image is scaled accordingly.
func resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	imgWidth, imgHeight := bounds.Dx(), bounds.Dy()

	// Return original image if no resizing is needed.
	if maxWidth == 0 && maxHeight == 0 {
		return img
	}

	var scale float64 = 1
	if maxWidth > 0 && maxHeight > 0 {
		scaleWidth := float64(maxWidth) / float64(imgWidth)
		scaleHeight := float64(maxHeight) / float64(imgHeight)
		// Choose the smaller scale to fit both constraints.
		scale = min(scaleWidth, scaleHeight)
	} else if maxWidth > 0 {
		scale = float64(maxWidth) / float64(imgWidth)
	} else if maxHeight > 0 {
		scale = float64(maxHeight) / float64(imgHeight)
	}

	newWidth := int(float64(imgWidth) * scale)
	newHeight := int(float64(imgHeight) * scale)

	// Resize the image by creating a new image and manually copying pixels.
	thumbnail := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / scale)
			srcY := int(float64(y) / scale)
			thumbnail.Set(x, y, img.At(srcX, srcY))
		}
	}

	return thumbnail
}

// min returns the smaller of x or y.
func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
