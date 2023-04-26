package utils

import "github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"

type Options map[string]bool
type OptionsOpt func(Options)

func NewOptions(opts ...OptionsOpt) Options {
	options := make(map[string]bool, 11)
	options[constant.IsNotification] = false
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

func WithOptions(options Options, opts ...OptionsOpt) Options {
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithNotification() OptionsOpt {
	return func(options Options) {
		options[constant.IsNotification] = true
	}
}

func WithHistory() OptionsOpt {
	return func(options Options) {
		options[constant.IsHistory] = true
	}
}

func WithPersistent() OptionsOpt {
	return func(options Options) {
		options[constant.IsPersistent] = true
	}
}

func WithOfflinePush() OptionsOpt {
	return func(options Options) {
		options[constant.IsOfflinePush] = true
	}
}

func WithUnreadCount() OptionsOpt {
	return func(options Options) {
		options[constant.IsUnreadCount] = true
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

func (o Options) IsNotification() bool {
	return o.Is(constant.IsNotification)
}

func (o Options) IsHistory(options Options) bool {
	return o.Is(constant.IsHistory)
}

func (o Options) IsPersistent(options Options) bool {
	return o.Is(constant.IsPersistent)
}

func (o Options) IsOfflinePush(options Options) bool {
	return o.Is(constant.IsOfflinePush)
}

func (o Options) IsUnreadCount(options Options) bool {
	return o.Is(constant.IsUnreadCount)
}

func (o Options) IsConversationUpdate(options Options) bool {
	return o.Is(constant.IsConversationUpdate)
}

func (o Options) IsSenderSync(options Options) bool {
	return o.Is(constant.IsSenderSync)
}

func (o Options) IsNotPrivate(options Options) bool {
	return o.Is(constant.IsNotPrivate)
}

func (o Options) IsSenderConversationUpdate(options Options) bool {
	return o.Is(constant.IsSenderConversationUpdate)
}

func (o Options) IsSenderNotificationPush(options Options) bool {
	return o.Is(constant.IsSenderNotificationPush)
}

func (o Options) IsReactionFromCache(options Options) bool {
	return o.Is(constant.IsReactionFromCache)
}
