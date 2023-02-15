package new

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type LongConn interface {
	//Close this connection
	Close() error
	//Write message to connection,messageType means data type,can be set binary(2) and text(1).
	WriteMessage(messageType int, message []byte) error
	//Read message from connection.
	ReadMessage() (int, []byte, error)
	//SetReadTimeout sets the read deadline on the underlying network connection,
	//after a read has timed out, will return an error.
	SetReadTimeout(timeout int) error
	//SetWriteTimeout sets the write deadline when send message,when read has timed out,will return error.
	SetWriteTimeout(timeout int) error
	//Try to dial a connection,url must set auth args,header can control compress data
	Dial(urlStr string, requestHeader http.Header) (*http.Response, error)
	//Whether the connection of the current long connection is nil
	IsNil() bool
	//Set the connection of the current long connection to nil
	SetConnNil()
	//Check the connection of the current and when it was sent are the same
	CheckSendConnDiffNow() bool
}
type GWebSocket struct {
	protocolType int
	conn         *websocket.Conn
}

func NewDefault(protocolType int) *GWebSocket {
	return &GWebSocket{protocolType: protocolType}
}
func (d *GWebSocket) Close() error {
	return d.conn.Close()
}

func (d *GWebSocket) WriteMessage(messageType int, message []byte) error {
	d.setSendConn(d.conn)
	return d.conn.WriteMessage(messageType, message)
}

func (d *GWebSocket) setSendConn(sendConn *websocket.Conn) {
	d.sendConn = sendConn
}

func (d *GWebSocket) ReadMessage() (int, []byte, error) {
	return d.conn.ReadMessage()
}
func (d *GWebSocket) SetReadTimeout(timeout int) error {
	return d.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (d *GWebSocket) SetWriteTimeout(timeout int) error {
	return d.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (d *GWebSocket) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	conn, httpResp, err := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if err == nil {
		d.conn = conn
	}
	return httpResp, err

}

func (d *GWebSocket) IsNil() bool {
	if d.conn != nil {
		return false
	}
	return true
}

func (d *GWebSocket) SetConnNil() {
	d.conn = nil
}
func (d *GWebSocket) CheckSendConnDiffNow() bool {
	return d.conn == d.sendConn
}
