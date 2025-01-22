package service

import (
	"fmt"

	"google.golang.org/grpc"
)

type GrpcServer struct {
}

func (x *GrpcServer) RegisterService(desc *grpc.ServiceDesc, impl any) {
	fmt.Println("RegisterService", desc, impl)
}
