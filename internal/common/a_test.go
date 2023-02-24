package common

import (
	"OpenIM/internal/api2rpc"
	"OpenIM/pkg/proto/group"
	"github.com/gin-gonic/gin"
	"testing"
)

type AReq struct {
}

type AResp struct {
}

func KickGroupMember(c *gin.Context) {
	// 默认 全部自动
	api2rpc.NewRpc(api2rpc.NewGin[AReq, AResp](c), group.NewGroupClient, group.GroupClient.KickGroupMember).Name("group").Call()

	//// 可以自定义编辑请求和响应
	//a := NewRpc(NewGin[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp](c), group.NewGroupClient, group.GroupClient.KickGroupMember)
	//a.Before(func(apiReq *apistruct.KickGroupMemberReq, rpcReq *group.KickGroupMemberReq, bind func() error) error {
	//	return bind()
	//}).After(func(rpcResp *group.KickGroupMemberResp, apiResp *apistruct.KickGroupMemberResp, bind func() error) error {
	//	return bind()
	//}).Name("group").Call()
}

func TestName(t *testing.T) {
	KickGroupMember(nil)
}
