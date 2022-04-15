package push

type OfflinePusher interface {
	Push(userIDList []string, alert, detailContent, operationID string) (resp string, err error)
}
