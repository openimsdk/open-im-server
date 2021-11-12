package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pb "Open_IM/pkg/proto/group"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// InviteUserToGroupReq struct
type InviteUserToGroupReq struct {
	GroupID     string   `json:"groupID" binding:"required"`
	UidList     []string `json:"uidList" binding:"required"`
	Reason      string   `json:"reason"`
	OperationID string   `json:"operationID" binding:"required"`
}

// GetJoinedGroupListReq struct
type GetJoinedGroupListReq struct {
	OperationID string `json:"operationID" binding:"required"`
}

// paramsKickGroupMember struct
type paramsKickGroupMember struct {
	GroupID     string `json:"groupID"`
	UidListInfo []struct {
		UserId   string `json:"userId,omitempty"`
		Role     int32  `json:"role,omitempty"`
		JoinTime uint64 `json:"joinTime,omitempty"`
		NickName string `json:"nickName,omitempty"`
		FaceUrl  string `json:"faceUrl,omitempty"`
	} `json:"uidListInfo" binding:"required"`
	Reason      string `json:"reason"`
	OperationID string `json:"operationID" binding:"required"`
}

type KickGroupMemberReq struct {
	GroupID     string                    `json:"groupID"`
	UidListInfo []*pb.GroupMemberFullInfo `json:"uidListInfo" binding:"required"`
	Reason      string                    `json:"reason"`
	OperationID string                    `json:"operationID" binding:"required"`
}

// @Summary
// @Schemes
// @Description kick member from group
// @Tags group
// @Accept json
// @Produce json
// @Param body body group.paramsKickGroupMember true "kick member from group params"
// @Param token header string true "token"
// @Success 200 {object} group.KickGroupMemberResp
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /group/kick_group [post]
func KickGroupMember(c *gin.Context) {
	log.Info("", "", "KickGroupMember start....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)

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

	RpcResp, err := client.KickGroupMember(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "GetGroupMemberList failed, err: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	// KickGroupMemberResp struct
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

// GetGroupMembersInfoReq struct
type GetGroupMembersInfoReq struct {
	GroupID     string   `json:"groupID"`
	MemberList  []string `json:"memberList"`
	OperationID string   `json:"operationID"`
}

// GetGroupMembersInfoResp struct
type GetGroupMembersInfoResp struct {
	ErrorCode int32          `json:"errCode"`
	ErrorMsg  string         `json:"errMsg"`
	Data      []MemberResult `json:"data"`
}

// @Summary
// @Schemes
// @Description get group members info
// @Tags group
// @Accept json
// @Produce json
// @Param body body group.GetGroupMembersInfoReq true "get group members info params"
// @Param token header string true "token"
// @Success 200 {object} group.GetGroupMembersInfoResp
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /group/get_group_members_info [post]
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

// GetGroupMemberListReq struct
type GetGroupMemberListReq struct {
	GroupID     string `json:"groupID"`
	Filter      int32  `json:"filter"`
	NextSeq     int32  `json:"nextSeq,omitempty"`
	OperationID string `json:"operationID"`
}

// getGroupAllMemberReq struct
type getGroupAllMemberReq struct {
	GroupID     string `json:"groupID"`
	OperationID string `json:"operationID"`
}

// MemberResult struct
type MemberResult struct {
	GroupId  string `json:"groupID"`
	UserId   string `json:"userId"`
	Role     int32  `json:"role"`
	JoinTime uint64 `json:"joinTime"`
	Nickname string `json:"nickName"`
	FaceUrl  string `json:"faceUrl"`
}

// @Summary
// @Schemes
// @Description get group member list
// @Tags group
// @Accept json
// @Produce json
// @Param body body group.GetGroupMemberListReq true "get group member list params"
// @Param token header string true "token"
// @Success 200 {object} group.GetGroupMemberListResp
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /group/get_group_member_list [post]
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

	// GetGroupMemberListResp struct
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

// @Summary
// @Schemes
// @Description get group all members
// @Tags group
// @Accept json
// @Produce json
// @Param body body group.getGroupAllMemberReq true "get group all members params"
// @Param token header string true "token"
// @Success 200 {object} group.GetGroupMemberListResp
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /group/get_group_all_member_list [post]
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

	// GetGroupMemberListResp struct
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

// groupResult struct
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

// @Summary
// @Schemes
// @Description get joined group list
// @Tags group
// @Accept json
// @Produce json
// @Param body body group.GetJoinedGroupListReq true "get joined group list params"
// @Param token header string true "token"
// @Success 200 {object} group.GetJoinedGroupListResp
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /group/get_joined_group_list [post]
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

	// GetJoinedGroupListResp struct
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

// Id2Result struct
type Id2Result struct {
	UId    string `json:"uid"`
	Result int32  `json:"result"`
}

// @Summary
// @Schemes
// @Description invite user to group
// @Tags group
// @Accept json
// @Produce json
// @Param body body group.InviteUserToGroupReq true "invite user to group params"
// @Param token header string true "token"
// @Success 200 {object} group.InviteUserToGroupResp
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /group/invite_user_to_group [post]
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

	// InviteUserToGroupResp struct
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
