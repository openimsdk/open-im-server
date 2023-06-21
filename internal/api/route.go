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
	//f, _ := os.Create("../logs/api.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	//gin.SetMode(gin.DebugMode)
	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("required_if", RequiredIf)
	}
	log.ZInfo(context.Background(), "load config", "config", config.Config)
	r.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID())
	if config.Config.Prometheus.Enable {
		prome.NewApiRequestCounter()
		prome.NewApiRequestFailedCounter()
		prome.NewApiRequestSuccessCounter()
		r.Use(prome.PrometheusMiddleware)
		r.GET("/metrics", prome.PrometheusHandler())
	}
	userRouterGroup := r.Group("/user")
	{
		u := NewUserApi(discov)
		userRouterGroupChild := mw.NewRouterGroup(userRouterGroup, "")
		userRouterGroupChildToken := mw.NewRouterGroup(userRouterGroup, "", mw.WithGinParseToken(rdb))
		userRouterGroupChild.POST("/user_register", u.UserRegister)
		userRouterGroupChildToken.POST("/update_user_info", u.UpdateUserInfo) //1
		userRouterGroupChildToken.POST("/set_global_msg_recv_opt", u.SetGlobalRecvMessageOpt)
		userRouterGroupChildToken.POST("/get_users_info", u.GetUsersPublicInfo) //1
		userRouterGroupChildToken.POST("/get_all_users_uid", u.GetAllUsersID)   // todo
		userRouterGroupChildToken.POST("/account_check", u.AccountCheck)        // todo
		userRouterGroupChildToken.POST("/get_users", u.GetUsers)
		userRouterGroupChildToken.POST("/get_users_online_status", u.GetUsersOnlineStatus)
	}
	//friend routing group
	friendRouterGroup := r.Group("/friend")
	{
		f := NewFriendApi(discov)
		friendRouterGroup.Use(mw.GinParseToken(rdb))
		friendRouterGroup.POST("/delete_friend", f.DeleteFriend)                  //1
		friendRouterGroup.POST("/get_friend_apply_list", f.GetFriendApplyList)    //1
		friendRouterGroup.POST("/get_self_friend_apply_list", f.GetSelfApplyList) //1
		friendRouterGroup.POST("/get_friend_list", f.GetFriendList)               //1
		friendRouterGroup.POST("/add_friend", f.ApplyToAddFriend)                 //1
		friendRouterGroup.POST("/add_friend_response", f.RespondFriendApply)      //1
		friendRouterGroup.POST("/set_friend_remark", f.SetFriendRemark)           //1
		friendRouterGroup.POST("/add_black", f.AddBlack)                          //1
		friendRouterGroup.POST("/get_black_list", f.GetPaginationBlacks)          //1
		friendRouterGroup.POST("/remove_black", f.RemoveBlack)                    //1
		friendRouterGroup.POST("/import_friend", f.ImportFriends)                 //1
		friendRouterGroup.POST("/is_friend", f.IsFriend)                          //1
	}
	g := NewGroupApi(discov)
	groupRouterGroup := r.Group("/group")
	{

		groupRouterGroup.Use(mw.GinParseToken(rdb))
		groupRouterGroup.POST("/create_group", g.CreateGroup)                                   //1
		groupRouterGroup.POST("/set_group_info", g.SetGroupInfo)                                //1
		groupRouterGroup.POST("/join_group", g.JoinGroup)                                       //1
		groupRouterGroup.POST("/quit_group", g.QuitGroup)                                       //1
		groupRouterGroup.POST("/group_application_response", g.ApplicationGroupResponse)        //1
		groupRouterGroup.POST("/transfer_group", g.TransferGroupOwner)                          //1
		groupRouterGroup.POST("/get_recv_group_applicationList", g.GetRecvGroupApplicationList) //1
		groupRouterGroup.POST("/get_user_req_group_applicationList", g.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_groups_info", g.GetGroupsInfo) //1
		groupRouterGroup.POST("/kick_group", g.KickGroupMember)    //1
		// groupRouterGroup.POST("/get_group_all_member_list", g.GetGroupAllMemberList) //1
		groupRouterGroup.POST("/get_group_members_info", g.GetGroupMembersInfo) //1
		groupRouterGroup.POST("/get_group_member_list", g.GetGroupMemberList)   //1
		groupRouterGroup.POST("/invite_user_to_group", g.InviteUserToGroup)     //1
		groupRouterGroup.POST("/get_joined_group_list", g.GetJoinedGroupList)
		groupRouterGroup.POST("/dismiss_group", g.DismissGroup) //
		groupRouterGroup.POST("/mute_group_member", g.MuteGroupMember)
		groupRouterGroup.POST("/cancel_mute_group_member", g.CancelMuteGroupMember) //MuteGroup
		groupRouterGroup.POST("/mute_group", g.MuteGroup)
		groupRouterGroup.POST("/cancel_mute_group", g.CancelMuteGroup)
		//groupRouterGroup.POST("/set_group_member_nickname", g.SetGroupMemberNickname)
		groupRouterGroup.POST("/set_group_member_info", g.SetGroupMemberInfo)
		groupRouterGroup.POST("/get_group_abstract_info", g.GetGroupAbstractInfo)
	}
	superGroupRouterGroup := r.Group("/super_group")
	{
		superGroupRouterGroup.Use(mw.GinParseToken(rdb))
		superGroupRouterGroup.POST("/get_joined_group_list", g.GetJoinedSuperGroupList)
		superGroupRouterGroup.POST("/get_groups_info", g.GetSuperGroupsInfo)
	}
	////certificate
	authRouterGroup := r.Group("/auth")
	{
		a := NewAuthApi(discov)
		u := NewUserApi(discov)
		authRouterGroupChild := mw.NewRouterGroup(authRouterGroup, "")
		authRouterGroupChildToken := mw.NewRouterGroup(authRouterGroup, "", mw.WithGinParseToken(rdb))
		authRouterGroupChild.POST("/user_register", u.UserRegister)    //1
		authRouterGroupChild.POST("/user_token", a.UserToken)          //1
		authRouterGroupChildToken.POST("/parse_token", a.ParseToken)   //1
		authRouterGroupChildToken.POST("/force_logout", a.ForceLogout) //1
	}
	////Third service
	thirdGroup := r.Group("/third")
	{
		t := NewThirdApi(discov)
		thirdGroup.Use(mw.GinParseToken(rdb))
		thirdGroup.POST("/fcm_update_token", t.FcmUpdateToken)
		thirdGroup.POST("/set_app_badge", t.SetAppBadge)

		thirdGroup.POST("/apply_put", t.ApplyPut)
		thirdGroup.POST("/get_put", t.GetPut)
		thirdGroup.POST("/confirm_put", t.ConfirmPut)
		thirdGroup.POST("/get_hash", t.GetHash)
		thirdGroup.POST("/object", t.GetURL)
		thirdGroup.GET("/object", t.GetURL)
	}
	////Message
	msgGroup := r.Group("/msg")
	{
		m := NewMessageApi(discov)
		msgGroup.Use(mw.GinParseToken(rdb))
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

		//msgGroup.POST("/set_message_reaction_extensions", msg.SetMessageReactionExtensions)
		//msgGroup.POST("/get_message_list_reaction_extensions", msg.GetMessageListReactionExtensions)
		//msgGroup.POST("/add_message_reaction_extensions", msg.AddMessageReactionExtensions)
		//msgGroup.POST("/delete_message_reaction_extensions", msg.DeleteMessageReactionExtensions)
	}
	////Conversation
	conversationGroup := r.Group("/conversation")
	{
		c := NewConversationApi(discov)
		conversationGroup.Use(mw.GinParseToken(rdb))
		conversationGroup.POST("/get_all_conversations", c.GetAllConversations)
		conversationGroup.POST("/get_conversation", c.GetConversation)
		conversationGroup.POST("/get_conversations", c.GetConversations)
		conversationGroup.POST("/set_conversation", c.SetConversation)
		conversationGroup.POST("/batch_set_conversation", c.BatchSetConversations)
		conversationGroup.POST("/set_recv_msg_opt", c.SetRecvMsgOpt)
		conversationGroup.POST("/modify_conversation_field", c.ModifyConversationField)
		conversationGroup.POST("/set_conversations", c.SetConversations)
	}

	statisticsGroup := r.Group("/statistics")
	{
		s := NewStatisticsApi(discov)
		conversationGroup.Use(mw.GinParseToken(rdb))
		statisticsGroup.POST("/user_register", s.UserRegister)
	}
	return r
}
