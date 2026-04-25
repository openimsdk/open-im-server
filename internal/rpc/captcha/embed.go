package captcha

import "embed"

// resourceFS embeds background images for the click captcha at compile time.
//
//go:embed resources/images/*.jpg
var resourceFS embed.FS
