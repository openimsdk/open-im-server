package friend

import (
	"Open_IM/internal/rpc/chat"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type friendServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewFriendServer(port int) *friendServer {
	log.NewPrivateLog("friend")
	return &friendServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImFriendName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *friendServer) Run() {
	log.Info("", "", fmt.Sprintf("rpc friend init...."))

	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.InfoByArgs(fmt.Sprintf("Failed to listen rpc friend network,err=%s", err.Error()))
		return
	}
	log.Info("", "", "listen network success, address = %s", registerAddress)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//User friend related services register to etcd
	pbFriend.RegisterFriendServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.ErrorByArgs("register rpc fiend service to etcd failed,err=%s", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.ErrorByArgs("listen rpc friend error,err=%s", err.Error())
		return
	}
}

func (s *friendServer) GetFriendsInfo(ctx context.Context, req *pbFriend.GetFriendsInfoReq) (*pbFriend.GetFriendInfoResp, error) {
	log.NewInfo(req.CommID.OperationID, "GetFriendsInfo args ", req.String())
	var (
		isInBlackList int32
		//	isFriend      int32
		comment string
	)

	friendShip, err := imdb.FindFriendRelationshipFromFriend(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "FindFriendRelationshipFromFriend failed ", err.Error())
		return &pbFriend.GetFriendInfoResp{ErrCode: constant.ErrSearchUserInfo.ErrCode, ErrMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
		//	isFriend = constant.FriendFlag
	}
	comment = friendShip.Remark

	friendUserInfo, err := imdb.FindUserByUID(req.CommID.ToUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "FindUserByUID failed ", err.Error())
		return &pbFriend.GetFriendInfoResp{ErrCode: constant.ErrSearchUserInfo.ErrCode, ErrMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
	}

	err = imdb.FindRelationshipFromBlackList(req.CommID.FromUserID, req.CommID.ToUserID)
	if err == nil {
		isInBlackList = constant.BlackListFlag
	}

	resp := pbFriend.GetFriendInfoResp{
		ErrCode: 0,
		ErrMsg:  "",
		Data: &pbFriend.FriendInfo{
			IsBlack: isInBlackList,
		},
	}
	utils.CopyStructFields(resp.Data.FriendUser, friendUserInfo)
	resp.Data.IsBlack = isInBlackList
	resp.Data.OwnerUserID = req.CommID.FromUserID
	resp.Data.Remark = comment
	resp.Data.CreateTime = friendUserInfo.CreateTime

	log.NewInfo(req.CommID.OperationID, "GetFriendsInfo ok ", resp)
	return &resp, nil

}

func (s *friendServer) AddBlacklist(ctx context.Context, req *pbFriend.AddBlacklistReq) (*pbFriend.CommonResp, error) {
	log.NewInfo(req.CommID.OperationID, "AddBlacklist args ", req.String())
	ok := token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID)
	if !ok {
		log.NewError(req.CommID.OperationID, "CheckAccess failed ", req.CommID.OpUserID, req.CommID.FromUserID)
	}

	err := imdb.InsertInToUserBlackList(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "InsertInToUserBlackList failed ", err.Error())
		return &pbFriend.CommonResp{ErrCode: constant.ErrMysql.ErrCode, ErrMsg: constant.ErrMysql.ErrMsg}, nil
	}
	log.NewInfo(req.CommID.OperationID, "InsertInToUserBlackList ok ", req.CommID.FromUserID, req.CommID.ToUserID)
	chat.BlackAddedNotification(req.CommID.OperationID, req.CommID.OpUserID, req.CommID.FromUserID, req.CommID.ToUserID)
	return &pbFriend.CommonResp{}, nil

}

func (s *friendServer) AddFriend(ctx context.Context, req *pbFriend.AddFriendReq) (*pbFriend.CommonResp, error) {
	log.NewInfo(req.CommID.OperationID, "AddFriend args ", req.String())
	ok := token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID)
	if !ok {
		log.NewError(req.CommID.OperationID, "CheckAccess failed ", req.CommID.OpUserID, req.CommID.FromUserID)
	}
	//Cannot add non-existent users
	if _, err := imdb.FindUserByUID(req.CommID.ToUserID); err != nil {
		log.NewError(req.CommID.OperationID, "FindUserByUID failed ", err.Error(), req.CommID.ToUserID)
		return &pbFriend.CommonResp{ErrCode: constant.ErrAddFriend.ErrCode, ErrMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
	}

	//Establish a latest relationship in the friend request table
	err := imdb.ReplaceIntoFriendReq(req.CommID.FromUserID, req.CommID.ToUserID, constant.ApplicationFriendFlag, req.ReqMessage)
	if err != nil {
		log.NewError(req.CommID.OperationID, "ReplaceIntoFriendReq failed ", err.Error())
		return &pbFriend.CommonResp{ErrCode: constant.ErrAddFriend.ErrCode, ErrMsg: constant.ErrAddFriend.ErrMsg}, nil
	}

	chat.FriendApplicationAddedNotification(req.CommID.OperationID, req.CommID.OpUserID, req.CommID.FromUserID, req.CommID.ToUserID, req.ReqMessage)
	return &pbFriend.CommonResp{}, nil
}

func (s *friendServer) ImportFriend(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.ImportFriendResp, error) {
	log.NewInfo(req.OperationID, "ImportFriend failed ", req.String())
	var resp pbFriend.ImportFriendResp
	var c pbFriend.CommonResp

	if !utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
		log.NewError(req.OperationID, "not authorized", req.OpUserID)
		c.ErrCode = constant.ErrAddFriend.ErrCode
		c.ErrMsg = "not authorized"
		return &pbFriend.ImportFriendResp{CommonResp: &c, FailedUidList: req.FriendUserIDList}, nil
	}

	if _, err := imdb.FindUserByUID(req.FromUserID); err != nil {
		log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), req.FromUserID)
		c.ErrCode = constant.ErrAddFriend.ErrCode
		c.ErrMsg = "this user not exists,cant not add friend"
		return &pbFriend.ImportFriendResp{CommonResp: &c, FailedUidList: req.FriendUserIDList}, nil
	}

	for _, v := range req.FriendUserIDList {
		if _, fErr := imdb.FindUserByUID(v); fErr != nil {
			c.ErrMsg = "some uid establish failed"
			c.ErrCode = 408
			resp.FailedUidList = append(resp.FailedUidList, v)
		} else {
			if _, err := imdb.FindFriendRelationshipFromFriend(req.FromUserID, v); err != nil {
				//Establish two single friendship
				err1 := imdb.InsertToFriend(req.FromUserID, v, 1)
				if err1 != nil {
					resp.FailedUidList = append(resp.FailedUidList, v)
					log.NewError(req.OperationID, "InsertToFriend failed", req.FromUserID, v, err1.Error())
					c.ErrMsg = "some uid establish failed"
					c.ErrCode = 408
					continue
				}
				err2 := imdb.InsertToFriend(v, req.FromUserID, 1)
				if err2 != nil {
					resp.FailedUidList = append(resp.FailedUidList, v)
					log.NewError(req.OperationID, "InsertToFriend failed", v, req.FromUserID, err2.Error())
					c.ErrMsg = "some uid establish failed"
					c.ErrCode = 408
					continue
				}
				chat.FriendAddedNotification(req.OperationID, req.OpUserID, req.FromUserID, v)
			}
		}
	}
	resp.CommonResp = &c
	log.NewInfo(req.OperationID, "ImportFriend rpc ok ", resp)
	return &resp, nil
}

func (s *friendServer) AddFriendResponse(ctx context.Context, req *pbFriend.AddFriendResponseReq) (*pbFriend.CommonResp, error) {
	log.NewInfo(req.CommID.OperationID, "AddFriendResponse args ", req.String())

	if !token_verify.CheckAccess(rreq.CommID.FromUserID, req.CommID.ToUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess failed ", req.CommID.FromUserID, req.CommID.ToUserID)
		return &pbFriend.CommonResp{ErrCode: constant.ErrAgreeToAddFriend.ErrCode, ErrMsg: constant.ErrAgreeToAddFriend.ErrMsg}, nil
	}

	//Check there application before agreeing or refuse to a friend's application
	//req.CommID.FromUserID process req.CommID.ToUserID
	if _, err := imdb.FindFriendApplyFromFriendReqByUid(req.CommID.ToUserID, req.CommID.FromUserID); err != nil {
		log.NewError(req.CommID.OperationID, "FindFriendApplyFromFriendReqByUid failed ", err.Error(), req.CommID.ToUserID, req.CommID.FromUserID)
		return &pbFriend.CommonResp{ErrCode: constant.ErrAgreeToAddFriend.ErrCode, ErrMsg: constant.ErrAgreeToAddFriend.ErrMsg}, nil
	}

	//Change friend request status flag
	err := imdb.UpdateFriendRelationshipToFriendReq(req.CommID.ToUserID, req.CommID.FromUserID, req.Flag)
	if err != nil {
		log.NewError(req.CommID.OperationID, "UpdateFriendRelationshipToFriendReq failed ", err.Error(), req.CommID.ToUserID, req.CommID.FromUserID, req.Flag)
		return &pbFriend.CommonResp{ErrCode: constant.ErrMysql.ErrCode, ErrMsg: constant.ErrAgreeToAddFriend.ErrMsg}, nil
	}
	log.NewInfo(req.CommID.OperationID, "rpc AddFriendResponse ok")
	//Change the status of the friend request form
	if req.Flag == constant.FriendFlag {
		//Establish friendship after find friend relationship not exists
		_, err := imdb.FindFriendRelationshipFromFriend(req.CommID.FromUserID, req.CommID.ToUserID)
		if err == nil {
			log.NewWarn(req.CommID.OperationID, "FindFriendRelationshipFromFriend exist", req.CommID.FromUserID, req.CommID.ToUserID)
		} else {
			//Establish two single friendship
			err = imdb.InsertToFriend(req.CommID.FromUserID, req.CommID.ToUserID, req.Flag)
			if err != nil {
				log.NewError(req.CommID.OperationID, "InsertToFriend failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID, req.Flag)
				return &pbFriend.CommonResp{ErrCode: constant.ErrAgreeToAddFriend.ErrCode, ErrMsg: constant.ErrAgreeToAddFriend.ErrMsg}, nil
			}
		}

		_, err = imdb.FindFriendRelationshipFromFriend(req.CommID.ToUserID, req.CommID.FromUserID)
		if err == nil {
			log.NewWarn(req.CommID.OperationID, "FindFriendRelationshipFromFriend exist", req.CommID.ToUserID, req.CommID.FromUserID)
			return &pbFriend.CommonResp{}, nil
		}
		err = imdb.InsertToFriend(req.CommID.ToUserID, req.CommID.FromUserID, req.Flag)
		if err != nil {
			log.NewError(req.CommID.OperationID, "InsertToFriend failed ", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID, req.Flag)
			return &pbFriend.CommonResp{ErrCode: constant.ErrAgreeToAddFriend.ErrCode, ErrMsg: constant.ErrAgreeToAddFriend.ErrMsg}, nil
		}
	}

	chat.FriendApplicationProcessedNotification(req.CommID.OperationID, req.CommID.OpUserID, req.CommID.FromUserID, req.CommID.ToUserID, req.Flag)
	return &pbFriend.CommonResp{}, nil
}

func (s *friendServer) DeleteFriend(ctx context.Context, req *pbFriend.DeleteFriendReq) (*pbFriend.CommonResp, error) {
	log.NewInfo(req.CommID.OperationID, "DeleteFriend args ", req.String())
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.CommonResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	err := imdb.DeleteSingleFriendInfo(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "DeleteSingleFriendInfo failed", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
		return &pbFriend.CommonResp{ErrCode: constant.ErrMysql.ErrCode, ErrMsg: constant.ErrMysql.ErrMsg}, nil
	}
	log.NewInfo(req.CommID.OperationID, "DeleteFriend rpc ok")
	chat.FriendDeletedNotification(req.CommID.OperationID, req.CommID.OpUserID, req.CommID.FromUserID, req.CommID.ToUserID)
	return &pbFriend.CommonResp{}, nil
}

func (s *friendServer) GetBlacklist(ctx context.Context, req *pbFriend.GetBlacklistReq) (*pbFriend.GetBlacklistResp, error) {
	log.NewInfo(req.CommID.OperationID, "GetBlacklist args ", req.String())

	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess failed", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.GetBlacklistResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	blackListInfo, err := imdb.GetBlackListByUID(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetBlackListByUID failed ", err.Error(), req.CommID.FromUserID)
		return &pbFriend.GetBlacklistResp{ErrCode: constant.ErrGetBlackList.ErrCode, ErrMsg: constant.ErrGetBlackList.ErrMsg}, nil
	}

	var (
		userInfoList []*sdk_ws.PublicUserInfo
	)
	for _, blackUser := range blackListInfo {
		var blackUserInfo sdk_ws.PublicUserInfo
		//Find black user information
		us, err := imdb.FindUserByUID(blackUser.BlockUserID)
		if err != nil {
			log.NewError(req.CommID.OperationID, "FindUserByUID failed ", err.Error(), blackUser.BlockUserID)
			continue
		}
		utils.CopyStructFields(&blackUserInfo, us)
		userInfoList = append(userInfoList, &blackUserInfo)
	}
	log.NewInfo(req.CommID.OperationID, "rpc GetBlacklist ok")
	return &pbFriend.GetBlacklistResp{Data: userInfoList}, nil
}

func (s *friendServer) SetFriendComment(ctx context.Context, req *pbFriend.SetFriendCommentReq) (*pbFriend.CommonResp, error) {
	log.NewInfo(req.CommID.OperationID, "SetFriendComment args ", req.String())
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.CommonResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	err := imdb.UpdateFriendComment(req.CommID.FromUserID, req.CommID.OpUserID, req.Remark)
	if err != nil {
		log.NewError(req.CommID.OperationID, "UpdateFriendComment failed ", err.Error(), req.CommID.FromUserID, req.CommID.OpUserID, req.Remark)
		return &pbFriend.CommonResp{ErrCode: constant.ErrSetFriendComment.ErrCode, ErrMsg: constant.ErrSetFriendComment.ErrMsg}, nil
	}
	log.NewInfo(req.CommID.OperationID, "rpc SetFriendComment ok")
	chat.FriendInfoChangedNotification(req.CommID.OperationID, req.CommID.OpUserID, req.CommID.FromUserID, req.CommID.ToUserID)
	return &pbFriend.CommonResp{}, nil
}

func (s *friendServer) RemoveBlacklist(ctx context.Context, req *pbFriend.RemoveBlacklistReq) (*pbFriend.CommonResp, error) {
	log.NewInfo(req.CommID.OperationID, "RemoveBlacklist args ", req.String())
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.CommonResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	err := imdb.RemoveBlackList(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "RemoveBlackList failed", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
		return &pbFriend.CommonResp{ErrCode: constant.ErrMysql.ErrCode, ErrMsg: constant.ErrMysql.ErrMsg}, nil
	}
	log.NewInfo(req.CommID.OperationID, "rpc RemoveBlacklist ok")
	chat.BlackDeletedNotification(req.CommID.OperationID, req.CommID.OpUserID, req.CommID.FromUserID, req.CommID.ToUserID)
	return &pbFriend.CommonResp{}, nil
}

func (s *friendServer) IsInBlackList(ctx context.Context, req *pbFriend.IsInBlackListReq) (*pbFriend.IsInBlackListResp, error) {
	log.NewInfo("IsInBlackList args ", req.String())
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.IsInBlackListResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	var isInBlacklist = false
	err := imdb.FindRelationshipFromBlackList(req.CommID.FromUserID, req.CommID.ToUserID)
	if err == nil {
		isInBlacklist = true
	}
	log.NewInfo(req.CommID.OperationID, "IsInBlackList rpc ok")
	return &pbFriend.IsInBlackListResp{Response: isInBlacklist}, nil
}

func (s *friendServer) IsFriend(ctx context.Context, req *pbFriend.IsFriendReq) (*pbFriend.IsFriendResp, error) {
	log.NewInfo("IsFriend args ", req.String())
	var isFriend int32
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.IsFriendResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	_, err := imdb.FindFriendRelationshipFromFriend(req.CommID.FromUserID, req.CommID.ToUserID)
	if err == nil {
		isFriend = constant.FriendFlag
	} else {
		isFriend = constant.ApplicationFriendFlag
	}
	log.NewInfo("IsFriend rpc ok")
	return &pbFriend.IsFriendResp{ShipType: isFriend}, nil
}

func (s *friendServer) GetFriendList(ctx context.Context, req *pbFriend.GetFriendListReq) (*pbFriend.GetFriendListResp, error) {
	log.NewInfo("GetFriendList args ", req.String())
	var userInfoList []*pbFriend.FriendInfo
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.GetFriendListResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	friends, err := imdb.FindUserInfoFromFriend(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "FindUserInfoFromFriend failed", err.Error(), req.CommID.FromUserID)
		return &pbFriend.GetFriendListResp{ErrCode: constant.ErrSearchUserInfo.ErrCode, ErrMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
	}
	for _, friendUser := range friends {
		var friendUserInfo pbFriend.FriendInfo
		//find user is in blackList
		err = imdb.FindRelationshipFromBlackList(req.CommID.FromUserID, friendUser.FriendUserID)
		if err == nil {
			friendUserInfo.IsBlack = constant.BlackListFlag
		} else {
			friendUserInfo.IsBlack = 0
		}
		//Find user information
		us, err := imdb.FindUserByUID(friendUser.FriendUserID)
		if err != nil {
			log.NewError(req.CommID.OperationID, "FindUserByUID failed", err.Error(), friendUser.FriendUserID)
			continue
		}
		utils.CopyStructFields(friendUserInfo.FriendUser, us)
		friendUserInfo.Remark = friendUser.Remark
		friendUserInfo.OwnerUserID = req.CommID.FromUserID
		friendUserInfo.CreateTime = friendUser.CreateTime
		userInfoList = append(userInfoList, &friendUserInfo)
	}
	log.NewInfo(req.CommID.OperationID, "rpc GetFriendList ok", pbFriend.GetFriendListResp{Data: userInfoList})
	return &pbFriend.GetFriendListResp{Data: userInfoList}, nil
}

func (s *friendServer) GetFriendApplyList(ctx context.Context, req *pbFriend.GetFriendApplyReq) (*pbFriend.GetFriendApplyResp, error) {
	log.NewInfo(req.CommID.OperationID, "GetFriendApplyList args ", req.String())
	var appleUserList []*pbFriend.ApplyUserInfo
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.GetFriendApplyResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	//	Find the  current user friend applications received
	ApplyUsersInfo, err := imdb.FindFriendsApplyFromFriendReq(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "FindFriendsApplyFromFriendReq ", err.Error(), req.CommID.FromUserID)
		return &pbFriend.GetFriendApplyResp{ErrCode: constant.ErrMysql.ErrCode, ErrMsg: constant.ErrMysql.ErrMsg}, nil
	}

	for _, applyUserInfo := range ApplyUsersInfo {
		var userInfo pbFriend.ApplyUserInfo
		//Find friend application status
		userInfo.Flag = applyUserInfo.Flag
		userInfo.ReqMessage = applyUserInfo.ReqMessage
		userInfo.ApplyTime = applyUserInfo.CreateTime

		//Find  user information
		us, err := imdb.FindUserByUID(applyUserInfo.ReqID)
		if err != nil {
			log.NewError(req.CommID.OperationID, "FindUserByUID failed ", err.Error(), applyUserInfo.ReqID)
			continue
		}
		utils.CopyStructFields(userInfo.UserInfo, us)
		appleUserList = append(appleUserList, &userInfo)

	}
	log.NewInfo(req.CommID.OperationID, "rpc GetFriendApplyList ok", pbFriend.GetFriendApplyResp{Data: appleUserList})
	return &pbFriend.GetFriendApplyResp{Data: appleUserList}, nil
}

func (s *friendServer) GetSelfApplyList(ctx context.Context, req *pbFriend.GetFriendApplyReq) (*pbFriend.GetFriendApplyResp, error) {
	log.NewInfo(req.CommID.OperationID, "GetSelfApplyList args ", req.String())
	var selfApplyOtherUserList []*pbFriend.ApplyUserInfo
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.GetFriendApplyResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	//	Find the self add other userinfo
	usersInfo, err := imdb.FindSelfApplyFromFriendReq(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "FindSelfApplyFromFriendReq failed ", err.Error(), req.CommID.FromUserID)
		return &pbFriend.GetFriendApplyResp{ErrCode: constant.ErrMysql.ErrCode, ErrMsg: constant.ErrMysql.ErrMsg}, nil
	}
	for _, selfApplyOtherUserInfo := range usersInfo {
		var userInfo pbFriend.ApplyUserInfo
		//Find friend application status
		userInfo.Flag = selfApplyOtherUserInfo.Flag
		userInfo.ReqMessage = selfApplyOtherUserInfo.ReqMessage
		userInfo.ApplyTime = selfApplyOtherUserInfo.CreateTime
		//Find  user information
		us, err := imdb.FindUserByUID(selfApplyOtherUserInfo.UserID)
		if err != nil {
			log.NewError(req.CommID.OperationID, "FindUserByUID failed", err.Error(), selfApplyOtherUserInfo.UserID)
			continue
		}
		utils.CopyStructFields(userInfo.UserInfo, us)
		selfApplyOtherUserList = append(selfApplyOtherUserList, &userInfo)
	}
	log.NewInfo(req.CommID.OperationID, "rpc GetSelfApplyList ok", pbFriend.GetFriendApplyResp{Data: selfApplyOtherUserList})
	return &pbFriend.GetFriendApplyResp{Data: selfApplyOtherUserList}, nil
}
