package group

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"strings"

	pbGroup "Open_IM/pkg/proto/group"
	"github.com/gin-gonic/gin"
)

func GetGroups(c *gin.Context) {
	var (
		req cms_api_struct.GetGroupsRequest
		resp cms_api_struct.GetGroupsResponse
		reqPb pbGroup.GetGroupsReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroups(context.Background(), &reqPb)
	fmt.Println(respPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)

}

func GetGroup(c *gin.Context) {
	var (
		req cms_api_struct.GetGroupRequest
		resp cms_api_struct.GetGroupResponse
		reqPb pbGroup.GetGroupReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError("0", "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	utils.CopyStructFields(&reqPb, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroup(context.Background(), &reqPb)
	fmt.Println(respPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func CreateGroup(c *gin.Context) {
	var (
		req cms_api_struct.CreateGroupRequest
		resp cms_api_struct.CreateGroupResponse
		reqPb pbGroup.CreateGroupReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupInfo.GroupName = req.GroupName
	reqPb.
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.CreateGroup(context.Background(), &reqPb)
	fmt.Println(respPb)
	if err != nil {
		log.NewError("s", "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func SearchGroupsMember(c *gin.Context) {

}



func AddUsers(c *gin.Context) {

}

func InquireMember(c *gin.Context) {

}

func InquireGroup(c *gin.Context) {

}

func AddGroupMember(c *gin.Context) {

}

func AddMembers(c *gin.Context) {

}

func SetMaster(c *gin.Context) {

}

func BlockUser(c *gin.Context) {

}

func RemoveUser(c *gin.Context) {

}

func BanPrivateChat(c *gin.Context) {

}

func Withdraw(c *gin.Context) {

}

func SearchMessage(g *gin.Context) {

}
