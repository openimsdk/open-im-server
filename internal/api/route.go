package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGinRouter(discov discoveryregistry.SvcDiscoveryRegistry, rdb redis.UniversalClient) *gin.Engine {
	discov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("required_if", RequiredIf)
	}
	log.ZInfo(context.Background(), "load config", "config", config.Config)
	r.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID())
	u := NewUserApi(discov)
	if config.Config.Prometheus.Enable {
		prome.NewApiRequestCounter()
		prome.NewApiRequestFailedCounter()
		prome.NewApiRequestSuccessCounter()
		r.Use(prome.PrometheusMiddleware)
		r.GET("/metrics", prome.PrometheusHandler())
	}
	ParseToken := mw.GinParseToken(rdb)
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/user_register", u.UserRegister)
		userRouterGroup.POST("/update_user_info", ParseToken, u.UpdateUserInfo)
		userRouterGroup.POST("/set_global_msg_recv_opt", ParseToken, u.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_users_info", ParseToken, u.GetUsersPublicInfo)
		userRouterGroup.POST("/get_all_users_uid", ParseToken, u.GetAllUsersID)
		userRouterGroup.POST("/account_check", ParseToken, u.AccountCheck)
		userRouterGroup.POST("/get_users", ParseToken, u.GetUsers)
		userRouterGroup.POST("/get_users_online_status", ParseToken, u.GetUsersOnlineStatus)
	}
	//friend routing group
	friendRouterGroup := r.Group("/friend", ParseToken)
	{
		f := NewFriendApi(discov)
		friendRouterGroup.POST("/delete_friend", f.DeleteFriend)
		friendRouterGroup.POST("/get_friend_apply_list", f.GetFriendApplyList)
		friendRouterGroup.POST("/get_self_friend_apply_list", f.GetSelfApplyList)
		friendRouterGroup.POST("/get_friend_list", f.GetFriendList)
		friendRouterGroup.POST("/add_friend", f.ApplyToAddFriend)
		friendRouterGroup.POST("/add_friend_response", f.RespondFriendApply)
		friendRouterGroup.POST("/set_friend_remark", f.SetFriendRemark)
		friendRouterGroup.POST("/add_black", f.AddBlack)
		friendRouterGroup.POST("/get_black_list", f.GetPaginationBlacks)
		friendRouterGroup.POST("/remove_black", f.RemoveBlack)
		friendRouterGroup.POST("/import_friend", f.ImportFriends)
		friendRouterGroup.POST("/is_friend", f.IsFriend)
	}
	g := NewGroupApi(discov)
	groupRouterGroup := r.Group("/group", ParseToken)
	{
		groupRouterGroup.POST("/create_group", g.CreateGroup)
		groupRouterGroup.POST("/set_group_info", g.SetGroupInfo)
		groupRouterGroup.POST("/join_group", g.JoinGroup)
		groupRouterGroup.POST("/quit_group", g.QuitGroup)
		groupRouterGroup.POST("/group_application_response", g.ApplicationGroupResponse)
		groupRouterGroup.POST("/transfer_group", g.TransferGroupOwner)
		groupRouterGroup.POST("/get_recv_group_applicationList", g.GetRecvGroupApplicationList)
		groupRouterGroup.POST("/get_user_req_group_applicationList", g.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_groups_info", g.GetGroupsInfo)
		groupRouterGroup.POST("/kick_group", g.KickGroupMember)
		groupRouterGroup.POST("/get_group_members_info", g.GetGroupMembersInfo)
		groupRouterGroup.POST("/get_group_member_list", g.GetGroupMemberList)
		groupRouterGroup.POST("/invite_user_to_group", g.InviteUserToGroup)
		groupRouterGroup.POST("/get_joined_group_list", g.GetJoinedGroupList)
		groupRouterGroup.POST("/dismiss_group", g.DismissGroup) //
		groupRouterGroup.POST("/mute_group_member", g.MuteGroupMember)
		groupRouterGroup.POST("/cancel_mute_group_member", g.CancelMuteGroupMember)
		groupRouterGroup.POST("/mute_group", g.MuteGroup)
		groupRouterGroup.POST("/cancel_mute_group", g.CancelMuteGroup)
		groupRouterGroup.POST("/set_group_member_info", g.SetGroupMemberInfo)
		groupRouterGroup.POST("/get_group_abstract_info", g.GetGroupAbstractInfo)
	}
	superGroupRouterGroup := r.Group("/super_group", ParseToken)
	{
		superGroupRouterGroup.POST("/get_joined_group_list", g.GetJoinedSuperGroupList)
		superGroupRouterGroup.POST("/get_groups_info", g.GetSuperGroupsInfo)
	}
	//certificate
	authRouterGroup := r.Group("/auth")
	{
		a := NewAuthApi(discov)
		authRouterGroup.POST("/user_register", u.UserRegister)
		authRouterGroup.POST("/user_token", a.UserToken)
		authRouterGroup.POST("/parse_token", a.ParseToken)
		authRouterGroup.POST("/force_logout", ParseToken, a.ForceLogout)
	}
	//Third service
	thirdGroup := r.Group("/third", ParseToken)
	{
		t := NewThirdApi(discov)
		thirdGroup.POST("/fcm_update_token", t.FcmUpdateToken)
		thirdGroup.POST("/set_app_badge", t.SetAppBadge)

		thirdGroup.POST("/apply_put", t.ApplyPut)
		thirdGroup.POST("/get_put", t.GetPut)
		thirdGroup.POST("/confirm_put", t.ConfirmPut)
		thirdGroup.POST("/get_hash", t.GetHash)
		thirdGroup.POST("/object", t.GetURL)
		thirdGroup.GET("/object", t.GetURL)
	}
	//Message
	msgGroup := r.Group("/msg", ParseToken)
	{
		m := NewMessageApi(discov)
		msgGroup.POST("/newest_seq", m.GetSeq)
		msgGroup.POST("/send_msg", m.SendMessage)
		msgGroup.POST("/pull_msg_by_seq", m.PullMsgBySeqs)
		msgGroup.POST("/revoke_msg", m.RevokeMsg)
		msgGroup.POST("/mark_msgs_as_read", m.MarkMsgsAsRead)
		msgGroup.POST("/mark_conversation_as_read", m.MarkConversationAsRead)
		msgGroup.POST("/get_conversations_has_read_and_max_seq", m.GetConversationsHasReadAndMaxSeq)
		msgGroup.POST("/set_conversation_has_read_seq", m.SetConversationHasReadSeq)

		msgGroup.POST("/clear_conversation_msg", m.ClearConversationsMsg)
		msgGroup.POST("/user_clear_all_msg", m.UserClearAllMsg)
		msgGroup.POST("/delete_msgs", m.DeleteMsgs)
		msgGroup.POST("/delete_msg_phsical_by_seq", m.DeleteMsgPhysicalBySeq)
		msgGroup.POST("/delete_msg_physical", m.DeleteMsgPhysical)

		msgGroup.POST("/batch_send_msg", m.ManagementBatchSendMsg)
		msgGroup.POST("/check_msg_is_send_success", m.CheckMsgIsSendSuccess)
	}
	//Conversation
	conversationGroup := r.Group("/conversation", ParseToken)
	{
		c := NewConversationApi(discov)
		conversationGroup.POST("/get_all_conversations", c.GetAllConversations)
		conversationGroup.POST("/get_conversation", c.GetConversation)
		conversationGroup.POST("/get_conversations", c.GetConversations)
		conversationGroup.POST("/batch_set_conversation", c.BatchSetConversations)
		conversationGroup.POST("/set_recv_msg_opt", c.SetRecvMsgOpt)
		conversationGroup.POST("/modify_conversation_field", c.ModifyConversationField)
		conversationGroup.POST("/set_conversations", c.SetConversations)
	}

	statisticsGroup := r.Group("/statistics", ParseToken)
	{
		statisticsGroup.POST("/user_register", u.UserRegisterCount)
	}
	return r
}
