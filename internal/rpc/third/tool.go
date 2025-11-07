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
	if !authverify.CheckUserIsAdmin(ctx, opUserID) {
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

func putUpdate[T any](update map[string]any, name string, val interface{ GetValuePtr() *T }) {
	ptrVal := val.GetValuePtr()
	if ptrVal == nil {
		return
	}
	update[name] = *ptrVal
}
