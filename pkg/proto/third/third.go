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

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func (x *ApplyPutReq) Check() error {
	if x.PutID == "" {
		return errs.ErrArgs.Wrap("PutID is empty")
	}
	if x.ContentType == "" {
		return errs.ErrArgs.Wrap("ContentType is empty")
	}
	return nil
}

func (x *ConfirmPutReq) Check() error {
	if x.PutID == "" {
		return errs.ErrArgs.Wrap("PutID is empty")
	}
	return nil
}

func (x *GetUrlReq) Check() error {
	if x.Name == "" {
		return errs.ErrArgs.Wrap("Name is empty")
	}
	return nil
}

func (x *GetPutReq) Check() error {
	if x.PutID == "" {
		return errs.ErrArgs.Wrap("PutID is empty")
	}
	return nil
}

func (x *GetHashInfoReq) Check() error {
	if x.Hash == "" {
		return errs.ErrArgs.Wrap("Hash is empty")
	}
	return nil
}

func (x *FcmUpdateTokenReq) Check() error {
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errs.ErrArgs.Wrap("PlatformID is invalid")
	}
	if x.FcmToken == "" {
		return errs.ErrArgs.Wrap("FcmToken is empty")
	}
	if x.Account == "" {
		return errs.ErrArgs.Wrap("Account is empty")
	}
	return nil
}

func (x *SetAppBadgeReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("UserID is empty")
	}
	return nil
}
