package statistics

import (
	"net/http"

	"Open_IM/pkg/req_resp"

	"Open_IM/test"

	"github.com/gin-gonic/gin"
)

func MessagesStatistics(c *gin.Context) {
	var (
		req req_resp.StatisticsRequest
		//resp req_resp.MessageStatisticsResponse
	)
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if _, err := test.RpcFake(); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	fakeData := test.GetUserStatisticsResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fakeData})
}

func UsersStatistics(c *gin.Context) {
	var (
		req req_resp.StatisticsRequest
		//resp req_resp.MessageStatisticsResponse
	)
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if _, err := test.RpcFake(); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	fakeData := test.GetUserStatisticsResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fakeData})
}

func GroupsStatistics(c *gin.Context) {
	var (
		req req_resp.StatisticsRequest
		//resp req_resp.MessageStatisticsResponse
	)
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if _, err := test.RpcFake(); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	fakeData := test.GetUserStatisticsResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fakeData})
}

func GetActiveUser(c *gin.Context) {
	if _, err := test.RpcFake(); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	fakeData := test.GetActiveUserResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fakeData})
}

func GetActiveGroup(c *gin.Context) {
	if _, err := test.RpcFake(); err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	fakeData := test.GetActiveGroupResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fakeData})
}
