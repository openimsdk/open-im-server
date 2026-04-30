package api

import (
	"github.com/gin-gonic/gin"
	pbredpacket "github.com/openimsdk/protocol/redpacket"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/log"
)

type RedPacketApi struct {
	Client pbredpacket.RedPacketClient
}

func NewRedPacketApi(client pbredpacket.RedPacketClient) *RedPacketApi {
	return &RedPacketApi{Client: client}
}

func (h *RedPacketApi) CreateOrder(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.CreateOrderReq](ctx)
	if err != nil {
		log.ZError(ctx, "redpacket create order parse failed", err)
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.CreateOrder(ctx, req)
	if err != nil {
		log.ZError(ctx, "redpacket create order rpc failed", err)
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) CreatedCallback(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.CreatedCallbackReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.CreatedCallback(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) GetDetail(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.GetDetailReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.GetDetail(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) IssueClaimSign(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.IssueClaimSignReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.IssueClaimSign(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) ClaimResult(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.ClaimResultReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.ClaimResult(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) IssueWalletBindChallenge(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.IssueWalletBindChallengeReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.IssueWalletBindChallenge(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) ConfirmWalletBind(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.ConfirmWalletBindReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.ConfirmWalletBind(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) GetWalletBinding(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.GetWalletBindingReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.GetWalletBinding(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

// Admin endpoints

func (h *RedPacketApi) AdminSetSigner(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.SetSignerReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.SetSigner(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) AdminSetToken(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.SetTokenReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.SetToken(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) AdminSetExpiry(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.SetExpiryReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.SetExpiry(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) AdminSetAllowAllTokens(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.SetAllowAllTokensReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.SetAllowAllTokens(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) AdminSetNativeTokenEnabled(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.SetNativeTokenEnabledReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.SetNativeTokenEnabled(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}

func (h *RedPacketApi) AdminParseTxEvents(ctx *gin.Context) {
	req, err := a2r.ParseRequestNotCheck[pbredpacket.ParseTxEventsReq](ctx)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	resp, err := h.Client.ParseTxEvents(ctx, req)
	if err != nil {
		apiresp.GinError(ctx, err)
		return
	}
	apiresp.GinSuccess(ctx, resp)
}
