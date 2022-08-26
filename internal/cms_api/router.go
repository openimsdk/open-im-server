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
		statisticsRouterGroup.GET("/get_messages_statistics", statistics.GetMessagesStatistics)
		statisticsRouterGroup.GET("/get_user_statistics", statistics.GetUserStatistics)
		statisticsRouterGroup.GET("/get_group_statistics", statistics.GetGroupStatistics)
		statisticsRouterGroup.GET("/get_active_user", statistics.GetActiveUser)
		statisticsRouterGroup.GET("/get_active_group", statistics.GetActiveGroup)
	}
	groupRouterGroup := r2.Group("/group")
	{
		groupRouterGroup.GET("/get_group_by_id", group.GetGroupByID)
		groupRouterGroup.GET("/get_groups", group.GetGroups)
		groupRouterGroup.GET("/get_group_by_name", group.GetGroupByName)
		groupRouterGroup.GET("/get_group_members", group.GetGroupMembers)
		groupRouterGroup.POST("/create_group", group.CreateGroup)
		groupRouterGroup.POST("/add_members", group.AddGroupMembers)
		groupRouterGroup.POST("/remove_members", group.RemoveGroupMembers)
		groupRouterGroup.POST("/get_members_in_group", group.GetGroupMembers)
		groupRouterGroup.POST("/set_group_master", group.SetGroupOwner)
		groupRouterGroup.POST("/set_group_ordinary_user", group.SetGroupOrdinaryUsers)
		groupRouterGroup.POST("/alter_group_info", group.AlterGroupInfo)
	}
	userRouterGroup := r2.Group("/user")
	{
		userRouterGroup.POST("/add_user", user.AddUser)
		userRouterGroup.POST("/unblock_user", user.UnblockUser)
		userRouterGroup.POST("/block_user", user.BlockUser)
		userRouterGroup.GET("/get_block_users", user.GetBlockUsers)
		userRouterGroup.GET("/get_block_user", user.GetBlockUserById)
	}
	messageCMSRouterGroup := r2.Group("/message")
	{
		messageCMSRouterGroup.GET("/get_chat_logs", messageCMS.GetChatLogs)
	}
	return baseRouter
}
