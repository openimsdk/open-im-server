package minio

import (
	"github.com/minio/minio-go/v7"
	"net/url"
	_ "unsafe"
)

//go:linkname makeTargetURL github.com/minio/minio-go/v7.(*Client).makeTargetURL
func makeTargetURL(client *minio.Client, bucketName, objectName, bucketLocation string, isVirtualHostStyle bool, queryValues url.Values) (*url.URL, error)
