package msggateway

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/stringutil"
)

var (
	ErrConnClosed                = errs.New("conn has closed")
	ErrNotSupportMessageProtocol = errs.New("not support message protocol")
	ErrClientClosed              = errs.New("client actively close the connection")
	ErrPanic                     = errs.New("panic error")
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
	SDKType        string `json:"sdkType"`
	SDKVersion     string `json:"sdkVersion"`
	Encoder        Encoder
	ctx            *UserConnContext
	longConnServer LongConnServer
	closed         atomic.Bool
	closedErr      error
	token          string
	hbCtx          context.Context
	hbCancel       context.CancelFunc
	subLock        *sync.Mutex
	subUserIDs     map[string]struct{} // client conn subscription list
}

// ResetClient updates the client's state with new connection and context information.
func (c *Client) ResetClient(ctx *UserConnContext, conn LongConn, longConnServer LongConnServer) {
	c.w = new(sync.Mutex)
	c.conn = conn
	c.PlatformID = stringutil.StringToInt(ctx.GetPlatformID())
	c.IsCompress = ctx.GetCompression()
	c.IsBackground = ctx.GetBackground()
	c.UserID = ctx.GetUserID()
	c.ctx = ctx
	c.longConnServer = longConnServer
	c.IsBackground = false
	c.closed.Store(false)
	c.closedErr = nil
	c.token = ctx.GetToken()
	c.SDKType = ctx.GetSDKType()
	c.SDKVersion = ctx.GetSDKVersion()
	c.hbCtx, c.hbCancel = context.WithCancel(c.ctx)
	c.subLock = new(sync.Mutex)
	if c.subUserIDs != nil {
		clear(c.subUserIDs)
	}
	if c.SDKType == GoSDK {
		c.Encoder = NewGobEncoder()
	} else {
		c.Encoder = NewJsonEncoder()
	}
	c.subUserIDs = make(map[string]struct{})
}

func (c *Client) pingHandler(appData string) error {
	if err := c.conn.SetReadDeadline(pongWait); err != nil {
		return err
	}

	log.ZDebug(c.ctx, "ping Handler Success.", "appData", appData)
	return c.writePongMsg(appData)
}

func (c *Client) pongHandler(_ string) error {
	if err := c.conn.SetReadDeadline(pongWait); err != nil {
		return err
	}
	return nil
}

// readMessage continuously reads messages from the connection.
func (c *Client) readMessage() {
	defer func() {
		if r := recover(); r != nil {
			c.closedErr = ErrPanic
			log.ZPanic(c.ctx, "socket have panic err:", errs.ErrPanic(r))
		}
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(pongWait)
	c.conn.SetPongHandler(c.pongHandler)
	c.conn.SetPingHandler(c.pingHandler)
	c.activeHeartbeat(c.hbCtx)

	for {
		log.ZDebug(c.ctx, "readMessage")
		messageType, message, returnErr := c.conn.ReadMessage()
		if returnErr != nil {
			log.ZWarn(c.ctx, "readMessage", returnErr, "messageType", messageType)
			c.closedErr = returnErr
			return
		}

		log.ZDebug(c.ctx, "readMessage", "messageType", messageType)
		if c.closed.Load() {
			// The scenario where the connection has just been closed, but the coroutine has not exited
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
			_ = c.conn.SetReadDeadline(pongWait)
			parseDataErr := c.handlerTextMessage(message)
			if parseDataErr != nil {
				c.closedErr = parseDataErr
				return
			}
		case PingMessage:
			err := c.writePongMsg("")
			log.ZError(c.ctx, "writePongMsg", err)

		case CloseMessage:
			c.closedErr = ErrClientClosed
			return

		default:
		}
	}
}

// handleMessage processes a single message received by the client.
func (c *Client) handleMessage(message []byte) error {
	if c.IsCompress {
		var err error
		message, err = c.longConnServer.DecompressWithPool(message)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	var binaryReq = getReq()
	defer freeReq(binaryReq)

	err := c.Encoder.Decode(message, binaryReq)
	if err != nil {
		return err
	}

	if err := c.longConnServer.Validate(binaryReq); err != nil {
		return err
	}

	if binaryReq.SendID != c.UserID {
		return errs.New("exception conn userID not same to req userID", "binaryReq", binaryReq.String())
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
	case WSPullMsg:
		resp, messageErr = c.longConnServer.GetSeqMessage(ctx, binaryReq)
	case WSGetConvMaxReadSeq:
		resp, messageErr = c.longConnServer.GetConversationsHasReadAndMaxSeq(ctx, binaryReq)
	case WsPullConvLastMessage:
		resp, messageErr = c.longConnServer.GetLastMessage(ctx, binaryReq)
	case WsLogoutMsg:
		resp, messageErr = c.longConnServer.UserLogout(ctx, binaryReq)
	case WsSetBackgroundStatus:
		resp, messageErr = c.setAppBackgroundStatus(ctx, binaryReq)
	case WsSubUserOnlineStatus:
		resp, messageErr = c.longConnServer.SubUserOnlineStatus(ctx, c, binaryReq)
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
	// TODO: callback
	return resp, nil
}

func (c *Client) close() {
	c.w.Lock()
	defer c.w.Unlock()
	if c.closed.Load() {
		return
	}
	c.closed.Store(true)
	c.conn.Close()
	c.hbCancel() // Close server-initiated heartbeat.
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
	t := time.Now()
	log.ZDebug(ctx, "gateway reply message", "resp", mReply.String())
	err = c.writeBinaryMsg(mReply)
	if err != nil {
		log.ZWarn(ctx, "wireBinaryMsg replyMessage", err, "resp", mReply.String())
	}
	log.ZDebug(ctx, "wireBinaryMsg end", "time cost", time.Since(t))

	if binaryReq.ReqIdentifier == WsLogoutMsg {
		return errs.New("user logout", "operationID", binaryReq.OperationID).Wrap()
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
	log.ZDebug(c.ctx, "KickOnlineMessage debug ")
	err := c.writeBinaryMsg(resp)
	c.close()
	return err
}

func (c *Client) PushUserOnlineStatus(data []byte) error {
	resp := Resp{
		ReqIdentifier: WsSubUserOnlineStatus,
		Data:          data,
	}
	return c.writeBinaryMsg(resp)
}

func (c *Client) writeBinaryMsg(resp Resp) error {
	if c.closed.Load() {
		return nil
	}

	encodedBuf, err := c.Encoder.Encode(resp)
	if err != nil {
		return err
	}

	c.w.Lock()
	defer c.w.Unlock()

	err = c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		return err
	}

	if c.IsCompress {
		resultBuf, compressErr := c.longConnServer.CompressWithPool(encodedBuf)
		if compressErr != nil {
			return compressErr
		}
		return c.conn.WriteMessage(MessageBinary, resultBuf)
	}

	return c.conn.WriteMessage(MessageBinary, encodedBuf)
}

// Actively initiate Heartbeat when platform in Web.
func (c *Client) activeHeartbeat(ctx context.Context) {
	if c.PlatformID == constant.WebPlatformID {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.ZPanic(ctx, "activeHeartbeat Panic", errs.ErrPanic(r))
				}
			}()
			log.ZDebug(ctx, "server initiative send heartbeat start.")
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := c.writePingMsg(); err != nil {
						log.ZWarn(c.ctx, "send Ping Message error.", err)
						return
					}
				case <-c.hbCtx.Done():
					return
				}
			}
		}()
	}
}
func (c *Client) writePingMsg() error {
	if c.closed.Load() {
		return nil
	}

	c.w.Lock()
	defer c.w.Unlock()

	err := c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(PingMessage, nil)
}

func (c *Client) writePongMsg(appData string) error {
	log.ZDebug(c.ctx, "write Pong Msg in Server", "appData", appData)
	if c.closed.Load() {
		log.ZWarn(c.ctx, "is closed in server", nil, "appdata", appData, "closed err", c.closedErr)
		return nil
	}

	c.w.Lock()
	defer c.w.Unlock()

	err := c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		log.ZWarn(c.ctx, "SetWriteDeadline in Server have error", errs.Wrap(err), "writeWait", writeWait, "appData", appData)
		return errs.Wrap(err)
	}
	err = c.conn.WriteMessage(PongMessage, []byte(appData))
	if err != nil {
		log.ZWarn(c.ctx, "Write Message have error", errs.Wrap(err), "Pong msg", PongMessage)
	}

	return errs.Wrap(err)
}

func (c *Client) handlerTextMessage(b []byte) error {
	var msg TextMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		return err
	}
	switch msg.Type {
	case TextPong:
		return nil
	case TextPing:
		msg.Type = TextPong
		msgData, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		c.w.Lock()
		defer c.w.Unlock()
		if err := c.conn.SetWriteDeadline(writeWait); err != nil {
			return err
		}
		return c.conn.WriteMessage(MessageText, msgData)
	default:
		return fmt.Errorf("not support message type %s", msg.Type)
	}
}
