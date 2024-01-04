package control

import (
	"call-back-http/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CallbackBeforeSendSingleMsgCommand(c *gin.Context) {
	var req model.CallbackBeforeSendSingleMsgReq
	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	fmt.Printf("CallbackBeforeSendSingleMsgCommand received:%#v\n", req)
	str := "callback return message"
	byte, err := json.Marshal(str)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.CallbackBeforeSendSingleMsgResp{
			CommonCallbackResp: model.CommonCallbackResp{
				ActionCode: 500,
				ErrCode:    5001,
				ErrMsg:     "callback error",
				ErrDlt:     err.Error(),
				NextCode:   2,
			},
		})
	}
	resp := &model.CallbackBeforeSendSingleMsgResp{
		CommonCallbackResp: model.CommonCallbackResp{
			ActionCode: 0,
			ErrCode:    2000,
			ErrMsg:     "Success",
			ErrDlt:     "Successful",
			NextCode:   2,
		},
		Content: byte,
	}
	fmt.Printf("CallbackBeforeSendSingleMsgCommand return:%#v\n", resp)
	c.JSON(http.StatusOK, resp)
}
