// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package third

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
)

func genLogID() string {
	const dataLen = 10
	data := make([]byte, dataLen)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}

func (t *thirdServer) UploadLogs(ctx context.Context, req *third.UploadLogsReq) (*third.UploadLogsResp, error) {
	var dbLogs []*relationtb.Log
	userID := ctx.Value(constant.OpUserID).(string)
	platform := constant.PlatformID2Name[int(req.Platform)]
	for _, fileURL := range req.FileURLs {
		log := relationtb.Log{
			Platform:     platform,
			UserID:       userID,
			CreateTime:   time.Now(),
			Url:          fileURL.URL,
			FileName:     fileURL.Filename,
			AppFramework: req.AppFramework,
			Version:      req.Version,
			Ex:           req.Ex,
		}
		for i := 0; i < 20; i++ {
			id := genLogID()
			logs, err := t.thirdDatabase.GetLogs(ctx, []string{id}, "")
			if err != nil {
				return nil, err
			}
			if len(logs) == 0 {
				log.LogID = id
				break
			}
		}
		if log.LogID == "" {
			return nil, servererrs.ErrData.WrapMsg("Log id gen error")
		}
		dbLogs = append(dbLogs, &log)
	}
	err := t.thirdDatabase.UploadLogs(ctx, dbLogs)
	if err != nil {
		return nil, err
	}
	return &third.UploadLogsResp{}, nil
}

func (t *thirdServer) DeleteLogs(ctx context.Context, req *third.DeleteLogsReq) (*third.DeleteLogsResp, error) {
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	userID := ""
	logs, err := t.thirdDatabase.GetLogs(ctx, req.LogIDs, userID)
	if err != nil {
		return nil, err
	}
	var logIDs []string
	for _, log := range logs {
		logIDs = append(logIDs, log.LogID)
	}
	if ids := datautil.Single(req.LogIDs, logIDs); len(ids) > 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("logIDs not found", "logIDs", ids)
	}
	err = t.thirdDatabase.DeleteLogs(ctx, req.LogIDs, userID)
	if err != nil {
		return nil, err
	}

	return &third.DeleteLogsResp{}, nil
}

func dbToPbLogInfos(logs []*relationtb.Log) []*third.LogInfo {
	db2pbForLogInfo := func(log *relationtb.Log) *third.LogInfo {
		return &third.LogInfo{
			Filename:   log.FileName,
			UserID:     log.UserID,
			Platform:   log.Platform,
			Url:        log.Url,
			CreateTime: log.CreateTime.UnixMilli(),
			LogID:      log.LogID,
			SystemType: log.SystemType,
			Version:    log.Version,
			Ex:         log.Ex,
		}
	}
	return datautil.Slice(logs, db2pbForLogInfo)
}

func (t *thirdServer) SearchLogs(ctx context.Context, req *third.SearchLogsReq) (*third.SearchLogsResp, error) {
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	var (
		resp    third.SearchLogsResp
		userIDs []string
	)
	if req.StartTime > req.EndTime {
		return nil, errs.ErrArgs.WrapMsg("startTime>endTime")
	}
	if req.StartTime == 0 && req.EndTime == 0 {
		t := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
		timestampMills := t.UnixNano() / int64(time.Millisecond)
		req.StartTime = timestampMills
		req.EndTime = time.Now().UnixNano() / int64(time.Millisecond)
	}

	total, logs, err := t.thirdDatabase.SearchLogs(ctx, req.Keyword, time.UnixMilli(req.StartTime), time.UnixMilli(req.EndTime), req.Pagination)
	if err != nil {
		return nil, err
	}
	pbLogs := dbToPbLogInfos(logs)
	for _, log := range logs {
		userIDs = append(userIDs, log.UserID)
	}
	userMap, err := t.userClient.GetUsersInfoMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	for _, pbLog := range pbLogs {
		if user, ok := userMap[pbLog.UserID]; ok {
			pbLog.Nickname = user.Nickname
		}
	}
	resp.LogsInfos = pbLogs
	resp.Total = uint32(total)
	return &resp, nil
}
