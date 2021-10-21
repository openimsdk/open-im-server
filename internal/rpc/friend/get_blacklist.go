package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
)

func (s *friendServer) GetBlacklist(ctx context.Context, req *pbFriend.GetBlacklistReq) (*pbFriend.GetBlacklistResp, error) {
	log.Info(req.Token, req.OperationID, "rpc get blacklist is server,args=%s", req.String())
	var (
		userInfoList []*pbFriend.UserInfo
		comment      string
	)
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.GetBlacklistResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	blackListInfo, err := im_mysql_model.GetBlackListByUID(claims.UID)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s get blacklist failed", err.Error())
		return &pbFriend.GetBlacklistResp{ErrorCode: config.ErrGetBlackList.ErrCode, ErrorMsg: config.ErrGetBlackList.ErrMsg}, nil
	}
	for _, blackUser := range blackListInfo {
		var blackUserInfo pbFriend.UserInfo
		//Find black user information
		us, err := im_mysql_model.FindUserByUID(blackUser.BlockId)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s search black list userInfo failed", err.Error())
			continue
		}
		friendShip, err := im_mysql_model.FindFriendRelationshipFromFriend(claims.UID, blackUser.BlockId)
		if err == nil {
			comment = friendShip.Comment
		}
		blackUserInfo.Uid = us.UID
		blackUserInfo.Icon = us.Icon
		blackUserInfo.Name = us.Name
		blackUserInfo.Gender = us.Gender
		blackUserInfo.Mobile = us.Mobile
		blackUserInfo.Birth = us.Birth
		blackUserInfo.Email = us.Email
		blackUserInfo.Ex = us.Ex
		blackUserInfo.Comment = comment

		userInfoList = append(userInfoList, &blackUserInfo)
	}
	log.Info(req.Token, req.OperationID, "rpc get blacklist success return")
	return &pbFriend.GetBlacklistResp{Data: userInfoList}, nil
}
