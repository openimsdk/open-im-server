package main

import (
	apiAuth "Open_IM/src/api/auth"
	apiChat "Open_IM/src/api/chat"
	"Open_IM/src/api/friend"
	apiThird "Open_IM/src/api/third"
	"Open_IM/src/api/user"
	"Open_IM/src/common/log"
	"Open_IM/src/utils"
	"flag"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
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
		friendRouterGroup.POST("/search_friend", friend.SearchFriend)
		friendRouterGroup.POST("/add_friend", friend.AddFriend)
		friendRouterGroup.POST("/get_friend_apply_list", friend.GetFriendApplyList)
		friendRouterGroup.POST("/get_friend_list", friend.GetFriendList)
		friendRouterGroup.POST("/add_blacklist", friend.AddBlacklist)
		friendRouterGroup.POST("/get_blacklist", friend.GetBlacklist)
		friendRouterGroup.POST("/remove_blacklist", friend.RemoveBlacklist)
		friendRouterGroup.POST("/delete_friend", friend.DeleteFriend)
		friendRouterGroup.POST("/add_friend_response", friend.AddFriendResponse)
		friendRouterGroup.POST("/set_friend_comment", friend.SetFriendComment)
	}
	//group related routing group
	/*groupRouterGroup := r.Group("/group")
	{
		groupRouterGroup.POST("/create_group", group.CreateGroup)
		groupRouterGroup.POST("/get_group_list", group.GetGroupList)
		groupRouterGroup.POST("/get_group_info", group.GetGroupInfo)
		groupRouterGroup.POST("/delete_group_member", group.DeleteGroupMember)
		groupRouterGroup.POST("/set_group_name", group.SetGroupName)
		groupRouterGroup.POST("/set_group_bulletin", group.SetGroupBulletin)
		groupRouterGroup.POST("/set_owner_group_nickname", group.SetOwnerGroupNickname)
		groupRouterGroup.POST("/set_group_head_image", group.SetGroupHeadImage)
		groupRouterGroup.POST("/member_exit_group", group.MemberExitGroup)
	}*/
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
