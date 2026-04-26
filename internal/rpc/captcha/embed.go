package captcha

import "embed"

// resourceFS embeds background images and tile images at compile time.
// Background images come from go-captcha-resources (sourcedata/images/image-{1..5}).
// Tile images come from go-captcha-resources (sourcedata/tiles/tile-{1..4}):
//   overlay.png  → GraphImage.OverlayImage
//   shadow.png   → GraphImage.ShadowImage
//   mask.png     → GraphImage.MaskImage
//
//go:embed resources/images/*.jpg resources/tiles/*/*.png
var resourceFS embed.FS
