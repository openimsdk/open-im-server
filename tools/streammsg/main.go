package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	getAdminToken   = "/auth/get_admin_token"
	sendMsgApi      = "/msg/send_msg"
	appendStreamMsg = "/msg/append_stream_msg"
)

var (
	ApiAddr = "http://127.0.0.1:10002"
	Token   string
)

func ApiCall[R any](api string, req any) (*R, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, ApiAddr+api, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	if Token != "" {
		request.Header.Set("token", Token)
	}
	request.Header.Set(constant.OperationID, uuid.New().String())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var resp R
	apiResponse := apiresp.ApiResponse{
		Data: &resp,
	}
	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}
	if apiResponse.ErrCode != 0 {
		return nil, errs.NewCodeError(apiResponse.ErrCode, apiResponse.ErrMsg)
	}
	return &resp, nil
}

func main() {
	resp, err := ApiCall[auth.GetAdminTokenResp](getAdminToken, &auth.GetAdminTokenReq{
		Secret: "openIM123",
		UserID: "imAdmin",
	})
	if err != nil {
		fmt.Println("get admin token failed", err)
		return
	}
	Token = resp.Token
	g := gin.Default()
	g.POST("/callbackExample/callbackAfterSendSingleMsgCommand", toGin(handlerUserMsg))
	if err := g.Run(":10006"); err != nil {
		panic(err)
	}
}

func toGin[R any](fn func(c *gin.Context, req *R) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Printf("HTTP %s %s %s\n", c.Request.Method, c.Request.URL, body)
		var req R
		if err := json.Unmarshal(body, &req); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if err := fn(c, &req); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, "{}")
	}
}

func handlerUserMsg(c *gin.Context, req *cbapi.CallbackAfterSendSingleMsgReq) error {
	if req.ContentType != constant.Text {
		return nil
	}
	if !strings.Contains(req.Content, "stream") {
		return nil
	}
	apiReq := apistruct.SendMsgReq{
		RecvID: req.SendID,
		SendMsg: apistruct.SendMsg{
			SendID:           req.RecvID,
			SenderNickname:   "xxx",
			SenderFaceURL:    "",
			SenderPlatformID: constant.AdminPlatformID,
			ContentType:      constant.Stream,
			SessionType:      req.SessionType,
			SendTime:         time.Now().UnixMilli(),
			Content: map[string]any{
				"type":    "xxx",
				"content": "server test stream msg",
			},
		},
	}
	go func() {
		if err := doPushStreamMsg(&apiReq); err != nil {
			fmt.Println("doPushStreamMsg failed", err)
			return
		}
		fmt.Println("doPushStreamMsg success")
	}()
	return nil
}

func doPushStreamMsg(sendReq *apistruct.SendMsgReq) error {
	resp, err := ApiCall[msg.SendMsgResp](sendMsgApi, sendReq)
	if err != nil {
		return err
	}
	const num = 5
	for i := 1; i <= num; i++ {
		_, err := ApiCall[msg.AppendStreamMsgResp](appendStreamMsg, &msg.AppendStreamMsgReq{
			ClientMsgID: resp.ClientMsgID,
			StartIndex:  int64(i - 1),
			Packets: []string{
				fmt.Sprintf("stream_msg_packet_%03d", i),
			},
			End: i == num,
		})
		if err != nil {
			fmt.Println("append stream msg failed", "clientMsgID", resp.ClientMsgID, "index", fmt.Sprintf("%d/%d", i, num), "error", err)
			return err
		}
		fmt.Println("append stream msg success", "clientMsgID", resp.ClientMsgID, "index", fmt.Sprintf("%d/%d", i, num))
		time.Sleep(time.Second * 10)
	}
	return nil
}
