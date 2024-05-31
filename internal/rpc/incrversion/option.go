package incrversion

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/utils/datautil"
)

func Limit(maxSync int, version uint64) int {
	if version == 0 {
		return 0
	}
	return maxSync
}

type Option[A, B any] struct {
	VersionID string
	Version   func() (*model.VersionLog, error)
	AllID     func() ([]string, error)
	Find      func(ids []string) ([]A, error)
	ID        func(elem A) string
	Resp      func(version *model.VersionLog, delIDs []string, list []A, full bool) *B
}

func (o *Option[A, B]) Build() (*B, error) {
	version, err := o.Version()
	if err != nil {
		return nil, err
	}
	var (
		deleteIDs []string
		changeIDs []string
	)
	full := o.VersionID != version.ID.Hex() || version.Full()
	if full {
		changeIDs, err = o.AllID()
		if err != nil {
			return nil, err
		}
	} else {
		deleteIDs, changeIDs = version.DeleteAndChangeIDs()
	}
	var list []A
	if len(changeIDs) > 0 {
		list, err = o.Find(changeIDs)
		if err != nil {
			return nil, err
		}
		if (!full) && o.ID != nil && len(changeIDs) != len(list) {
			foundIDs := datautil.SliceSetAny(list, o.ID)
			for _, id := range changeIDs {
				if _, ok := foundIDs[id]; !ok {
					deleteIDs = append(deleteIDs, id)
				}
			}
		}
	}
	return o.Resp(version, deleteIDs, list, full), nil
}
