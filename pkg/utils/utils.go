package utils

import (
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"runtime"
	"strconv"
)

// copy a by b  b->a
func CopyStructFields(a interface{}, b interface{}, fields ...string) (err error) {
	return copier.Copy(a, b)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}
