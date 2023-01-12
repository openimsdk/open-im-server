package constant

var (
	ErrNone             = ErrInfo{0, "", ""}
	ErrArgs             = ErrInfo{ArgsError, "ArgsError", ""}
	ErrDatabase         = ErrInfo{DatabaseError, "DatabaseError", ""}
	ErrInternalServer   = ErrInfo{ServerInternalError, "ServerInternalError", ""}
	ErrNetwork          = ErrInfo{NetworkError, "NetworkError", ""}
	ErrNoPermission     = ErrInfo{NoPermissionError, "NoPermissionError", ""}
	ErrIdentity         = ErrInfo{IdentityError, "IdentityError", ""}
	ErrCallback         = ErrInfo{ErrMsg: "CallbackError"}
	ErrCallbackContinue = ErrInfo{ErrMsg: "CallbackContinueError"}

	ErrUserIDNotFound  = ErrInfo{UserIDNotFoundError, "UserIDNotFoundError", ""}
	ErrGroupIDNotFound = ErrInfo{GroupIDNotFoundError, "GroupIDNotFoundError", ""}

	ErrRecordNotFound = ErrInfo{RecordNotFoundError, "RecordNotFoundError", ""}

	ErrRelationshipAlready = ErrInfo{RelationshipAlreadyError, "RelationshipAlreadyError", ""}
	ErrNotRelationshipYet  = ErrInfo{NotRelationshipYetError, "NotRelationshipYetError", ""}

	ErrOnlyOneOwner        = ErrInfo{OnlyOneOwnerError, "OnlyOneOwnerError", ""}
	ErrInGroupAlready      = ErrInfo{InGroupAlreadyError, "InGroupAlreadyError", ""}
	ErrNotInGroupYet       = ErrInfo{NotInGroupYetError, "NotInGroupYetError", ""}
	ErrDismissedAlready    = ErrInfo{DismissedAlreadyError, "DismissedAlreadyError", ""}
	ErrOwnerNotAllowedQuit = ErrInfo{OwnerNotAllowedQuitError, "OwnerNotAllowedQuitError", ""}
	ErrRegisteredAlready   = ErrInfo{RegisteredAlreadyError, "RegisteredAlreadyError", ""}
	ErrGroupTypeNotSupport = ErrInfo{GroupTypeNotSupport, "", ""}
	ErrGroupNoOwner        = ErrInfo{GroupNoOwner, "ErrGroupNoOwner", ""}

	ErrDefaultOther             = ErrInfo{DefaultOtherError, "DefaultOtherError", ""}
	ErrData                     = ErrInfo{DataError, "DataError", ""}
	ErrTokenExpired             = ErrInfo{TokenExpiredError, "TokenExpiredError", ""}
	ErrTokenInvalid             = ErrInfo{TokenInvalidError, "TokenInvalidError", ""}         //
	ErrTokenMalformed           = ErrInfo{TokenMalformedError, "TokenMalformedError", ""}     //格式错误
	ErrTokenNotValidYet         = ErrInfo{TokenNotValidYetError, "TokenNotValidYetError", ""} //还未生效
	ErrTokenUnknown             = ErrInfo{TokenUnknownError, "TokenUnknownError", ""}         //未知错误
	ErrTokenKicked              = ErrInfo{TokenKickedError, "TokenKickedError", ""}
	ErrTokenNotExist            = ErrInfo{TokenNotExistError, "TokenNotExistError", ""} //在redis中不存在
	ErrTokenDifferentPlatformID = ErrInfo{TokenDifferentPlatformIDError, "TokenDifferentPlatformIDError", ""}
	ErrTokenDifferentUserID     = ErrInfo{TokenDifferentUserIDError, "TokenDifferentUserIDError", ""}

	ErrMessageHasReadDisable = ErrInfo{MessageHasReadDisable, "MessageHasReadDisable", ""}

	ErrDB        = ErrDatabase
	ErrSendLimit = ErrInternalServer
)

const (
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
	RegisterLimit        = 10012
	LoginLimit           = 10013
	InvitationError      = 10014
)

// 通用错误码
const (
	NoError             = 0     //无错误
	ArgsError           = 90001 //输入参数错误
	DatabaseError       = 90002 //redis/mysql等db错误
	ServerInternalError = 90003 //服务器内部错误
	NetworkError        = 90004 //网络错误
	NoPermissionError   = 90005 //权限不足
	GRPCConnIsNil       = 90006 //grpc连接空

	DefaultOtherError = 90006 //其他错误
	DataError         = 90007 //数据错误

	IdentityError = 90008 // 身份错误 非管理员token，且token中userID与请求userID不一致
)

// 账号错误码
const (
	UserIDNotFoundError  = 91001 //UserID不存在 或未注册
	GroupIDNotFoundError = 91002 //GroupID不存在
	RecordNotFoundError  = 91002 //记录不存在
)

// 关系链错误码
const (
	RelationshipAlreadyError = 92001 //已经是好友关系（或者黑名单）
	NotRelationshipYetError  = 92002 //不是好友关系（或者黑名单）
)

// 群组错误码
const (
	OnlyOneOwnerError        = 93001 //只能有一个群主
	InGroupAlreadyError      = 93003 //已在群组中
	NotInGroupYetError       = 93004 //不在群组中
	DismissedAlreadyError    = 93004 //群组已经解散
	OwnerNotAllowedQuitError = 93004 //群主不能退群
	GroupTypeNotSupport      = 93005
	GroupNoOwner             = 93006
)

// 用户错误码
const (
	RegisteredAlreadyError = 94001 //用户已经注册过了
)

// token错误码
const (
	TokenExpiredError             = 95001
	TokenInvalidError             = 95002
	TokenMalformedError           = 95003
	TokenNotValidYetError         = 95004
	TokenUnknownError             = 95005
	TokenKickedError              = 95006
	TokenDifferentPlatformIDError = 95007
	TokenDifferentUserIDError     = 95008
	TokenNotExistError            = 95009
)

// 消息错误码
const (
	MessageHasReadDisable = 96001
)

// temp

var (
	ErrServer            = ErrInfo{500, "server error", ""}
	ErrTencentCredential = ErrInfo{400, "ErrTencentCredential", ""}
)
