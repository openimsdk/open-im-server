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

package specialerror

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

var handlers []func(err error) errs.CodeError

func AddErrHandler(h func(err error) errs.CodeError) {
	if h == nil {
		panic("nil handler")
	}
	handlers = append(handlers, h)
}

func AddReplace(target error, codeErr errs.CodeError) {
	AddErrHandler(func(err error) errs.CodeError {
		if err == target {
			return codeErr
		}
		return nil
	})
}

func ErrCode(err error) errs.CodeError {
	if codeErr, ok := err.(errs.CodeError); ok {
		return codeErr
	}
	for i := 0; i < len(handlers); i++ {
		if codeErr := handlers[i](err); codeErr != nil {
			return codeErr
		}
	}
	return nil
}
