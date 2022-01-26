package group

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"strings"

	pbGroup "Open_IM/pkg/proto/group"

	"github.com/gin-gonic/gin"
)

func GetGroups(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupsRequest
		resp  cms_api_struct.GetGroupsResponse
		reqPb pbGroup.GetGroupsReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.Pagination = &commonPb.RequestPagination{}
	utils.CopyStructFields(&reqPb.Pagination, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroups(context.Background(), &reqPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	for _, v := range respPb.GroupInfo {
		resp.Groups = append(resp.Groups, cms_api_struct.GroupResponse{
			GroupName:        v.GroupName,
			GroupID:          v.GroupID,
			GroupMasterName:  v.OwnerUserID,
			GroupMasterId:    v.OwnerUserID,
			CreateTime:       (utils.UnixSecondToTime(int64(v.CreateTime))).String(),
			IsBanChat:        false,
			IsBanPrivateChat: false,
			ProfilePhoto:     v.FaceURL,
		})
	}
	resp.GroupNums = int(respPb.GroupNum)
	resp.CurrentPage = int(respPb.Pagination.PageNumber)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetGroup(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupRequest
		resp  cms_api_struct.GetGroupResponse
		reqPb pbGroup.GetGroupReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupName = req.GroupName
	fmt.Println(reqPb)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroup(context.Background(), &reqPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	for _, v := range respPb.GroupInfo {
		resp.Groups = append(resp.Groups, cms_api_struct.GroupResponse{
			GroupName:        v.GroupName,
			GroupID:          v.GroupID,
			GroupMasterName:  v.OwnerUserID,
			GroupMasterId:    v.OwnerUserID,
			CreateTime:       (utils.UnixSecondToTime(int64(v.CreateTime))).String(),
			IsBanChat:        false,
			IsBanPrivateChat: false,
			ProfilePhoto:     v.FaceURL,
		})
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func CreateGroup(c *gin.Context) {
	var (
		req   cms_api_struct.CreateGroupRequest
		resp  cms_api_struct.CreateGroupResponse
		reqPb pbGroup.CreateGroupReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupInfo.GroupName = req.GroupName
	reqPb.GroupInfo.CreatorUserID = ""
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.CreateGroup(context.Background(), &reqPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	fmt.Println(respPb)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func BanGroupChat(c *gin.Context) {
	var (
		req   cms_api_struct.BanGroupChatRequest
		reqPb pbGroup.BanGroupChatReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.BanGroupChat(context.Background(), &reqPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)

}

func BanPrivateChat(c *gin.Context) {
	var (
		req   cms_api_struct.BanPrivateChatRequest
		reqPb pbGroup.BanPrivateChatReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.BanPrivateChat(context.Background(), &reqPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func GetGroupsMember(c *gin.Context) {
	var (
		req cms_api_struct.GetGroupMembersRequest
		_   cms_api_struct.GetGroupMembersResponse
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
}

func InquireMember(c *gin.Context) {

}

func InquireGroup(c *gin.Context) {

}

func AddMembers(c *gin.Context) {

}

func RemoveUser(c *gin.Context) {

}

func Withdraw(c *gin.Context) {

}

func SearchMessage(g *gin.Context) {

}
