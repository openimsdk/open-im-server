package log

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
)

func NewLfsHook(rotationTime time.Duration, maxRemainNum uint, moduleName string) logrus.Hook {
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.InfoLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.WarnLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
		logrus.ErrorLevel: initRotateLogs(rotationTime, maxRemainNum, "all", moduleName),
	}, &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        false,
		FieldsOrder:     []string{"PID", "FilePath", "OperationID", "Msg"},
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
