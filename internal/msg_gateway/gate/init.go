package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"

	"Open_IM/pkg/statistics"
	"fmt"
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	rwLock              *sync.RWMutex
	validate            *validator.Validate
	ws                  WServer
	rpcSvr              RPCServer
	sendMsgAllCount     uint64
	sendMsgFailedCount  uint64
	sendMsgSuccessCount uint64
	userCount           uint64

	sendMsgAllCountLock sync.RWMutex
)

func Init(rpcPort, wsPort int) {
	//log initialization

	rwLock = new(sync.RWMutex)
	validate = validator.New()
	statistics.NewStatistics(&sendMsgAllCount, config.Config.ModuleName.LongConnSvrName, fmt.Sprintf("%d second recv to msg_gateway sendMsgCount", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	statistics.NewStatistics(&userCount, config.Config.ModuleName.LongConnSvrName, fmt.Sprintf("%d second add user conn", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	ws.onInit(wsPort)
	rpcSvr.onInit(rpcPort)
}

func Run() {
	go ws.run()
	go rpcSvr.run()
}
