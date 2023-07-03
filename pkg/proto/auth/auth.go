package auth

import "github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

func (x *UserTokenReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.PlatformID > 9 || x.PlatformID < 1 {
		return errs.ErrArgs.Wrap("platform is invalidate")
	}
	return nil
}

func (x *ForceLogoutReq) Check() error {
	if x.UserID == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	if x.PlatformID > 9 || x.PlatformID < 1 {
		return errs.ErrArgs.Wrap("platformID is invalidate")
	}
	return nil
}

func (x *ParseTokenReq) Check() error {
	if x.Token == "" {
		return errs.ErrArgs.Wrap("userID is empty")
	}
	return nil
}
