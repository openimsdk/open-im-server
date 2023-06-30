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

package errs

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type CodeError interface {
	Code() int
	Msg() string
	Detail() string
	WithDetail(detail string) CodeError
	// Is 判断是否是某个错误, loose为false时, 只有错误码相同就认为是同一个错误, 默认为true
	Is(err error, loose ...bool) bool
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

func (e *codeError) Detail() string {
	return e.detail
}

func (e *codeError) WithDetail(detail string) CodeError {
	var d string
	if e.detail == "" {
		d = detail
	} else {
		d = e.detail + ", " + detail
	}
	return &codeError{
		code:   e.code,
		msg:    e.msg,
		detail: d,
	}
}

func (e *codeError) Wrap(w ...string) error {
	return errors.Wrap(e, strings.Join(w, ", "))
}

func (e *codeError) Is(err error, loose ...bool) bool {
	if err == nil {
		return false
	}
	var allowSubclasses bool
	if len(loose) == 0 {
		allowSubclasses = true
	} else {
		allowSubclasses = loose[0]
	}
	codeErr, ok := Unwrap(err).(CodeError)
	if ok {
		if allowSubclasses {
			return Relation.Is(e.code, codeErr.Code())
		} else {
			return codeErr.Code() == e.code
		}
	}
	return false
}

func (e *codeError) Error() string {
	return fmt.Sprintf("%s", e.msg)
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

func Wrap(err error, msg ...string) error {
	if err == nil {
		return nil
	}
	if len(msg) == 0 {
		return errors.WithStack(err)
	}
	return errors.Wrap(err, strings.Join(msg, ", "))
}
