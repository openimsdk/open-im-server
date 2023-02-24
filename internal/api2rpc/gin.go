package api2rpc

import (
	"context"
	"github.com/gin-gonic/gin"
)

//func KickGroupMember(c *gin.Context) {
//	// 默认 全部自动
//	//var api ApiBind[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp] = NewGin[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp](c)
//	var api ApiBind[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp] = nil
//	var client func(conn *grpc.ClientConn) group.GroupClient = nil
//	var rpc func(ctx context.Context, in *group.KickGroupMemberReq, opts ...grpc.CallOption) (*group.KickGroupMemberResp, error) = nil
//	//NewRpc(api, client, rpc).Name("group").Call()
//	NewRpc(api, client, rpc).Name("group").Call()
//
//	// 可以自定义编辑请求和响应
//	//a := NewRpc(NewGin[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp](c), "", group.NewGroupClient, group.GroupClient.KickGroupMember)
//	//a.Before(func(apiReq *apistruct.KickGroupMemberReq, rpcReq *group.KickGroupMemberReq, bind func() error) error {
//	//	return bind()
//	//}).After(func(rpcResp *group.KickGroupMemberResp, apiResp *apistruct.KickGroupMemberResp, bind func() error) error {
//	//	return bind()
//	//}).Execute()
//}
//

func NewGin[A, B any](c *gin.Context) ApiBind[A, B] {
	return &ginApiBind[A, B]{
		c: c,
	}
}

type ginApiBind[A, B any] struct {
	c *gin.Context
}

func (g *ginApiBind[A, B]) OperationID() string {
	return g.c.GetHeader("operationID")
}

func (g *ginApiBind[A, B]) OpUserID() (string, error) {
	return "", nil
}

func (g *ginApiBind[A, B]) Bind(a *A) error {
	return g.c.BindJSON(a)
}

func (g *ginApiBind[A, B]) Resp(resp *B, err error) {
	if err == nil {
		g.Write(resp)
	} else {
		g.Error(err)
	}
}

func (g *ginApiBind[A, B]) Error(err error) {
	//TODO implement me
}

func (g *ginApiBind[A, B]) Write(b *B) {
	//TODO implement me
}

func (g *ginApiBind[A, B]) Context() context.Context {
	return g.c
}
