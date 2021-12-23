package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pb "Open_IM/pkg/proto/group"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type InviteUserToGroupReq struct {
	GroupID     string   `json:"groupID" binding:"required"`
	UidList     []string `json:"uidList" binding:"required"`
	Reason      string   `json:"reason"`
	OperationID string   `json:"operationID" binding:"required"`
}

type GetJoinedGroupListReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

type KickGroupMemberReq struct {
	GroupID     string                    `json:"groupID"`
	UidListInfo []*pb.GroupMemberFullInfo `json:"uidListInfo" binding:"required"`
	Reason      string                    `json:"reason"`
	OperationID string                    `json:"operationID" binding:"required"`
}

func KickGroupMember(c *gin.Context) {

	params := KickGroupMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &pb.KickGroupMemberReq{
		OperationID: params.OperationID,
		GroupID:     params.GroupID,
		Token:       c.Request.Header.Get("token"),

		UidListInfo: params.UidListInfo,
	}
	log.Info(req.Token, req.OperationID, "recv req: ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	RpcResp, err := client.KickGroupMember(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "GetGroupMemberList failed, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	type KickGroupMemberResp struct {
		ErrorCode int32       `json:"errCode"`
		ErrorMsg  string      `json:"errMsg"`
		Data      []Id2Result `json:"data"`
	}

	var memberListResp KickGroupMemberResp
	memberListResp.ErrorMsg = RpcResp.ErrorMsg
	memberListResp.ErrorCode = RpcResp.ErrorCode
	for _, v := range RpcResp.Id2Result {
		memberListResp.Data = append(memberListResp.Data,
			Id2Result{UId: v.UId,
				Result: v.Result})
	}
	c.JSON(http.StatusOK, memberListResp)
}

type GetGroupMembersInfoReq struct {
	GroupID     string   `json:"groupID"`
	MemberList  []string `json:"memberList"`
	OperationID string   `json:"operationID"`
}
type GetGroupMembersInfoResp struct {
	ErrorCode int32          `json:"errCode"`
	ErrorMsg  string         `json:"errMsg"`
	Data      []MemberResult `json:"data"`
}

func GetGroupMembersInfo(c *gin.Context) {
	log.Info("", "", "GetGroupMembersInfo start....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)

	params := GetGroupMembersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &pb.GetGroupMembersInfoReq{
		OperationID: params.OperationID,
		GroupID:     params.GroupID,
		MemberList:  params.MemberList,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "recv req: ", len(params.MemberList))

	RpcResp, err := client.GetGroupMembersInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "GetGroupMemberList failed, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	var memberListResp GetGroupMembersInfoResp
	memberListResp.ErrorMsg = RpcResp.ErrorMsg
	memberListResp.ErrorCode = RpcResp.ErrorCode
	for _, v := range RpcResp.MemberList {
		memberListResp.Data = append(memberListResp.Data,
			MemberResult{GroupId: req.GroupID,
				UserId:   v.UserId,
				Role:     v.Role,
				JoinTime: uint64(v.JoinTime),
				Nickname: v.NickName,
				FaceUrl:  v.FaceUrl})
	}
	c.JSON(http.StatusOK, memberListResp)
}

type GetGroupMemberListReq struct {
	GroupID     string `json:"groupID"`
	Filter      int32  `json:"filter"`
	NextSeq     int32  `json:"nextSeq"`
	OperationID string `json:"operationID"`
}
type getGroupAllMemberReq struct {
	GroupID     string `json:"groupID"`
	OperationID string `json:"operationID"`
}

type MemberResult struct {
	GroupId  string `json:"groupID"`
	UserId   string `json:"userId"`
	Role     int32  `json:"role"`
	JoinTime uint64 `json:"joinTime"`
	Nickname string `json:"nickName"`
	FaceUrl  string `json:"faceUrl"`
}

func GetGroupMemberList(c *gin.Context) {
	log.Info("", "", "GetGroupMemberList start....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)

	params := GetGroupMemberListReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.GetGroupMemberListReq{
		OperationID: params.OperationID,
		Filter:      params.Filter,
		NextSeq:     params.NextSeq,
		GroupID:     params.GroupID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "recv req: ", req.String())
	RpcResp, err := client.GetGroupMemberList(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "GetGroupMemberList failed, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	type GetGroupMemberListResp struct {
		ErrorCode int32          `json:"errCode"`
		ErrorMsg  string         `json:"errMsg"`
		NextSeq   int32          `json:"nextSeq"`
		Data      []MemberResult `json:"data"`
	}

	var memberListResp GetGroupMemberListResp
	memberListResp.ErrorMsg = RpcResp.ErrorMsg
	memberListResp.ErrorCode = RpcResp.ErrorCode
	memberListResp.NextSeq = RpcResp.NextSeq
	for _, v := range RpcResp.MemberList {
		memberListResp.Data = append(memberListResp.Data,
			MemberResult{GroupId: req.GroupID,
				UserId:   v.UserId,
				Role:     v.Role,
				JoinTime: uint64(v.JoinTime),
				Nickname: v.NickName,
				FaceUrl:  v.FaceUrl})
	}
	c.JSON(http.StatusOK, memberListResp)

}

func GetGroupAllMember(c *gin.Context) {
	log.Info("", "", "GetGroupAllMember start....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)

	params := getGroupAllMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.GetGroupAllMemberReq{
		GroupID:     params.GroupID,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "recv req: ", req.String())
	RpcResp, err := client.GetGroupAllMember(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "GetGroupAllMember failed, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	type GetGroupMemberListResp struct {
		ErrorCode int32          `json:"errCode"`
		ErrorMsg  string         `json:"errMsg"`
		Data      []MemberResult `json:"data"`
	}

	var memberListResp GetGroupMemberListResp
	memberListResp.ErrorMsg = RpcResp.ErrorMsg
	memberListResp.ErrorCode = RpcResp.ErrorCode
	for _, v := range RpcResp.MemberList {
		memberListResp.Data = append(memberListResp.Data,
			MemberResult{GroupId: req.GroupID,
				UserId:   v.UserId,
				Role:     v.Role,
				JoinTime: uint64(v.JoinTime),
				Nickname: v.NickName,
				FaceUrl:  v.FaceUrl})
	}
	c.JSON(http.StatusOK, memberListResp)
}

type groupResult struct {
	GroupId      string `json:"groupId"`
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	FaceUrl      string `json:"faceUrl"`
	OwnerId      string `json:"ownerId"`
	CreateTime   uint64 `json:"createTime"`
	MemberCount  uint32 `json:"memberCount"`
}

func GetJoinedGroupList(c *gin.Context) {
	log.Info("", "", "GetJoinedGroupList start....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	fmt.Println("config:    ", etcdConn, config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)

	params := GetJoinedGroupListReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.GetJoinedGroupListReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "recv req: ", req.String())

	RpcResp, err := client.GetJoinedGroupList(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "GetJoinedGroupList failed, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info(req.Token, req.OperationID, "GetJoinedGroupList: ", RpcResp)

	type GetJoinedGroupListResp struct {
		ErrorCode int32         `json:"errCode"`
		ErrorMsg  string        `json:"errMsg"`
		Data      []groupResult `json:"data"`
	}

	var GroupListResp GetJoinedGroupListResp
	GroupListResp.ErrorCode = RpcResp.ErrorCode
	GroupListResp.ErrorMsg = RpcResp.ErrorMsg
	for _, v := range RpcResp.GroupList {
		GroupListResp.Data = append(GroupListResp.Data,
			groupResult{GroupId: v.GroupId, GroupName: v.GroupName,
				Notification: v.Notification,
				Introduction: v.Introduction,
				FaceUrl:      v.FaceUrl,
				OwnerId:      v.OwnerId,
				CreateTime:   v.CreateTime,
				MemberCount:  v.MemberCount})
	}
	c.JSON(http.StatusOK, GroupListResp)
}

type Id2Result struct {
	UId    string `json:"uid"`
	Result int32  `json:"result"`
}

func InviteUserToGroup(c *gin.Context) {
	log.Info("", "", "InviteUserToGroup start....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)

	params := InviteUserToGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.InviteUserToGroupReq{
		OperationID: params.OperationID,
		GroupID:     params.GroupID,
		Reason:      params.Reason,
		UidList:     params.UidList,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "recv req: ", req.String())

	RpcResp, err := client.InviteUserToGroup(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "InviteUserToGroup failed, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	type InviteUserToGroupResp struct {
		ErrorCode int32       `json:"errCode"`
		ErrorMsg  string      `json:"errMsg"`
		I2R       []Id2Result `json:"data"`
	}

	var iResp InviteUserToGroupResp
	iResp.ErrorMsg = RpcResp.ErrorMsg
	iResp.ErrorCode = RpcResp.ErrorCode
	for _, v := range RpcResp.Id2Result {
		iResp.I2R = append(iResp.I2R, Id2Result{UId: v.UId, Result: v.Result})
	}

	//resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg, "data": RpcResp.Id2Result}
	c.JSON(http.StatusOK, iResp)
}

type paramsCreateGroupStruct struct {
	MemberList   []*pb.GroupAddMemberInfo `json:"memberList"`
	GroupName    string                   `json:"groupName"`
	Introduction string                   `json:"introduction"`
	Notification string                   `json:"notification"`
	FaceUrl      string                   `json:"faceUrl"`
	OperationID  string                   `json:"operationID" binding:"required"`
	Ex           string                   `json:"ex"`
}

func CreateGroup(c *gin.Context) {
	log.Info("", "", "api create group init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsCreateGroupStruct{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.CreateGroupReq{
		MemberList:   params.MemberList,
		GroupName:    params.GroupName,
		Introduction: params.Introduction,
		Notification: params.Notification,
		FaceUrl:      params.FaceUrl,
		OperationID:  params.OperationID,
		Ex:           params.Ex,
		Token:        c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api create group is server,params=%s", req.String())
	RpcResp, err := client.CreateGroup(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call create group  rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call create group  rpc server success,args=%s", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg, "data": gin.H{"groupID": RpcResp.GroupID}}
		c.JSON(http.StatusOK, resp)
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	}
	log.InfoByArgs("api create group success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}

type paramsGroupApplicationList struct {
	OperationID string `json:"operationID" binding:"required"`
}

func newUserRegisterReq(params *paramsGroupApplicationList) *group.GetGroupApplicationListReq {
	pbData := group.GetGroupApplicationListReq{
		OperationID: params.OperationID,
	}
	return &pbData
}

type paramsGroupApplicationListRet struct {
	ID               string `json:"id"`
	GroupID          string `json:"groupID"`
	FromUserID       string `json:"fromUserID"`
	ToUserID         string `json:"toUserID"`
	Flag             int32  `json:"flag"`
	RequestMsg       string `json:"reqMsg"`
	HandledMsg       string `json:"handledMsg"`
	AddTime          int64  `json:"createTime"`
	FromUserNickname string `json:"fromUserNickName"`
	ToUserNickname   string `json:"toUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceURL"`
	ToUserFaceUrl    string `json:"toUserFaceURL"`
	HandledUser      string `json:"handledUser"`
	Type             int32  `json:"type"`
	HandleStatus     int32  `json:"handleStatus"`
	HandleResult     int32  `json:"handleResult"`
}

func GetGroupApplicationList(c *gin.Context) {
	log.Info("", "", "api GetGroupApplicationList init ....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := group.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsGroupApplicationList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	pbData := newUserRegisterReq(&params)

	token := c.Request.Header.Get("token")
	if claims, err := token_verify.ParseToken(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	} else {
		pbData.UID = claims.UID
	}

	log.Info("", "", "api GetGroupApplicationList is server, [data: %s]", pbData.String())
	reply, err := client.GetGroupApplicationList(context.Background(), pbData)
	if err != nil {
		log.Error("", "", "api GetGroupApplicationList call rpc fail, [data: %s] [err: %s]", pbData.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info("", "", "api GetGroupApplicationList call rpc success, [data: %s] [reply: %s]", pbData.String(), reply.String())

	unProcessCount := 0
	userReq := make([]paramsGroupApplicationListRet, 0)
	if reply != nil && reply.Data != nil && reply.Data.User != nil {
		for i := 0; i < len(reply.Data.User); i++ {
			req := paramsGroupApplicationListRet{}
			req.ID = reply.Data.User[i].ID
			req.GroupID = reply.Data.User[i].GroupID
			req.FromUserID = reply.Data.User[i].FromUserID
			req.ToUserID = reply.Data.User[i].ToUserID
			req.Flag = reply.Data.User[i].Flag
			req.RequestMsg = reply.Data.User[i].RequestMsg
			req.HandledMsg = reply.Data.User[i].HandledMsg
			req.AddTime = reply.Data.User[i].AddTime
			req.FromUserNickname = reply.Data.User[i].FromUserNickname
			req.ToUserNickname = reply.Data.User[i].ToUserNickname
			req.FromUserFaceUrl = reply.Data.User[i].FromUserFaceUrl
			req.ToUserFaceUrl = reply.Data.User[i].ToUserFaceUrl
			req.HandledUser = reply.Data.User[i].HandledUser
			req.Type = reply.Data.User[i].Type
			req.HandleStatus = reply.Data.User[i].HandleStatus
			req.HandleResult = reply.Data.User[i].HandleResult
			userReq = append(userReq, req)

			if req.Flag == 0 {
				unProcessCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
		"data": gin.H{
			"count": unProcessCount,
			"user":  userReq,
		},
	})

}

type paramsGetGroupInfo struct {
	GroupIDList []string `json:"groupIDList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}

func GetGroupsInfo(c *gin.Context) {
	log.Info("", "", "api get groups info init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsGetGroupInfo{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.GetGroupsInfoReq{
		GroupIDList: params.GroupIDList,
		Token:       c.Request.Header.Get("token"),
		OperationID: params.OperationID,
	}
	log.Info(req.Token, req.OperationID, "get groups info is server,params=%s", req.String())
	RpcResp, err := client.GetGroupsInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "call get groups info rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call get groups info rpc server success", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		groupsInfo := make([]pb.GroupInfo, 0)
		for _, v := range RpcResp.Data {
			var groupInfo pb.GroupInfo
			groupInfo.GroupId = v.GroupId
			groupInfo.GroupName = v.GroupName
			groupInfo.Notification = v.Notification
			groupInfo.Introduction = v.Introduction
			groupInfo.FaceUrl = v.FaceUrl
			groupInfo.CreateTime = v.CreateTime
			groupInfo.OwnerId = v.OwnerId
			groupInfo.MemberCount = v.MemberCount

			groupsInfo = append(groupsInfo, groupInfo)
		}
		c.JSON(http.StatusOK, gin.H{
			"errCode": RpcResp.ErrorCode,
			"errMsg":  RpcResp.ErrorMsg,
			"data":    groupsInfo,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	}
}

type paramsGroupApplicationResponse struct {
	OperationID      string `json:"operationID" binding:"required"`
	GroupID          string `json:"groupID" binding:"required"`
	FromUserID       string `json:"fromUserID" binding:"required"`
	FromUserNickName string `json:"fromUserNickName"`
	FromUserFaceUrl  string `json:"fromUserFaceUrl"`
	ToUserID         string `json:"toUserID" binding:"required"`
	ToUserNickName   string `json:"toUserNickName"`
	ToUserFaceUrl    string `json:"toUserFaceUrl"`
	AddTime          int64  `json:"addTime"`
	RequestMsg       string `json:"requestMsg"`
	HandledMsg       string `json:"handledMsg"`
	Type             int32  `json:"type"`
	HandleStatus     int32  `json:"handleStatus"`
	HandleResult     int32  `json:"handleResult"`

	UserID string `json:"userID"`
}

func newGroupApplicationResponse(params *paramsGroupApplicationResponse) *group.GroupApplicationResponseReq {
	pbData := group.GroupApplicationResponseReq{
		OperationID:      params.OperationID,
		GroupID:          params.GroupID,
		FromUserID:       params.FromUserID,
		FromUserNickName: params.FromUserNickName,
		FromUserFaceUrl:  params.FromUserFaceUrl,
		ToUserID:         params.ToUserID,
		ToUserNickName:   params.ToUserNickName,
		ToUserFaceUrl:    params.ToUserFaceUrl,
		AddTime:          params.AddTime,
		RequestMsg:       params.RequestMsg,
		HandledMsg:       params.HandledMsg,
		Type:             params.Type,
		HandleStatus:     params.HandleStatus,
		HandleResult:     params.HandleResult,
	}
	return &pbData
}

func ApplicationGroupResponse(c *gin.Context) {
	log.Info("", "", "api GroupApplicationResponse init ....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := group.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsGroupApplicationResponse{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	pbData := newGroupApplicationResponse(&params)

	token := c.Request.Header.Get("token")
	if claims, err := token_verify.ParseToken(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	} else {
		pbData.OwnerID = claims.UID
	}

	log.Info("", "", "api GroupApplicationResponse is server, [data: %s]", pbData.String())
	reply, err := client.GroupApplicationResponse(context.Background(), pbData)
	if err != nil {
		log.Error("", "", "api GroupApplicationResponse call rpc fail, [data: %s] [err: %s]", pbData.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info("", "", "api GroupApplicationResponse call rpc success, [data: %s] [reply: %s]", pbData.String(), reply.String())

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
	})

}

type paramsJoinGroup struct {
	GroupID     string `json:"groupID" binding:"required"`
	Message     string `json:"message"`
	OperationID string `json:"operationID" binding:"required"`
}

func JoinGroup(c *gin.Context) {
	log.Info("", "", "api join group init....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsJoinGroup{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.JoinGroupReq{
		GroupID:     params.GroupID,
		Message:     params.Message,
		Token:       c.Request.Header.Get("token"),
		OperationID: params.OperationID,
	}
	log.Info(req.Token, req.OperationID, "api join group is server,params=%s", req.String())
	RpcResp, err := client.JoinGroup(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "call join group  rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call join group rpc server success,args=%s", RpcResp.String())
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
}

type paramsQuitGroup struct {
	GroupID     string `json:"groupID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

func QuitGroup(c *gin.Context) {
	log.Info("", "", "api quit group init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsQuitGroup{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.QuitGroupReq{
		GroupID:     params.GroupID,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api quit group is server,params=%s", req.String())
	RpcResp, err := client.QuitGroup(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "call quit group rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call quit group rpc server success,args=%s", RpcResp.String())
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	log.InfoByArgs("api quit group success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}

type paramsSetGroupInfo struct {
	GroupID      string `json:"groupId"  binding:"required"`
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	FaceUrl      string `json:"faceUrl"`
	OperationID  string `json:"operationID"  binding:"required"`
}

func SetGroupInfo(c *gin.Context) {
	log.Info("", "", "api set group info init...")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsSetGroupInfo{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.SetGroupInfoReq{
		GroupID:      params.GroupID,
		GroupName:    params.GroupName,
		Notification: params.Notification,
		Introduction: params.Introduction,
		FaceUrl:      params.FaceUrl,
		Token:        c.Request.Header.Get("token"),
		OperationID:  params.OperationID,
	}
	log.Info(req.Token, req.OperationID, "api set group info is server,params=%s", req.String())
	RpcResp, err := client.SetGroupInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "call set group info rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call set group info rpc server success,args=%s", RpcResp.String())
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
}

type paramsTransferGroupOwner struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
}

func newTransferGroupOwnerReq(params *paramsTransferGroupOwner) *group.TransferGroupOwnerReq {
	pbData := group.TransferGroupOwnerReq{
		OperationID: params.OperationID,
		GroupID:     params.GroupID,
		NewOwner:    params.UID,
	}
	return &pbData
}

func TransferGroupOwner(c *gin.Context) {
	log.Info("", "", "api TransferGroupOwner init ....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := group.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsTransferGroupOwner{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	pbData := newTransferGroupOwnerReq(&params)

	token := c.Request.Header.Get("token")
	if claims, err := token_verify.ParseToken(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	} else {
		pbData.OldOwner = claims.UID
	}

	log.Info("", "", "api TransferGroupOwner is server, [data: %s]", pbData.String())
	reply, err := client.TransferGroupOwner(context.Background(), pbData)
	if err != nil {
		log.Error("", "", "api TransferGroupOwner call rpc fail, [data: %s] [err: %s]", pbData.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info("", "", "api TransferGroupOwner call rpc success, [data: %s] [reply: %s]", pbData.String(), reply.String())

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
	})

}
