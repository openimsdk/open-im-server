package cms_api

import (
	"Open_IM/internal/cms_api/admin"
	"Open_IM/internal/cms_api/group"
	messageCMS "Open_IM/internal/cms_api/message_cms"
	"Open_IM/internal/cms_api/middleware"
	"Open_IM/internal/cms_api/organization"
	"Open_IM/internal/cms_api/statistics"
	"Open_IM/internal/cms_api/user"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	baseRouter := gin.Default()
	router := baseRouter.Group("/api")
	router.Use(middleware.CorsHandler())
	adminRouterGroup := router.Group("/admin")
	{
		adminRouterGroup.POST("/login", admin.AdminLogin)
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
	organizationRouterGroup := r2.Group("/organization")
	{
		organizationRouterGroup.GET("/get_staffs", organization.GetStaffs)
		organizationRouterGroup.GET("/get_organizations", organization.GetOrganizations)
		organizationRouterGroup.GET("/get_squad", organization.GetSquads)
		organizationRouterGroup.POST("/add_organization", organization.AddOrganization)
		organizationRouterGroup.POST("/alter_staff", organization.AlterStaff)
		organizationRouterGroup.GET("/inquire_organization", organization.InquireOrganization)
		organizationRouterGroup.POST("/alter_organization", organization.AlterOrganization)
		organizationRouterGroup.POST("/delete_organization", organization.DeleteOrganization)
		organizationRouterGroup.POST("/get_organization_squad", organization.GetOrganizationSquads)
		organizationRouterGroup.PATCH("/alter_corps_info", organization.AlterStaffsInfo)
		organizationRouterGroup.POST("/add_child_org", organization.AddChildOrganization)
	}
	groupRouterGroup := r2.Group("/group")
	{
		groupRouterGroup.GET("/get_group_by_id", group.GetGroupById)
		groupRouterGroup.GET("/get_groups", group.GetGroups)
		groupRouterGroup.GET("/get_group_by_name", group.GetGroupByName)
		groupRouterGroup.GET("/get_group_members", group.GetGroupMembers)
		groupRouterGroup.POST("/create_group", group.CreateGroup)
		groupRouterGroup.POST("/add_members", group.AddGroupMembers)
		groupRouterGroup.POST("/remove_members", group.RemoveGroupMembers)
		groupRouterGroup.POST("/ban_group_private_chat", group.BanPrivateChat)
		groupRouterGroup.POST("/open_group_private_chat", group.OpenPrivateChat)
		groupRouterGroup.POST("/ban_group_chat", group.BanGroupChat)
		groupRouterGroup.POST("/open_group_chat", group.OpenGroupChat)
		groupRouterGroup.POST("/delete_group", group.DeleteGroup)
		groupRouterGroup.POST("/get_members_in_group", group.GetGroupMembers)
		groupRouterGroup.POST("/set_group_master", group.SetGroupMaster)
		groupRouterGroup.POST("/set_group_ordinary_user", group.SetGroupOrdinaryUsers)
		groupRouterGroup.POST("/alter_group_info", group.AlterGroupInfo)
	}
	userRouterGroup := r2.Group("/user")
	{
		userRouterGroup.POST("/resign", user.ResignUser)
		userRouterGroup.GET("/get_user", user.GetUserById)
		userRouterGroup.POST("/alter_user", user.AlterUser)
		userRouterGroup.GET("/get_users", user.GetUsers)
		userRouterGroup.POST("/add_user", user.AddUser)
		userRouterGroup.POST("/unblock_user", user.UnblockUser)
		userRouterGroup.POST("/block_user", user.BlockUser)
		userRouterGroup.GET("/get_block_users", user.GetBlockUsers)
		userRouterGroup.GET("/get_block_user", user.GetBlockUserById)
		userRouterGroup.POST("/delete_user", user.DeleteUser)
		userRouterGroup.GET("/get_users_by_name", user.GetUsersByName)
	}
	friendRouterGroup := r2.Group("/friend")
	{
		friendRouterGroup.POST("/get_friends_by_id")
		friendRouterGroup.POST("/set_friend")
		friendRouterGroup.POST("/remove_friend")
	}
	messageCMSRouterGroup := r2.Group("/message")
	{
		messageCMSRouterGroup.GET("/get_chat_logs", messageCMS.GetChatLogs)
		messageCMSRouterGroup.POST("/broadcast_message", messageCMS.BroadcastMessage)
		messageCMSRouterGroup.POST("/mass_send_message", messageCMS.MassSendMassage)
		messageCMSRouterGroup.POST("/withdraw_message", messageCMS.WithdrawMessage)
	}
	return baseRouter
}
