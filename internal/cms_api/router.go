package cms_api

import (
	"Open_IM/internal/cms_api/admin"
	"Open_IM/internal/cms_api/friend"
	"Open_IM/internal/cms_api/group"
	messageCMS "Open_IM/internal/cms_api/message_cms"
	"Open_IM/internal/cms_api/middleware"
	"Open_IM/internal/cms_api/statistics"
	"Open_IM/internal/cms_api/user"
	"Open_IM/internal/demo/register"
	"Open_IM/pkg/common/config"

	promePkg "Open_IM/pkg/common/prometheus"

	"github.com/gin-gonic/gin"
)

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	baseRouter := gin.New()
	baseRouter.Use(gin.Recovery())
	baseRouter.Use(middleware.CorsHandler())
	if config.Config.Prometheus.Enable {
		baseRouter.GET("/metrics", promePkg.PrometheusHandler())
	}
	router := baseRouter.Group("/cms")
	adminRouterGroup := router.Group("/admin")
	{
		adminRouterGroup.POST("/login", admin.AdminLogin)
		adminRouterGroup.Use(middleware.JWTAuth())
		adminRouterGroup.POST("/get_user_token", admin.GetUserToken)

		adminRouterGroup.POST("/add_user_register_add_friend_id", admin.AddUserRegisterAddFriendIDList)
		adminRouterGroup.POST("/reduce_user_register_reduce_friend_id", admin.ReduceUserRegisterAddFriendIDList)
		adminRouterGroup.POST("/get_user_register_reduce_friend_id_list", admin.GetUserRegisterAddFriendIDList)

		adminRouterGroup.POST("/generate_invitation_code", register.GenerateInvitationCode)
		adminRouterGroup.POST("/query_invitation_code", register.QueryInvitationCode)
		adminRouterGroup.POST("/get_invitation_codes", register.GetInvitationCodes)

		adminRouterGroup.POST("/query_user_ip_limit_login", register.QueryUserIDLimitLogin)
		adminRouterGroup.POST("/add_user_ip_limit_login", register.AddUserIPLimitLogin)
		adminRouterGroup.POST("/remove_user_ip_limit_login", register.RemoveUserIPLimitLogin)

		adminRouterGroup.POST("/query_ip_register", register.QueryIPRegister)
		adminRouterGroup.POST("/add_ip_limit", register.AddIPLimit)
		adminRouterGroup.POST("/remove_ip_Limit", register.RemoveIPLimit)
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
		userRouterGroup.POST("/get_user_id_by_email_phone", user.GetUserIDByEmailAndPhoneNumber)

		userRouterGroup.POST("/add_user", user.AddUser)
		userRouterGroup.POST("/unblock_user", user.UnblockUser)
		userRouterGroup.POST("/block_user", user.BlockUser)
		userRouterGroup.POST("/get_block_users", user.GetBlockUsers)
	}
	messageCMSRouterGroup := r2.Group("/message")
	{
		messageCMSRouterGroup.POST("/get_chat_logs", messageCMS.GetChatLogs)
	}
	friendCMSRouterGroup := r2.Group("/friend")
	{
		friendCMSRouterGroup.POST("/get_friends", friend.GetUserFriends)
	}

	return baseRouter
}
