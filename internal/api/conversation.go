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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/a2r"
)

type ConversationApi struct{}

func NewConversationApi() ConversationApi {
	return ConversationApi{}
}

func (o *ConversationApi) GetAllConversations(c *gin.Context) {
	a2r.CallV2(conversation.GetAllConversationsCaller.Invoke, c)
}

func (o *ConversationApi) GetSortedConversationList(c *gin.Context) {
	a2r.CallV2(conversation.GetSortedConversationListCaller.Invoke, c)
}

func (o *ConversationApi) GetConversation(c *gin.Context) {
	a2r.CallV2(conversation.GetConversationCaller.Invoke, c)
}

func (o *ConversationApi) GetConversations(c *gin.Context) {
	a2r.CallV2(conversation.GetConversationsCaller.Invoke, c)
}

func (o *ConversationApi) SetConversations(c *gin.Context) {
	a2r.CallV2(conversation.SetConversationsCaller.Invoke, c)
}

func (o *ConversationApi) GetConversationOfflinePushUserIDs(c *gin.Context) {
	a2r.CallV2(conversation.GetConversationOfflinePushUserIDsCaller.Invoke, c)
}

func (o *ConversationApi) GetFullOwnerConversationIDs(c *gin.Context) {
	a2r.CallV2(conversation.GetFullOwnerConversationIDsCaller.Invoke, c)
}

func (o *ConversationApi) GetIncrementalConversation(c *gin.Context) {
	a2r.CallV2(conversation.GetIncrementalConversationCaller.Invoke, c)
}

func (o *ConversationApi) GetOwnerConversation(c *gin.Context) {
	a2r.CallV2(conversation.GetOwnerConversationCaller.Invoke, c)
}

func (o *ConversationApi) GetNotNotifyConversationIDs(c *gin.Context) {
	a2r.CallV2(conversation.GetNotNotifyConversationIDsCaller.Invoke, c)
}

func (o *ConversationApi) GetPinnedConversationIDs(c *gin.Context) {
	a2r.CallV2(conversation.GetPinnedConversationIDsCaller.Invoke, c)
}
