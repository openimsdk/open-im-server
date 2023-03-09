package msggateway

import (
	"net/http"
)

func httpError(ctx *UserConnContext, err error) {
	code := http.StatusUnauthorized
	ctx.SetHeader("Sec-Websocket-Version", "13")
	ctx.SetHeader("ws_err_msg", err.Error())
	//if errors.Is(err, errs.ErrTokenExpired) {
	//	code = errs.ErrTokenExpired.Code()
	//}
	//if errors.Is(err, errs.ErrTokenInvalid) {
	//	code = errs.ErrTokenInvalid.Code()
	//}
	//if errors.Is(err, errs.ErrTokenMalformed) {
	//	code = errs.ErrTokenMalformed.Code()
	//}
	//if errors.Is(err, errs.ErrTokenNotValidYet) {
	//	code = errs.ErrTokenNotValidYet.Code()
	//}
	//if errors.Is(err, errs.ErrTokenUnknown) {
	//	code = errs.ErrTokenUnknown.Code()
	//}
	//if errors.Is(err, errs.ErrTokenKicked) {
	//	code = errs.ErrTokenKicked.Code()
	//}
	//if errors.Is(err, errs.ErrTokenDifferentPlatformID) {
	//	code = errs.ErrTokenDifferentPlatformID.Code()
	//}
	//if errors.Is(err, errs.ErrTokenDifferentUserID) {
	//	code = errs.ErrTokenDifferentUserID.Code()
	//}
	//if errors.Is(err, errs.ErrConnOverMaxNumLimit) {
	//	code = errs.ErrConnOverMaxNumLimit.Code()
	//}
	//if errors.Is(err, errs.ErrConnArgsErr) {
	//	code = errs.ErrConnArgsErr.Code()
	//}
	ctx.ErrReturn(err.Error(), code)
}
