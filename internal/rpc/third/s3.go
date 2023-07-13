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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3/cont"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"time"
)

func (t *thirdServer) PartLimit(ctx context.Context, req *third.PartLimitReq) (*third.PartLimitResp, error) {
	limit := t.s3dataBase.PartLimit()
	return &third.PartLimitResp{
		MinPartSize: limit.MinPartSize,
		MaxPartSize: limit.MaxPartSize,
		MaxNumSize:  int32(limit.MaxNumSize),
	}, nil
}

func (t *thirdServer) PartSize(ctx context.Context, req *third.PartSizeReq) (*third.PartSizeResp, error) {
	size, err := t.s3dataBase.PartSize(ctx, req.Size)
	if err != nil {
		return nil, err
	}
	return &third.PartSizeResp{Size: size}, nil
}

func (t *thirdServer) InitiateMultipartUpload(ctx context.Context, req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	defer log.ZDebug(ctx, "return")
	if err := checkUploadName(ctx, req.Name); err != nil {
		return nil, err
	}
	expireTime := time.Now().Add(t.defaultExpire)
	result, err := t.s3dataBase.InitiateMultipartUpload(ctx, req.Hash, req.Size, t.defaultExpire, int(req.MaxParts))
	if err != nil {
		if haErr, ok := errs.Unwrap(err).(*cont.HashAlreadyExistsError); ok {
			obj := &relation.ObjectModel{
				Name:        req.Name,
				UserID:      mcontext.GetOpUserID(ctx),
				Hash:        req.Hash,
				Key:         haErr.Object.Key,
				Size:        haErr.Object.Size,
				ContentType: req.ContentType,
				Cause:       req.Cause,
				CreateTime:  time.Now(),
			}
			if err := t.s3dataBase.SetObject(ctx, obj); err != nil {
				return nil, err
			}
			return &third.InitiateMultipartUploadResp{
				Url: t.apiAddress(obj.Name),
			}, nil
		}
		return nil, err
	}
	var sign *third.AuthSignParts
	if result.Sign != nil && len(result.Sign.Parts) > 0 {
		sign = &third.AuthSignParts{
			Url:    result.Sign.URL,
			Query:  toPbMapArray(result.Sign.Query),
			Header: toPbMapArray(result.Sign.Header),
			Parts:  make([]*third.SignPart, len(result.Sign.Parts)),
		}
		for i, part := range result.Sign.Parts {
			sign.Parts[i] = &third.SignPart{
				PartNumber: int32(part.PartNumber),
				Url:        part.URL,
				Query:      toPbMapArray(part.Query),
				Header:     toPbMapArray(part.Header),
			}
		}
	}
	return &third.InitiateMultipartUploadResp{
		Upload: &third.UploadInfo{
			UploadID:   result.UploadID,
			PartSize:   result.PartSize,
			Sign:       sign,
			ExpireTime: expireTime.UnixMilli(),
		},
	}, nil
}

func (t *thirdServer) AuthSign(ctx context.Context, req *third.AuthSignReq) (*third.AuthSignResp, error) {
	defer log.ZDebug(ctx, "return")
	partNumbers := utils.Slice(req.PartNumbers, func(partNumber int32) int { return int(partNumber) })
	result, err := t.s3dataBase.AuthSign(ctx, req.UploadID, partNumbers)
	if err != nil {
		return nil, err
	}
	resp := &third.AuthSignResp{
		Url:    result.URL,
		Query:  toPbMapArray(result.Query),
		Header: toPbMapArray(result.Header),
		Parts:  make([]*third.SignPart, len(result.Parts)),
	}
	for i, part := range result.Parts {
		resp.Parts[i] = &third.SignPart{
			PartNumber: int32(part.PartNumber),
			Url:        part.URL,
			Query:      toPbMapArray(part.Query),
			Header:     toPbMapArray(part.Header),
		}
	}
	return resp, nil
}

func (t *thirdServer) CompleteMultipartUpload(ctx context.Context, req *third.CompleteMultipartUploadReq) (*third.CompleteMultipartUploadResp, error) {
	defer log.ZDebug(ctx, "return")
	if err := checkUploadName(ctx, req.Name); err != nil {
		return nil, err
	}
	result, err := t.s3dataBase.CompleteMultipartUpload(ctx, req.UploadID, req.Parts)
	if err != nil {
		return nil, err
	}
	obj := &relation.ObjectModel{
		Name:        req.Name,
		UserID:      mcontext.GetOpUserID(ctx),
		Hash:        result.Hash,
		Key:         result.Key,
		Size:        result.Size,
		ContentType: req.ContentType,
		Cause:       req.Cause,
		CreateTime:  time.Now(),
	}
	if err := t.s3dataBase.SetObject(ctx, obj); err != nil {
		return nil, err
	}
	return &third.CompleteMultipartUploadResp{
		Url: t.apiAddress(obj.Name),
	}, nil
}

func (t *thirdServer) AccessURL(ctx context.Context, req *third.AccessURLReq) (*third.AccessURLResp, error) {
	expireTime, rawURL, err := t.s3dataBase.AccessURL(ctx, req.Name, t.defaultExpire)
	if err != nil {
		return nil, err
	}
	return &third.AccessURLResp{
		Url:        rawURL,
		ExpireTime: expireTime.UnixMilli(),
	}, nil
}

func (t *thirdServer) apiAddress(name string) string {
	return t.apiURL + name
}
