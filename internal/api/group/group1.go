package group

import (
	"OpenIM/internal/a2r"
	"OpenIM/pkg/apistruct"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/proto/group"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
)

func _() {
	context.Background()
}

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
	a2r.Call(&apistruct.KickGroupMemberReq{}, &apistruct.KickGroupMemberResp{}, group.GroupClient.KickGroupMember, g.getGroupClient, c, nil, nil)
}

func (g *Group) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call1(group.GroupClient.GetGroupMembersInfo, g.getGroupClient, c)
}
