package errs

import (
	"OpenIM/pkg/utils"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Code interface {
	Code() int
	Msg() string
}

type Coderr interface {
	Code
	Detail() string
	WithDetail(string) Coderr
	Warp(...string) error
	error
}

func NewCodeError(code int, msg string) Coderr {
	return &errInfo{
		code: code,
		msg:  msg,
	}
}

type errInfo struct {
	code   int
	msg    string
	detail string
}

func (e *errInfo) WithDetail(s string) Coderr {
	if e.detail == "" {
		e.detail = s
	} else {
		e.detail = s + ", " + e.detail
	}
	return e
}

func (e *errInfo) Code() int {
	return e.code
}

func (e *errInfo) Msg() string {
	return e.msg
}

func (e *errInfo) Detail() string {
	return e.detail
}

func (e *errInfo) Warp(w ...string) error {
	return errors.Wrap(e, strings.Join(w, ", "))
}

func (e *errInfo) Error() string {
	return fmt.Sprintf("[%d]%s", e.code, e.msg)
}

func Unwrap(err error) error {
	return utils.Unwrap(err)
}

func GetCode(err error) Code {
	if err == nil {
		return NewCodeError(UnknownCode, "nil")
	}
	if code, ok := Unwrap(err).(Code); ok {
		if code.Code() == 0 {
			return NewCodeError(UnknownCode, "code == 0")
		}
		return code
	}
	return NewCodeError(UnknownCode, "unknown code")
}
