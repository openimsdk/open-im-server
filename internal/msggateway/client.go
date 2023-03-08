package msggateway

import (
	"OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
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

type Client struct {
	w              *sync.Mutex
	conn           LongConn
	platformID     int
	isCompress     bool
	userID         string
	isBackground   bool
	connID         string
	onlineAt       int64 // 上线时间戳（毫秒）
	longConnServer LongConnServer
	closed         bool
}

func newClient(ctx *UserConnContext, conn LongConn, isCompress bool) *Client {
	return &Client{
		w:          new(sync.Mutex),
		conn:       conn,
		platformID: utils.StringToInt(ctx.GetPlatformID()),
		isCompress: isCompress,
		userID:     ctx.GetUserID(),
		connID:     ctx.GetConnID(),
		onlineAt:   utils.GetCurrentTimestampByMill(),
	}
}
func (c *Client) ResetClient(ctx *UserConnContext, conn LongConn, isCompress bool, longConnServer LongConnServer) {
	c.w = new(sync.Mutex)
	c.conn = conn
	c.platformID = utils.StringToInt(ctx.GetPlatformID())
	c.isCompress = isCompress
	c.userID = ctx.GetUserID()
	c.connID = ctx.GetConnID()
	c.onlineAt = utils.GetCurrentTimestampByMill()
	c.longConnServer = longConnServer
}
func (c *Client) readMessage() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("socket have panic err:", r, string(debug.Stack()))
		}
		//c.close()
	}()
	//var returnErr error
	for {
		messageType, message, returnErr := c.conn.ReadMessage()
		if returnErr != nil {
			break
		}
		if c.closed == true { //连接刚置位已经关闭，但是协程还没退出的场景
			break
		}
		switch messageType {
		case PingMessage:
		case PongMessage:
		case CloseMessage:
			return
		case MessageText:
		case MessageBinary:
			if len(message) == 0 {
				continue
			}
			returnErr = c.handleMessage(message)
			if returnErr != nil {
				break
			}

		}
	}

}
func (c *Client) handleMessage(message []byte) error {
	if c.isCompress {
		var decompressErr error
		message, decompressErr = c.longConnServer.DeCompress(message)
		if decompressErr != nil {
			return utils.Wrap(decompressErr, "")
		}
	}
	var binaryReq Req
	err := c.longConnServer.Decode(message, &binaryReq)
	if err != nil {
		return utils.Wrap(err, "")
	}
	if err := c.longConnServer.Validate(binaryReq); err != nil {
		return utils.Wrap(err, "")
	}
	if binaryReq.SendID != c.userID {
		return errors.New("exception conn userID not same to req userID")
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, ConnID, c.connID)
	ctx = context.WithValue(ctx, OperationID, binaryReq.OperationID)
	ctx = context.WithValue(ctx, CommonUserID, binaryReq.SendID)
	ctx = context.WithValue(ctx, PlatformID, c.platformID)
	var messageErr error
	var resp []byte
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
		return errors.New(fmt.Sprintf("ReqIdentifier failed,sendID:%d,msgIncr:%s,reqIdentifier:%s", binaryReq.SendID, binaryReq.MsgIncr, binaryReq.ReqIdentifier))
	}
	c.replyMessage(&binaryReq, messageErr, resp)
	return nil

}
func (c *Client) setAppBackgroundStatus(ctx context.Context, req Req) ([]byte, error) {
	resp, isBackground, messageErr := c.longConnServer.SetUserDeviceBackground(ctx, req)
	if messageErr != nil {
		return nil, messageErr
	}
	c.isBackground = isBackground
	//todo callback
	return resp, nil

}
func (c *Client) close() {
	c.w.Lock()
	defer c.w.Unlock()
	c.conn.Close()
	c.longConnServer.UnRegister(c)

}
func (c *Client) replyMessage(binaryReq *Req, err error, resp []byte) {
	mReply := Resp{
		ReqIdentifier: binaryReq.ReqIdentifier,
		MsgIncr:       binaryReq.MsgIncr,
		OperationID:   binaryReq.OperationID,
		Data:          resp,
	}
	_ = c.writeMsg(mReply)
}
func (c *Client) PushMessage(ctx context.Context, msgData *sdkws.MsgData) error {
	return nil
}

func (c *Client) KickOnlineMessage(ctx context.Context) error {
	return nil
}

func (c *Client) writeMsg(resp Resp) error {
	c.w.Lock()
	defer c.w.Unlock()
	if c.closed == true {
		return nil
	}
	encodedBuf := bufferPool.Get().([]byte)
	resultBuf := bufferPool.Get().([]byte)
	encodeBuf, err := c.longConnServer.Encode(resp)
	if err != nil {
		return utils.Wrap(err, "")
	}
	_ = c.conn.SetWriteTimeout(60)
	if c.isCompress {
		var compressErr error
		resultBuf, compressErr = c.longConnServer.Compress(encodeBuf)
		if compressErr != nil {
			return utils.Wrap(compressErr, "")
		}
		return c.conn.WriteMessage(MessageBinary, resultBuf)
	} else {
		return c.conn.WriteMessage(MessageBinary, encodedBuf)
	}
}
