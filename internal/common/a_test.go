package common

import (
	"OpenIM/internal/api2rpc"
	"OpenIM/pkg/proto/group"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
	"time"
)

type AReq struct {
}

type AResp struct {
}

func KickGroupMember(c *gin.Context) {
	// 默认 全部自动
	api2rpc.NewRpc(api2rpc.NewGin[AReq, AResp](c), group.NewGroupClient, group.GroupClient.KickGroupMember).Call()

	//// 可以自定义编辑请求和响应
	//a := NewRpc(NewGin[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp](c), group.NewGroupClient, group.GroupClient.KickGroupMember)
	//a.Before(func(apiReq *apistruct.KickGroupMemberReq, rpcReq *group.KickGroupMemberReq, bind func() error) error {
	//	return bind()
	//}).After(func(rpcResp *group.KickGroupMemberResp, apiResp *apistruct.KickGroupMemberResp, bind func() error) error {
	//	return bind()
	//}).Name("group").Call()
}

//func getInterfaceName(handler PackerHandler) string {
//	funcInfo := runtime.FuncForPC(reflect.ValueOf(handler).Pointer())
//	name := funcInfo.Name()
//	names := strings.Split(name, "/")
//	if len(names) == 0 {
//		return ""
//	}
//
//	return names[len(names)-1]
//}

func TestName(t *testing.T) {
	n := 100000000
	start := time.Now().UnixNano()
	for i := 0; i < n; i++ {
		var val group.GroupClient
		reflect.TypeOf(&val).Elem().String()
	}
	end := time.Now().UnixNano()
	fmt.Println(time.Duration(end-start) / time.Duration(n))
}
