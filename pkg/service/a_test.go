package service

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	_ user.UnimplementedUserServer
	_ group.UnimplementedGroupServer
)

func TestName1(t *testing.T) {
	cc := newStandaloneConn()
	user.RegisterUserServer(cc.Registry(), &user.UnimplementedUserServer{})
	group.RegisterGroupServer(cc.Registry(), &group.UnimplementedGroupServer{})
	ctx := context.Background()
	resp, err := user.NewUserClient(cc).GetUserStatus(ctx, &user.GetUserStatusReq{UserID: "imAdmin", UserIDs: []string{"10000", "20000"}})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp)
}

func newStandaloneConn() *standaloneConn {
	return &standaloneConn{
		registry:   newStandaloneRegistry(),
		serializer: NewProtoSerializer(),
	}
}

type standaloneConn struct {
	registry   *standaloneRegistry
	serializer Serializer
}

func (x *standaloneConn) Registry() grpc.ServiceRegistrar {
	return x.registry
}

func (x *standaloneConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	handler := x.registry.getMethod(method)
	if handler == nil {
		return fmt.Errorf("service %s not found", method)
	}
	resp, err := handler(ctx, args, nil)
	if err != nil {
		return err
	}
	tmp, err := x.serializer.Marshal(resp)
	if err != nil {
		return err
	}
	return x.serializer.Unmarshal(tmp, reply)
}

func (x *standaloneConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Errorf(codes.Unimplemented, "method stream not implemented")
}

type serverHandler func(ctx context.Context, req any, interceptor grpc.UnaryServerInterceptor) (any, error)

func newStandaloneRegistry() *standaloneRegistry {
	return &standaloneRegistry{
		methods:    make(map[string]serverHandler),
		serializer: NewProtoSerializer(),
	}
}

type standaloneRegistry struct {
	lock       sync.RWMutex
	methods    map[string]serverHandler
	serializer Serializer
}

func (x *standaloneConn) emptyDec(req any) error {
	return nil
}

func (x *standaloneRegistry) RegisterService(desc *grpc.ServiceDesc, impl any) {
	x.lock.Lock()
	defer x.lock.Unlock()
	for i := range desc.Methods {
		method := desc.Methods[i]
		name := fmt.Sprintf("/%s/%s", desc.ServiceName, method.MethodName)
		if _, ok := x.methods[name]; ok {
			panic(fmt.Errorf("service %s already registered, method %s", desc.ServiceName, method.MethodName))
		}
		x.methods[name] = func(ctx context.Context, req any, interceptor grpc.UnaryServerInterceptor) (any, error) {
			return method.Handler(impl, ctx, func(in any) error {
				tmp, err := x.serializer.Marshal(req)
				if err != nil {
					return err
				}
				return x.serializer.Unmarshal(tmp, in)
			}, interceptor)
		}
	}
}

func (x *standaloneRegistry) getMethod(name string) serverHandler {
	x.lock.RLock()
	defer x.lock.RUnlock()
	return x.methods[name]
}

type Serializer interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}

func NewProtoSerializer() Serializer {
	return protoSerializer{}
}

type protoSerializer struct{}

func (protoSerializer) Marshal(in any) ([]byte, error) {
	return proto.Marshal(in.(proto.Message))
}

func (protoSerializer) Unmarshal(b []byte, out any) error {
	return proto.Unmarshal(b, out.(proto.Message))
}
