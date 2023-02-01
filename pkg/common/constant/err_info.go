package constant

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

type ErrInfo struct {
	ErrCode      int32
	ErrMsg       string
	DetailErrMsg string
}

func NewErrInfo(code int32, msg, detail string) *ErrInfo {
	return &ErrInfo{
		ErrCode:      code,
		ErrMsg:       msg,
		DetailErrMsg: detail,
	}
}

func (e *ErrInfo) Error() string {
	return "errMsg: " + e.ErrMsg + " detail errMsg: " + e.DetailErrMsg
}

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}

func (e *ErrInfo) Wrap(msg ...string) error {
	return errors.Wrap(e, strings.Join(msg, "--"))
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
