package dataver

import "github.com/openimsdk/tools/utils/datautil"

type SyncResult struct {
	Version   uint
	DeleteEID []string
	Changes   []string
	Full      bool
}

func NewSyncResult(wl *WriteLog, fullIDs []string) *SyncResult {
	var findEIDs []string
	var res SyncResult
	if wl.Full() {
		res.Changes = fullIDs
		res.Full = true
	} else {
		idSet := datautil.SliceSet(fullIDs)
		for _, l := range wl.Logs {
			if l.Deleted {
				res.DeleteEID = append(res.DeleteEID, l.EID)
			} else {
				if _, ok := idSet[l.EID]; ok {
					findEIDs = append(findEIDs, l.EID)
				}
			}
		}
	}
	return &res
}
