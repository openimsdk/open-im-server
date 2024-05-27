package dataver

import (
	"github.com/openimsdk/tools/utils/datautil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SyncResult struct {
	Version   uint
	VersionID string
	DeleteEID []string
	Changes   []string
	Full      bool
}

func VersionIDStr(id primitive.ObjectID) string {
	if id.IsZero() {
		return ""
	}
	return id.String()
}

func NewSyncResult(wl *WriteLog, fullIDs []string, versionID string) *SyncResult {
	var findEIDs []string
	var res SyncResult
	if wl.Full() || VersionIDStr(wl.ID) != versionID {
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
