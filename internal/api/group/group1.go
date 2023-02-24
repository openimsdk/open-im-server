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
	api2rpc.Rpc(
		&apistruct.KickGroupMemberReq{},
		&apistruct.KickGroupMemberResp{},
		group.GroupClient.KickGroupMember,
	).Must(c, g.getGroupClient).Call()
}

//func (g *Group) KickGroupMember1(c *gin.Context) {
//	var fn func(client group.GroupClient, ctx context.Context, in *group.KickGroupMemberReq, opts ...grpc.CallOption) (*group.KickGroupMemberResp, error) = group.GroupClient.KickGroupMember
//	api2rpc.Rpc(&apistruct.KickGroupMemberReq{}, &apistruct.KickGroupMemberResp{}, fn).Must(c, g.getGroupClient).Call()
//}
