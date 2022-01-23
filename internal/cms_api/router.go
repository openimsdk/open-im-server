package cms_api

import (
	"Open_IM/internal/cms_api/admin"
	"Open_IM/internal/cms_api/group"
	"Open_IM/internal/cms_api/message"
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
	router.Use(middleware.JWTAuth())
	router.Use(middleware.CorsHandler())
	adminRouterGroup := router.Group("/admin")
	{
		adminRouterGroup.POST("/register", admin.UserRegister)
		adminRouterGroup.POST("/login", admin.UserLogin)
		adminRouterGroup.GET("/get_user_settings", admin.GetUserSettings)
		adminRouterGroup.POST("/alter_user_settings", admin.AlterUserSettings)
	}
	statisticsRouterGroup := router.Group("/statistics")
	{
		statisticsRouterGroup.GET("/get_messages_statistics", statistics.MessagesStatistics)
		statisticsRouterGroup.GET("/get_users_statistics", statistics.UsersStatistics)
		statisticsRouterGroup.GET("/get_groups_statistics", statistics.GroupsStatistics)
		statisticsRouterGroup.GET("/get_active_user", statistics.GetActiveUser)
		statisticsRouterGroup.GET("/get_active_group", statistics.GetActiveGroup)
	}
	organizationRouterGroup := router.Group("/organization")
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
	messageRouterGroup := router.Group("/message")
	{
		messageRouterGroup.POST("/broadcast", message.Broadcast)
		messageRouterGroup.GET("/search_message_by_user", message.SearchMessageByUser)
		messageRouterGroup.POST("/message_mass_send", message.MassSendMassage)
		messageRouterGroup.GET("/search_message_by_group", message.SearchMessageByGroup)
		messageRouterGroup.POST("/withdraw_message", message.Withdraw)
	}
	groupRouterGroup := router.Group("/groups")
	{
		groupRouterGroup.GET("/search_groups", group.SearchGroups)
		groupRouterGroup.GET("/search_groups_member", group.SearchGroupsMember)
		groupRouterGroup.POST("/create_group", group.CreateGroup)
		groupRouterGroup.GET("/inquire_group", group.InquireGroup)
		groupRouterGroup.GET("/inquireMember_by_group", group.InquireMember)
		groupRouterGroup.POST("/add_members", group.AddMembers)
		groupRouterGroup.POST("/set_master", group.SetMaster)
		groupRouterGroup.POST("/block_user", group.BlockUser)
		groupRouterGroup.POST("/remove_user", group.RemoveUser)
		groupRouterGroup.POST("/ban_private_chat", group.BanPrivateChat)
		groupRouterGroup.POST("/withdraw_message", group.Withdraw)
		groupRouterGroup.POST("/search_group_message", group.SearchMessage)
	}
	userRouterGroup := router.Group("/user")
	{
		userRouterGroup.POST("/resign", user.ResignUser)
		userRouterGroup.GET("/get_user", user.GetUser)
		userRouterGroup.POST("/alter_user", user.AlterUser)
		userRouterGroup.GET("/get_users", user.GetUsers)
		userRouterGroup.POST("/add_user", user.AddUser)
		userRouterGroup.POST("/unblock_user", user.UnblockUser)
		userRouterGroup.POST("/block_user", user.BlockUser)
		userRouterGroup.GET("/block_users", user.GetBlockUsers)
	}
	return baseRouter
}
