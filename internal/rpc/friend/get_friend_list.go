package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
)

func (s *friendServer) GetFriendList(ctx context.Context, req *pbFriend.GetFriendListReq) (*pbFriend.GetFriendListResp, error) {
	log.Info(req.Token, req.OperationID, "rpc get friend list is server,args=%s", req.String())
	var userInfoList []*pbFriend.UserInfo
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.GetFriendListResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	friends, err := im_mysql_model.FindUserInfoFromFriend(claims.UID)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s search friendInfo failed", err.Error())
		return &pbFriend.GetFriendListResp{ErrorCode: config.ErrSearchUserInfo.ErrCode, ErrorMsg: config.ErrSearchUserInfo.ErrMsg}, nil
	}
	for _, friendUser := range friends {
		var friendUserInfo pbFriend.UserInfo

		//find user is in blackList
		err = im_mysql_model.FindRelationshipFromBlackList(claims.UID, friendUser.FriendId)
		if err == nil {
			friendUserInfo.IsInBlackList = constant.BlackListFlag
		} else {
			friendUserInfo.IsInBlackList = 0
		}
		//Find user information
		us, err := im_mysql_model.FindUserByUID(friendUser.FriendId)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s search userInfo failed", err.Error())
			continue
		}
		friendUserInfo.Uid = friendUser.FriendId
		friendUserInfo.Comment = friendUser.Comment
		friendUserInfo.Icon = us.Icon
		friendUserInfo.Name = us.Name
		friendUserInfo.Gender = us.Gender
		friendUserInfo.Mobile = us.Mobile
		friendUserInfo.Birth = us.Birth
		friendUserInfo.Email = us.Email
		friendUserInfo.Ex = us.Ex

		userInfoList = append(userInfoList, &friendUserInfo)

	}
	log.Info(req.Token, req.OperationID, "rpc get friend list success return")
	return &pbFriend.GetFriendListResp{Data: userInfoList}, nil
}
