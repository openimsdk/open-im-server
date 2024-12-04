package incrversion

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type BatchOption[A, B any] struct {
	Ctx            context.Context
	TargetKeys     []string
	VersionIDs     []string
	VersionNumbers []uint64
	//SyncLimit       int
	Versions         func(ctx context.Context, dIds []string, versions []uint64, limits []int) (map[string]*model.VersionLog, error)
	CacheMaxVersions func(ctx context.Context, dIds []string) (map[string]*model.VersionLog, error)
	Find             func(ctx context.Context, dId string, ids []string) (A, error)
	Resp             func(versionsMap map[string]*model.VersionLog, deleteIdsMap map[string][]string, insertListMap, updateListMap map[string]A, fullMap map[string]bool) *B
}

func (o *BatchOption[A, B]) newError(msg string) error {
	return errs.ErrInternalServer.WrapMsg(msg)
}

func (o *BatchOption[A, B]) check() error {
	if o.Ctx == nil {
		return o.newError("opt ctx is nil")
	}
	if len(o.TargetKeys) == 0 {
		return o.newError("targetKeys is empty")
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
		valids[i] = (err == nil && (!objID.IsZero()) && o.VersionNumbers[i] > 0)
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
	var dIDs []string
	var versionNums []uint64
	var limits []int

	valids := o.validVersions()

	if o.CacheMaxVersions == nil {
		for i, valid := range valids {
			if valid {
				(*tags)[i] = tagQuery
				dIDs = append(dIDs, o.TargetKeys[i])
				versionNums = append(versionNums, o.VersionNumbers[i])
				limits = append(limits, syncLimit)
			} else {
				(*tags)[i] = tagFull
				dIDs = append(dIDs, o.TargetKeys[i])
				versionNums = append(versionNums, 0)
				limits = append(limits, 0)
			}
		}

		versions, err = o.Versions(o.Ctx, dIDs, versionNums, limits)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		return versions, nil

	} else {
		caches, err := o.CacheMaxVersions(o.Ctx, o.TargetKeys)
		if err != nil {
			return nil, errs.Wrap(err)
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
			} else if !equals[i] {
				(*tags)[i] = tagFull
			} else if o.VersionNumbers[i] == uint64(caches[o.TargetKeys[i]].Version) {
				(*tags)[i] = tagEqual
			} else {
				(*tags)[i] = tagQuery
				dIDs = append(dIDs, o.TargetKeys[i])
				versionNums = append(versionNums, o.VersionNumbers[i])
				limits = append(limits, syncLimit)

				delete(caches, o.TargetKeys[i])
			}
		}

		if dIDs != nil {
			versionMap, err := o.Versions(o.Ctx, dIDs, versionNums, limits)
			if err != nil {
				return nil, errs.Wrap(err)
			}

			for k, v := range versionMap {
				caches[k] = v
			}
		}

		versions = caches
	}
	return versions, nil
}

func (o *BatchOption[A, B]) Build() (*B, error) {
	if err := o.check(); err != nil {
		return nil, errs.Wrap(err)
	}

	tags := make([]int, len(o.TargetKeys))
	versions, err := o.getVersions(&tags)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	fullMap := make(map[string]bool)
	for i, tag := range tags {
		switch tag {
		case tagQuery:
			vLog := versions[o.TargetKeys[i]]
			fullMap[o.TargetKeys[i]] = vLog.ID.Hex() != o.VersionIDs[i] || uint64(vLog.Version) < o.VersionNumbers[i] || len(vLog.Logs) != vLog.LogLen
		case tagFull:
			fullMap[o.TargetKeys[i]] = true
		case tagEqual:
			fullMap[o.TargetKeys[i]] = false
		default:
			panic(fmt.Errorf("undefined tag %d", tag))
		}
	}

	var (
		insertIdsMap = make(map[string][]string)
		deleteIdsMap = make(map[string][]string)
		updateIdsMap = make(map[string][]string)
	)

	for _, targetKey := range o.TargetKeys {
		if !fullMap[targetKey] {
			version := versions[targetKey]
			insertIds, deleteIds, updateIds := version.DeleteAndChangeIDs()
			insertIdsMap[targetKey] = insertIds
			deleteIdsMap[targetKey] = deleteIds
			updateIdsMap[targetKey] = updateIds
		}
	}

	var (
		insertListMap = make(map[string]A)
		updateListMap = make(map[string]A)
	)

	for targetKey, insertIds := range insertIdsMap {
		if len(insertIds) > 0 {
			insertList, err := o.Find(o.Ctx, targetKey, insertIds)
			if err != nil {
				return nil, errs.Wrap(err)
			}
			insertListMap[targetKey] = insertList
		}
	}

	for targetKey, updateIds := range updateIdsMap {
		if len(updateIds) > 0 {
			updateList, err := o.Find(o.Ctx, targetKey, updateIds)
			if err != nil {
				return nil, errs.Wrap(err)
			}
			updateListMap[targetKey] = updateList
		}
	}

	return o.Resp(versions, deleteIdsMap, insertListMap, updateListMap, fullMap), nil
}
