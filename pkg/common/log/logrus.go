package log

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"bufio"
	"context"
	"runtime"
	"strings"

	//"bufio"
	"fmt"
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var logger *Logger
var ctxLogger *Logger

type Logger struct {
	*logrus.Logger
	Pid  int
	Type string
}

func init() {
	logger = loggerInit("")
	ctxLogger = ctxLoggerInit("")
}

func NewPrivateLog(moduleName string) {
	logger = loggerInit(moduleName)
	ctxLogger = ctxLoggerInit(moduleName)
}

func ctxLoggerInit(moduleName string) *Logger {
	var ctxLogger = logrus.New()
	ctxLogger.SetLevel(logrus.Level(config.Config.Log.RemainLogLevel))
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err.Error())
	}
	writer := bufio.NewWriter(src)
	ctxLogger.SetOutput(writer)
	ctxLogger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	if config.Config.Log.ElasticSearchSwitch {
		ctxLogger.AddHook(newEsHook(moduleName))
	}
	//Log file segmentation hook
	hook := NewLfsHook(time.Duration(config.Config.Log.RotationTime)*time.Hour, config.Config.Log.RemainRotationCount, moduleName)
	ctxLogger.AddHook(hook)
	return &Logger{
		ctxLogger,
		os.Getpid(),
		"ctxLogger",
	}
}

func loggerInit(moduleName string) *Logger {
	var logger = logrus.New()
	//All logs will be printed
	logger.SetLevel(logrus.Level(config.Config.Log.RemainLogLevel))
	//Close std console output
	//os.O_WRONLY | os.O_CREATE | os.O_APPEND
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err.Error())
	}
	writer := bufio.NewWriter(src)
	logger.SetOutput(writer)
	// logger.SetOutput(os.Stdout)
	//Log Console Print Style Setting

	logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})

	//File name and line number display hook
	logger.AddHook(newFileHook())

	//Send logs to elasticsearch hook
	if config.Config.Log.ElasticSearchSwitch {
		logger.AddHook(newEsHook(moduleName))
	}
	//Log file segmentation hook
	hook := NewLfsHook(time.Duration(config.Config.Log.RotationTime)*time.Hour, config.Config.Log.RemainRotationCount, moduleName)
	logger.AddHook(hook)
	return &Logger{
		logger,
		os.Getpid(),
		"",
	}
}
func NewLfsHook(rotationTime time.Duration, maxRemainNum uint, moduleName string) logrus.Hook {
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.InfoLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.WarnLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.ErrorLevel: initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
	}, &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID"},
	})
	return lfsHook
}
func initRotateLogs(rotationTime time.Duration, maxRemainNum uint, level string, moduleName string) *rotatelogs.RotateLogs {
	if moduleName != "" {
		moduleName = moduleName + "."
	}
	writer, err := rotatelogs.New(
		config.Config.Log.StorageLocation+moduleName+level+"."+"%Y-%m-%d",
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	if err != nil {
		panic(err.Error())
	} else {
		return writer
	}
}

func Info(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Infoln(args)
}

func Error(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Errorln(args)
}

func Debug(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Debugln(args)
}

//Deprecated
func Warning(token, OperationID, format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"PID":         logger.Pid,
		"OperationID": OperationID,
	}).Warningf(format, args...)

}

//Deprecated
func InfoByArgs(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Infof(format, args)
}

//Deprecated
func ErrorByArgs(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Errorf(format, args...)
}

//Print log information in k, v format,
//kv is best to appear in pairs. tipInfo is the log prompt information for printing,
//and kv is the key and value for printing.
//Deprecated
func InfoByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Info(tipInfo)
}

//Deprecated
func ErrorByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Error(tipInfo)
}

//Deprecated
func DebugByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Debug(tipInfo)
}

//Deprecated
func WarnByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Warn(tipInfo)
}

//internal method
func argsHandle(OperationID string, fields logrus.Fields, args []interface{}) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[fmt.Sprintf("%v", args[i])] = args[i+1]
		} else {
			fields[fmt.Sprintf("%v", args[i])] = ""
		}
	}
	fields["OperationID"] = OperationID
	fields["PID"] = logger.Pid
}
func NewInfo(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Infoln(args)
}
func NewError(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Errorln(args)
}
func NewDebug(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Debugln(args)
}
func NewWarn(OperationID string, args ...interface{}) {
	logger.WithFields(logrus.Fields{
		"OperationID": OperationID,
		"PID":         logger.Pid,
	}).Warnln(args)
}

func ShowLog(ctx context.Context) {
	t := ctx.Value(trace_log.TraceLogKey).(*trace_log.ApiInfo)
	OperationID := trace_log.GetOperationID(ctx)
	if ctx.Value(trace_log.TraceLogKey).(*trace_log.ApiInfo).GinCtx != nil {
		ctxLogger.WithFields(logrus.Fields{
			"OperationID": OperationID,
			"PID":         ctxLogger.Pid,
		}).Infoln("api: ", t.ApiName)
	} else {
		ctxLogger.WithFields(logrus.Fields{
			"OperationID": OperationID,
			"PID":         ctxLogger.Pid,
		}).Infoln("rpc: ", t.ApiName)
	}
	for _, v := range *t.Funcs {
		if v.File == "" {
			_, file, line, _ := runtime.Caller(1)
			var s string
			i := strings.SplitAfter(file, "/")
			if len(i) > 3 {
				s = i[len(i)-3] + i[len(i)-2] + i[len(i)-1] + ":" + utils.IntToString(line)
			}
			v.File = s
		}
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
	t := ctx.Value(trace_log.TraceLogKey).(*trace_log.ApiInfo)
	OperationID := trace_log.GetOperationID(ctx)
	for _, v := range *t.Funcs {
		logger.WithFields(logrus.Fields{
			"OperationID": OperationID,
			"PID":         logger.Pid,
		}).Infoln(v.Args, args)
	}
}

func DebugWithCtx(ctx context.Context, args ...interface{}) {
	t := ctx.Value(trace_log.TraceLogKey).(*trace_log.ApiInfo)
	OperationID := trace_log.GetOperationID(ctx)
	for _, v := range *t.Funcs {
		logger.WithFields(logrus.Fields{
			"OperationID": OperationID,
			"PID":         logger.Pid,
		}).Debugln(v.Args, args)
	}
}

func ErrorWithCtx(ctx context.Context, args ...interface{}) {
	t := ctx.Value(trace_log.TraceLogKey).(*trace_log.ApiInfo)
	OperationID := trace_log.GetOperationID(ctx)
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
	t := ctx.Value(trace_log.TraceLogKey).(*trace_log.ApiInfo)
	OperationID := trace_log.GetOperationID(ctx)
	for _, v := range *t.Funcs {
		logger.WithFields(logrus.Fields{
			"OperationID": OperationID,
			"PID":         logger.Pid,
		}).Warnln(v.Args, args)
	}
}
