package constant

import (
	sdkws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ErrInfo struct {
	ErrCode      int32
	ErrMsg       string
	DetailErrMsg string
}

func (e *ErrInfo) Error() string {
	return "errMsg: " + e.ErrMsg + " detail errMsg: " + e.DetailErrMsg
}

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}

func (e *ErrInfo) Warp() error {
	return utils.Wrap(e, "")
}

func (e *ErrInfo) WarpMessage(msg string) error {
	return utils.Wrap(e, msg)
}

func NewErrNetwork(err error) error {
	return toDetail(err, ErrNetwork)
}

func NewErrData(err error) error {
	return toDetail(err, ErrData)
}

func toDetail(err error, info *ErrInfo) *ErrInfo {
	errInfo := *info
	errInfo.DetailErrMsg = err.Error()
	return &errInfo
}

func ToAPIErrWithErr(err error) *ErrInfo {
	errComm := errors.New("")
	var marshalErr *json.MarshalerError
	errInfo := &ErrInfo{}
	switch {
	case errors.As(err, &errComm):
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return toDetail(err, ErrRecordNotFound)
		}
		return toDetail(err, ErrData)
	case errors.As(err, &marshalErr):
		return toDetail(err, ErrData)
	case errors.As(err, &errInfo):
		return toDetail(err, errInfo)
	}
	return toDetail(err, ErrDefaultOther)
}

func SetErrorForResp(err error, commonResp *sdkws.CommonResp) {
	errInfo := ToAPIErrWithErr(err)
	commonResp.ErrCode = errInfo.ErrCode
	commonResp.ErrMsg = errInfo.ErrMsg
	commonResp.DetailErrMsg = err.Error()
}

func CommonResp2Err(resp *sdkws.CommonResp) error {
	if resp.ErrCode != NoError {
		return errors.New(fmt.Sprintf("call rpc error, errCode is %d, errMsg is %s, detailErrMsg is %s", resp.ErrCode, resp.ErrMsg, resp.DetailErrMsg))
	}
	return nil
}

func Error2CommResp(ctx context.Context, info ErrInfo, detailErrMsg string) *sdkws.CommonResp {
	err := &sdkws.CommonResp{
		ErrCode: info.ErrCode,
		ErrMsg:  info.ErrMsg,
	}
	if detailErrMsg != "" {
		err.DetailErrMsg = detailErrMsg
	}
	return err
}
