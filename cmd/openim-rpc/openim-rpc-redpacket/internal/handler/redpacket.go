package handler

import (
	"net/http"

	"redpacket/internal/authctx"
	"redpacket/internal/service"
	"redpacket/pkg/resp"

	"github.com/gin-gonic/gin"
)

type RedPacketHandler struct {
	rpSvc *service.RedPacketService
}

func NewRedPacketHandler(rpSvc *service.RedPacketService) *RedPacketHandler {
	return &RedPacketHandler{rpSvc: rpSvc}
}

func (h *RedPacketHandler) CreateOrder(c *gin.Context) {
	if err := authctx.BindCurrentUserID(c); err != nil {
		resp.Forbidden(c, err.Error())
		return
	}

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	result, err := h.rpSvc.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	resp.OK(c, result)
}

func (h *RedPacketHandler) CreatedCallback(c *gin.Context) {
	var req service.CreatedCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.rpSvc.CreatedCallback(c.Request.Context(), &req); err != nil {
		resp.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	resp.OK(c, gin.H{"ok": true})
}

func (h *RedPacketHandler) Detail(c *gin.Context) {
	packetID := c.Query("packet_id")
	if packetID == "" {
		resp.BadRequest(c, "packet_id is required")
		return
	}

	detail, err := h.rpSvc.GetDetail(c.Request.Context(), packetID)
	if err != nil {
		resp.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}

	resp.OK(c, detail)
}

func (h *RedPacketHandler) ClaimSign(c *gin.Context) {
	if err := authctx.BindCurrentUserID(c); err != nil {
		resp.Forbidden(c, err.Error())
		return
	}

	var req struct {
		PacketID   string `json:"packet_id" binding:"required"`
		Claimer    string `json:"claimer" binding:"required"`
		RandomSeed string `json:"random_seed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	result, err := h.rpSvc.IssueClaimSign(c.Request.Context(), req.PacketID, req.Claimer, req.RandomSeed)
	if err != nil {
		resp.InternalError(c, "failed to issue claim signature: "+err.Error())
		return
	}

	resp.OK(c, result)
}

func (h *RedPacketHandler) ClaimResult(c *gin.Context) {
	if err := authctx.BindCurrentUserID(c); err != nil {
		resp.Forbidden(c, err.Error())
		return
	}

	var req service.ClaimResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.rpSvc.ClaimResult(c.Request.Context(), &req); err != nil {
		resp.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	resp.OK(c, gin.H{"ok": true})
}

func (h *RedPacketHandler) WalletBindChallenge(c *gin.Context) {
	if err := authctx.BindCurrentUserID(c); err != nil {
		resp.Forbidden(c, err.Error())
		return
	}

	var req service.WalletBindChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	result, err := h.rpSvc.IssueWalletBindChallenge(c.Request.Context(), &req)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	resp.OK(c, result)
}

func (h *RedPacketHandler) WalletBindConfirm(c *gin.Context) {
	var req service.WalletBindConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	result, err := h.rpSvc.ConfirmWalletBind(c.Request.Context(), &req)
	if err != nil {
		resp.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	resp.OK(c, result)
}

func (h *RedPacketHandler) WalletBindDetail(c *gin.Context) {
	if err := authctx.BindCurrentUserID(c); err != nil {
		resp.Forbidden(c, err.Error())
		return
	}

	chainType := c.Query("chain_type")
	walletAddress := c.Query("wallet_address")
	if chainType == "" || walletAddress == "" {
		resp.BadRequest(c, "chain_type and wallet_address are required")
		return
	}

	result, err := h.rpSvc.GetWalletBinding(c.Request.Context(), "", chainType, walletAddress)
	if err != nil {
		resp.Fail(c, http.StatusNotFound, 404, err.Error())
		return
	}

	resp.OK(c, result)
}
