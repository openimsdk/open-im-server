package push

import "Open_IM/pkg/common/constant"

var PushTerminal = []int{constant.IOSPlatformID, constant.AndroidPlatformID}

type OfflinePusher interface {
	Push(userIDList []string, title, detailContent, operationID string, opts PushOpts) (resp string, err error)
}

type PushOpts struct {
	Signal        Signal
	IOSPushSound  string
	IOSBadgeCount bool
	Data          string
}

type Signal struct {
	ClientMsgID string
}
