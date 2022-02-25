package main

import (
	"Open_IM/cmd/rpc/open_im_msg/widget"
	rpcChat "Open_IM/internal/rpc/msg"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 10300, "rpc listening port")
	flag.Parse()
	rpcServer := rpcChat.NewRpcChatServer(*rpcPort)
	// register widgets

	// mock 注册发送前的拦截器
	rpcServer.UseWidgetBeforSend(widget.MockBeforeSendHandler)

	//
	rpcServer.Run()
}
