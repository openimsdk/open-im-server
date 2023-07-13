package third

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
)

func toPbMapArray(m map[string][]string) []*third.KeyValues {
	res := make([]*third.KeyValues, 0, len(m))
	for key := range m {
		res = append(res, &third.KeyValues{
			Key:    key,
			Values: m[key],
		})
	}
	return res
}

func checkUploadName(ctx context.Context, name string) error {
	if name == "" {
		return errs.ErrArgs.Wrap("name is empty")
	}
	if name[0] == '/' {
		return errs.ErrArgs.Wrap("name cannot start with `/`")
	}
	if err := checkValidObjectName(name); err != nil {
		return errs.ErrArgs.Wrap(err.Error())
	}
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		return errs.ErrNoPermission.Wrap("opUserID is empty")
	}
	if !tokenverify.IsManagerUserID(opUserID) {
		if !strings.HasPrefix(name, opUserID+"/") {
			return errs.ErrNoPermission.Wrap(fmt.Sprintf("name must start with `%s/`", opUserID))
		}
	}
	return nil
}

func checkValidObjectNamePrefix(objectName string) error {
	if len(objectName) > 1024 {
		return errors.New("object name cannot be longer than 1024 characters")
	}
	if !utf8.ValidString(objectName) {
		return errors.New("object name with non UTF-8 strings are not supported")
	}
	return nil
}

func checkValidObjectName(objectName string) error {
	if strings.TrimSpace(objectName) == "" {
		return errors.New("object name cannot be empty")
	}
	return checkValidObjectNamePrefix(objectName)
}
