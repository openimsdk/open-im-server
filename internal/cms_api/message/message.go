package message

import (
	"net/http"

	"Open_IM_CMS/test"

	"github.com/gin-gonic/gin"
)

func Broadcast(c *gin.Context) {

}

func SearchMessageByUser(c *gin.Context) {
	fake := test.GetSearchUserMsgFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fake})
}

func SearchMessageByGroup(c *gin.Context) {
	fake := test.GetSearchGroupMsgFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fake})
}

func MassSendMassage(c *gin.Context) {

}

func Withdraw(c *gin.Context) {

}
