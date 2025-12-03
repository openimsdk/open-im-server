package cachekey

import "strconv"

const (
	object         = "OBJECT:"
	s3             = "S3:"
	minioImageInfo = "MINIO:IMAGE:"
	minioThumbnail = "MINIO:THUMBNAIL:"
)

func GetObjectKey(engine string, name string) string {
	return object + engine + ":" + name
}

func GetS3Key(engine string, name string) string {
	return s3 + engine + ":" + name
}

func GetObjectImageInfoKey(key string) string {
	return minioImageInfo + key
}

func GetMinioImageThumbnailKey(key string, format string, width int, height int) string {
	return minioThumbnail + format + ":w" + strconv.Itoa(width) + ":h" + strconv.Itoa(height) + ":" + key
}
