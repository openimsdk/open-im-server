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

package utils

import "github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"

type Options map[string]bool
type OptionsOpt func(Options)

// new option
func NewOptions(opts ...OptionsOpt) Options {
	options := make(map[string]bool, 11)
	options[constant.IsNotNotification] = false
	options[constant.IsSendMsg] = false
	options[constant.IsHistory] = false
	options[constant.IsPersistent] = false
	options[constant.IsOfflinePush] = false
	options[constant.IsUnreadCount] = false
	options[constant.IsConversationUpdate] = false
	options[constant.IsSenderSync] = false
	options[constant.IsNotPrivate] = false
	options[constant.IsSenderConversationUpdate] = false
	options[constant.IsSenderNotificationPush] = false
	options[constant.IsReactionFromCache] = false
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// new message option
func NewMsgOptions() Options {
	options := make(map[string]bool, 11)
	options[constant.IsOfflinePush] = false
	return make(map[string]bool)
}

// WithOptions
func WithOptions(options Options, opts ...OptionsOpt) Options {
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// WithNotNotification
func WithNotNotification(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsNotNotification] = b
	}
}

// WithSendMsg
func WithSendMsg(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsSendMsg] = b
	}
}

// WithHistory
func WithHistory(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsHistory] = b
	}
}

// WithPersistent
func WithPersistent() OptionsOpt {
	return func(options Options) {
		options[constant.IsPersistent] = true
	}
}

// WithOfflinePush
func WithOfflinePush(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsOfflinePush] = b
	}
}

// unread count
func WithUnreadCount(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsUnreadCount] = b
	}
}

// WithConversationUpdate
func WithConversationUpdate() OptionsOpt {
	return func(options Options) {
		options[constant.IsConversationUpdate] = true
	}
}

// WithSenderSync
func WithSenderSync() OptionsOpt {
	return func(options Options) {
		options[constant.IsSenderSync] = true
	}
}

// WithNotPrivate
func WithNotPrivate() OptionsOpt {
	return func(options Options) {
		options[constant.IsNotPrivate] = true
	}
}

// WithSenderConversationUpdate
func WithSenderConversationUpdate() OptionsOpt {
	return func(options Options) {
		options[constant.IsSenderConversationUpdate] = true
	}
}

// WithSenderNotificationPush
func WithSenderNotificationPush() OptionsOpt {
	return func(options Options) {
		options[constant.IsSenderNotificationPush] = true
	}
}

// react from cache is or not
func WithReactionFromCache() OptionsOpt {
	return func(options Options) {
		options[constant.IsReactionFromCache] = true
	}
}

// is or not exit in map named o
func (o Options) Is(notification string) bool {
	v, ok := o[notification]
	if !ok || v {
		return true
	}
	return false
}

// is or nit notification
func (o Options) IsNotNotification() bool {
	return o.Is(constant.IsNotNotification)
}

// is or not send msg
func (o Options) IsSendMsg() bool {
	return o.Is(constant.IsSendMsg)
}

// is or not a history
func (o Options) IsHistory() bool {
	return o.Is(constant.IsHistory)
}

// is or not persistent
func (o Options) IsPersistent() bool {
	return o.Is(constant.IsPersistent)
}

// is oor not push offline
func (o Options) IsOfflinePush() bool {
	return o.Is(constant.IsOfflinePush)
}

// unread count
func (o Options) IsUnreadCount() bool {
	return o.Is(constant.IsUnreadCount)
}

// is or not conversation update
func (o Options) IsConversationUpdate() bool {
	return o.Is(constant.IsConversationUpdate)
}

// is or not send async
func (o Options) IsSenderSync() bool {
	return o.Is(constant.IsSenderSync)
}

// is or not private
func (o Options) IsNotPrivate() bool {
	return o.Is(constant.IsNotPrivate)
}

// is or not notification push update
func (o Options) IsSenderConversationUpdate() bool {
	return o.Is(constant.IsSenderConversationUpdate)
}

// is or not notification push sender
func (o Options) IsSenderNotificationPush() bool {
	return o.Is(constant.IsSenderNotificationPush)
}

// reaction is or not from cache
func (o Options) IsReactionFromCache() bool {
	return o.Is(constant.IsReactionFromCache)
}
