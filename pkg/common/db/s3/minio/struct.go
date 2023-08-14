package minio

type minioImageInfo struct {
	NotImage bool   `json:"notImage,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Format   string `json:"format,omitempty"`
}
