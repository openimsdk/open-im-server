package api

import (
	"github.com/gin-gonic/gin"
	pbcaptcha "github.com/openimsdk/protocol/captcha"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/log"
)

type CaptchaApi struct {
	Client pbcaptcha.CaptchaClient
}

func NewCaptchaApi(client pbcaptcha.CaptchaClient) *CaptchaApi {
	return &CaptchaApi{Client: client}
}

func (c *CaptchaApi) GenerateCaptcha(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbcaptcha.GenerateCaptchaReq](ctx)
	if err != nil {
		log.ZError(ctx, "captcha generate request parse failed", err)
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := c.Client.GenerateCaptcha(ctx, req)
	if err != nil {
		log.ZError(ctx, "captcha generate rpc failed", err)
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (c *CaptchaApi) VerifyCaptcha(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbcaptcha.VerifyCaptchaReq](ctx)
	if err != nil {
		log.ZError(ctx, "captcha verify request parse failed", err)
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := c.Client.VerifyCaptcha(ctx, req)
	if err != nil {
		log.ZError(ctx, "captcha verify rpc failed", err, "captchaID", req.GetCaptchaID(), "clickPoints", req.GetClickPoints())
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}
