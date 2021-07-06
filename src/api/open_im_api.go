package main

import (
	apiAuth "Open_IM/src/api/auth"
	apiChat "Open_IM/src/api/chat"
	"Open_IM/src/api/friend"
	"Open_IM/src/api/group"
	apiThird "Open_IM/src/api/third"
	"Open_IM/src/api/user"
	"Open_IM/src/common/log"
	"Open_IM/src/utils"
	"flag"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"syscall"
)

func main() {

	logFile, err := os.OpenFile("./fatal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {

		return
	}
	syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
	//gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = ioutil.Discard
	log.Info("", "", "api server running...")
	r := gin.Default()
	r.Use(utils.CorsHandler())
	// user routing group, which handles user registration and login services
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/update_user_info", user.UpdateUserInfo)
		userRouterGroup.POST("/get_user_info", user.GetUserInfo)
	}
	//friend routing group
	friendRouterGroup := r.Group("/friend")
	{
		friendRouterGroup.POST("/get_friends_info", friend.GetFriendsInfo)
		friendRouterGroup.POST("/add_friend", friend.AddFriend)
		friendRouterGroup.POST("/get_friend_apply_list", friend.GetFriendApplyList)
		friendRouterGroup.POST("/get_self_apply_list", friend.GetSelfApplyList)
		friendRouterGroup.POST("/get_friend_list", friend.GetFriendList)
		friendRouterGroup.POST("/add_blacklist", friend.AddBlacklist)
		friendRouterGroup.POST("/get_blacklist", friend.GetBlacklist)
		friendRouterGroup.POST("/remove_blacklist", friend.RemoveBlacklist)
		friendRouterGroup.POST("/delete_friend", friend.DeleteFriend)
		friendRouterGroup.POST("/add_friend_response", friend.AddFriendResponse)
		friendRouterGroup.POST("/set_friend_comment", friend.SetFriendComment)
		friendRouterGroup.POST("/is_friend", friend.IsFriend)
	}
	//group related routing group
	groupRouterGroup := r.Group("/group")
	{
		groupRouterGroup.POST("/create_group", group.CreateGroup)
		groupRouterGroup.POST("/set_group_info", group.SetGroupInfo)
		groupRouterGroup.POST("join_group", group.JoinGroup)
		groupRouterGroup.POST("/quit_group", group.QuitGroup)
		groupRouterGroup.POST("/group_application_response", group.ApplicationGroupResponse)
		groupRouterGroup.POST("/transfer_group", group.TransferGroupOwner)
		groupRouterGroup.POST("/get_group_applicationList", group.GetGroupApplicationList)
		groupRouterGroup.POST("/get_groups_info", group.GetGroupsInfo)
		groupRouterGroup.POST("/kick_group", group.KickGroupMember)
		groupRouterGroup.POST("/get_group_member_list", group.GetGroupMemberList)
		groupRouterGroup.POST("/get_group_all_member_list", group.GetGroupAllMember)
		groupRouterGroup.POST("/get_group_members_info", group.GetGroupMembersInfo)
		groupRouterGroup.POST("/invite_user_to_group", group.InviteUserToGroup)
		groupRouterGroup.POST("/get_joined_group_list", group.GetJoinedGroupList)
	}
	//certificate
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_register", apiAuth.UserRegister)
		authRouterGroup.POST("/user_token", apiAuth.UserToken)
	}
	//Third service
	thirdGroup := r.Group("/third")
	{
		thirdGroup.POST("/tencent_cloud_storage_credential", apiThird.TencentCloudStorageCredential)
	}
	//Message
	chatGroup := r.Group("/chat")
	{
		chatGroup.POST("/newest_seq", apiChat.UserNewestSeq)
		chatGroup.POST("/pull_msg", apiChat.UserPullMsg)
		chatGroup.POST("/send_msg", apiChat.UserSendMsg)
	}

	ginPort := flag.Int("port", 10000, "get ginServerPort from cmd,default 10000 as port")
	flag.Parse()
	r.Run(utils.ServerIP + ":" + strconv.Itoa(*ginPort))
}
