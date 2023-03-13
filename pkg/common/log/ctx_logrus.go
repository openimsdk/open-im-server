package log

import (
	"OpenIM/pkg/common/tracelog"
	"context"
	"github.com/sirupsen/logrus"
)

func ShowLog(ctx context.Context) {
	t := ctx.Value(tracelog.TraceLogKey).(*tracelog.FuncInfos)
	OperationID := tracelog.GetOperationID(ctx)
	for _, v := range *t.Funcs {
		if v.Err != nil {
			ctxLogger.WithFields(logrus.Fields{
				"OperationID": OperationID,
				"PID":         ctxLogger.Pid,
				"FilePath":    v.File,
			}).Errorln("func: ", v.FuncName, " args: ", v.Args, v.Err.Error())
		} else {
			switch v.LogLevel {
			case logrus.InfoLevel:
				ctxLogger.WithFields(logrus.Fields{
					"OperationID": OperationID,
					"PID":         ctxLogger.Pid,
					"FilePath":    v.File,
				}).Infoln("func: ", v.FuncName, " args: ", v.Args)
			case logrus.DebugLevel:
				ctxLogger.WithFields(logrus.Fields{
					"OperationID": OperationID,
					"PID":         ctxLogger.Pid,
					"FilePath":    v.File,
				}).Debugln("func: ", v.FuncName, " args: ", v.Args)
			case logrus.WarnLevel:
				ctxLogger.WithFields(logrus.Fields{
					"OperationID": OperationID,
					"PID":         ctxLogger.Pid,
					"FilePath":    v.File,
				}).Warnln("func: ", v.FuncName, " args: ", v.Args)
			}
		}
	}
}

func InfoWithCtx(ctx context.Context, args ...interface{}) {
	t := ctx.Value(tracelog.TraceLogKey).(*tracelog.FuncInfos)
	OperationID := tracelog.GetOperationID(ctx)
	for _, v := range *t.Funcs {
		logger.WithFields(logrus.Fields{
			"OperationID": OperationID,
			"PID":         logger.Pid,
		}).Infoln(v.Args, args)
	}
}

func DebugWithCtx(ctx context.Context, args ...interface{}) {
	t := ctx.Value(tracelog.TraceLogKey).(*tracelog.FuncInfos)
	OperationID := tracelog.GetOperationID(ctx)
	for _, v := range *t.Funcs {
		logger.WithFields(logrus.Fields{
			"OperationID": OperationID,
			"PID":         logger.Pid,
		}).Debugln(v.Args, args)
	}
}

func ErrorWithCtx(ctx context.Context, args ...interface{}) {
	t := ctx.Value(tracelog.TraceLogKey).(*tracelog.FuncInfos)
	OperationID := tracelog.GetOperationID(ctx)
	for _, v := range *t.Funcs {
		if v.Err != nil {
			logger.WithFields(logrus.Fields{
				"OperationID": OperationID,
				"PID":         logger.Pid,
			}).Errorln(v.Err, v.Args, args)
		}
	}
}

func WarnWithCtx(ctx context.Context, args ...interface{}) {
	t := ctx.Value(tracelog.TraceLogKey).(*tracelog.FuncInfos)
	OperationID := tracelog.GetOperationID(ctx)
	for _, v := range *t.Funcs {
		logger.WithFields(logrus.Fields{
			"OperationID": OperationID,
			"PID":         logger.Pid,
		}).Warnln(v.Args, args)
	}
}
