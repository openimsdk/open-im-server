package common

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/getcdv3"
	utils2 "Open_IM/pkg/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"net/http"
	"reflect"
	"strings"
)

func ApiToRpc(c *gin.Context, apiReq, apiResp interface{}, rpcName string, rpcClientFunc interface{}, rpcFuncName string) {
	if rpcName == "" {
		rpcName = utils2.GetFuncName(1)
	}
	logFuncName := fmt.Sprintf("[ApiToRpc: %s]%s", utils2.GetFuncName(1), rpcFuncName)
	ctx := tracelog.NewCtx(c, rpcFuncName)
	defer log.ShowLog(ctx)
	if err := c.BindJSON(apiReq); err != nil {
		WriteErrorResponse(ctx, "BindJSON", err)
		return
	}
	tracelog.SetCtxInfo(ctx, logFuncName, nil, "apiReq", apiReq)
	etcdConn, err := getcdv3.GetConn(ctx, rpcName)
	if err != nil {
		WriteErrorResponse(ctx, "GetConn", err)
		return
	}
	rpcClient := reflect.ValueOf(rpcClientFunc).Call([]reflect.Value{
		reflect.ValueOf(etcdConn),
	})[0].MethodByName(rpcFuncName) // rpcClient func
	rpcReqPtr := reflect.New(rpcClient.Type().In(1).Elem()) // *req
	CopyAny(apiReq, rpcReqPtr.Interface())
	tracelog.SetCtxInfo(ctx, logFuncName, nil, "opUserID", c.GetString("opUserID"), "callRpcReq", rpcString(rpcReqPtr.Elem().Interface()))
	respArr := rpcClient.Call([]reflect.Value{
		reflect.ValueOf(context.Context(c)), // context.Context (ctx operationID. opUserID)
		rpcReqPtr,                           // rpcClient apiReq
	}) // respArr => (apiResp, error)
	if !respArr[1].IsNil() { // rpcClient err != nil
		err := respArr[1].Interface().(error)
		WriteErrorResponse(ctx, rpcFuncName, err, "callRpcResp", "error")
		return
	}
	rpcResp := respArr[0].Elem()
	tracelog.SetCtxInfo(ctx, rpcFuncName, nil, "callRpcResp", rpcString(rpcResp.Interface()))
	if apiResp != nil {
		CopyAny(rpcResp.Interface(), apiResp)
	}
	SetSuccess(ctx, rpcFuncName, apiResp)
}

func rpcString(v interface{}) string {
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%+v", v)
}

type baseResp struct {
	ErrCode int32       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	ErrDtl  string      `json:"errDtl"`
	Data    interface{} `json:"data"`
}

func WriteErrorResponse(ctx context.Context, funcName string, err error, args ...interface{}) {
	tracelog.SetCtxInfo(ctx, funcName, err, args)
	e := tracelog.Unwrap(err)
	switch t := e.(type) {
	case *constant.ErrInfo:
		ctx.Value(tracelog.TraceLogKey).(*tracelog.ApiInfo).GinCtx.JSON(http.StatusOK, baseResp{ErrCode: t.ErrCode, ErrMsg: t.ErrMsg, ErrDtl: t.DetailErrMsg})
		//ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": t.ErrCode, "errMsg": t.ErrMsg, "errDtl": t.DetailErrMsg})
		return
	default:
		s, ok := status.FromError(e)
		if !ok {
			ctx.Value(tracelog.TraceLogKey).(*tracelog.ApiInfo).GinCtx.JSON(http.StatusOK, &baseResp{ErrCode: constant.ErrDefaultOther.ErrCode, ErrMsg: err.Error(), ErrDtl: fmt.Sprintf("%+v", err)})
			//ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": constant.ErrDefaultOther.ErrCode, "errMsg": err.Error(), "errDtl": fmt.Sprintf("%+v", err)})
			return
		}
		var details []string
		if err != e {
			details = append(details, fmt.Sprintf("%+v", err))
		}
		for _, s := range s.Details() {
			details = append(details, fmt.Sprintf("%+v", s))
		}
		ctx.Value(tracelog.TraceLogKey).(*tracelog.ApiInfo).GinCtx.JSON(http.StatusOK, &baseResp{ErrCode: int32(s.Code()), ErrMsg: s.Message(), ErrDtl: strings.Join(details, "\n")})
		//ctx.Value(TraceLogKey).(*ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": s.Code(), "errMsg": s.Message(), "errDtl": strings.Join(details, "\n")})
		return
	}
}

func SetSuccess(ctx context.Context, funcName string, data interface{}) {
	tracelog.SetCtxInfo(ctx, funcName, nil, "data", data)
	ctx.Value(tracelog.TraceLogKey).(*tracelog.ApiInfo).GinCtx.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "errDtl": "", "data": data})
}
