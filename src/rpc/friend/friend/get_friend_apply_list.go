package friend

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbFriend "Open_IM/src/proto/friend"
	"Open_IM/src/utils"
	"context"
	"fmt"
	"strconv"
)

func (s *friendServer) GetFriendApplyList(ctx context.Context, req *pbFriend.GetFriendApplyReq) (*pbFriend.GetFriendApplyResp, error) {
	log.Info(req.Token, req.OperationID, "rpc get friend apply list is server,args=%s", req.String())
	var appleUserList []*pbFriend.ApplyUserInfo
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.GetFriendApplyResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	//	Find the  current user friend applications received
	ApplyUsersInfo, err := im_mysql_model.FindFriendsApplyFromFriendReq(claims.UID)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,search applyInfo failed", err.Error())
		return &pbFriend.GetFriendApplyResp{ErrorCode: config.ErrMysql.ErrCode, ErrorMsg: config.ErrMysql.ErrMsg}, nil
	}
	for _, applyUserInfo := range ApplyUsersInfo {
		var userInfo pbFriend.ApplyUserInfo
		//Find friend application status
		friendReqStatus, err := im_mysql_model.FindFriendRelationshipFromFriendReq(applyUserInfo.ReqId, applyUserInfo.UserId)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s,search friendRelationshipStatus failed", err.Error())
		} else {
			userInfo.Flag = friendReqStatus.Flag
			userInfo.ReqMessage = friendReqStatus.ReqMessage
			userInfo.ApplyTime = strconv.FormatInt(friendReqStatus.CreateTime.Unix(), 10)
		}
		//Find  user information
		us, err := im_mysql_model.FindUserByUID(applyUserInfo.ReqId)
		if err != nil {
			log.Error(req.Token, req.OperationID, fmt.Sprintf("err=%s,search userInfo failed", err.Error()))
		} else {
			userInfo.Uid = us.UID
			userInfo.Icon = us.Icon
			userInfo.Name = us.Name
			userInfo.Gender = us.Gender
			userInfo.Mobile = us.Mobile
			userInfo.Birth = us.Birth
			userInfo.Email = us.Email
			userInfo.Ex = us.Ex
			appleUserList = append(appleUserList, &userInfo)
		}
	}
	log.Info(req.Token, req.OperationID, fmt.Sprintf("rpc get friendapplylist success return"))
	return &pbFriend.GetFriendApplyResp{Data: appleUserList}, nil
}
