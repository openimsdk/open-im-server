package group

import (
	"OpenIM/internal/a2r"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/proto/group"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
)

var _ context.Context = nil // 解决goland编辑器bug

func NewGroup(zk *openKeeper.ZkClient) *Group {
	return &Group{zk: zk}
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
	a2r.Call(group.GroupClient.KickGroupMember, g.getGroupClient, c)
}

func (g *Group) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMembersInfo, g.getGroupClient, c)
}
