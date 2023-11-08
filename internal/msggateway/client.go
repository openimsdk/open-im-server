// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msggateway

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"

	"google.golang.org/protobuf/proto"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/utils"
)

var (
	ErrConnClosed                = errors.New("conn has closed")
	ErrNotSupportMessageProtocol = errors.New("not support message protocol")
	ErrClientClosed              = errors.New("client actively close the connection")
	ErrPanic                     = errors.New("panic error")
)

const (
	// MessageText is for UTF-8 encoded text messages like JSON.
	MessageText = iota + 1
	// MessageBinary is for binary messages like protobufs.
	MessageBinary
	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

type PingPongHandler func(string) error

type Client struct {
	w              *sync.Mutex
	conn           LongConn
	PlatformID     int    `json:"platformID"`
	IsCompress     bool   `json:"isCompress"`
	UserID         string `json:"userID"`
	IsBackground   bool   `json:"isBackground"`
	ctx            *UserConnContext
	longConnServer LongConnServer
	closed         atomic.Bool
	closedErr      error
	token          string
}

func newClient(ctx *UserConnContext, conn LongConn, isCompress bool) *Client {
	return &Client{
		w:          new(sync.Mutex),
		conn:       conn,
		PlatformID: utils.StringToInt(ctx.GetPlatformID()),
		IsCompress: isCompress,
		UserID:     ctx.GetUserID(),
		ctx:        ctx,
	}
}

func (c *Client) ResetClient(
	ctx *UserConnContext,
	conn LongConn,
	isBackground, isCompress bool,
	longConnServer LongConnServer,
	token string,
) {
	c.w = new(sync.Mutex)
	c.conn = conn
	c.PlatformID = utils.StringToInt(ctx.GetPlatformID())
	c.IsCompress = isCompress
	c.IsBackground = isBackground
	c.UserID = ctx.GetUserID()
	c.ctx = ctx
	c.longConnServer = longConnServer
	c.IsBackground = false
	c.closed.Store(false)
	c.closedErr = nil
	c.token = token
}

func (c *Client) pingHandler(_ string) error {
	_ = c.conn.SetReadDeadline(pongWait)
	return c.writePongMsg()
}

func (c *Client) readMessage() {
	defer func() {
		if r := recover(); r != nil {
			c.closedErr = ErrPanic
			fmt.Println("socket have panic err:", r, string(debug.Stack()))
		}
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(pongWait)
	c.conn.SetPingHandler(c.pingHandler)

	for {
		messageType, message, returnErr := c.conn.ReadMessage()
		if returnErr != nil {
			log.ZWarn(c.ctx, "readMessage", returnErr, "messageType", messageType)
			c.closedErr = returnErr
			return
		}

		log.ZDebug(c.ctx, "readMessage", "messageType", messageType)
		if c.closed.Load() { // 连接刚置位已经关闭，但是协程还没退出的场景
			c.closedErr = ErrConnClosed
			return
		}

		switch messageType {
		case MessageBinary:
			_ = c.conn.SetReadDeadline(pongWait)
			parseDataErr := c.handleMessage(message)
			if parseDataErr != nil {
				c.closedErr = parseDataErr
				return
			}
		case MessageText:
			c.closedErr = ErrNotSupportMessageProtocol
			return

		case PingMessage:
			err := c.writePongMsg()
			log.ZError(c.ctx, "writePongMsg", err)

		case CloseMessage:
			c.closedErr = ErrClientClosed
			return
		default:
		}
	}
}

func (c *Client) handleMessage(message []byte) error {
	if c.IsCompress {
		var err error
		message, err = c.longConnServer.DeCompress(message)
		if err != nil {
			return utils.Wrap(err, "")
		}
	}

	var binaryReq = getReq()
	defer freeReq(binaryReq)

	err := c.longConnServer.Decode(message, binaryReq)
	if err != nil {
		return utils.Wrap(err, "")
	}

	if err := c.longConnServer.Validate(binaryReq); err != nil {
		return utils.Wrap(err, "")
	}

	if binaryReq.SendID != c.UserID {
		return utils.Wrap(errors.New("exception conn userID not same to req userID"), binaryReq.String())
	}

	ctx := mcontext.WithMustInfoCtx(
		[]string{binaryReq.OperationID, binaryReq.SendID, constant.PlatformIDToName(c.PlatformID), c.ctx.GetConnID()},
	)

	log.ZDebug(ctx, "gateway req message", "req", binaryReq.String())

	var (
		resp       []byte
		messageErr error
	)

	switch binaryReq.ReqIdentifier {
	case WSGetNewestSeq:
		resp, messageErr = c.longConnServer.GetSeq(ctx, binaryReq)
	case WSSendMsg:
		resp, messageErr = c.longConnServer.SendMessage(ctx, binaryReq)
	case WSSendSignalMsg:
		resp, messageErr = c.longConnServer.SendSignalMessage(ctx, binaryReq)
	case WSPullMsgBySeqList:
		resp, messageErr = c.longConnServer.PullMessageBySeqList(ctx, binaryReq)
	case WsLogoutMsg:
		resp, messageErr = c.longConnServer.UserLogout(ctx, binaryReq)
	case WsSetBackgroundStatus:
		resp, messageErr = c.setAppBackgroundStatus(ctx, binaryReq)
	default:
		return fmt.Errorf(
			"ReqIdentifier failed,sendID:%s,msgIncr:%s,reqIdentifier:%d",
			binaryReq.SendID,
			binaryReq.MsgIncr,
			binaryReq.ReqIdentifier,
		)
	}

	return c.replyMessage(ctx, binaryReq, messageErr, resp)
}

func (c *Client) setAppBackgroundStatus(ctx context.Context, req *Req) ([]byte, error) {
	resp, isBackground, messageErr := c.longConnServer.SetUserDeviceBackground(ctx, req)
	if messageErr != nil {
		return nil, messageErr
	}

	c.IsBackground = isBackground
	// todo callback
	return resp, nil
}

func (c *Client) close() {
	if c.closed.Load() {
		return
	}

	c.w.Lock()
	defer c.w.Unlock()

	c.closed.Store(true)
	c.conn.Close()
	c.longConnServer.UnRegister(c)
}

func (c *Client) replyMessage(ctx context.Context, binaryReq *Req, err error, resp []byte) error {
	errResp := apiresp.ParseError(err)
	mReply := Resp{
		ReqIdentifier: binaryReq.ReqIdentifier,
		MsgIncr:       binaryReq.MsgIncr,
		OperationID:   binaryReq.OperationID,
		ErrCode:       errResp.ErrCode,
		ErrMsg:        errResp.ErrMsg,
		Data:          resp,
	}
	log.ZDebug(ctx, "gateway reply message", "resp", mReply.String())
	err = c.writeBinaryMsg(mReply)
	if err != nil {
		log.ZWarn(ctx, "wireBinaryMsg replyMessage", err, "resp", mReply.String())
	}

	if binaryReq.ReqIdentifier == WsLogoutMsg {
		return errors.New("user logout")
	}
	return nil
}

func (c *Client) PushMessage(ctx context.Context, msgData *sdkws.MsgData) error {
	var msg sdkws.PushMessages
	conversationID := msgprocessor.GetConversationIDByMsg(msgData)
	m := map[string]*sdkws.PullMsgs{conversationID: {Msgs: []*sdkws.MsgData{msgData}}}
	if msgprocessor.IsNotification(conversationID) {
		msg.NotificationMsgs = m
	} else {
		msg.Msgs = m
	}
	log.ZDebug(ctx, "PushMessage", "msg", &msg)
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}
	resp := Resp{
		ReqIdentifier: WSPushMsg,
		OperationID:   mcontext.GetOperationID(ctx),
		Data:          data,
	}
	return c.writeBinaryMsg(resp)
}

func (c *Client) KickOnlineMessage() error {
	resp := Resp{
		ReqIdentifier: WSKickOnlineMsg,
	}
	err := c.writeBinaryMsg(resp)
	c.close()
	return err
}

func (c *Client) writeBinaryMsg(resp Resp) error {
	if c.closed.Load() {
		return nil
	}

	encodedBuf, err := c.longConnServer.Encode(resp)
	if err != nil {
		return utils.Wrap(err, "")
	}

	c.w.Lock()
	defer c.w.Unlock()

	_ = c.conn.SetWriteDeadline(writeWait)
	if c.IsCompress {
		resultBuf, compressErr := c.longConnServer.Compress(encodedBuf)
		if compressErr != nil {
			return utils.Wrap(compressErr, "")
		}
		return c.conn.WriteMessage(MessageBinary, resultBuf)
	}

	return c.conn.WriteMessage(MessageBinary, encodedBuf)
}

func (c *Client) writePongMsg() error {
	if c.closed.Load() {
		return nil
	}

	c.w.Lock()
	defer c.w.Unlock()

	err := c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		return utils.Wrap(err, "")
	}

	return c.conn.WriteMessage(PongMessage, nil)
}
