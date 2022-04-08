package push

type OfflinePusher interface {
	Push(userIDList []string, alert, detailContent, platform, operationID string) (resp string, err error)
}
