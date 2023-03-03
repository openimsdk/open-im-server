package api

import (
	"OpenIM/internal/api/conversation"
	"OpenIM/internal/api/manage"
	"OpenIM/internal/api/msg"
	"OpenIM/internal/api/third"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/middleware"
	"OpenIM/pkg/common/prome"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func NewGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("../logs/api.log")
	gin.DefaultWriter = io.MultiWriter(f)
	//	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.GinParseOperationID)
	log.Info("load config: ", config.Config)
	if config.Config.Prometheus.Enable {
		prome.NewApiRequestCounter()
		prome.NewApiRequestFailedCounter()
		prome.NewApiRequestSuccessCounter()
		r.Use(prome.PrometheusMiddleware)
		r.GET("/metrics", prome.PrometheusHandler())
	}

	userRouterGroup := r.Group("/user")
	{
		u := NewUser(nil)
		userRouterGroup.POST("/user_register", u.UserRegister)
		userRouterGroup.POST("/update_user_info", u.UpdateUserInfo) //1
		userRouterGroup.POST("/set_global_msg_recv_opt", u.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_users_info", u.GetUsersPublicInfo)            //1
		userRouterGroup.POST("/get_self_user_info", u.GetSelfUserInfo)           //1
		userRouterGroup.POST("/get_users_online_status", u.GetUsersOnlineStatus) //1
		userRouterGroup.POST("/get_users_info_from_cache", u.GetUsersInfoFromCache)
		userRouterGroup.POST("/get_user_friend_from_cache", u.GetFriendIDListFromCache)
		userRouterGroup.POST("/get_black_list_from_cache", u.GetBlackIDListFromCache)
		//userRouterGroup.POST("/get_all_users_uid", manage.GetAllUsersUid) // todo
		//userRouterGroup.POST("/account_check", manage.AccountCheck)       // todo
		userRouterGroup.POST("/get_users", u.GetUsers)
	}
	////friend routing group
	friendRouterGroup := r.Group("/friend")
	{
		f := NewFriend(nil)
		friendRouterGroup.POST("/add_friend", f.AddFriend)                        //1
		friendRouterGroup.POST("/delete_friend", f.DeleteFriend)                  //1
		friendRouterGroup.POST("/get_friend_apply_list", f.GetFriendApplyList)    //1
		friendRouterGroup.POST("/get_self_friend_apply_list", f.GetSelfApplyList) //1
		friendRouterGroup.POST("/get_friend_list", f.GetFriendList)               //1
		friendRouterGroup.POST("/add_friend_response", f.AddFriendResponse)       //1
		friendRouterGroup.POST("/set_friend_remark", f.SetFriendRemark)           //1
		friendRouterGroup.POST("/add_black", f.AddBlack)                          //1
		friendRouterGroup.POST("/get_black_list", f.GetBlacklist)                 //1
		friendRouterGroup.POST("/remove_black", f.RemoveBlacklist)                //1
		friendRouterGroup.POST("/import_friend", f.ImportFriend)                  //1
		friendRouterGroup.POST("/is_friend", f.IsFriend)                          //1

	}
	groupRouterGroup := r.Group("/group")
	g := NewGroup(nil)
	{
		groupRouterGroup.POST("/create_group", g.NewCreateGroup)                                //1
		groupRouterGroup.POST("/set_group_info", g.NewSetGroupInfo)                             //1
		groupRouterGroup.POST("/join_group", g.JoinGroup)                                       //1
		groupRouterGroup.POST("/quit_group", g.QuitGroup)                                       //1
		groupRouterGroup.POST("/group_application_response", g.ApplicationGroupResponse)        //1
		groupRouterGroup.POST("/transfer_group", g.TransferGroupOwner)                          //1
		groupRouterGroup.POST("/get_recv_group_applicationList", g.GetRecvGroupApplicationList) //1
		groupRouterGroup.POST("/get_user_req_group_applicationList", g.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_groups_info", g.GetGroupsInfo) //1
		groupRouterGroup.POST("/kick_group", g.KickGroupMember)    //1
		//groupRouterGroup.POST("/get_group_all_member_list", g.GetGroupAllMemberList) //1
		groupRouterGroup.POST("/get_group_members_info", g.GetGroupMembersInfo) //1
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
		superGroupRouterGroup.POST("/get_joined_group_list", g.GetJoinedSuperGroupList)
		superGroupRouterGroup.POST("/get_groups_info", g.GetSuperGroupsInfo)
	}
	////certificate
	authRouterGroup := r.Group("/auth")
	{
		a := NewAuth(nil)
		u := NewUser(nil)
		authRouterGroup.POST("/user_register", u.UserRegister) //1
		authRouterGroup.POST("/user_token", a.UserToken)       //1
		authRouterGroup.POST("/parse_token", a.ParseToken)     //1
		authRouterGroup.POST("/force_logout", a.ForceLogout)   //1
	}
	////Third service
	thirdGroup := r.Group("/third")
	{
		thirdGroup.POST("/tencent_cloud_storage_credential", third.TencentCloudStorageCredential)
		thirdGroup.POST("/ali_oss_credential", third.AliOSSCredential)
		thirdGroup.POST("/minio_storage_credential", third.MinioStorageCredential)
		thirdGroup.POST("/minio_upload", third.MinioUploadFile)
		thirdGroup.POST("/upload_update_app", third.UploadUpdateApp)
		thirdGroup.POST("/get_download_url", third.GetDownloadURL)
		thirdGroup.POST("/get_rtc_invitation_info", third.GetRTCInvitationInfo)
		thirdGroup.POST("/get_rtc_invitation_start_app", third.GetRTCInvitationInfoStartApp)
		thirdGroup.POST("/fcm_update_token", third.FcmUpdateToken)
		thirdGroup.POST("/aws_storage_credential", third.AwsStorageCredential)
		thirdGroup.POST("/set_app_badge", third.SetAppBadge)
	}
	////Message
	chatGroup := r.Group("/msg")
	{
		chatGroup.POST("/newest_seq", msg.GetSeq)
		chatGroup.POST("/send_msg", msg.SendMsg)
		chatGroup.POST("/pull_msg_by_seq", msg.PullMsgBySeqList)
		chatGroup.POST("/del_msg", msg.DelMsg)
		chatGroup.POST("/del_super_group_msg", msg.DelSuperGroupMsg)
		chatGroup.POST("/clear_msg", msg.ClearMsg)
		chatGroup.POST("/manage_send_msg", manage.ManagementSendMsg)
		chatGroup.POST("/batch_send_msg", manage.ManagementBatchSendMsg)
		chatGroup.POST("/check_msg_is_send_success", manage.CheckMsgIsSendSuccess)
		chatGroup.POST("/set_msg_min_seq", msg.SetMsgMinSeq)

		chatGroup.POST("/set_message_reaction_extensions", msg.SetMessageReactionExtensions)
		chatGroup.POST("/get_message_list_reaction_extensions", msg.GetMessageListReactionExtensions)
		chatGroup.POST("/add_message_reaction_extensions", msg.AddMessageReactionExtensions)
		chatGroup.POST("/delete_message_reaction_extensions", msg.DeleteMessageReactionExtensions)
	}
	////Conversation
	conversationGroup := r.Group("/conversation")
	{ //1
		conversationGroup.POST("/get_all_conversations", conversation.GetAllConversations)
		conversationGroup.POST("/get_conversation", conversation.GetConversation)
		conversationGroup.POST("/get_conversations", conversation.GetConversations)
		conversationGroup.POST("/set_conversation", conversation.SetConversation)
		conversationGroup.POST("/batch_set_conversation", conversation.BatchSetConversations)
		conversationGroup.POST("/set_recv_msg_opt", conversation.SetRecvMsgOpt)
		conversationGroup.POST("/modify_conversation_field", conversation.ModifyConversationField)
	}
	return r
}
