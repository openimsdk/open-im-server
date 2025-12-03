package msggateway

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/openimsdk/tools/apiresp"

	"github.com/gorilla/websocket"
	"github.com/openimsdk/tools/errs"
)

type LongConn interface {
	// Close this connection
	Close() error
	// WriteMessage Write message to connection,messageType means data type,can be set binary(2) and text(1).
	WriteMessage(messageType int, message []byte) error
	// ReadMessage Read message from connection.
	ReadMessage() (int, []byte, error)
	// SetReadDeadline sets the read deadline on the underlying network connection,
	// after a read has timed out, will return an error.
	SetReadDeadline(timeout time.Duration) error
	// SetWriteDeadline sets to write deadline when send message,when read has timed out,will return error.
	SetWriteDeadline(timeout time.Duration) error
	// Dial Try to dial a connection,url must set auth args,header can control compress data
	Dial(urlStr string, requestHeader http.Header) (*http.Response, error)
	// IsNil Whether the connection of the current long connection is nil
	IsNil() bool
	// SetConnNil Set the connection of the current long connection to nil
	SetConnNil()
	// SetReadLimit sets the maximum size for a message read from the peer.bytes
	SetReadLimit(limit int64)
	SetPongHandler(handler PingPongHandler)
	SetPingHandler(handler PingPongHandler)
	// GenerateLongConn Check the connection of the current and when it was sent are the same
	GenerateLongConn(w http.ResponseWriter, r *http.Request) error
}
type GWebSocket struct {
	protocolType     int
	conn             *websocket.Conn
	handshakeTimeout time.Duration
	writeBufferSize  int
}

func newGWebSocket(protocolType int, handshakeTimeout time.Duration, wbs int) *GWebSocket {
	return &GWebSocket{protocolType: protocolType, handshakeTimeout: handshakeTimeout, writeBufferSize: wbs}
}

func (d *GWebSocket) Close() error {
	return d.conn.Close()
}

func (d *GWebSocket) GenerateLongConn(w http.ResponseWriter, r *http.Request) error {
	upgrader := &websocket.Upgrader{
		HandshakeTimeout: d.handshakeTimeout,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	if d.writeBufferSize > 0 { // default is 4kb.
		upgrader.WriteBufferSize = d.writeBufferSize
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// The upgrader.Upgrade method usually returns enough error messages to diagnose problems that may occur during the upgrade
		return errs.WrapMsg(err, "GenerateLongConn: WebSocket upgrade failed")
	}
	d.conn = conn
	return nil
}

func (d *GWebSocket) WriteMessage(messageType int, message []byte) error {
	// d.setSendConn(d.conn)
	return d.conn.WriteMessage(messageType, message)
}

// func (d *GWebSocket) setSendConn(sendConn *websocket.Conn) {
//	d.sendConn = sendConn
//}

func (d *GWebSocket) ReadMessage() (int, []byte, error) {
	return d.conn.ReadMessage()
}

func (d *GWebSocket) SetReadDeadline(timeout time.Duration) error {
	return d.conn.SetReadDeadline(time.Now().Add(timeout))
}

func (d *GWebSocket) SetWriteDeadline(timeout time.Duration) error {
	if timeout <= 0 {
		return errs.New("timeout must be greater than 0")
	}

	// TODO SetWriteDeadline Future add error handling
	if err := d.conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return errs.WrapMsg(err, "GWebSocket.SetWriteDeadline failed")
	}
	return nil
}

func (d *GWebSocket) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	conn, httpResp, err := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if err != nil {
		return httpResp, errs.WrapMsg(err, "GWebSocket.Dial failed", "url", urlStr)
	}
	d.conn = conn
	return httpResp, nil
}

func (d *GWebSocket) IsNil() bool {
	return d.conn == nil
}

func (d *GWebSocket) SetConnNil() {
	d.conn = nil
}

func (d *GWebSocket) SetReadLimit(limit int64) {
	d.conn.SetReadLimit(limit)
}

func (d *GWebSocket) SetPongHandler(handler PingPongHandler) {
	d.conn.SetPongHandler(handler)
}

func (d *GWebSocket) SetPingHandler(handler PingPongHandler) {
	d.conn.SetPingHandler(handler)
}

func (d *GWebSocket) RespondWithError(err error, w http.ResponseWriter, r *http.Request) error {
	if err := d.GenerateLongConn(w, r); err != nil {
		return err
	}
	data, err := json.Marshal(apiresp.ParseError(err))
	if err != nil {
		_ = d.Close()
		return errs.WrapMsg(err, "json marshal failed")
	}

	if err := d.WriteMessage(MessageText, data); err != nil {
		_ = d.Close()
		return errs.WrapMsg(err, "WriteMessage failed")
	}
	_ = d.Close()
	return nil
}

func (d *GWebSocket) RespondWithSuccess() error {
	data, err := json.Marshal(apiresp.ParseError(nil))
	if err != nil {
		_ = d.Close()
		return errs.WrapMsg(err, "json marshal failed")
	}

	if err := d.WriteMessage(MessageText, data); err != nil {
		_ = d.Close()
		return errs.WrapMsg(err, "WriteMessage failed")
	}
	return nil
}
