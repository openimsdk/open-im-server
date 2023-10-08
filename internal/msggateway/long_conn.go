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
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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
}

func newGWebSocket(protocolType int, handshakeTimeout time.Duration) *GWebSocket {
	return &GWebSocket{protocolType: protocolType, handshakeTimeout: handshakeTimeout}
}

func (d *GWebSocket) Close() error {
	return d.conn.Close()
}

func (d *GWebSocket) GenerateLongConn(w http.ResponseWriter, r *http.Request) error {
	upgrader := &websocket.Upgrader{
		HandshakeTimeout: d.handshakeTimeout,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	d.conn = conn
	return nil
}

func (d *GWebSocket) WriteMessage(messageType int, message []byte) error {
	// d.setSendConn(d.conn)
	return d.conn.WriteMessage(messageType, message)
}

//func (d *GWebSocket) setSendConn(sendConn *websocket.Conn) {
//	d.sendConn = sendConn
//}

func (d *GWebSocket) ReadMessage() (int, []byte, error) {
	return d.conn.ReadMessage()
}

func (d *GWebSocket) SetReadDeadline(timeout time.Duration) error {
	return d.conn.SetReadDeadline(time.Now().Add(timeout))
}

func (d *GWebSocket) SetWriteDeadline(timeout time.Duration) error {
	return d.conn.SetWriteDeadline(time.Now().Add(timeout))
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

func (d *GWebSocket) SetReadLimit(limit int64) {
	d.conn.SetReadLimit(limit)
}

func (d *GWebSocket) SetPongHandler(handler PingPongHandler) {
	d.conn.SetPongHandler(handler)
}

func (d *GWebSocket) SetPingHandler(handler PingPongHandler) {
	d.conn.SetPingHandler(handler)
}

//func (d *GWebSocket) CheckSendConnDiffNow() bool {
//	return d.conn == d.sendConn
//}
