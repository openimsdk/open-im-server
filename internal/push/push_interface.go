package push

import "Open_IM/internal/push/logic"

type OfflinePusher interface {
	Push(userIDList []string, alert, detailContent, operationID string, opts logic.PushOpts) (resp string, err error)
}
