package msggateway

import (
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/log"
)

func httpError(ctx *UserConnContext, err error) {
	log.ZWarn(ctx, "ws connection error", err)
	apiresp.HttpError(ctx.RespWriter, err)
}
