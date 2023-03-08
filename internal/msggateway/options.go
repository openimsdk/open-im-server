package msggateway

import "time"

type Option func(opt *configs)
type configs struct {
	//长连接监听端口
	port int
	//长连接允许最大链接数
	maxConnNum int64
	//连接握手超时时间
	handshakeTimeout time.Duration
	//允许消息最大长度
	messageMaxMsgLength int
}

func WithPort(port int) Option {
	return func(opt *configs) {
		opt.port = port
	}
}
func WithMaxConnNum(num int64) Option {
	return func(opt *configs) {
		opt.maxConnNum = num
	}
}
func WithHandshakeTimeout(t time.Duration) Option {
	return func(opt *configs) {
		opt.handshakeTimeout = t
	}
}
func WithMessageMaxMsgLength(length int) Option {
	return func(opt *configs) {
		opt.messageMaxMsgLength = length
	}
}
