package push

type OfflinePusher interface {
	Push(userIDList []string, alert, detailContent, operationID string, opts PushOpts) (resp string, err error)
}

type PushOpts struct {
	Signal Signal
}

type Signal struct {
	ClientMsgID string
}
