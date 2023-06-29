package zookeeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"google.golang.org/grpc"
)

func newClientConnInterface(cc grpc.ClientConnInterface) grpc.ClientConnInterface {
	return &clientConnInterface{cc: cc}
}

type clientConnInterface struct {
	cc grpc.ClientConnInterface
}

func (c *clientConnInterface) callOptionToString(opts []grpc.CallOption) string {
	arr := make([]string, 0, len(opts)+1)
	arr = append(arr, fmt.Sprintf("opts len: %d", len(opts)))
	for i, opt := range opts {
		arr = append(arr, fmt.Sprintf("[%d:%T]", i, opt))
	}
	return strings.Join(arr, ", ")
}

func (c *clientConnInterface) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	log.ZDebug(ctx, "grpc.ClientConnInterface.Invoke in", "method", method, "args", args, "reply", reply, "opts", c.callOptionToString(opts))
	if err := c.cc.Invoke(ctx, method, args, reply, opts...); err != nil {
		log.ZError(ctx, "grpc.ClientConnInterface.Invoke error", err, "method", method, "args", args, "reply", reply)
		return err
	}
	log.ZDebug(ctx, "grpc.ClientConnInterface.Invoke success", "method", method, "args", args, "reply", reply)
	return nil
}

func (c *clientConnInterface) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.ZDebug(ctx, "grpc.ClientConnInterface.NewStream in", "desc", desc, "method", method, "opts", c.callOptionToString(opts))
	cs, err := c.cc.NewStream(ctx, desc, method, opts...)
	if err != nil {
		log.ZError(ctx, "grpc.ClientConnInterface.NewStream error", err, "desc", desc, "method", method, "opts", len(opts))
		return nil, err
	}
	log.ZDebug(ctx, "grpc.ClientConnInterface.NewStream success", "desc", desc, "method", method, "opts", len(opts))
	return cs, nil
}
