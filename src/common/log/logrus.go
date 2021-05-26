package log

import (
	"Open_IM/src/common/config"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var logger *Logger

type Logger struct {
	*logrus.Logger
	Pid int
}

func init() {
	logger = loggerInit("")

}
func NewPrivateLog(moduleName string) {
	logger = loggerInit(moduleName)
}

func loggerInit(moduleName string) *Logger {
	var logger = logrus.New()
	//All logs will be printed
	logger.SetLevel(logrus.TraceLevel)
	//Log Style Setting
	logger.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		HideKeys:        false,
		FieldsOrder:     []string{"PID"},
	})
	//File name and line number display hook
	logger.AddHook(newFileHook())

	//Send logs to elasticsearch hook
	if config.Config.Log.ElasticSearchSwitch == true {
		logger.AddHook(newEsHook(moduleName))
	}
	//Log file segmentation hook
	hook := NewLfsHook(config.Config.Log.StorageLocation+time.Now().Format("2006-01-02")+".log", 0, 5, moduleName)
	logger.AddHook(hook)
	return &Logger{
		logger,
		os.Getpid(),
	}
}
func NewLfsHook(logName string, rotationTime time.Duration, maxRemainNum uint, moduleName string) logrus.Hook {
	var fileNameSuffix string
	if GetCurrentTimestamp() >= GetCurDayZeroTimestamp() && GetCurrentTimestamp() <= GetCurDayHalfTimestamp() {
		fileNameSuffix = time.Now().Format("2006-01-02") + ".log"
	} else {
		fileNameSuffix = time.Now().Format("2006-01-02") + ".log"
	}
	writer, err := rotatelogs.New(
		logName,
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	if err != nil {
		panic(err)
	}
	writeInfo, err := rotatelogs.New(
		config.Config.Log.StorageLocation+moduleName+"/info."+fileNameSuffix,
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	writeError, err := rotatelogs.New(
		config.Config.Log.StorageLocation+moduleName+"/error."+fileNameSuffix,
		rotatelogs.WithRotationTime(time.Minute),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	writeDebug, err := rotatelogs.New(
		config.Config.Log.StorageLocation+moduleName+"/debug."+fileNameSuffix,
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	writeWarn, err := rotatelogs.New(
		config.Config.Log.StorageLocation+moduleName+"/warn."+fileNameSuffix,
		rotatelogs.WithRotationTime(time.Minute),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	if err != nil {
		panic(err)
	}
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writeDebug,
		logrus.InfoLevel:  writeInfo,
		logrus.WarnLevel:  writeWarn,
		logrus.ErrorLevel: writeError,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		HideKeys:        false,
		FieldsOrder:     []string{"PID"},
	})

	return lfsHook
}

func Info(token, OperationID, format string, args ...interface{}) {
	if token == "" && OperationID == "" {
		logger.WithFields(logrus.Fields{}).Infof(format, args...)
	} else {
		logger.WithFields(logrus.Fields{
			"token":       token,
			"OperationID": OperationID,
		}).Infof(format, args...)
	}
}

func Error(token, OperationID, format string, args ...interface{}) {
	if token == "" && OperationID == "" {
		logger.WithFields(logrus.Fields{}).Errorf(format, args...)
	} else {
		logger.WithFields(logrus.Fields{
			"token":       token,
			"OperationID": OperationID,
		}).Errorf(format, args...)
	}
}

func Debug(token, OperationID, format string, args ...interface{}) {
	if token == "" && OperationID == "" {
		logger.WithFields(logrus.Fields{}).Debugf(format, args...)
	} else {
		logger.WithFields(logrus.Fields{
			"token":       token,
			"OperationID": OperationID,
		}).Debugf(format, args...)
	}
}

func Warning(token, OperationID, format string, args ...interface{}) {
	if token == "" && OperationID == "" {
		logger.WithFields(logrus.Fields{}).Warningf(format, args...)
	} else {
		logger.WithFields(logrus.Fields{
			"token":       token,
			"OperationID": OperationID,
		}).Warningf(format, args...)
	}
}

func InfoByArgs(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Infof(format, args)
}

func ErrorByArgs(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Errorf(format, args...)
}

//Print log information in k, v format,
//kv is best to appear in pairs. tipInfo is the log prompt information for printing,
//and kv is the key and value for printing.
func InfoByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Info(tipInfo)
}
func ErrorByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Error(tipInfo)
}
func DebugByKv(tipInfo, OperationID string, args ...interface{}) {
	fields := make(logrus.Fields)
	argsHandle(OperationID, fields, args)
	logger.WithFields(fields).Debug(tipInfo)
}
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
	fields["operationID"] = OperationID
	fields["PID"] = logger.Pid
}
