// Copyright © 2024 OpenIM. All rights reserved.
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

package msg

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

func genReportID() string {
	const dataLen = 12
	data := make([]byte, dataLen)
	rand.Read(data)
	chars := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < len(data); i++ {
		data[i] = chars[data[i]%byte(len(chars))]
	}
	return string(data)
}

func (m *msgServer) ReportSpam(ctx context.Context, req *msg.ReportSpamReq) (*msg.ReportSpamResp, error) {
	if req.ReportedUserID == "" {
		return nil, errs.ErrArgs.WrapMsg("reportedUserID is required")
	}
	if req.ReasonType <= 0 {
		return nil, errs.ErrArgs.WrapMsg("reasonType must be positive")
	}

	reporterUserID := mcontext.GetOpUserID(ctx)

	report := &model.SpamReport{
		ReporterUserID: reporterUserID,
		ReportedUserID: req.ReportedUserID,
		ConversationID: req.ConversationID,
		ClientMsgID:    req.ClientMsgID,
		Seq:            req.Seq,
		ReasonType:     req.ReasonType,
		Reason:         req.Reason,
		Status:         model.SpamReportStatusPending,
		CreateTime:     time.Now(),
		Ex:             req.Ex,
	}

	// Generate a unique reportID.
	for i := 0; i < 20; i++ {
		id := genReportID()
		existing, err := m.spamReportDB.Get(ctx, id)
		if err == nil && existing != nil {
			continue
		}
		report.ReportID = id
		break
	}
	if report.ReportID == "" {
		return nil, errs.ErrInternalServer.WrapMsg("failed to generate report ID")
	}

	if err := m.spamReportDB.Create(ctx, report); err != nil {
		return nil, err
	}
	return &msg.ReportSpamResp{ReportID: report.ReportID}, nil
}

func (m *msgServer) GetSpamReports(ctx context.Context, req *msg.GetSpamReportsReq) (*msg.GetSpamReportsResp, error) {
	if err := authverify.CheckAdmin(ctx, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}

	var start, end time.Time
	if req.StartTime > 0 {
		start = time.UnixMilli(req.StartTime)
	}
	if req.EndTime > 0 {
		end = time.UnixMilli(req.EndTime)
	}

	total, reports, err := m.spamReportDB.Find(ctx, req.Status, req.ReportedUserID, req.ReporterUserID,
		start, end, req.Pagination)
	if err != nil {
		return nil, err
	}

	pbReports := datautil.Slice(reports, func(r *model.SpamReport) *msg.SpamReportInfo {
		return &msg.SpamReportInfo{
			ReportID:       r.ReportID,
			ReporterUserID: r.ReporterUserID,
			ReportedUserID: r.ReportedUserID,
			ConversationID: r.ConversationID,
			ClientMsgID:    r.ClientMsgID,
			Seq:            r.Seq,
			ReasonType:     r.ReasonType,
			Reason:         r.Reason,
			Status:         r.Status,
			CreateTime:     r.CreateTime.UnixMilli(),
			HandleTime:     r.HandleTime.UnixMilli(),
			HandlerUserID:  r.HandlerUserID,
			Ex:             r.Ex,
		}
	})

	return &msg.GetSpamReportsResp{
		Reports: pbReports,
		Total:   uint32(total),
	}, nil
}

func (m *msgServer) HandleSpamReport(ctx context.Context, req *msg.HandleSpamReportReq) (*msg.HandleSpamReportResp, error) {
	if err := authverify.CheckAdmin(ctx, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	if req.ReportID == "" {
		return nil, errs.ErrArgs.WrapMsg("reportID is required")
	}
	if req.Status != model.SpamReportStatusHandled && req.Status != model.SpamReportStatusIgnored {
		return nil, errs.ErrArgs.WrapMsg("status must be 1 (handled) or 2 (ignored)")
	}

	handlerUserID := mcontext.GetOpUserID(ctx)
	if err := m.spamReportDB.UpdateStatus(ctx, req.ReportID, req.Status, handlerUserID, time.Now()); err != nil {
		return nil, err
	}
	return &msg.HandleSpamReportResp{}, nil
}
