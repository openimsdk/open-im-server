package utils

import (
	"errors"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

func GenSmallImage(src, dst string) error {
	fIn, _ := os.Open(src)
	defer fIn.Close()

	fOut, _ := os.Create(dst)
	defer fOut.Close()

	if err := scale(fIn, fOut, 0, 0, 0); err != nil {
		return err
	}
	return nil
}

func scale(in io.Reader, out io.Writer, width, height, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X / 2
		height = origin.Bounds().Max.Y / 2
	}
	if quality == 0 {
		quality = 25
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	switch fm {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	case "bmp":
		return bmp.Encode(out, canvas)
	default:
		return errors.New("ERROR FORMAT")
	}

	return nil
}
