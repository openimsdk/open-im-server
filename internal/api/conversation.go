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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

type ConversationApi rpcclient.Conversation

func NewConversationApi(discov discoveryregistry.SvcDiscoveryRegistry) ConversationApi {
	return ConversationApi(*rpcclient.NewConversation(discov))
}

func (o *ConversationApi) GetAllConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetAllConversations, o.Client, c)
}

func (o *ConversationApi) GetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversation, o.Client, c)
}

func (o *ConversationApi) GetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversations, o.Client, c)
}

// deprecated
func (o *ConversationApi) SetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversation, o.Client, c)
}

// deprecated
func (o *ConversationApi) BatchSetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.BatchSetConversations, o.Client, c)
}

func (o *ConversationApi) SetRecvMsgOpt(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetRecvMsgOpt, o.Client, c)
}

func (o *ConversationApi) ModifyConversationField(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.ModifyConversationField, o.Client, c)
}

func (o *ConversationApi) SetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversations, o.Client, c)
}
