package handler

import (
	"redpacket/internal/service"
	"redpacket/pkg/resp"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminSvc *service.AdminService
}

func NewAdminHandler(adminSvc *service.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc}
}

// SetSigner sets the signer address in the contract
func (h *AdminHandler) SetSigner(c *gin.Context) {
	var req struct {
		SignerAddress string `json:"signer_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.adminSvc.SetSigner(c.Request.Context(), req.SignerAddress); err != nil {
		resp.InternalError(c, "failed to set signer: "+err.Error())
		return
	}

	resp.OK(c, gin.H{"message": "signer address updated successfully"})
}

// SetToken configures allowed token
func (h *AdminHandler) SetToken(c *gin.Context) {
	var req struct {
		TokenAddress string `json:"token_address" binding:"required"`
		Allowed      bool   `json:"allowed"`
		MinAmount    string `json:"min_amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.adminSvc.SetToken(c.Request.Context(), req.TokenAddress, req.Allowed, req.MinAmount); err != nil {
		resp.InternalError(c, "failed to set token: "+err.Error())
		return
	}

	resp.OK(c, gin.H{"message": "token configuration updated"})
}

// SetExpiry sets default expiry duration
func (h *AdminHandler) SetExpiry(c *gin.Context) {
	var req struct {
		ExpirySeconds int64 `json:"expiry_seconds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.adminSvc.SetExpiry(c.Request.Context(), req.ExpirySeconds); err != nil {
		resp.InternalError(c, "failed to set expiry: "+err.Error())
		return
	}

	resp.OK(c, gin.H{"message": "expiry duration updated"})
}

// SetAllowAllTokens sets whether all tokens are allowed
func (h *AdminHandler) SetAllowAllTokens(c *gin.Context) {
	var req struct {
		AllowAll bool `json:"allow_all"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.adminSvc.SetAllowAllTokens(c.Request.Context(), req.AllowAll); err != nil {
		resp.InternalError(c, "failed to update allow all tokens: "+err.Error())
		return
	}

	resp.OK(c, gin.H{"message": "allow all tokens setting updated"})
}

// SetNativeTokenEnabled enables/disables native token (ETH/TRX)
func (h *AdminHandler) SetNativeTokenEnabled(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	if err := h.adminSvc.SetNativeTokenEnabled(c.Request.Context(), req.Enabled); err != nil {
		resp.InternalError(c, "failed to update native token setting: "+err.Error())
		return
	}

	resp.OK(c, gin.H{"message": "native token setting updated"})
}

// ParseTxEvents manually parses events from a transaction hash (for debugging)
func (h *AdminHandler) ParseTxEvents(c *gin.Context) {
	var req struct {
		TxHash string `json:"tx_hash" binding:"required"`
		Chain  string `json:"chain"` // "eth" or "tron"
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	result, err := h.adminSvc.ParseTxEvents(c.Request.Context(), req.TxHash, req.Chain)
	if err != nil {
		resp.InternalError(c, "failed to parse tx events: "+err.Error())
		return
	}

	resp.OK(c, result)
}
