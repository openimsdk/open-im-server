package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// register
func AdminLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "someJSON", "status": 200})
}

func AdminRegister(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "someJSON", "status": 200})
}

func GetAdminSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "someJSON", "status": 200})
}

func AlterAdminSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "someJSON", "status": 200})
}
