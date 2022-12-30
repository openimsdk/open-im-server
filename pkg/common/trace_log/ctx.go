package trace_log

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

const TraceLogKey = "trace_log"
const GinContextKey = "gin_context"

func ToErrInfoWithErr(errCode int32, errMsg string) *ErrInfo {
	return nil
}

func getAPIErrorResponse(c *gin.Context, errInfo ErrInfo) {

}

func NewCtx(c *gin.Context, api string) context.Context {
	req := ApiInfo{ApiName: api}
	return context.WithValue(c, GinContextKey, req)
}

func ShowLog(ctx context.Context) {
	t := ctx.Value(TraceLogKey).(ApiInfo)
	log.Info(t.OperationID, "api: ", t.ApiName)
	for _, v := range t.Funcs {
		if v.Err != nil {
			log.Error(v.FuncName, v.Err, v.Args)
		} else {
			log.Info(v.FuncName, v.Err, v.Args)
		}
	}

}

func WriteErrorResponse(c *gin.Context, err error) {
	e := new(constant.ErrInfo)
	if errors.As(err, &e) {

	}
}

type FuncInfo struct {
	FuncName string
	Args     map[string]interface{}
	Err      error
}

type ApiInfo struct {
	ApiName     string
	OperationID string
	Funcs       []FuncInfo
}

func SetContextInfo(ctx context.Context, funcName string, err error, args ...interface{}) {
	var req ReqInfo
	t := ctx.Value("f").([]ReqInfo)
	argsHandle(args, req.Args)
	req.FuncName = funcName
	req.Err = err
	t = append(t, req)
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
