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
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

func toPbMapArray(m map[string][]string) []*third.KeyValues {
	if len(m) == 0 {
		return nil
	}
	res := make([]*third.KeyValues, 0, len(m))
	for key := range m {
		res = append(res, &third.KeyValues{
			Key:    key,
			Values: m[key],
		})
	}
	return res
}

func (t *thirdServer) checkUploadName(ctx context.Context, name string) error {
	if name == "" {
		return errs.ErrArgs.WrapMsg("name is empty")
	}
	if name[0] == '/' {
		return errs.ErrArgs.WrapMsg("name cannot start with `/`")
	}
	if err := checkValidObjectName(name); err != nil {
		return errs.ErrArgs.WrapMsg(err.Error())
	}
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		return errs.ErrNoPermission.WrapMsg("opUserID is empty")
	}
	if !authverify.IsManagerUserID(opUserID, t.config.Share.IMAdminUserID) {
		if !strings.HasPrefix(name, opUserID+"/") {
			return errs.ErrNoPermission.WrapMsg(fmt.Sprintf("name must start with `%s/`", opUserID))
		}
	}
	return nil
}

func checkValidObjectNamePrefix(objectName string) error {
	if len(objectName) > 1024 {
		return errs.New("object name cannot be longer than 1024 characters")
	}
	if !utf8.ValidString(objectName) {
		return errs.New("object name with non UTF-8 strings are not supported")
	}
	return nil
}

func checkValidObjectName(objectName string) error {
	if strings.TrimSpace(objectName) == "" {
		return errs.New("object name cannot be empty")
	}
	return checkValidObjectNamePrefix(objectName)
}

func (t *thirdServer) IsManagerUserID(opUserID string) bool {
	return authverify.IsManagerUserID(opUserID, t.config.Share.IMAdminUserID)
}

func putUpdate[T any](update map[string]any, name string, val interface{ GetValuePtr() *T }) {
	ptrVal := val.GetValuePtr()
	if ptrVal == nil {
		return
	}
	update[name] = *ptrVal
}
