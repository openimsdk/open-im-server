
### mongo
```go
type Seq struct {
	ConversationID string `bson:"conversation_id"`
	MaxSeq         int64  `bson:"max_seq"`
	MinSeq         int64  `bson:"min_seq"`
}
```

```go
type Seq interface {
  Malloc(ctx context.Context, conversationID string, size int64) (int64, error)
  GetMaxSeq(ctx context.Context, conversationID string) (int64, error)
  GetMinSeq(ctx context.Context, conversationID string) (int64, error)
  SetMinSeq(ctx context.Context, conversationID string, seq int64) error
}
```

1. Malloc 申请seq数量，返回的已有seq的最大值，消息用的第一个seq是返回值+1
2. GetMaxSeq 获取申请的seq的最大值，在发消息的seq小于这个值
3. GetMinSeq 获取最小的seq，用于拉取历史消息
4. SetMinSeq 设置最小的seq，用于拉取历史消息

### redis
```go
type RedisSeq struct {
	Curr int64 // 当前的最大seq
	Last int64 // mongodb中申请的最大seq
	Lock *int64 // 锁，用于在mongodb中申请seq
}
```

1. Malloc 申请seq数量，返回的已有seq的最大值，消息用的第一个seq是返回值+1，如果redis中申请数量够用，直接返回，并自增对应数量。如果redis中申请数量不够用，加锁，从mongodb中申请seq。
2. GetMaxSeq 获取已发消息的最大seq就是Curr的值。如果redis中缓存不存在就通过mongodb获取最大seq。存储在redis中。其中Curr和Last都是这个seq值。
3. GetMinSeq, SetMinSeq用之前rockscache的方案。