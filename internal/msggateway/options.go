package msggateway

import "time"

type (
	Option  func(opt *configs)
	configs struct {
		// Long connection listening port
		port int
		// Maximum number of connections allowed for long connection
		maxConnNum int64
		// Connection handshake timeout
		handshakeTimeout time.Duration
		// Maximum length allowed for messages
		messageMaxMsgLength int
		// Websocket write buffer, default: 4096, 4kb.
		writeBufferSize int
	}
)

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

func WithWriteBufferSize(size int) Option {
	return func(opt *configs) {
		opt.writeBufferSize = size
	}
}
