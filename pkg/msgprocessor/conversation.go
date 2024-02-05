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

package msgprocessor

import (
	"sort"
	"strings"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/sdkws"
	"google.golang.org/protobuf/proto"
)

func GetNotificationConversationIDByMsg(msg *sdkws.MsgData) string {
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		return "n_" + strings.Join(l, "_")
	case constant.GroupChatType:
		return "n_" + msg.GroupID
	case constant.SuperGroupChatType:
		return "n_" + msg.GroupID
	case constant.NotificationChatType:
		return "n_" + msg.SendID + "_" + msg.RecvID
	}
	return ""
}

func GetChatConversationIDByMsg(msg *sdkws.MsgData) string {
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		return "si_" + strings.Join(l, "_")
	case constant.GroupChatType:
		return "g_" + msg.GroupID
	case constant.SuperGroupChatType:
		return "sg_" + msg.GroupID
	case constant.NotificationChatType:
		return "sn_" + msg.SendID + "_" + msg.RecvID
	}

	return ""
}

func GenConversationUniqueKey(msg *sdkws.MsgData) string {
	switch msg.SessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		return strings.Join(l, "_")
	case constant.SuperGroupChatType:
		return msg.GroupID
	}
	return ""
}

func GetConversationIDByMsg(msg *sdkws.MsgData) string {
	options := Options(msg.Options)
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		if !options.IsNotNotification() {
			return "n_" + strings.Join(l, "_")
		}
		return "si_" + strings.Join(l, "_") // single chat
	case constant.GroupChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.GroupID // group chat
		}
		return "g_" + msg.GroupID // group chat
	case constant.SuperGroupChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.GroupID // super group chat
		}
		return "sg_" + msg.GroupID // super group chat
	case constant.NotificationChatType:
		if !options.IsNotNotification() {
			return "n_" + msg.SendID + "_" + msg.RecvID // super group chat
		}
		return "sn_" + msg.SendID + "_" + msg.RecvID // server notification chat
	}
	return ""
}

func GetConversationIDBySessionType(sessionType int, ids ...string) string {
	sort.Strings(ids)
	if len(ids) > 2 || len(ids) < 1 {
		return ""
	}
	switch sessionType {
	case constant.SingleChatType:
		return "si_" + strings.Join(ids, "_") // single chat
	case constant.GroupChatType:
		return "g_" + ids[0] // group chat
	case constant.SuperGroupChatType:
		return "sg_" + ids[0] // super group chat
	case constant.NotificationChatType:
		return "sn_" + ids[0] // server notification chat
	}
	return ""
}

func GetNotificationConversationIDByConversationID(conversationID string) string {
	l := strings.Split(conversationID, "_")
	if len(l) > 1 {
		l[0] = "n"
		return strings.Join(l, "_")
	}

	return ""
}

func GetNotificationConversationID(sessionType int, ids ...string) string {
	sort.Strings(ids)
	if len(ids) > 2 || len(ids) < 1 {
		return ""
	}
	switch sessionType {
	case constant.SingleChatType:
		return "n_" + strings.Join(ids, "_") // single chat
	case constant.SuperGroupChatType:
		return "n_" + ids[0] // super group chat
	}
	return ""
}

func IsNotification(conversationID string) bool {
	return strings.HasPrefix(conversationID, "n_")
}

func IsNotificationByMsg(msg *sdkws.MsgData) bool {
	return !Options(msg.Options).IsNotNotification()
}

func ParseConversationID(msg *sdkws.MsgData) (isNotification bool, conversationID string) {
	options := Options(msg.Options)
	switch msg.SessionType {
	case constant.SingleChatType:
		l := []string{msg.SendID, msg.RecvID}
		sort.Strings(l)
		if !options.IsNotNotification() {
			return true, "n_" + strings.Join(l, "_")
		}
		return false, "si_" + strings.Join(l, "_") // single chat
	case constant.SuperGroupChatType:
		if !options.IsNotNotification() {
			return true, "n_" + msg.GroupID // super group chat
		}
		return false, "sg_" + msg.GroupID // super group chat
	case constant.NotificationChatType:
		if !options.IsNotNotification() {
			return true, "n_" + msg.SendID + "_" + msg.RecvID // super group chat
		}
		return false, "sn_" + msg.SendID + "_" + msg.RecvID // server notification chat
	}
	return false, ""
}

type MsgBySeq []*sdkws.MsgData

func (s MsgBySeq) Len() int {
	return len(s)
}

func (s MsgBySeq) Less(i, j int) bool {
	return s[i].Seq < s[j].Seq
}

func (s MsgBySeq) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func Pb2String(pb proto.Message) (string, error) {
	s, err := proto.Marshal(pb)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func String2Pb(s string, pb proto.Message) error {
	return proto.Unmarshal([]byte(s), pb)
}
