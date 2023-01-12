package trace_log

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/utils"
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"runtime"
	"strings"

	//"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const TraceLogKey = "trace_log"

func NewCtx(c *gin.Context, api string) context.Context {
	req := &ApiInfo{ApiName: api, GinCtx: c, Funcs: &[]FuncInfo{}}
	return context.WithValue(c, TraceLogKey, req)
}

func NewCtx1(c *gin.Context, api, operationID string) context.Context {
	req := &ApiInfo{ApiName: api, GinCtx: c, OperationID: operationID, Funcs: &[]FuncInfo{}}
	return context.WithValue(c, TraceLogKey, req)
}

func NewRpcCtx(c context.Context, rpc string, operationID string) context.Context {
	req := &ApiInfo{ApiName: rpc, Funcs: &[]FuncInfo{}}
	ctx := context.WithValue(c, TraceLogKey, req)
	SetOperationID(ctx, operationID)
	return ctx
}

func SetOperationID(ctx context.Context, operationID string) {
	ctx.Value(TraceLogKey).(*ApiInfo).OperationID = operationID
}

func GetOperationID(ctx context.Context) string {
	return ctx.Value(TraceLogKey).(*ApiInfo).OperationID
}

//func ShowLog(ctx context.Context) {
//	t := ctx.Value(TraceLogKey).(*ApiInfo)
//	if ctx.Value(TraceLogKey).(*ApiInfo).GinCtx != nil {
//		log.Info(t.OperationID, "api: ", t.ApiName)
//	} else {
//		log.Info(t.OperationID, "rpc: ", t.ApiName)
//	}
//	for _, v := range *t.Funcs {
//		if v.Err != nil {
//			log.Error(t.OperationID, "func: ", v.FuncName, " args: ", v.Args, v.Err.Error())
//		} else {
//			switch v.LogLevel {
//			case logrus.InfoLevel:
//				log.Info(t.OperationID, "func: ", v.FuncName, " args: ", v.Args)
//			case logrus.DebugLevel:
//				log.Debug(t.OperationID, "func: ", v.FuncName, " args: ", v.Args)
//			case logrus.WarnLevel:
//				log.Debug(t.OperationID, "func: ", v.FuncName, " args: ", v.Args)
//			}
//		}
//	}
//}

func WriteErrorResponse(ctx context.Context, funcName string, err error, args ...interface{}) {
	SetCtxInfo(ctx, funcName, err, args)
	e := Unwrap(err)
	switch t := e.(type) {
	case *constant.ErrInfo:
		ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": t.ErrCode, "errMsg": t.ErrMsg, "errDtl": t.DetailErrMsg})
		return
	default:
		s, ok := status.FromError(e)
		if !ok {
			ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": constant.ErrDefaultOther.ErrCode, "errMsg": err.Error(), "errDtl": fmt.Sprintf("%+v", err)})
			return
		}
		var details []string
		if err != e {
			details = append(details, fmt.Sprintf("%+v", err))
		}
		for _, s := range s.Details() {
			details = append(details, fmt.Sprintf("%+v", s))
		}
		ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": s.Code(), "errMsg": s.Message(), "errDtl": strings.Join(details, "\n")})
		return
	}
}

//func WriteErrorResponse(ctx context.Context, funcName string, err error, args ...interface{}) {
//	SetCtxInfo(ctx, funcName, err, args)
//	e := new(constant.ErrInfo)
//	switch {
//	case errors.As(err, &e):
//		ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": e.ErrCode, "errMsg": e.ErrMsg})
//		return
//	default:
//		ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": constant.ErrDefaultOther.ErrCode, "errMsg": constant.ErrDefaultOther.ErrMsg, "errDtl": err.Error()})
//		return
//	}
//}

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

type ApiInfo struct {
	ApiName     string
	OperationID string
	Funcs       *[]FuncInfo
	GinCtx      *gin.Context
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
	t := ctx.Value(TraceLogKey).(*ApiInfo)
	var funcInfo FuncInfo
	funcInfo.Args = make(map[string]interface{})
	argsHandle(args, funcInfo.Args)
	funcInfo.FuncName = funcName
	funcInfo.Err = err
	funcInfo.LogLevel = logLevel
	_, file, line, _ := runtime.Caller(3)
	i := strings.SplitAfter(file, "/")
	if len(i) > 3 {
		file = i[len(i)-3] + i[len(i)-2] + i[len(i)-1] + ":" + utils.IntToString(line)
	}
	funcInfo.File = file
	*t.Funcs = append(*t.Funcs, funcInfo)
}

func SetRpcReqInfo(ctx context.Context, funcName string, req string) {
	t := ctx.Value(TraceLogKey).(*ApiInfo)
	var funcInfo FuncInfo
	funcInfo.Args = make(map[string]interface{})
	var args []interface{}
	args = append(args, " rpc req ", req)
	argsHandle(args, funcInfo.Args)
	funcInfo.FuncName = funcName
	*t.Funcs = append(*t.Funcs, funcInfo)
}

func SetRpcRespInfo(ctx context.Context, funcName string, resp string) {
	t := ctx.Value(TraceLogKey).(*ApiInfo)
	var funcInfo FuncInfo
	funcInfo.Args = make(map[string]interface{})
	var args []interface{}
	args = append(args, " rpc resp ", resp)
	argsHandle(args, funcInfo.Args)
	funcInfo.FuncName = funcName
	*t.Funcs = append(*t.Funcs, funcInfo)
}

func SetSuccess(ctx context.Context, funcName string, data interface{}) {
	SetCtxInfo(ctx, funcName, nil, "data", data)
	ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "errDtl": "", "data": data})
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
