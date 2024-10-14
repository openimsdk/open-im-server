package jssdk

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/checker"
	"github.com/openimsdk/tools/errs"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"strings"
)

func field[A, B, C any](ctx context.Context, fn func(ctx context.Context, req *A, opts ...grpc.CallOption) (*B, error), req *A, get func(*B) C) (C, error) {
	resp, err := fn(ctx, req)
	if err != nil {
		var c C
		return c, err
	}
	return get(resp), nil
}

func call[A, B any](c *gin.Context, fn func(ctx context.Context, req *A) (*B, error)) {
	var isJSON bool
	switch contentType := c.GetHeader("Content-Type"); {
	case contentType == "":
		isJSON = true
	case strings.Contains(contentType, "application/json"):
		isJSON = true
	case strings.Contains(contentType, "application/protobuf"):
	case strings.Contains(contentType, "application/x-protobuf"):
	default:
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("unsupported content type"))
		return
	}
	var req *A
	if isJSON {
		var err error
		req, err = a2r.ParseRequest[A](c)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
	} else {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			apiresp.GinError(c, err)
			return
		}
		req = new(A)
		if err := proto.Unmarshal(body, any(req).(proto.Message)); err != nil {
			apiresp.GinError(c, err)
			return
		}
		if err := checker.Validate(&req); err != nil {
			apiresp.GinError(c, err)
			return
		}
	}
	resp, err := fn(c, req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if isJSON {
		apiresp.GinSuccess(c, resp)
		return
	}
	body, err := proto.Marshal(any(resp).(proto.Message))
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, body)
}
