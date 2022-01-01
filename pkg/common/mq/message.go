package mq

import "time"

type Message struct {
	Key, Value     []byte
	Topic          string
	Partition      int32
	Offset         int64
	Timestamp      time.Time
	BlockTimestamp time.Time
	Headers        []*RecordHeader
}

type RecordHeader struct {
	Key   []byte
	Value []byte
}
