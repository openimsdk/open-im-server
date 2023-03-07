package errs

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type CodeError interface {
	Code() int
	Msg() string
	Wrap(msg ...string) error
	error
}

func NewCodeError(code int, msg string) CodeError {
	return &codeError{
		code: code,
		msg:  msg,
	}
}

type codeError struct {
	code   int
	msg    string
	detail string
}

func (e *codeError) Code() int {
	return e.code
}

func (e *codeError) Msg() string {
	return e.msg
}

func (e *codeError) Wrap(w ...string) error {
	return errors.Wrap(e, strings.Join(w, ", "))
}

func (e *codeError) Error() string {
	return fmt.Sprintf("[%d]%s", e.code, e.msg)
}

func Unwrap(err error) error {
	for err != nil {
		unwrap, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			break
		}
		err = unwrap.Unwrap()
	}
	return err
}
