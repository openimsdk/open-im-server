package trace_log

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const TraceLogKey = "trace_log"
const GinContextKey = "gin_context"

func NewCtx(c *gin.Context, api string) context.Context {
	req := ApiInfo{ApiName: api, GinCtx: c}
	req.OperationID = c.PostForm("operationID")
	return context.WithValue(c, GinContextKey, req)
}

func ShowLog(ctx context.Context) {
	t := ctx.Value(TraceLogKey).(ApiInfo)
	log.Info(t.OperationID, "api: ", t.ApiName)
	for _, v := range t.Funcs {
		if v.Err != nil {
			log.Error(v.FuncName, v.Err, v.Args)
		} else {
			log.Info(v.FuncName, v.Args)
		}
	}
}

func WriteErrorResponse(ctx context.Context, funcName string, err error) {
	SetContextInfo(ctx, funcName, err)
	e := new(constant.ErrInfo)
	switch {
	case errors.As(err, &e):
		ctx.Value(GinContextKey).(ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": e.ErrCode, "errMsg": e.ErrMsg})
		return
	default:
		ctx.Value(GinContextKey).(ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": constant.ErrDefaultOther.ErrCode, "errMsg": constant.ErrDefaultOther.ErrMsg})
		return
	}
}

type ApiInfo struct {
	ApiName     string
	OperationID string
	Funcs       []FuncInfo
	GinCtx      *gin.Context
}

type FuncInfo struct {
	FuncName string
	Args     map[string]interface{}
	Err      error
}

func SetContextInfo(ctx context.Context, funcName string, err error, args ...interface{}) {
	t := ctx.Value(TraceLogKey).(ApiInfo)
	var funcInfo FuncInfo
	funcInfo.Args = make(map[string]interface{})
	argsHandle(args, funcInfo.Args)
	funcInfo.FuncName = funcName
	funcInfo.Err = err
	t.Funcs = append(t.Funcs, funcInfo)
}

func argsHandle(args []interface{}, fields map[string]interface{}) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[fmt.Sprintf("%v", args[i])] = args[i+1]
		} else {
			fields[fmt.Sprintf("%v", args[i])] = ""
		}
	}
}
