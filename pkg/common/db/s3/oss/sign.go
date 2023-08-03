package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"net/http"
	_ "unsafe"
)

//go:linkname ossSignHeader github.com/aliyun/aliyun-oss-go-sdk/oss.(*Conn).signHeader
func ossSignHeader(c *oss.Conn, req *http.Request, canonicalizedResource string)
