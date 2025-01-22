package service

import (
	"crypto/tls"
	"net"
	"testing"
	"time"
)

func TestName1(t *testing.T) {

	tls.Client(&testConn{}, &tls.Config{}).Handshake()

	time.Sleep(time.Hour)
}

type testConn struct {
}

func (testConn) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (testConn) Write(b []byte) (n int, err error) {
	panic("implement me")
}

func (testConn) Close() error {
	//TODO implement me
	panic("implement me")
}

func (testConn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (testConn) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (testConn) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (testConn) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (testConn) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}
