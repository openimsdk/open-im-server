package common

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/getcdv3"
	utils2 "Open_IM/pkg/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

func ApiToRpc(c *gin.Context, apiReq, apiResp interface{}, rpcName string, rpcClientFunc interface{}, rpcFuncName string) {
	if rpcName == "" {
		rpcName = utils2.GetFuncName(1)
	}
	logFuncName := fmt.Sprintf("[ApiToRpc: %s]%s", utils2.GetFuncName(1), rpcFuncName)
	ctx := trace_log.NewCtx1(c, rpcFuncName)
	defer log.ShowLog(ctx)
	if err := c.BindJSON(apiReq); err != nil {
		trace_log.WriteErrorResponse(ctx, "BindJSON", err)
		return
	}
	trace_log.SetCtxInfo(ctx, logFuncName, nil, "apiReq", apiReq)
	etcdConn, err := getcdv3.GetConn(ctx, rpcName)
	if err != nil {
		trace_log.WriteErrorResponse(ctx, "GetConn", err)
		return
	}
	rpcClient := reflect.ValueOf(rpcClientFunc).Call([]reflect.Value{
		reflect.ValueOf(etcdConn),
	})[0].MethodByName(rpcFuncName) // rpcClient func
	rpcReqPtr := reflect.New(rpcClient.Type().In(1).Elem()) // *req
	CopyAny(apiReq, rpcReqPtr.Interface())
	trace_log.SetCtxInfo(ctx, logFuncName, nil, "opUserID", c.GetString("opUserID"), "callRpcReq", rpcString(rpcReqPtr.Elem().Interface()))
	respArr := rpcClient.Call([]reflect.Value{
		reflect.ValueOf(context.Context(c)), // context.Context (ctx operationID. opUserID)
		rpcReqPtr,                           // rpcClient apiReq
	}) // respArr => (apiResp, error)
	if !respArr[1].IsNil() { // rpcClient err != nil
		err := respArr[1].Interface().(error)
		trace_log.WriteErrorResponse(ctx, rpcFuncName, err, "callRpcResp", "error")
		return
	}
	rpcResp := respArr[0].Elem()
	trace_log.SetCtxInfo(ctx, rpcFuncName, nil, "callRpcResp", rpcString(rpcResp.Interface()))
	if apiResp != nil {
		CopyAny(rpcResp.Interface(), apiResp)
	}
	trace_log.SetSuccess(ctx, rpcFuncName, apiResp)
}

func rpcString(v interface{}) string {
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%+v", v)
}
