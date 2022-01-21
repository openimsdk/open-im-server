package organization

import (
	"Open_IM_CMS/pkg/common/config"
	"Open_IM_CMS/pkg/errno"
	"Open_IM_CMS/pkg/etcdv3"
	commonProto "Open_IM_CMS/pkg/proto/common"
	proto "Open_IM_CMS/pkg/proto/organization"
	"Open_IM_CMS/pkg/req_resp"
	"Open_IM_CMS/test"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetStaffs(c *gin.Context) {
	var (
		req    req_resp.GetStaffsResponse
		resp   req_resp.GetStaffsResponse
		reqPb  commonProto.Pagination
		respPb *proto.GetStaffsResp
	)
	fmt.Println(resp, req)
	fakeData := test.GetStaffsResponseFake()
	etcdConn := etcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCMSApiOrganizationName)
	client := proto.NewOrganizationClient(etcdConn)
	fmt.Println(client, reqPb)
	respPb, err := client.GetStaffs(context.Background(), &reqPb)
	fmt.Println(respPb, err)
	fmt.Println(etcdConn)
	req_resp.RespHttp200(c, errno.RespOK, fakeData)
}

func GetOrganizations(c *gin.Context) {
	var (
		req  req_resp.GetOrganizationsResponse
		resp req_resp.GetStaffsResponse
	)
	fmt.Println(resp, req)
	fakeData := test.GetOrganizationsResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fakeData})
}

func GetSquads(c *gin.Context) {
	fakeData := test.GetSquadResponseFake()
	c.JSON(http.StatusOK, gin.H{"code": "0", "data": fakeData})
}

func AlterStaff(c *gin.Context) {

}

func AddOrganization(c *gin.Context) {

}

func InquireOrganization(g *gin.Context) {

}

func AlterOrganization(c *gin.Context) {

}

func DeleteOrganization(g *gin.Context) {

}

func GetOrganizationSquads(c *gin.Context) {

}

func AlterStaffsInfo(c *gin.Context) {

}

func AddChildOrganization(c *gin.Context) {

}
