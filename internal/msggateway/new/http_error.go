package new

import (
	"OpenIM/pkg/common/constant"
	"errors"
	"net/http"
)

func httpError(ctx *UserConnContext, err error) {
	code := http.StatusUnauthorized
	ctx.SetHeader("Sec-Websocket-Version", "13")
	ctx.SetHeader("ws_err_msg", err.Error())
	if errors.Is(err, constant.ErrTokenExpired) {
		code = int(constant.ErrTokenExpired.ErrCode)
	}
	if errors.Is(err, constant.ErrTokenInvalid) {
		code = int(constant.ErrTokenInvalid.ErrCode)
	}
	if errors.Is(err, constant.ErrTokenMalformed) {
		code = int(constant.ErrTokenMalformed.ErrCode)
	}
	if errors.Is(err, constant.ErrTokenNotValidYet) {
		code = int(constant.ErrTokenNotValidYet.ErrCode)
	}
	if errors.Is(err, constant.ErrTokenUnknown) {
		code = int(constant.ErrTokenUnknown.ErrCode)
	}
	if errors.Is(err, constant.ErrTokenKicked) {
		code = int(constant.ErrTokenKicked.ErrCode)
	}
	if errors.Is(err, constant.ErrTokenDifferentPlatformID) {
		code = int(constant.ErrTokenDifferentPlatformID.ErrCode)
	}
	if errors.Is(err, constant.ErrTokenDifferentUserID) {
		code = int(constant.ErrTokenDifferentUserID.ErrCode)
	}
	if errors.Is(err, constant.ErrConnOverMaxNumLimit) {
		code = int(constant.ErrConnOverMaxNumLimit.ErrCode)
	}
	if errors.Is(err, constant.ErrConnArgsErr) {
		code = int(constant.ErrConnArgsErr.ErrCode)
	}
	ctx.ErrReturn(err.Error(), code)
}
