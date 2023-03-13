package tracelog

import (
	"OpenIM/pkg/utils"
	"context"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"

	//"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

const TraceLogKey = "tracelog"

func SetFuncInfos(c context.Context, rootFuncName string, operationID string) context.Context {
	req := &FuncInfos{RootFuncName: rootFuncName, Funcs: &[]FuncInfo{}}
	ctx := context.WithValue(c, TraceLogKey, req)
	SetOperationID(ctx, operationID)
	return ctx
}

func NewCtx(rootFuncName string, operationID string) context.Context {
	c := context.Background()
	req := &FuncInfos{RootFuncName: rootFuncName, Funcs: &[]FuncInfo{}}
	ctx := context.WithValue(c, TraceLogKey, req)
	SetOperationID(ctx, operationID)
	return ctx
}

func SetOperationID(ctx context.Context, operationID string) {
	ctx.Value(TraceLogKey).(*FuncInfos).OperationID = operationID
}

func GetOperationID(ctx context.Context) string {
	f, ok := ctx.Value(TraceLogKey).(*FuncInfos)
	if ok {
		return f.OperationID
	}
	return utils.GetFuncName(1)
}

func GetOpUserID(ctx context.Context) string {
	s, _ := ctx.Value("opUserID").(string)
	return s
}

func Unwrap(err error) error {
	for err != nil {
		unwrap, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			break
		}
		err = unwrap.Unwrap()
	}
	return err
}

type FuncInfos struct {
	RootFuncName string
	OperationID  string
	Funcs        *[]FuncInfo
	GinCtx       *gin.Context
}

type FuncInfo struct {
	FuncName string
	Args     Args
	Err      error
	LogLevel logrus.Level
	File     string
}

type Args map[string]interface{}

func (a Args) String() string {
	var s string
	var hasElement bool
	for k, v := range a {
		if !hasElement {
			s += "{"
			hasElement = true
		}
		s += fmt.Sprintf("%s: %v", k, v)
	}
	if hasElement {
		s += "}"
	}
	return s
}

func SetCtxDebug(ctx context.Context, funcName string, err error, args ...interface{}) {
	SetContextInfo(ctx, funcName, logrus.DebugLevel, err, args)
}

func SetCtxInfo(ctx context.Context, funcName string, err error, args ...interface{}) {
	SetContextInfo(ctx, funcName, logrus.InfoLevel, err, args)
}

func SetCtxWarn(ctx context.Context, funcName string, err error, args ...interface{}) {
	SetContextInfo(ctx, funcName, logrus.WarnLevel, err, args)
}

func SetContextInfo(ctx context.Context, funcName string, logLevel logrus.Level, err error, args ...interface{}) {
	var t *FuncInfos
	var ok bool
	if t, ok = ctx.Value(TraceLogKey).(*FuncInfos); !ok {
		return
	}
	var funcInfo FuncInfo
	funcInfo.Args = make(map[string]interface{})
	argsHandle(args, funcInfo.Args)
	funcInfo.FuncName = funcName
	funcInfo.Err = err
	funcInfo.LogLevel = logLevel
	_, file, line, _ := runtime.Caller(3)
	var s string
	i := strings.SplitAfter(file, "/")
	if len(i) > 3 {
		s = i[len(i)-3] + i[len(i)-2] + i[len(i)-1] + ":" + utils.IntToString(line)
	}
	funcInfo.File = s
	*t.Funcs = append(*t.Funcs, funcInfo)
}

func SetRpcReqInfo(ctx context.Context, funcName string, req string) {
	t := ctx.Value(TraceLogKey).(*FuncInfos)
	var funcInfo FuncInfo
	funcInfo.Args = make(map[string]interface{})
	var args []interface{}
	args = append(args, " rpc req ", req)
	argsHandle(args, funcInfo.Args)
	funcInfo.FuncName = funcName
	*t.Funcs = append(*t.Funcs, funcInfo)
}

func SetRpcRespInfo(ctx context.Context, funcName string, resp string) {
	t := ctx.Value(TraceLogKey).(*FuncInfos)
	var funcInfo FuncInfo
	funcInfo.Args = make(map[string]interface{})
	var args []interface{}
	args = append(args, " rpc resp ", resp)
	argsHandle(args, funcInfo.Args)
	funcInfo.FuncName = funcName
	*t.Funcs = append(*t.Funcs, funcInfo)
}

func argsHandle(args []interface{}, fields map[string]interface{}) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[fmt.Sprintf("%v", args[i])] = fmt.Sprintf("%+v", args[i+1])
		} else {
			fields[fmt.Sprintf("%v", args[i])] = ""
		}
	}
}
