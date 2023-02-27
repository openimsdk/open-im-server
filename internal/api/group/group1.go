package group

import (
	"OpenIM/internal/api2rpc"
	"OpenIM/pkg/apistruct"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/proto/group"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
)

type Group struct {
	zk *openKeeper.ZkClient
}

func (g *Group) getGroupClient() (group.GroupClient, error) {
	conn, err := g.zk.GetConn(config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	return group.NewGroupClient(conn), nil
}

func (g *Group) KickGroupMember(c *gin.Context) {
	api2rpc.New(&apistruct.KickGroupMemberReq{}, &apistruct.KickGroupMemberResp{}, group.GroupClient.KickGroupMember).Call(c, g.getGroupClient)
}

func (g *Group) GetGroupMembersInfo(c *gin.Context) {
	api2rpc.New(&apistruct.GetGroupMembersInfoReq{}, &apistruct.GetGroupMembersInfoResp{}, group.GroupClient.GetGroupMembersInfo).Call(c, g.getGroupClient)
}

//func GetGroupMembersInfo(c *gin.Context) {
//	api2rpc.New(&apistruct.GetGroupMembersInfoReq{}, &apistruct.GetGroupMembersInfoResp{}, group.GroupClient.GetGroupMembersInfo).Call(c, g.getGroupClient)
//}

//func New[A, B, C, D any, E any](apiReq *A, apiResp *B, rpc func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)) {
//
//}
