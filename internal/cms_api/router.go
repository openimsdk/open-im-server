package cms_api

import (
	"Open_IM/internal/cms_api/admin"
	"Open_IM/internal/cms_api/group"
	messageCMS "Open_IM/internal/cms_api/message_cms"
	"Open_IM/internal/cms_api/middleware"
	"Open_IM/internal/cms_api/statistics"
	"Open_IM/internal/cms_api/user"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	baseRouter := gin.Default()
	router := baseRouter.Group("/cms")
	router.Use(middleware.CorsHandler())
	adminRouterGroup := router.Group("/admin")
	{
		adminRouterGroup.POST("/login", admin.AdminLogin)
		adminRouterGroup.Use(middleware.JWTAuth())
		adminRouterGroup.POST("/add_user_register_add_friend_id", admin.AddUserRegisterAddFriendIDList)
		adminRouterGroup.POST("/reduce_user_register_reduce_friend_id", admin.ReduceUserRegisterAddFriendIDList)
		adminRouterGroup.POST("/get_user_register_reduce_friend_id_list", admin.GetUserRegisterAddFriendIDList)
	}
	r2 := router.Group("")
	r2.Use(middleware.JWTAuth())
	statisticsRouterGroup := r2.Group("/statistics")
	{
		statisticsRouterGroup.POST("/get_messages_statistics", statistics.GetMessagesStatistics)
		statisticsRouterGroup.POST("/get_user_statistics", statistics.GetUserStatistics)
		statisticsRouterGroup.POST("/get_group_statistics", statistics.GetGroupStatistics)
		statisticsRouterGroup.POST("/get_active_user", statistics.GetActiveUser)
		statisticsRouterGroup.POST("/get_active_group", statistics.GetActiveGroup)
	}
	groupRouterGroup := r2.Group("/group")
	{
		groupRouterGroup.POST("/get_groups", group.GetGroups)
		groupRouterGroup.POST("/get_group_members", group.GetGroupMembers)
	}
	userRouterGroup := r2.Group("/user")
	{
		userRouterGroup.POST("/add_user", user.AddUser)
		userRouterGroup.POST("/unblock_user", user.UnblockUser)
		userRouterGroup.POST("/block_user", user.BlockUser)
		userRouterGroup.POST("/get_block_users", user.GetBlockUsers)
	}
	messageCMSRouterGroup := r2.Group("/message")
	{
		messageCMSRouterGroup.POST("/get_chat_logs", messageCMS.GetChatLogs)
	}
	return baseRouter
}
