package captcha

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
	"golang.org/x/image/font/gofont/goregular"
)

// buildClickCaptcha constructs a click.Captcha instance configured with
// alphanumeric characters, a bundled Go font, and the embedded background images.
func buildClickCaptcha() (click.Captcha, error) {
	font, err := loadGoRegularFont()
	if err != nil {
		return nil, fmt.Errorf("load font: %w", err)
	}
	backgrounds, err := loadBackgrounds()
	if err != nil {
		return nil, fmt.Errorf("load captcha backgrounds: %w", err)
	}

	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 6, Max: 8}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 3, Max: 4}),
		click.WithRangeSize(option.RangeVal{Min: 26, Max: 34}),
		click.WithDisplayShadow(true),
	)
	builder.SetResources(
		click.WithChars(alphanumChars),
		click.WithFonts([]*truetype.Font{font}),
		click.WithBackgrounds(backgrounds),
	)
	return builder.Make(), nil
}

// loadGoRegularFont parses the bundled Go Regular TTF font.
func loadGoRegularFont() (*truetype.Font, error) {
	return freetype.ParseFont(goregular.TTF)
}

// loadBackgrounds decodes the embedded JPEG background images.
func loadBackgrounds() ([]image.Image, error) {
	const count = 5
	images := make([]image.Image, 0, count)
	for i := 1; i <= count; i++ {
		path := fmt.Sprintf("resources/images/image-%d.jpg", i)
		img, err := decodeEmbedImage(path)
		if err != nil {
			return nil, fmt.Errorf("decode %s: %w", path, err)
		}
		images = append(images, img)
	}
	return images, nil
}

func decodeEmbedImage(path string) (image.Image, error) {
	f, err := resourceFS.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
