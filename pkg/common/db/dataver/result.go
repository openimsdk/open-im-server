package dataver

type SyncResult[T any] struct {
	Version   uint
	DeleteEID []string
	Changes   []T
	Full      bool
}

func NewSyncResult[T any](wl *WriteLog, find func(eIds []string) ([]T, error)) (*SyncResult[T], error) {
	var findEIDs []string
	var res SyncResult[T]
	if wl.Full() {
		res.Full = true
	} else {
		for _, l := range wl.Logs {
			if l.Deleted {
				res.DeleteEID = append(res.DeleteEID, l.EID)
			} else {
				findEIDs = append(findEIDs, l.EID)
			}
		}
	}
	if res.Full || len(findEIDs) > 0 {
		var err error
		res.Changes, err = find(findEIDs)
		if err != nil {
			return nil, err
		}
	}
	return &res, nil
}
