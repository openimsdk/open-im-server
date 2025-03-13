package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/openimsdk/open-im-server/v3/internal/api/jssdk"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	pbAuth "github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"
)

const (
	NoCompression      = -1
	DefaultCompression = 0
	BestCompression    = 1
	BestSpeed          = 2
)

func prommetricsGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		path := c.FullPath()
		if c.Writer.Status() == http.StatusNotFound {
			prommetrics.HttpCall("<404>", c.Request.Method, c.Writer.Status())
		} else {
			prommetrics.HttpCall(path, c.Request.Method, c.Writer.Status())
		}
		if resp := apiresp.GetGinApiResponse(c); resp != nil {
			prommetrics.APICall(path, c.Request.Method, resp.ErrCode)
		}
	}
}

func newGinRouter(ctx context.Context, client discovery.Conn, cfg *Config) (*gin.Engine, error) {
	authConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Auth)
	if err != nil {
		return nil, err
	}
	userConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.User)
	if err != nil {
		return nil, err
	}
	groupConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Group)
	if err != nil {
		return nil, err
	}
	friendConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Friend)
	if err != nil {
		return nil, err
	}
	conversationConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Conversation)
	if err != nil {
		return nil, err
	}
	thirdConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Third)
	if err != nil {
		return nil, err
	}
	msgConn, err := client.GetConn(ctx, cfg.Discovery.RpcService.Msg)
	if err != nil {
		return nil, err
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("required_if", RequiredIf)
	}
	switch cfg.API.Api.CompressionLevel {
	case NoCompression:
	case DefaultCompression:
		r.Use(gzip.Gzip(gzip.DefaultCompression))
	case BestCompression:
		r.Use(gzip.Gzip(gzip.BestCompression))
	case BestSpeed:
		r.Use(gzip.Gzip(gzip.BestSpeed))
	}
	r.Use(prommetricsGin(), gin.RecoveryWithWriter(gin.DefaultErrorWriter, mw.GinPanicErr), mw.CorsHandler(),
		mw.GinParseOperationID(), GinParseToken(rpcli.NewAuthClient(authConn)))

	u := NewUserApi(user.NewUserClient(userConn), client, cfg.Discovery.RpcService)
	{
		userRouterGroup := r.Group("/user")
		userRouterGroup.POST("/user_register", u.UserRegister)
		userRouterGroup.POST("/update_user_info", u.UpdateUserInfo)
		userRouterGroup.POST("/update_user_info_ex", u.UpdateUserInfoEx)
		userRouterGroup.POST("/set_global_msg_recv_opt", u.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_users_info", u.GetUsersPublicInfo)
		userRouterGroup.POST("/get_all_users_uid", u.GetAllUsersID)
		userRouterGroup.POST("/account_check", u.AccountCheck)
		userRouterGroup.POST("/get_users", u.GetUsers)
		userRouterGroup.POST("/get_users_online_status", u.GetUsersOnlineStatus)
		userRouterGroup.POST("/get_users_online_token_detail", u.GetUsersOnlineTokenDetail)
		userRouterGroup.POST("/subscribe_users_status", u.SubscriberStatus)
		userRouterGroup.POST("/get_users_status", u.GetUserStatus)
		userRouterGroup.POST("/get_subscribe_users_status", u.GetSubscribeUsersStatus)

		userRouterGroup.POST("/process_user_command_add", u.ProcessUserCommandAdd)
		userRouterGroup.POST("/process_user_command_delete", u.ProcessUserCommandDelete)
		userRouterGroup.POST("/process_user_command_update", u.ProcessUserCommandUpdate)
		userRouterGroup.POST("/process_user_command_get", u.ProcessUserCommandGet)
		userRouterGroup.POST("/process_user_command_get_all", u.ProcessUserCommandGetAll)

		userRouterGroup.POST("/add_notification_account", u.AddNotificationAccount)
		userRouterGroup.POST("/update_notification_account", u.UpdateNotificationAccountInfo)
		userRouterGroup.POST("/search_notification_account", u.SearchNotificationAccount)
	}
	// friend routing group
	{
		f := NewFriendApi(relation.NewFriendClient(friendConn))
		friendRouterGroup := r.Group("/friend")
		friendRouterGroup.POST("/delete_friend", f.DeleteFriend)
		friendRouterGroup.POST("/get_friend_apply_list", f.GetFriendApplyList)
		friendRouterGroup.POST("/get_designated_friend_apply", f.GetDesignatedFriendsApply)
		friendRouterGroup.POST("/get_self_friend_apply_list", f.GetSelfApplyList)
		friendRouterGroup.POST("/get_friend_list", f.GetFriendList)
		friendRouterGroup.POST("/get_designated_friends", f.GetDesignatedFriends)
		friendRouterGroup.POST("/add_friend", f.ApplyToAddFriend)
		friendRouterGroup.POST("/add_friend_response", f.RespondFriendApply)
		friendRouterGroup.POST("/set_friend_remark", f.SetFriendRemark)
		friendRouterGroup.POST("/add_black", f.AddBlack)
		friendRouterGroup.POST("/get_black_list", f.GetPaginationBlacks)
		friendRouterGroup.POST("/get_specified_blacks", f.GetSpecifiedBlacks)
		friendRouterGroup.POST("/remove_black", f.RemoveBlack)
		friendRouterGroup.POST("/get_incremental_blacks", f.GetIncrementalBlacks)
		friendRouterGroup.POST("/import_friend", f.ImportFriends)
		friendRouterGroup.POST("/is_friend", f.IsFriend)
		friendRouterGroup.POST("/get_friend_id", f.GetFriendIDs)
		friendRouterGroup.POST("/get_specified_friends_info", f.GetSpecifiedFriendsInfo)
		friendRouterGroup.POST("/update_friends", f.UpdateFriends)
		friendRouterGroup.POST("/get_incremental_friends", f.GetIncrementalFriends)
		friendRouterGroup.POST("/get_full_friend_user_ids", f.GetFullFriendUserIDs)
	}

	g := NewGroupApi(group.NewGroupClient(groupConn))
	{
		groupRouterGroup := r.Group("/group")
		groupRouterGroup.POST("/create_group", g.CreateGroup)
		groupRouterGroup.POST("/set_group_info", g.SetGroupInfo)
		groupRouterGroup.POST("/set_group_info_ex", g.SetGroupInfoEx)
		groupRouterGroup.POST("/join_group", g.JoinGroup)
		groupRouterGroup.POST("/quit_group", g.QuitGroup)
		groupRouterGroup.POST("/group_application_response", g.ApplicationGroupResponse)
		groupRouterGroup.POST("/transfer_group", g.TransferGroupOwner)
		groupRouterGroup.POST("/get_recv_group_applicationList", g.GetRecvGroupApplicationList)
		groupRouterGroup.POST("/get_user_req_group_applicationList", g.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_group_users_req_application_list", g.GetGroupUsersReqApplicationList)
		groupRouterGroup.POST("/get_specified_user_group_request_info", g.GetSpecifiedUserGroupRequestInfo)
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
		groupRouterGroup.POST("/get_groups", g.GetGroups)
		groupRouterGroup.POST("/get_group_member_user_id", g.GetGroupMemberUserIDs)
		groupRouterGroup.POST("/get_incremental_join_groups", g.GetIncrementalJoinGroup)
		groupRouterGroup.POST("/get_incremental_group_members", g.GetIncrementalGroupMember)
		groupRouterGroup.POST("/get_incremental_group_members_batch", g.GetIncrementalGroupMemberBatch)
		groupRouterGroup.POST("/get_full_group_member_user_ids", g.GetFullGroupMemberUserIDs)
		groupRouterGroup.POST("/get_full_join_group_ids", g.GetFullJoinGroupIDs)
	}
	// certificate
	{
		a := NewAuthApi(pbAuth.NewAuthClient(authConn))
		authRouterGroup := r.Group("/auth")
		authRouterGroup.POST("/get_admin_token", a.GetAdminToken)
		authRouterGroup.POST("/get_user_token", a.GetUserToken)
		authRouterGroup.POST("/parse_token", a.ParseToken)
		authRouterGroup.POST("/force_logout", a.ForceLogout)

	}
	// Third service
	{
		t := NewThirdApi(third.NewThirdClient(thirdConn), cfg.API.Prometheus.GrafanaURL)
		thirdGroup := r.Group("/third")
		thirdGroup.GET("/prometheus", t.GetPrometheus)
		thirdGroup.POST("/fcm_update_token", t.FcmUpdateToken)
		thirdGroup.POST("/set_app_badge", t.SetAppBadge)

		logs := thirdGroup.Group("/logs")
		logs.POST("/upload", t.UploadLogs)
		logs.POST("/delete", t.DeleteLogs)
		logs.POST("/search", t.SearchLogs)

		objectGroup := r.Group("/object")

		objectGroup.POST("/part_limit", t.PartLimit)
		objectGroup.POST("/part_size", t.PartSize)
		objectGroup.POST("/initiate_multipart_upload", t.InitiateMultipartUpload)
		objectGroup.POST("/auth_sign", t.AuthSign)
		objectGroup.POST("/complete_multipart_upload", t.CompleteMultipartUpload)
		objectGroup.POST("/access_url", t.AccessURL)
		objectGroup.POST("/initiate_form_data", t.InitiateFormData)
		objectGroup.POST("/complete_form_data", t.CompleteFormData)
		objectGroup.GET("/*name", t.ObjectRedirect)
	}
	// Message
	m := NewMessageApi(msg.NewMsgClient(msgConn), rpcli.NewUserClient(userConn), cfg.Share.IMAdminUserID)
	{
		msgGroup := r.Group("/msg")
		msgGroup.POST("/newest_seq", m.GetSeq)
		msgGroup.POST("/search_msg", m.SearchMsg)
		msgGroup.POST("/send_msg", m.SendMessage)
		msgGroup.POST("/send_business_notification", m.SendBusinessNotification)
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

		msgGroup.POST("/batch_send_msg", m.BatchSendMsg)
		msgGroup.POST("/check_msg_is_send_success", m.CheckMsgIsSendSuccess)
		msgGroup.POST("/get_server_time", m.GetServerTime)
		msgGroup.POST("/get_stream_msg", m.GetStreamMsg)
		msgGroup.POST("/append_stream_msg", m.AppendStreamMsg)
	}
	// Conversation
	{
		c := NewConversationApi(conversation.NewConversationClient(conversationConn))
		conversationGroup := r.Group("/conversation")
		conversationGroup.POST("/get_sorted_conversation_list", c.GetSortedConversationList)
		conversationGroup.POST("/get_all_conversations", c.GetAllConversations)
		conversationGroup.POST("/get_conversation", c.GetConversation)
		conversationGroup.POST("/get_conversations", c.GetConversations)
		conversationGroup.POST("/set_conversations", c.SetConversations)
		conversationGroup.POST("/get_conversation_offline_push_user_ids", c.GetConversationOfflinePushUserIDs)
		conversationGroup.POST("/get_full_conversation_ids", c.GetFullOwnerConversationIDs)
		conversationGroup.POST("/get_incremental_conversations", c.GetIncrementalConversation)
		conversationGroup.POST("/get_owner_conversation", c.GetOwnerConversation)
		conversationGroup.POST("/get_not_notify_conversation_ids", c.GetNotNotifyConversationIDs)
		conversationGroup.POST("/get_pinned_conversation_ids", c.GetPinnedConversationIDs)
		conversationGroup.POST("/update_conversations_by_user", c.UpdateConversationsByUser)
	}

	{
		statisticsGroup := r.Group("/statistics")
		statisticsGroup.POST("/user/register", u.UserRegisterCount)
		statisticsGroup.POST("/user/active", m.GetActiveUser)
		statisticsGroup.POST("/group/create", g.GroupCreateCount)
		statisticsGroup.POST("/group/active", m.GetActiveGroup)
	}

	{
		j := jssdk.NewJSSdkApi(rpcli.NewUserClient(userConn), rpcli.NewRelationClient(friendConn),
			rpcli.NewGroupClient(groupConn), rpcli.NewConversationClient(conversationConn), rpcli.NewMsgClient(msgConn))
		jssdk := r.Group("/jssdk")
		jssdk.POST("/get_conversations", j.GetConversations)
		jssdk.POST("/get_active_conversations", j.GetActiveConversations)
	}
	{
		pd := NewPrometheusDiscoveryApi(cfg, client)
		proDiscoveryGroup := r.Group("/prometheus_discovery")
		proDiscoveryGroup.GET("/api", pd.Api)
		proDiscoveryGroup.GET("/user", pd.User)
		proDiscoveryGroup.GET("/group", pd.Group)
		proDiscoveryGroup.GET("/msg", pd.Msg)
		proDiscoveryGroup.GET("/friend", pd.Friend)
		proDiscoveryGroup.GET("/conversation", pd.Conversation)
		proDiscoveryGroup.GET("/third", pd.Third)
		proDiscoveryGroup.GET("/auth", pd.Auth)
		proDiscoveryGroup.GET("/push", pd.Push)
		proDiscoveryGroup.GET("/msg_gateway", pd.MessageGateway)
		proDiscoveryGroup.GET("/msg_transfer", pd.MessageTransfer)
	}

	var etcdClient *clientv3.Client
	if cfg.Discovery.Enable == config.ETCD {
		etcdClient = client.(*etcd.SvcDiscoveryRegistryImpl).GetClient()
	}
	cm := NewConfigManager(cfg.Share.IMAdminUserID, &cfg.AllConfig, etcdClient, string(cfg.ConfigPath))
	{
		configGroup := r.Group("/config", cm.CheckAdmin)
		configGroup.POST("/get_config_list", cm.GetConfigList)
		configGroup.POST("/get_config", cm.GetConfig)
		configGroup.POST("/set_config", cm.SetConfig)
		configGroup.POST("/set_configs", cm.SetConfigs)
		configGroup.POST("/reset_config", cm.ResetConfig)
		configGroup.POST("/set_enable_config_manager", cm.SetEnableConfigManager)
		configGroup.POST("/get_enable_config_manager", cm.GetEnableConfigManager)
	}
	{
		r.POST("/restart", cm.CheckAdmin, cm.Restart)
	}
	return r, nil
}

func GinParseToken(authClient *rpcli.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost:
			for _, wApi := range Whitelist {
				if strings.HasPrefix(c.Request.URL.Path, wApi) {
					c.Next()
					return
				}
			}

			token := c.Request.Header.Get(constant.Token)
			if token == "" {
				log.ZWarn(c, "header get token error", servererrs.ErrArgs.WrapMsg("header must have token"))
				apiresp.GinError(c, servererrs.ErrArgs.WrapMsg("header must have token"))
				c.Abort()
				return
			}
			resp, err := authClient.ParseToken(c, token)
			if err != nil {
				apiresp.GinError(c, err)
				c.Abort()
				return
			}
			c.Set(constant.OpUserPlatform, constant.PlatformIDToName(int(resp.PlatformID)))
			c.Set(constant.OpUserID, resp.UserID)
			c.Next()
		}
	}
}

// Whitelist api not parse token
var Whitelist = []string{
	"/auth/get_admin_token",
	"/auth/parse_token",
}
