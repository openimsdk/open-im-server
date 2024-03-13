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
	"fmt"
	"time"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/third"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"
	utils2 "github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
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
	var DBlogs []*relationtb.LogModel
	userID := ctx.Value(constant.OpUserID).(string)
	platform := constant.PlatformID2Name[int(req.Platform)]
	for _, fileURL := range req.FileURLs {
		log := relationtb.LogModel{
			Version:    req.Version,
			SystemType: req.SystemType,
			Platform:   platform,
			UserID:     userID,
			CreateTime: time.Now(),
			Url:        fileURL.URL,
			FileName:   fileURL.Filename,
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
			return nil, errs.ErrData.Wrap("LogModel id gen error")
		}
		DBlogs = append(DBlogs, &log)
	}
	err := t.thirdDatabase.UploadLogs(ctx, DBlogs)
	if err != nil {
		return nil, err
	}
	return &third.UploadLogsResp{}, nil
}

func (t *thirdServer) DeleteLogs(ctx context.Context, req *third.DeleteLogsReq) (*third.DeleteLogsResp, error) {
	if err := authverify.CheckAdmin(ctx, t.config); err != nil {
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
	if ids := utils2.Single(req.LogIDs, logIDs); len(ids) > 0 {
		return nil, errs.ErrRecordNotFound.Wrap(fmt.Sprintf("logIDs not found%#v", ids))
	}
	err = t.thirdDatabase.DeleteLogs(ctx, req.LogIDs, userID)
	if err != nil {
		return nil, err
	}

	return &third.DeleteLogsResp{}, nil
}

func dbToPbLogInfos(logs []*relationtb.LogModel) []*third.LogInfo {
	db2pbForLogInfo := func(log *relationtb.LogModel) *third.LogInfo {
		return &third.LogInfo{
			Filename:   log.FileName,
			UserID:     log.UserID,
			Platform:   utils.StringToInt32(log.Platform),
			Url:        log.Url,
			CreateTime: log.CreateTime.UnixMilli(),
			LogID:      log.LogID,
			SystemType: log.SystemType,
			Version:    log.Version,
			Ex:         log.Ex,
		}
	}
	return utils.Slice(logs, db2pbForLogInfo)
}

func (t *thirdServer) SearchLogs(ctx context.Context, req *third.SearchLogsReq) (*third.SearchLogsResp, error) {
	if err := authverify.CheckAdmin(ctx, t.config); err != nil {
		return nil, err
	}
	var (
		resp    third.SearchLogsResp
		userIDs []string
	)
	if req.StartTime > req.EndTime {
		return nil, errs.ErrArgs.Wrap("startTime>endTime")
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
	userMap, err := t.userRpcClient.GetUsersInfoMap(ctx, userIDs)
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
