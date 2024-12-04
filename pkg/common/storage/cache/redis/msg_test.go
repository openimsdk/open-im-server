// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/openimsdk/protocol/sdkws"
)

func Test_msgCache_SetMessagesToCache(t *testing.T) {
	type fields struct {
		rdb redis.UniversalClient
	}
	type args struct {
		ctx            context.Context
		conversationID string
		msgs           []*sdkws.MsgData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{"test1", fields{rdb: redis.NewClient(&redis.Options{Addr: "localhost:16379", Username: "", Password: "openIM123", DB: 0})}, args{context.Background(),
			"cid", []*sdkws.MsgData{{Seq: 1}, {Seq: 2}, {Seq: 3}}}, 3, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &msgCache{
				rdb: tt.fields.rdb,
			}
			got, err := c.SetMessagesToCache(tt.args.ctx, tt.args.conversationID, tt.args.msgs)
			if !tt.wantErr(t, err, fmt.Sprintf("SetMessagesToCache(%v, %v, %v)", tt.args.ctx, tt.args.conversationID, tt.args.msgs)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SetMessagesToCache(%v, %v, %v)", tt.args.ctx, tt.args.conversationID, tt.args.msgs)
		})
	}
}

func Test_msgCache_GetMessagesBySeq(t *testing.T) {
	type fields struct {
		rdb redis.UniversalClient
	}
	type args struct {
		ctx            context.Context
		conversationID string
		seqs           []int64
	}
	var failedSeq []int64
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantSeqMsgs    []*sdkws.MsgData
		wantFailedSeqs []int64
		wantErr        assert.ErrorAssertionFunc
	}{
		{"test1", fields{rdb: redis.NewClient(&redis.Options{Addr: "localhost:16379", Password: "openIM123", DB: 0})},
			args{context.Background(), "cid", []int64{1, 2, 3}},
			[]*sdkws.MsgData{{Seq: 1}, {Seq: 2}, {Seq: 3}}, failedSeq, assert.NoError},
		{"test2", fields{rdb: redis.NewClient(&redis.Options{Addr: "localhost:16379", Password: "openIM123", DB: 0})},
			args{context.Background(), "cid", []int64{4, 5, 6}},
			nil, []int64{4, 5, 6}, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &msgCache{
				rdb: tt.fields.rdb,
			}
			gotSeqMsgs, gotFailedSeqs, err := c.GetMessagesBySeq(tt.args.ctx, tt.args.conversationID, tt.args.seqs)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMessagesBySeq(%v, %v, %v)", tt.args.ctx, tt.args.conversationID, tt.args.seqs)) {
				return
			}
			equalMsgDataSlices(t, tt.wantSeqMsgs, gotSeqMsgs)
			assert.Equalf(t, tt.wantFailedSeqs, gotFailedSeqs, "GetMessagesBySeq(%v, %v, %v)", tt.args.ctx, tt.args.conversationID, tt.args.seqs)
		})
	}
}

func equalMsgDataSlices(t *testing.T, expected, actual []*sdkws.MsgData) {
	assert.Equal(t, len(expected), len(actual), "Slices have different lengths")
	for i := range expected {
		assert.True(t, proto.Equal(expected[i], actual[i]), "Element %d not equal: expected %v, got %v", i, expected[i], actual[i])
	}
}

func Test_msgCache_DeleteMessagesFromCache(t *testing.T) {
	type fields struct {
		rdb redis.UniversalClient
	}
	type args struct {
		ctx            context.Context
		conversationID string
		seqs           []int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"test1", fields{rdb: redis.NewClient(&redis.Options{Addr: "localhost:16379", Password: "openIM123"})},
			args{context.Background(), "cid", []int64{1, 2, 3}}, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &msgCache{
				rdb: tt.fields.rdb,
			}
			tt.wantErr(t, c.DeleteMessagesFromCache(tt.args.ctx, tt.args.conversationID, tt.args.seqs),
				fmt.Sprintf("DeleteMessagesFromCache(%v, %v, %v)", tt.args.ctx, tt.args.conversationID, tt.args.seqs))
		})
	}
}
