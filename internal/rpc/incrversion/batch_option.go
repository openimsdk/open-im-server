package incrversion

import (
	"context"
	"fmt"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BatchOption[A, B any] struct {
	Ctx            context.Context
	VersionKeys    []string
	VersionIDs     []string
	VersionNumbers []uint64
	//SyncLimit       int
	Versions         func(ctx context.Context, dIds []string, versions []uint64, limits []int) (map[string]*model.VersionLog, error)
	CacheMaxVersions func(ctx context.Context, dIds []string) (map[string]*model.VersionLog, error)
	//SortID          func(ctx context.Context, dId string) ([]string, error)
	Find func(ctx context.Context, ids []string) ([]A, error)
	// Resp func(version map[string]*model.VersionLog, deleteIds, insertList, updateList []A, full bool) []*B
	Resp func(versions map[string]*model.VersionLog, deleteIdsMap map[string][]string, insertListMap, updateListMap map[string][]A, fullMap map[string]bool) *B
}

func (o *BatchOption[A, B]) newError(msg string) error {
	return errs.ErrInternalServer.WrapMsg(msg)
}

func (o *BatchOption[A, B]) check() error {
	if o.Ctx == nil {
		return o.newError("opt ctx is nil")
	}
	if len(o.VersionKeys) == 0 {
		return o.newError("versionKeys is empty")
	}
	if o.Versions == nil {
		return o.newError("func versions is nil")
	}
	if o.Find == nil {
		return o.newError("func find is nil")
	}
	if o.Resp == nil {
		return o.newError("func resp is nil")
	}
	return nil
}

func (o *BatchOption[A, B]) validVersions() []bool {
	valids := make([]bool, len(o.VersionIDs))
	for i, versionID := range o.VersionIDs {
		objID, err := primitive.ObjectIDFromHex(versionID)
		valids[i] = err == nil && (!objID.IsZero()) && o.VersionNumbers[i] > 0
	}
	return valids
}

func (o *BatchOption[A, B]) equalIDs(objIDs []primitive.ObjectID) []bool {
	equals := make([]bool, len(o.VersionIDs))
	for i, versionID := range o.VersionIDs {
		equals[i] = versionID == objIDs[i].Hex()
	}
	return equals
}

func (o *BatchOption[A, B]) getVersions(tags *[]int) (versions map[string]*model.VersionLog, err error) {
	valids := o.validVersions()

	var dIDs []string
	var versionNums []uint64
	var limits []int

	if o.CacheMaxVersions == nil {
		for i, valid := range valids {
			if valid {
				(*tags)[i] = tagQuery
				dIDs = append(dIDs, o.VersionKeys[i])
				versionNums = append(versionNums, o.VersionNumbers[i])
				limits = append(limits, syncLimit)

				// version, err := o.Versions(o.Ctx, []string{o.VersionKeys[i]}, []uint64{o.VersionNumbers[i]}, syncLimit)
				// if err != nil {
				// 	return nil, err
				// }
				// versions[o.VersionKeys[i]] = version[o.VersionKeys[i]]
			} else {
				(*tags)[i] = tagFull
				dIDs = append(dIDs, o.VersionKeys[i])
				versionNums = append(versionNums, 0)
				limits = append(limits, 0)

				// version, err := o.Versions(o.Ctx, []string{o.VersionKeys[i]}, []uint64{0}, 0)
				// if err != nil {
				// 	return nil, err
				// }
				// versions[o.VersionKeys[i]] = version[o.VersionKeys[i]]
			}
		}
		versions, err = o.Versions(o.Ctx, dIDs, versionNums, limits)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		return versions, nil
	} else {
		caches, err := o.CacheMaxVersions(o.Ctx, o.VersionKeys)
		if err != nil {
			return nil, err
		}
		objIDs := make([]primitive.ObjectID, len(o.VersionIDs))
		for i, versionID := range o.VersionIDs {
			objID, _ := primitive.ObjectIDFromHex(versionID)
			objIDs[i] = objID
		}
		equals := o.equalIDs(objIDs)
		for i, valid := range valids {
			if !valid {
				(*tags)[i] = tagFull
				// versions[o.VersionKeys[i]] = caches[o.VersionKeys[i]]
			} else if !equals[i] {
				(*tags)[i] = tagFull
				// versions[o.VersionKeys[i]] = caches[o.VersionKeys[i]]
			} else if o.VersionNumbers[i] == uint64(caches[o.VersionKeys[i]].Version) {
				(*tags)[i] = tagEqual
				// versions[o.VersionKeys[i]] = caches[o.VersionKeys[i]]
			} else {
				(*tags)[i] = tagQuery
				dIDs = append(dIDs, o.VersionKeys[i])
				versionNums = append(versionNums, o.VersionNumbers[i])
				limits = append(limits, syncLimit)

				delete(caches, o.VersionKeys[i])
				// versions[o.VersionKeys[i]] = version[o.VersionKeys[i]]
			}
		}
		version, err := o.Versions(o.Ctx, dIDs, versionNums, limits)
		if err != nil {
			return nil, err
		}
		for k, v := range version {
			caches[k] = v
		}
		versions = caches
	}
	return versions, nil
}

func (o *BatchOption[A, B]) Build() (*B, error) {
	if err := o.check(); err != nil {
		return nil, err
	}

	tags := make([]int, len(o.VersionKeys))
	versions, err := o.getVersions(&tags)
	if err != nil {
		return nil, err
	}

	fullMap := make(map[string]bool)
	for i, tag := range tags {
		switch tag {
		case tagQuery:
			version := versions[o.VersionKeys[i]]
			fullMap[o.VersionKeys[i]] = version.ID.Hex() != o.VersionIDs[i] || uint64(version.Version) < o.VersionNumbers[i] || len(version.Logs) != version.LogLen
		case tagFull:
			fullMap[o.VersionKeys[i]] = true
		case tagEqual:
			fullMap[o.VersionKeys[i]] = false
		default:
			panic(fmt.Errorf("undefined tag %d", tag))
		}
	}

	var (
		insertIdsMap = make(map[string][]string)
		deleteIdsMap = make(map[string][]string)
		updateIdsMap = make(map[string][]string)
	)

	for _, versionKey := range o.VersionKeys {
		if !fullMap[versionKey] {
			version := versions[versionKey]
			insertIds, deleteIds, updateIds := version.DeleteAndChangeIDs()
			insertIdsMap[versionKey] = insertIds
			deleteIdsMap[versionKey] = deleteIds
			updateIdsMap[versionKey] = updateIds
		}
	}

	var (
		insertListMap = make(map[string][]A)
		updateListMap = make(map[string][]A)
	)

	for versionKey, insertIds := range insertIdsMap {
		if len(insertIds) > 0 {
			insertList, err := o.Find(o.Ctx, insertIds)
			if err != nil {
				return nil, err
			}
			insertListMap[versionKey] = insertList
		}
	}

	for versionKey, updateIds := range updateIdsMap {
		if len(updateIds) > 0 {
			updateList, err := o.Find(o.Ctx, updateIds)
			if err != nil {
				return nil, err
			}
			updateListMap[versionKey] = updateList
		}
	}

	return o.Resp(versions, deleteIdsMap, insertListMap, updateListMap, fullMap), nil
}

// for _, versionLog := range versionLogs {
// 	if versionLog != nil {
// 		if !full {

// 		}
// 		insertIds, deleteIds, updateIds = append(insertIds, versionLog.InsertID...), append(deleteIds, versionLog.DeleteIDs...), append(updateIds, versionLog.UpdateIDs...)
// 	}
// }

// insertList, err := o.Find(o.Ctx, insertIds)
// if err != nil {
// 	return nil, err
// }

// updateList, err := o.Find(o.Ctx, updateIds)
// if err != nil {
// 	return nil, err
// }

// full := len(insertIds) > 0 || len(updateIds) > 0

// return o.Resp(versionLogs, deleteIds, insertList, updateList, full), nil
