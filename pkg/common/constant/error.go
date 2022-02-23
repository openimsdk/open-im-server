package constant

import "errors"

// key = errCode, string = errMsg
type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

var (
	OK = ErrInfo{0, ""}
	ErrServer = ErrInfo{500, "server error"}

	//	ErrMysql             = ErrInfo{100, ""}
	//	ErrMongo             = ErrInfo{110, ""}
	//	ErrRedis             = ErrInfo{120, ""}
	ErrParseToken = ErrInfo{700, ParseTokenMsg.Error()}
	//	ErrCreateToken       = ErrInfo{201, "Create token failed"}
	//	ErrAppServerKey      = ErrInfo{300, "key error"}
	ErrTencentCredential = ErrInfo{400, ThirdPartyMsg.Error()}

	//	ErrorUserRegister             = ErrInfo{600, "User registration failed"}
	//	ErrAccountExists              = ErrInfo{601, "The account is already registered and cannot be registered again"}
	//	ErrUserPassword               = ErrInfo{602, "User password error"}
	//	ErrRefreshToken               = ErrInfo{605, "Failed to refresh token"}
	//	ErrAddFriend                  = ErrInfo{606, "Failed to add friends"}
	//	ErrAgreeToAddFriend           = ErrInfo{607, "Failed to agree application"}
	//	ErrAddFriendToBlack           = ErrInfo{608, "Failed to add friends to the blacklist"}
	//	ErrGetBlackList               = ErrInfo{609, "Failed to get blacklist"}
	//	ErrDeleteFriend               = ErrInfo{610, "Failed to delete friend"}
	//	ErrGetFriendApplyList         = ErrInfo{611, "Failed to get friend application list"}
	//	ErrGetFriendList              = ErrInfo{612, "Failed to get friend list"}
	//	ErrRemoveBlackList            = ErrInfo{613, "Failed to remove blacklist"}
	//	ErrSearchUserInfo             = ErrInfo{614, "Can't find the user information"}
	//	ErrDelAppleDeviceToken        = ErrInfo{615, ""}
	//	ErrModifyUserInfo             = ErrInfo{616, "update user some attribute failed"}
	//	ErrSetFriendComment           = ErrInfo{617, "set friend comment failed"}
	//	ErrSearchUserInfoFromTheGroup = ErrInfo{618, "There is no such group or the user not in the group"}
	//	ErrCreateGroup                = ErrInfo{619, "create group chat failed"}
	//	ErrJoinGroupApplication       = ErrInfo{620, "Failed to apply to join the group"}
	//	ErrQuitGroup                  = ErrInfo{621, "Failed to quit the group"}
	//	ErrSetGroupInfo               = ErrInfo{622, "Failed to set group info"}
	//	ErrParam                      = ErrInfo{700, "param failed"}
	ErrTokenExpired     = ErrInfo{701, TokenExpiredMsg.Error()}
	ErrTokenInvalid     = ErrInfo{702, TokenInvalidMsg.Error()}
	ErrTokenMalformed   = ErrInfo{703, TokenMalformedMsg.Error()}
	ErrTokenNotValidYet = ErrInfo{704, TokenNotValidYetMsg.Error()}
	ErrTokenUnknown     = ErrInfo{705, TokenUnknownMsg.Error()}

	ErrAccess = ErrInfo{ErrCode: 801, ErrMsg: AccessMsg.Error()}
	ErrDB     = ErrInfo{ErrCode: 802, ErrMsg: DBMsg.Error()}
	ErrArgs   = ErrInfo{ErrCode: 8003, ErrMsg: ArgsMsg.Error()}
)

var (
	ParseTokenMsg       = errors.New("parse token failed")
	TokenExpiredMsg     = errors.New("token is timed out, please log in again")
	TokenInvalidMsg     = errors.New("token has been invalidated")
	TokenNotValidYetMsg = errors.New("token not active yet")
	TokenMalformedMsg   = errors.New("that's not even a token")
	TokenUnknownMsg     = errors.New("couldn't handle this token")

	AccessMsg = errors.New("no permission")
	DBMsg     = errors.New("db failed")
	ArgsMsg   = errors.New("args failed")

	ThirdPartyMsg = errors.New("third party error")
)

const (
	NoError              = 0
	FormattingError      = 10001
	HasRegistered        = 10002
	NotRegistered        = 10003
	PasswordErr          = 10004
	GetIMTokenErr        = 10005
	RepeatSendCode       = 10006
	MailSendCodeErr      = 10007
	SmsSendCodeErr       = 10008
	CodeInvalidOrExpired = 10009
	RegisterFailed       = 10010
	ResetPasswordFailed  = 10011
	DatabaseError        = 10002
	ServerError          = 10004
	HttpError            = 10005
	IoError              = 10006
	IntentionalError     = 10007
)

func (e ErrInfo) Error() string {
	return e.ErrMsg
}

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}
