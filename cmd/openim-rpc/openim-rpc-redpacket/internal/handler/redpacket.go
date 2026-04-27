package handler

import (
	"net/http"

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
	var req struct {
		PacketID   string `json:"packet_id" binding:"required"`
		Claimer    string `json:"claimer" binding:"required"`
		UserID     string `json:"user_id" binding:"required"`
		RandomSeed string `json:"random_seed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.rpSvc.CanClaim(c.Request.Context(), req.PacketID, req.Claimer, req.UserID); err != nil {
		resp.Forbidden(c, err.Error())
		return
	}

	result, err := h.rpSvc.IssueClaimSign(c.Request.Context(), req.PacketID, req.Claimer, req.UserID, req.RandomSeed)
	if err != nil {
		resp.InternalError(c, "failed to issue claim signature: "+err.Error())
		return
	}

	resp.OK(c, result)
}

func (h *RedPacketHandler) ClaimResult(c *gin.Context) {
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
