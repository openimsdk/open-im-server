package msggateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"github.com/openimsdk/tools/log"
)

var ErrWriteFull = fmt.Errorf("websocket write buffer full,close connection")

type ClientConn interface {
	ReadMessage() ([]byte, error)
	WriteMessage(message []byte) error
	Close() error
}

type websocketMessage struct {
	MessageType int
	Data        []byte
}

func NewWebSocketClientConn(conn *websocket.Conn, readLimit int64, readTimeout time.Duration, pingInterval time.Duration) ClientConn {
	c := &websocketClientConn{
		readTimeout: readTimeout,
		conn:        conn,
		writer:      make(chan *websocketMessage, 256),
		done:        make(chan struct{}),
	}
	if readLimit > 0 {
		c.conn.SetReadLimit(readLimit)
	}
	c.conn.SetPingHandler(c.pingHandler)
	c.conn.SetPongHandler(c.pongHandler)

	go c.loopSend()
	if pingInterval > 0 {
		go c.doPing(pingInterval)
	}
	return c
}

type websocketClientConn struct {
	readTimeout time.Duration
	conn        *websocket.Conn
	writer      chan *websocketMessage
	done        chan struct{}
	err         atomic.Pointer[error]
}

func (c *websocketClientConn) ReadMessage() ([]byte, error) {
	buf, err := c.readMessage()
	if err != nil {
		return nil, c.closeBy(fmt.Errorf("read message %w", err))
	}
	return buf, nil
}

func (c *websocketClientConn) WriteMessage(message []byte) error {
	return c.writeMessage(websocket.BinaryMessage, message)
}

func (c *websocketClientConn) Close() error {
	_ = c.closeBy(fmt.Errorf("websocket connection closed"))
	return nil
}

func (c *websocketClientConn) closeBy(err error) error {
	if !c.err.CompareAndSwap(nil, &err) {
		return *c.err.Load()
	}
	close(c.done)
	log.ZWarn(context.Background(), "websocket connection closed", err, "remoteAddr", c.conn.RemoteAddr(),
		"chan length", len(c.writer))
	_ = c.conn.Close()
	return err
}

func (c *websocketClientConn) writeMessage(messageType int, data []byte) error {
	if errPtr := c.err.Load(); errPtr != nil {
		return *errPtr
	}
	select {
	case c.writer <- &websocketMessage{MessageType: messageType, Data: data}:
		return nil
	default:
		return c.closeBy(ErrWriteFull)
	}
}

func (c *websocketClientConn) loopSend() {
	var err error
	for {
		select {
		case <-c.done:
			return
		case msg := <-c.writer:
			switch msg.MessageType {
			case websocket.TextMessage, websocket.BinaryMessage:
				err = c.conn.WriteMessage(msg.MessageType, msg.Data)
			default:
				err = c.conn.WriteControl(msg.MessageType, msg.Data, time.Time{})
			}
			if err != nil {
				_ = c.closeBy(err)
				return
			}
		}
	}
}

func (c *websocketClientConn) setReadDeadline() error {
	deadline := time.Now().Add(c.readTimeout)
	return c.conn.SetReadDeadline(deadline)
}

func (c *websocketClientConn) readMessage() ([]byte, error) {
	for {
		if err := c.setReadDeadline(); err != nil {
			return nil, err
		}
		messageType, buf, err := c.conn.ReadMessage()
		if err != nil {
			return nil, err
		}
		switch messageType {
		case websocket.BinaryMessage:
			return buf, nil
		case websocket.TextMessage:
			if err := c.onReadTextMessage(buf); err != nil {
				return nil, err
			}
		case websocket.PingMessage:
			if err := c.pingHandler(string(buf)); err != nil {
				return nil, err
			}
		case websocket.PongMessage:
			if err := c.pongHandler(string(buf)); err != nil {
				return nil, err
			}
		case websocket.CloseMessage:
			if len(buf) == 0 {
				return nil, errors.New("websocket connection closed by peer")
			}
			return nil, fmt.Errorf("websocket connection closed by peer, data %s", string(buf))
		default:
			return nil, fmt.Errorf("unknown websocket message type %d", messageType)
		}
	}
}

func (c *websocketClientConn) onReadTextMessage(buf []byte) error {
	var msg struct {
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
	}
	if err := json.Unmarshal(buf, &msg); err != nil {
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
		return c.writeMessage(websocket.TextMessage, msgData)
	default:
		return fmt.Errorf("not support text message type %s", msg.Type)
	}
}

func (c *websocketClientConn) pingHandler(appData string) error {
	//log.ZWarn(context.Background(), "ping handler recv ping", nil, "remoteAddr", c.conn.RemoteAddr(), "appData", appData)
	if err := c.setReadDeadline(); err != nil {
		return err
	}
	err := c.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second*1))
	if err != nil {
		log.ZWarn(context.Background(), "ping handler write pong error", err, "remoteAddr", c.conn.RemoteAddr(), "appData", appData)
	}
	//log.ZWarn(context.Background(), "ping handler write pong success", nil, "remoteAddr", c.conn.RemoteAddr(), "appData", appData)
	return nil
}

func (c *websocketClientConn) pongHandler(string) error {
	return nil
}

func (c *websocketClientConn) doPing(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			if err := c.writeMessage(websocket.PingMessage, nil); err != nil {
				_ = c.closeBy(fmt.Errorf("send ping %w", err))
				return
			}
		}
	}
}
