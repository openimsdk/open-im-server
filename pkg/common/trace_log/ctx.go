package trace_log

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const TraceLogKey = "trace_log"

func NewCtx(c *gin.Context, api string) context.Context {
	req := &ApiInfo{ApiName: api, GinCtx: c, Funcs: &[]FuncInfo{}}
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

func ShowLog(ctx context.Context) {
	t := ctx.Value(TraceLogKey).(*ApiInfo)
	if ctx.Value(TraceLogKey).(*ApiInfo).GinCtx != nil {
		log.Info(t.OperationID, "api: ", t.ApiName)
	} else {
		log.Info(t.OperationID, "rpc: ", t.ApiName)
	}

	for _, v := range *t.Funcs {
		if v.Err != nil {
			log.Error(t.OperationID, "func: ", v.FuncName, " args: ", v.Args, v.Err.Error())
		} else {
			log.Info(t.OperationID, "func: ", v.FuncName, " args: ", v.Args)
		}
	}
}

func WriteErrorResponse(ctx context.Context, funcName string, err error, args ...interface{}) {
	SetContextInfo(ctx, funcName, err, args)
	e := new(constant.ErrInfo)
	switch {
	case errors.As(err, &e):
		ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": e.ErrCode, "errMsg": e.ErrMsg})
		return
	default:
		ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": constant.ErrDefaultOther.ErrCode, "errMsg": constant.ErrDefaultOther.ErrMsg})
		return
	}
}

type ApiInfo struct {
	ApiName     string
	OperationID string
	Funcs       *[]FuncInfo
	GinCtx      *gin.Context
}

type FuncInfo struct {
	FuncName string
	Args     map[string]interface{}
	Err      error
}

func SetContextInfo(ctx context.Context, funcName string, err error, args ...interface{}) {
	t := ctx.Value(TraceLogKey).(*ApiInfo)
	var funcInfo FuncInfo
	funcInfo.Args = make(map[string]interface{})
	argsHandle(args, funcInfo.Args)
	funcInfo.FuncName = utils.GetSelfFuncName()
	funcInfo.Err = err
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
	SetContextInfo(ctx, funcName, nil, "data", data)
	ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": data})
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

func GetApiErr(errCode int32, errMsg string) constant.ErrInfo {
	return constant.ErrInfo{ErrCode: errCode, ErrMsg: errMsg}
}
