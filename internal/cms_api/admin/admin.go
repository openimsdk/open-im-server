package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// register
func UserLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "someJSON", "status": 200})
}

func UserRegister(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "someJSON", "status": 200})
}

func GetUserSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "someJSON", "status": 200})
}

func AlterUserSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "someJSON", "status": 200})
}
