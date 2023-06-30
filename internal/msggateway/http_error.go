package msggateway

import "github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"

func httpError(ctx *UserConnContext, err error) {
	apiresp.HttpError(ctx.RespWriter, err)
}
