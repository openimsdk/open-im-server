package common

import (
	"Open_IM/internal/fault_tolerant"
	"Open_IM/pkg/common/trace_log"
	"context"
	"examination/pkg/api_struct"
	utils "github.com/OpenIMSDK/open_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func ApiToRpc(c *gin.Context, apiReq, apiResp interface{}, rpcName string, fn interface{}, rpcFuncName string, tokenFunc func(token string) (string, error)) {
	nCtx := trace_log.NewCtx(c, rpcFuncName)
	defer trace_log.ShowLog(nCtx)
	if err := c.BindJSON(apiReq); err != nil {
		trace_log.WriteErrorResponse(nCtx, "BindJSON", err)
		return
	}
	reqValue := reflect.ValueOf(apiReq).Elem()
	operationID := reqValue.FieldByName("OperationID").String()
	trace_log.SetOperationID(nCtx, operationID)

	trace_log.SetContextInfo(nCtx, "BindJSON", nil, "params", apiReq)

	etcdConn, err := fault_tolerant.GetDefaultConn(rpcName, operationID)
	if err != nil {
		trace_log.WriteErrorResponse(nCtx, "GetDefaultConn", err)
		return
	}
	rpc := reflect.ValueOf(fn).Call([]reflect.Value{
		reflect.ValueOf(etcdConn),
	})[0].MethodByName(rpcFuncName) // rpc func
	rpcReqPtr := reflect.New(rpc.Type().In(1).Elem()) // *req参数
	var opUserID string
	if tokenFunc != nil {
		var err error
		opUserID, err = tokenFunc(c.GetHeader("token"))
		if err != nil {
			trace_log.WriteErrorResponse(nCtx, "TokenFunc", err)
			return
		}
	}
	if opID := rpcReqPtr.Elem().FieldByName("OperationID"); opID.IsValid() {
		opID.SetString(operationID)
		if opU := rpcReqPtr.Elem().FieldByName("OpUserID"); opU.IsValid() {
			opU.SetString(opUserID)
		}
	} else {
		op := rpcReqPtr.Elem().FieldByName("Operation").Elem()
		op.FieldByName("OperationID").SetString(operationID)
		op.FieldByName("OpUserID").SetString(opUserID)
	}
	if err := utils.CopyStructFields(rpcReqPtr.Interface(), apiReq); err != nil {
		trace_log.WriteErrorResponse(nCtx, "CopyStructFields_RpcReq", err)
		return
	}
	respArr := rpc.Call([]reflect.Value{
		reflect.ValueOf(context.Context(c)), // context.Context
		rpcReqPtr,                           // rpc apiReq
	}) // respArr => (apiResp, error)
	if !respArr[1].IsNil() { // rpc err != nil
		err := respArr[1].Interface().(error)
		trace_log.WriteErrorResponse(nCtx, rpcFuncName, err, "rpc req", rpcReqPtr.Interface())
		return
	}
	rpcResp := respArr[0].Elem()
	trace_log.SetContextInfo(nCtx, rpcFuncName, nil, "rpc req", rpcReqPtr.Interface(), "resp", rpcResp.Interface())
	commonResp := rpcResp.FieldByName("CommonResp").Elem()
	errCodeVal := commonResp.FieldByName("ErrCode")
	errMsgVal := commonResp.FieldByName("ErrMsg").Interface().(string)
	errCode := errCodeVal.Interface().(int32)
	if errCode != 0 {
		trace_log.WriteErrorResponse(nCtx, "RpcErrCode", &api_struct.ErrInfo{
			ErrCode: errCode,
			ErrMsg:  errMsgVal,
		})
		return
	}
	if apiResp != nil {
		if err := utils.CopyStructFields(apiResp, rpcResp.Interface()); err != nil {
			trace_log.WriteErrorResponse(nCtx, "CopyStructFields_RpcResp", err)
			return
		}
	}
	trace_log.SetSuccess(nCtx, rpcFuncName, apiResp)
}

func WriteDataJson(c *gin.Context, code int32, msg string, data interface{}) {
	c.JSON(http.StatusOK, &BaseResp{
		ErrCode: code,
		ErrMsg:  msg,
		Data:    data,
	})
}

//func WriteErr(c *gin.Context, err error) {
//	if cerr, ok := err.(interface{ Code() int32 }); ok {
//		WriteDataJson(c, cerr.Code(), err.Error(), nil)
//	} else {
//		WriteDataJson(c, ErrDefault.ErrCode, err.Error(), nil)
//	}
//}

func WriteSuccess(c *gin.Context, msg string, data interface{}) {
	WriteDataJson(c, 0, msg, data)
}
