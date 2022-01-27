package statistics

import (
	"Open_IM/pkg/cms_api_struct"
	"github.com/gin-gonic/gin"
	statisticsPb "Open_IM/pkg/proto/statistics"
)

func GetMessagesStatistics(c *gin.Context) {
	var (
		req cms_api_struct.GetGroupMembersRequest
		resp cms_api_struct.GetGroupMembersResponse
		reqPb statisticsPb.GetMessageStatisticsReq
	)
}

func GetUsersStatistics(c *gin.Context) {

}

func GetGroupsStatistics(c *gin.Context) {

}

func GetActiveUser(c *gin.Context) {

}

func GetActiveGroup(c *gin.Context) {

}
