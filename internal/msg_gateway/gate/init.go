package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	rwLock   *sync.RWMutex
	validate *validator.Validate
	ws       WServer
	rpcSvr   RPCServer
)

func Init(rpcPort, wsPort int) {
	//log initialization
	log.NewPrivateLog(config.Config.ModuleName.LongConnSvrName)
	rwLock = new(sync.RWMutex)
	validate = validator.New()
	ws.onInit(wsPort)
	rpcSvr.onInit(rpcPort)
}

func Run() {
	go ws.run()
	go rpcSvr.run()
}
