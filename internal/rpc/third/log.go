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
	var DBlogs []*relationtb.Log
	userID := ctx.Value(constant.OpUserID).(string)
	platform := constant.PlatformID2Name[int(req.Platform)]
	for _, fileURL := range req.FileURLs {
		log := relationtb.Log{
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
			return nil, errs.ErrData.Wrap("Log id gen error")
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
	if ids := utils2.Single(req.LogIDs, logIDs); len(ids) > 0 {
		return nil, errs.ErrRecordNotFound.Wrap(fmt.Sprintf("logIDs not found%#v", ids))
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
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	var (
		resp    third.SearchLogsResp
		userIDs []string
	)
	if req.StartTime > req.EndTime {
		return nil, errs.ErrArgs.Wrap("startTime>endTime")
	}
	total, logs, err := t.thirdDatabase.SearchLogs(ctx, req.Keyword, time.UnixMilli(req.StartTime), time.UnixMilli(req.EndTime), req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	pbLogs := dbToPbLogInfos(logs)
	for _, log := range logs {
		userIDs = append(userIDs, log.UserID)
	}
	users, err := t.thirdDatabase.FindUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	IDtoName := make(map[string]string)
	for _, user := range users {
		IDtoName[user.UserID] = user.Nickname
	}
	for _, pbLog := range pbLogs {
		pbLog.Nickname = IDtoName[pbLog.UserID]
	}
	resp.LogsInfos = pbLogs
	resp.Total = total
	return &resp, nil
}
