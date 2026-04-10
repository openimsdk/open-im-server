// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

type PhoneSNApi struct {
	db database.PhoneSN
}

func NewPhoneSNApi(db database.PhoneSN) *PhoneSNApi {
	return &PhoneSNApi{db: db}
}

type phoneGetSNInfoReq struct {
	Phone string `json:"phone" binding:"required"`
}

type phoneGetSNInfoResp struct {
	IsSnd  bool  `json:"is_snd"`
	UserID int64 `json:"userID"`
}

type phoneSetSNInfoReq struct {
	Phone  string `json:"phone" binding:"required"`
	UserID int64  `json:"userID"`
	IsSnd  bool   `json:"is_snd"`
}

// GetSNInfo POST /phone/get_sn_info
func (a *PhoneSNApi) GetSNInfo(c *gin.Context) {
	var req phoneGetSNInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		log.ZError(c, "GetSNInfo", err)
		return
	}
	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("phone is empty"))
		log.ZError(c, "GetSNInfo", errs.ErrArgs.WrapMsg("phone is empty"))
		return
	}
	info, err := a.db.GetByPhone(c, phone)
	if err != nil {
		apiresp.GinError(c, err)
		log.ZError(c, "GetSNInfo", err)
		return
	}
	resp := phoneGetSNInfoResp{IsSnd: false, UserID: 0}
	if info != nil {
		resp.IsSnd = info.IsSnd
		resp.UserID = info.UserID
	}
	apiresp.GinSuccess(c, resp)
}

// SetSNInfo POST /phone/set_sn_info
func (a *PhoneSNApi) SetSNInfo(c *gin.Context) {
	var req phoneSetSNInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("phone is empty"))
		return
	}
	if err := a.db.Upsert(c, phone, req.UserID, req.IsSnd); err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, nil)
}
