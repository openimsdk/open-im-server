package msggateway

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"

	"OpenIM/pkg/statistics"
	"fmt"
	"sync"

	prome "OpenIM/pkg/common/prome"

	"github.com/go-playground/validator/v10"
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
	rwLock = new(sync.RWMutex)
	validate = validator.New()
	statistics.NewStatistics(&sendMsgAllCount, config.Config.ModuleName.LongConnSvrName, fmt.Sprintf("%d second recv to msg_gateway sendMsgCount", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	statistics.NewStatistics(&userCount, config.Config.ModuleName.LongConnSvrName, fmt.Sprintf("%d second add user conn", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	ws.onInit(wsPort)
	rpcSvr.onInit(rpcPort)
	initPrometheus()
}

func Run(prometheusPort int) {
	go ws.run()
	go rpcSvr.run()
	go func() {
		err := prome.StartPromeSrv(prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
}
