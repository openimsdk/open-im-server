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

// export ConversationAPI
type ConversationAPI rpcclient.Conversation

// NewConversationAPI creates a new Conversation
func NewConversationAPI(discov discoveryregistry.SvcDiscoveryRegistry) ConversationAPI {
	return ConversationAPI(*rpcclient.NewConversation(discov))
}

// get all conversation
func (o *ConversationAPI) GetAllConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetAllConversations, o.Client, c)
}

// get conversation
func (o *ConversationAPI) GetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversation, o.Client, c)
}

// getConversations returns a list of conversations.
func (o *ConversationAPI) GetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversations, o.Client, c)
}

// set list of conversations
func (o *ConversationAPI) BatchSetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.BatchSetConversations, o.Client, c)
}

// setRecvMsgOpt sets the option for receiving messages.
func (o *ConversationAPI) SetRecvMsgOpt(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetRecvMsgOpt, o.Client, c)
}

// modifyConversationField modifies the conversation field.
func (o *ConversationAPI) ModifyConversationField(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.ModifyConversationField, o.Client, c)
}

// set conversations
func (o *ConversationAPI) SetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversations, o.Client, c)
}
