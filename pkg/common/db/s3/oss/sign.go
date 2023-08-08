package oss

import (
	"net/http"
	_ "unsafe"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//go:linkname ossSignHeader github.com/aliyun/aliyun-oss-go-sdk/oss.(*Conn).signHeader
func ossSignHeader(c *oss.Conn, req *http.Request, canonicalizedResource string)
