package captcha

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/wenlng/go-captcha/v2/slide"
)

// loadResources reads the embedded files and returns slide.Resource options
// ready to be passed to slide.NewBuilder().SetResources(...).
func loadResources() ([]slide.Resource, error) {
	backgrounds, err := loadBackgrounds()
	if err != nil {
		return nil, fmt.Errorf("load captcha backgrounds: %w", err)
	}
	graphImages, err := loadGraphImages()
	if err != nil {
		return nil, fmt.Errorf("load captcha graph images: %w", err)
	}
	return []slide.Resource{
		slide.WithBackgrounds(backgrounds),
		slide.WithGraphImages(graphImages),
	}, nil
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

// loadGraphImages decodes the 4 sets of tile overlay/shadow/mask PNG images.
func loadGraphImages() ([]*slide.GraphImage, error) {
	const count = 4
	graphs := make([]*slide.GraphImage, 0, count)
	for i := 1; i <= count; i++ {
		overlay, err := decodeEmbedImage(fmt.Sprintf("resources/tiles/tile-%d/overlay.png", i))
		if err != nil {
			return nil, fmt.Errorf("decode tile-%d overlay: %w", i, err)
		}
		shadow, err := decodeEmbedImage(fmt.Sprintf("resources/tiles/tile-%d/shadow.png", i))
		if err != nil {
			return nil, fmt.Errorf("decode tile-%d shadow: %w", i, err)
		}
		mask, err := decodeEmbedImage(fmt.Sprintf("resources/tiles/tile-%d/mask.png", i))
		if err != nil {
			return nil, fmt.Errorf("decode tile-%d mask: %w", i, err)
		}
		graphs = append(graphs, &slide.GraphImage{
			OverlayImage: overlay,
			ShadowImage:  shadow,
			MaskImage:    mask,
		})
	}
	return graphs, nil
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
