package utils

import "github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"

type Options map[string]bool
type OptionsOpt func(Options)

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

func NewMsgOptions() Options {
	options := make(map[string]bool, 11)
	options[constant.IsOfflinePush] = false
	return make(map[string]bool)
}

func WithOptions(options Options, opts ...OptionsOpt) Options {
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithNotNotification(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsNotNotification] = b
	}
}

func WithSendMsg(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsSendMsg] = b
	}
}

func WithHistory(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsHistory] = b
	}
}

func WithPersistent() OptionsOpt {
	return func(options Options) {
		options[constant.IsPersistent] = true
	}
}

func WithOfflinePush(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsOfflinePush] = b
	}
}

func WithUnreadCount(b bool) OptionsOpt {
	return func(options Options) {
		options[constant.IsUnreadCount] = b
	}
}

func WithConversationUpdate() OptionsOpt {
	return func(options Options) {
		options[constant.IsConversationUpdate] = true
	}
}

func WithSenderSync() OptionsOpt {
	return func(options Options) {
		options[constant.IsSenderSync] = true
	}
}

func WithNotPrivate() OptionsOpt {
	return func(options Options) {
		options[constant.IsNotPrivate] = true
	}
}

func WithSenderConversationUpdate() OptionsOpt {
	return func(options Options) {
		options[constant.IsSenderConversationUpdate] = true
	}
}

func WithSenderNotificationPush() OptionsOpt {
	return func(options Options) {
		options[constant.IsSenderNotificationPush] = true
	}
}

func WithReactionFromCache() OptionsOpt {
	return func(options Options) {
		options[constant.IsReactionFromCache] = true
	}
}

func (o Options) Is(notification string) bool {
	v, ok := o[notification]
	if !ok || v {
		return true
	}
	return false
}

func (o Options) IsNotNotification() bool {
	return o.Is(constant.IsNotNotification)
}

func (o Options) IsSendMsg() bool {
	return o.Is(constant.IsSendMsg)
}

func (o Options) IsHistory() bool {
	return o.Is(constant.IsHistory)
}

func (o Options) IsPersistent() bool {
	return o.Is(constant.IsPersistent)
}

func (o Options) IsOfflinePush() bool {
	return o.Is(constant.IsOfflinePush)
}

func (o Options) IsUnreadCount() bool {
	return o.Is(constant.IsUnreadCount)
}

func (o Options) IsConversationUpdate() bool {
	return o.Is(constant.IsConversationUpdate)
}

func (o Options) IsSenderSync() bool {
	return o.Is(constant.IsSenderSync)
}

func (o Options) IsNotPrivate() bool {
	return o.Is(constant.IsNotPrivate)
}

func (o Options) IsSenderConversationUpdate() bool {
	return o.Is(constant.IsSenderConversationUpdate)
}

func (o Options) IsSenderNotificationPush() bool {
	return o.Is(constant.IsSenderNotificationPush)
}

func (o Options) IsReactionFromCache() bool {
	return o.Is(constant.IsReactionFromCache)
}
