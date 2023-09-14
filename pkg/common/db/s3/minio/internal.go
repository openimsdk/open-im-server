package minio

import (
	"net/url"
	_ "unsafe"

	"github.com/minio/minio-go/v7"
)

//go:linkname makeTargetURL github.com/minio/minio-go/v7.(*Client).makeTargetURL
func makeTargetURL(client *minio.Client, bucketName, objectName, bucketLocation string, isVirtualHostStyle bool, queryValues url.Values) (*url.URL, error)
