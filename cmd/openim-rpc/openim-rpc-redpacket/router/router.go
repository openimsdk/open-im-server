package router

import (
	"redpacket/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, rpHandler *handler.RedPacketHandler, adminHandler *handler.AdminHandler) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// User-facing red packet APIs
	api := r.Group("/api/redpacket")
	{
		api.POST("/create-order", rpHandler.CreateOrder)
		api.POST("/created-callback", rpHandler.CreatedCallback)
		api.GET("/detail", rpHandler.Detail)
		api.POST("/claim-sign", rpHandler.ClaimSign)
		api.POST("/claim-result", rpHandler.ClaimResult)
	}

	// Admin APIs - should be protected with authentication in production
	admin := r.Group("/admin/redpacket")
	{
		admin.POST("/set-signer", adminHandler.SetSigner)
		admin.POST("/set-token", adminHandler.SetToken)
		admin.POST("/set-expiry", adminHandler.SetExpiry)
		admin.POST("/set-allow-all-tokens", adminHandler.SetAllowAllTokens)
		admin.POST("/set-native-token", adminHandler.SetNativeTokenEnabled)
		admin.POST("/parse-tx-events", adminHandler.ParseTxEvents)
	}
}
