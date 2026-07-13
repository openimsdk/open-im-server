package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
)

func NewStreamMsg(rdb redis.UniversalClient) cache.StreamMsgCache {
	return &streamMsg{rdb: rdb}
}

type streamMsg struct {
	rdb redis.UniversalClient
}

func (x *streamMsg) getMsgKey(conversationID string, clientMsgID string) string {
	return cachekey.GetStreamMsgKey(conversationID, clientMsgID)
}

func (x *streamMsg) CreateStreamMsg(ctx context.Context, conversationID string, clientMsgID string, msg *cache.StreamMsg) error {
	key := x.getMsgKey(conversationID, clientMsgID)
	pipeline := x.rdb.Pipeline()
	pipeline.HSet(ctx, key, "sendUserID", msg.SendUserID)
	pipeline.HSet(ctx, key, "recvID", msg.RecvID)
	pipeline.HSet(ctx, key, "sessionType", strconv.Itoa(int(msg.SessionType)))
	pipeline.HSet(ctx, key, "updateTime", time.Now().UnixMilli())
	pipeline.HSet(ctx, key, "isEnd", false)
	pipeline.HSet(ctx, key, "streamType", msg.StreamType)
	pipeline.HSet(ctx, key, "streamContent", msg.StreamContent)
	pipeline.Expire(ctx, key, 24*time.Hour)
	_, err := pipeline.Exec(ctx)
	return err
}

func (x *streamMsg) AppendStreamMsg(ctx context.Context, conversationID string, clientMsgID string, startIndex int, packets []string, end bool, retPacket bool) (*cache.StreamMsg, error) {
	key := x.getMsgKey(conversationID, clientMsgID)
	var mapCmd *redis.MapStringStringCmd
	var sliceCmd *redis.SliceCmd
	pipeline := x.rdb.Pipeline()
	for i, packet := range packets {
		pipeline.HSet(ctx, key, "i_"+strconv.Itoa(startIndex+i), packet)
	}
	pipeline.HSet(ctx, key, "isEnd", end)
	pipeline.HSet(ctx, key, "updateTime", time.Now().UnixMilli())
	pipeline.Expire(ctx, key, 24*time.Hour)
	if retPacket {
		mapCmd = pipeline.HGetAll(ctx, key)
	} else {
		sliceCmd = pipeline.HMGet(ctx, key, "sendUserID", "recvID", "sessionType")
	}
	if _, err := pipeline.Exec(ctx); err != nil {
		return nil, err
	}
	var data map[string]string
	var err error
	if retPacket {
		data, err = mapCmd.Result()
	} else {
		arr, resultErr := sliceCmd.Result()
		if resultErr != nil {
			return nil, resultErr
		}
		if len(arr) != 3 || arr[0] == nil || arr[1] == nil || arr[2] == nil {
			return nil, errs.ErrRecordNotFound.WrapMsg("stream message not found")
		}
		data = map[string]string{
			"sendUserID": fmt.Sprint(arr[0]), "recvID": fmt.Sprint(arr[1]), "sessionType": fmt.Sprint(arr[2]),
		}
	}
	if err != nil {
		return nil, err
	}
	return x.mapToStreamMsg(data, retPacket)
}

func (x *streamMsg) GetStreamMsg(ctx context.Context, conversationID string, clientMsgID string) (*cache.StreamMsg, error) {
	data, err := x.rdb.HGetAll(ctx, x.getMsgKey(conversationID, clientMsgID)).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("stream message not found")
	}
	return x.mapToStreamMsg(data, true)
}

func (x *streamMsg) mapToStreamMsg(data map[string]string, full bool) (*cache.StreamMsg, error) {
	sessionType, err := strconv.ParseInt(data["sessionType"], 10, 32)
	if err != nil {
		return nil, err
	}
	if !full {
		return &cache.StreamMsg{SendUserID: data["sendUserID"], RecvID: data["recvID"], SessionType: int32(sessionType)}, nil
	}
	end, err := strconv.ParseBool(data["isEnd"])
	if err != nil {
		return nil, err
	}
	updateTime, err := strconv.ParseInt(data["updateTime"], 10, 64)
	if err != nil {
		return nil, err
	}
	msg := cache.StreamMsg{
		SendUserID: data["sendUserID"], RecvID: data["recvID"], SessionType: int32(sessionType),
		StreamType: data["streamType"], StreamContent: data["streamContent"], UpdateTime: updateTime,
		Packets: make(map[int64]string), End: end,
	}
	var maxIndex int64 = -1
	for indexStr, value := range data {
		if !strings.HasPrefix(indexStr, "i_") {
			continue
		}
		index, err := strconv.ParseInt(strings.TrimPrefix(indexStr, "i_"), 10, 64)
		if err != nil || index < 0 {
			return nil, errs.ErrInternalServer.WrapMsg("packet index is invalid", "index", indexStr)
		}
		msg.Packets[index] = value
		if maxIndex < index {
			maxIndex = index
		}
	}
	for i := int64(0); i <= maxIndex; i++ {
		if _, ok := msg.Packets[i]; !ok {
			return nil, errs.ErrInternalServer.WrapMsg("packet index is not continuous", "index", i)
		}
	}
	return &msg, nil
}

func (x *streamMsg) GetStreamMsgEnd(ctx context.Context, conversationID string, clientMsgID string) (bool, error) {
	return x.rdb.HGet(ctx, x.getMsgKey(conversationID, clientMsgID), "isEnd").Bool()
}
