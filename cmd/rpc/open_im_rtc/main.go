package main

import (
	"Open_IM/internal/rpc/rtc"
	"Open_IM/pkg/common/log"
	rtcPb "Open_IM/pkg/proto/rtc"
	"Open_IM/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:11300")
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), err.Error())
	} // 创建 RPC 服务容器
	grpcServer := grpc.NewServer()
	rtcPb.RegisterRtcServiceServer(grpcServer, &rtc.RtcService{})

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.NewError("", utils.GetSelfFuncName(), err.Error())
	}
	log.NewInfo("", utils.GetSelfFuncName(), "start success")
}