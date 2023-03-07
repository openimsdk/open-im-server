package errs

import (
	"OpenIM/pkg/utils"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Coderr interface {
	Code() int
	Msg() string
	Warp(msg ...string) error
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

func (e *errInfo) Code() int {
	return e.code
}

func (e *errInfo) Msg() string {
	return e.msg
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
