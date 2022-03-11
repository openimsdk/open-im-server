/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/8 14:35).
 */
package main

import (
	"flag"
	"fmt"
	"open_im_sdk/sdk_struct"

	//"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/ws_wrapper/utils"
	"open_im_sdk/ws_wrapper/ws_local_server"
	"runtime"
	"sync"
)

func main() {
	var sdkWsPort, openIMApiPort, openIMWsPort *int
	var openIMWsAddress, openIMApiAddress *string
	//
	//openIMTerminalType := flag.String("terminal_type", "web", "different terminal types")

	sdkWsPort = flag.Int("sdk_ws_port", 30000, "openIMSDK ws listening port")
	openIMApiPort = flag.Int("openIM_api_port", 10000, "openIM api listening port")
	openIMWsPort = flag.Int("openIM_ws_port", 17778, "openIM ws listening port")
	flag.Parse()
	//switch *openIMTerminalType {
	//case "pc":
	//	openIMWsAddress = flag.String("openIM_ws_address", "web", "different terminal types")
	//	openIMApiAddress = flag.String("openIM_api_address", "web", "different terminal types")
	//	flag.Parse()
	//case "web":
	//	openIMApiPort = flag.Int("openIM_api_port", 0, "openIM api listening port")
	//	openIMWsPort = flag.Int("openIM_ws_port", 0, "openIM ws listening port")
	//	flag.Parse()
	//}
	APIADDR := "http://1.14.194.38:10000"
	WSADDR := "ws://1.14.194.38:17778"

	sysType := runtime.GOOS
	switch sysType {

	case "darwin":
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: *openIMApiAddress,
			WsAddr: *openIMWsAddress, Platform: utils.OSXPlatformID, DataDir: "./"})
	case "linux":
		//sdkDBDir:= flag.String("sdk_db_dir","","openIMSDK initialization path")
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: "http://" + utils.ApiServerIP + ":" + utils.IntToString(*openIMApiPort),
			WsAddr: "ws://" + utils.WsServerIP + ":" + utils.IntToString(*openIMWsPort), Platform: utils.WebPlatformID, DataDir: "../db/sdk/"})

	case "windows":
		//	sdkWsPort = flag.Int("sdk_ws_port", 7799, "openIM ws listening port")
		//flag.Parse()
		ws_local_server.InitServer(&sdk_struct.IMConfig{ApiAddr: APIADDR,
			WsAddr: WSADDR, Platform: utils.WebPlatformID, DataDir: "./"})
	default:
		fmt.Println("this os not support", sysType)

	}
	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog("sdk", 6)
	fmt.Println("ws server is starting")
	ws_local_server.WS.OnInit(*sdkWsPort)
	ws_local_server.WS.Run()
	wg.Wait()

}
