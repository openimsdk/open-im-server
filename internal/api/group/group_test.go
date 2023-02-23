package group

import (
	"OpenIM/pkg/apistruct"
	"OpenIM/pkg/proto/group"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"reflect"
	"testing"
)

type Ignore struct{}

func temp(client group.GroupClient, ctx context.Context, in *group.KickGroupMemberReq, opts ...grpc.CallOption) (*group.KickGroupMemberResp, error) {

	return nil, nil
}

type ApiBind[A, B any] interface {
	OperationID() string
	OpUserID() (string, error)
	Bind(*A) error
	Error(error)
	Write(*B)
}

func NewApiBind[A, B any](c *gin.Context) ApiBind[A, B] {
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

func (g *ginApiBind[A, B]) Error(err error) {
	//TODO implement me
}

func (g *ginApiBind[A, B]) Write(b *B) {
	//TODO implement me
}

func TestName(t *testing.T) {
	//var bind ApiBind[int, int]
	//NewRpc(bind, "", group.NewGroupClient, temp)

	var c *gin.Context
	NewRpc(NewApiBind[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp](c), "", group.NewGroupClient, group.GroupClient.KickGroupMember)

}

func KickGroupMember(c *gin.Context) {
	// 默认 全部自动
	NewRpc(NewApiBind[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp](c), "", group.NewGroupClient, group.GroupClient.KickGroupMember).Execute()
	// 可以自定义编辑请求和响应
	a := NewRpc(NewApiBind[apistruct.KickGroupMemberReq, apistruct.KickGroupMemberResp](c), "", group.NewGroupClient, group.GroupClient.KickGroupMember)
	a.Before(func(apiReq *apistruct.KickGroupMemberReq, rpcReq *group.KickGroupMemberReq, bind func() error) error {
		return bind()
	}).After(func(rpcResp *group.KickGroupMemberResp, apiResp *apistruct.KickGroupMemberResp, bind func() error) error {
		return bind()
	}).Execute()
}

func NewRpc[A, B any, C, D any, Z any](bind ApiBind[A, B], name string, client func(conn *grpc.ClientConn) Z, rpc func(client Z, ctx context.Context, req C, options ...grpc.CallOption) (D, error)) *RpcRun[A, B, C, D, Z] {
	return &RpcRun[A, B, C, D, Z]{
		bind:   bind,
		name:   name,
		client: client,
		rpc:    rpc,
	}
}

type RpcRun[A, B any, C, D any, Z any] struct {
	bind   ApiBind[A, B]
	name   string
	client func(conn *grpc.ClientConn) Z
	rpc    func(client Z, ctx context.Context, req C, options ...grpc.CallOption) (D, error)
	before func(apiReq *A, rpcReq C, bind func() error) error
	after  func(rpcResp D, apiResp *B, bind func() error) error
}

func (a *RpcRun[A, B, C, D, Z]) Before(fn func(apiReq *A, rpcReq C, bind func() error) error) *RpcRun[A, B, C, D, Z] {
	a.before = fn
	return a
}

func (a *RpcRun[A, B, C, D, Z]) After(fn func(rpcResp D, apiResp *B, bind func() error) error) *RpcRun[A, B, C, D, Z] {
	a.after = fn
	return a
}

func (a *RpcRun[A, B, C, D, Z]) execute() (*B, error) {
	userID, err := a.bind.OpUserID()
	if err != nil {
		return nil, err
	}
	opID := a.bind.OperationID()
	var rpcReq C // C type => *Struct
	rpcReq = reflect.New(reflect.TypeOf(rpcReq).Elem()).Interface().(C)

	return nil, nil
}

func (a *RpcRun[A, B, C, D, Z]) Execute() {

}

func GetGrpcConn(name string) (*grpc.ClientConn, error) {
	return nil, errors.New("todo")
}
