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

package msg

import (
	"context"
	"errors"
	"sort"

	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/redis/go-redis/v9"
)

func (m *msgServer) GetConversationMaxSeq(ctx context.Context, req *pbmsg.GetConversationMaxSeqReq) (*pbmsg.GetConversationMaxSeqResp, error) {
	maxSeq, err := m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	return &pbmsg.GetConversationMaxSeqResp{MaxSeq: maxSeq}, nil
}

func (m *msgServer) GetMaxSeqs(ctx context.Context, req *pbmsg.GetMaxSeqsReq) (*pbmsg.SeqsInfoResp, error) {
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	return &pbmsg.SeqsInfoResp{MaxSeqs: maxSeqs}, nil
}

func (m *msgServer) GetHasReadSeqs(ctx context.Context, req *pbmsg.GetHasReadSeqsReq) (*pbmsg.SeqsInfoResp, error) {
	hasReadSeqs, err := m.MsgDatabase.GetHasReadSeqs(ctx, req.UserID, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	return &pbmsg.SeqsInfoResp{MaxSeqs: hasReadSeqs}, nil
}

func (m *msgServer) GetMsgByConversationIDs(ctx context.Context, req *pbmsg.GetMsgByConversationIDsReq) (*pbmsg.GetMsgByConversationIDsResp, error) {
	Msgs, err := m.MsgDatabase.FindOneByDocIDs(ctx, req.ConversationIDs, req.MaxSeqs)
	if err != nil {
		return nil, err
	}
	return &pbmsg.GetMsgByConversationIDsResp{MsgDatas: Msgs}, nil
}

func (m *msgServer) SetUserConversationsMinSeq(ctx context.Context, req *pbmsg.SetUserConversationsMinSeqReq) (*pbmsg.SetUserConversationsMinSeqResp, error) {
	for _, userID := range req.UserIDs {
		if err := m.MsgDatabase.SetUserConversationsMinSeqs(ctx, userID, map[string]int64{req.ConversationID: req.Seq}); err != nil {
			return nil, err
		}
	}
	return &pbmsg.SetUserConversationsMinSeqResp{}, nil
}

func (m *msgServer) GetActiveConversation(ctx context.Context, req *pbmsg.GetActiveConversationReq) (*pbmsg.GetActiveConversationResp, error) {
	res, err := m.MsgDatabase.GetCacheMaxSeqWithTime(ctx, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	conversations := make([]*pbmsg.ActiveConversation, 0, len(res))
	for conversationID, val := range res {
		conversations = append(conversations, &pbmsg.ActiveConversation{
			MaxSeq:         val.Seq,
			LastTime:       val.Time,
			ConversationID: conversationID,
		})
	}
	if req.Limit > 0 {
		sort.Sort(activeConversations(conversations))
		if len(conversations) > int(req.Limit) {
			conversations = conversations[:req.Limit]
		}
	}
	return &pbmsg.GetActiveConversationResp{Conversations: conversations}, nil
}

func (m *msgServer) SetUserConversationMaxSeq(ctx context.Context, req *pbmsg.SetUserConversationMaxSeqReq) (*pbmsg.SetUserConversationMaxSeqResp, error) {
	for _, userID := range req.OwnerUserID {
		if err := m.MsgDatabase.SetUserConversationsMaxSeq(ctx, req.ConversationID, userID, req.MaxSeq); err != nil {
			return nil, err
		}
	}
	return &pbmsg.SetUserConversationMaxSeqResp{}, nil
}

func (m *msgServer) SetUserConversationMinSeq(ctx context.Context, req *pbmsg.SetUserConversationMinSeqReq) (*pbmsg.SetUserConversationMinSeqResp, error) {
	for _, userID := range req.OwnerUserID {
		if err := m.MsgDatabase.SetUserConversationsMinSeq(ctx, req.ConversationID, userID, req.MinSeq); err != nil {
			return nil, err
		}
	}
	return &pbmsg.SetUserConversationMinSeqResp{}, nil
}
