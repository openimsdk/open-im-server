package log

import (
	"context"
	"fmt"
)

type ZkLogger struct{}

func (l *ZkLogger) Printf(format string, a ...interface{}) {
	ZInfo(context.Background(), "zookeeper output", "msg", fmt.Sprintf(format, a...))
}
