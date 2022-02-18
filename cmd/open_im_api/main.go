package main

import (
	apiAuth "Open_IM/internal/api/auth"
	apiChat "Open_IM/internal/api/chat"
	"Open_IM/internal/api/conversation"
	"Open_IM/internal/api/friend"
	"Open_IM/internal/api/group"
	"Open_IM/internal/api/manage"
	apiThird "Open_IM/internal/api/third"
	"Open_IM/internal/api/user"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"flag"
	"strconv"

	"github.com/gin-gonic/gin"
	//"syscall"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(utils.CorsHandler())
	// user routing group, which handles user registration and login services
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/update_user_info", user.UpdateUserInfo)    //1
		userRouterGroup.POST("/get_users_info", user.GetUsersInfo)        //1
		userRouterGroup.POST("/get_self_user_info", user.GetSelfUserInfo) //1
	}
	//friend routing group
	friendRouterGroup := r.Group("/friend")
	{
		//	friendRouterGroup.POST("/get_friends_info", friend.GetFriendsInfo)
		friendRouterGroup.POST("/add_friend", friend.AddFriend)                              //1
		friendRouterGroup.POST("/delete_friend", friend.DeleteFriend)                        //1
		friendRouterGroup.POST("/get_friend_apply_list", friend.GetFriendApplyList)          //1
		friendRouterGroup.POST("/get_self_friend_apply_list", friend.GetSelfFriendApplyList) //1
		friendRouterGroup.POST("/get_friend_list", friend.GetFriendList)                     //1
		friendRouterGroup.POST("/add_friend_response", friend.AddFriendResponse)             //1
		friendRouterGroup.POST("/set_friend_remark", friend.SetFriendRemark)                 //1

		friendRouterGroup.POST("/add_black", friend.AddBlack)          //1
		friendRouterGroup.POST("/get_black_list", friend.GetBlacklist) //1
		friendRouterGroup.POST("/remove_black", friend.RemoveBlack)    //1

		friendRouterGroup.POST("/import_friend", friend.ImportFriend) //1
		friendRouterGroup.POST("/is_friend", friend.IsFriend)         //1
	}
	//group related routing group
	groupRouterGroup := r.Group("/group")
	{
		groupRouterGroup.POST("/create_group", group.CreateGroup)                                   //1
		groupRouterGroup.POST("/set_group_info", group.SetGroupInfo)                                //1
		groupRouterGroup.POST("join_group", group.JoinGroup)                                        //1
		groupRouterGroup.POST("/quit_group", group.QuitGroup)                                       //1
		groupRouterGroup.POST("/group_application_response", group.ApplicationGroupResponse)        //1
		groupRouterGroup.POST("/transfer_group", group.TransferGroupOwner)                          //1
		groupRouterGroup.POST("/get_recv_group_applicationList", group.GetRecvGroupApplicationList) //1
		groupRouterGroup.POST("/get_user_req_group_applicationList", group.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_groups_info", group.GetGroupsInfo)                              //1
		groupRouterGroup.POST("/kick_group", group.KickGroupMember)                                 //1
		groupRouterGroup.POST("/get_group_member_list", group.GetGroupMemberList)                   //no use
		groupRouterGroup.POST("/get_group_all_member_list", group.GetGroupAllMemberList)            //1
		groupRouterGroup.POST("/get_group_members_info", group.GetGroupMembersInfo)                 //1
		groupRouterGroup.POST("/invite_user_to_group", group.InviteUserToGroup)                     //1
		groupRouterGroup.POST("/get_joined_group_list", group.GetJoinedGroupList)                   //1
	}
	//certificate
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_register", apiAuth.UserRegister) //1
		authRouterGroup.POST("/user_token", apiAuth.UserToken)       //1
	}
	//Third service
	thirdGroup := r.Group("/third")
	{
		thirdGroup.POST("/tencent_cloud_storage_credential", apiThird.TencentCloudStorageCredential)
		thirdGroup.POST("/minio_storage_credential", apiThird.MinioStorageCredential)
	}
	//Message
	chatGroup := r.Group("/msg")
	{
		chatGroup.POST("/newest_seq", apiChat.GetSeq)
		chatGroup.POST("/send_msg", apiChat.SendMsg)
		chatGroup.POST("/pull_msg_by_seq", apiChat.PullMsgBySeqList)
	}
	//Manager
	managementGroup := r.Group("/manager")
	{
		managementGroup.POST("/delete_user", manage.DeleteUser) //1
		managementGroup.POST("/send_msg", manage.ManagementSendMsg)
		managementGroup.POST("/get_all_users_uid", manage.GetAllUsersUid)             //1
		managementGroup.POST("/account_check", manage.AccountCheck)                   //1
		managementGroup.POST("/get_users_online_status", manage.GetUsersOnlineStatus) //1
	}
	//Conversation
	conversationGroup := r.Group("/conversation")
	{
		conversationGroup.POST("/set_receive_message_opt", conversation.SetReceiveMessageOpt)                  //1
		conversationGroup.POST("/get_receive_message_opt", conversation.GetReceiveMessageOpt)                  //1
		conversationGroup.POST("/get_all_conversation_message_opt", conversation.GetAllConversationMessageOpt) //1
	}

	log.NewPrivateLog("api")
	ginPort := flag.Int("port", 10000, "get ginServerPort from cmd,default 10000 as port")
	flag.Parse()
	r.Run(":" + strconv.Itoa(*ginPort))
}
