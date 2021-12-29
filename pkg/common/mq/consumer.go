package mq

import "time"

type Consumer interface {
	// RegisterMessageHandler is used to register message handler
	// any received messages will be passed to handler to process
	// once the Consumer started, it is forbidden to register handlers.
	RegisterMessageHandler(topic string, handler MessageHandler)

	// Start to consume messages
	Start() error
}

type MessageHandler interface {
	// HandleMessage process received messages,
	// if returned error is nil, the message will be auto committed.
	HandleMessage(msg *Message) error
}

type MessageHandleFunc func(msg *Message) error

func (fn MessageHandleFunc) HandleMessage(msg *Message) error {
	return fn(msg)
}

type Message struct {
	Key, Value     []byte
	Topic          string
	Partition      int32
	Offset         int64
	Timestamp      time.Time       // only set if kafka is version 0.10+, inner message timestamp
	BlockTimestamp time.Time       // only set if kafka is version 0.10+, outer (compressed) block timestamp
	Headers        []*RecordHeader // only set if kafka is version 0.11+
}

type RecordHeader struct {
	Key   []byte
	Value []byte
}
