package errs

// UnknownCode 没有解析到code或解析的code=0
const UnknownCode = 1000

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

	ConfigError = 90009

	CallbackError = 80000
)

const DuplicateKeyError = 12345

// 账号错误码
const (
	UserIDNotFoundError  = 91001 //UserID不存在 或未注册
	GroupIDNotFoundError = 91002 //GroupID不存在
	RecordNotFoundError  = 91002 //记录不存在
	GroupIDIDExisted     = 91002 //GroupID已存在
	UserIDExisted        = 91002 //UserID已存在
)

// 关系链错误码
const (
	RelationshipAlreadyError = 92001 //已经是好友关系（或者黑名单）
	NotRelationshipYetError  = 92002 //不是好友关系（或者黑名单）
	CanNotAddYourselfError   = 92003 //不能添加自己为好友
	BlockedByPeer            = 92004 //被对方拉黑
	NotPeersFriend           = 92005 //不是对方的好友
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

	MutedInGroup = 93007 //群成员被禁言
	MutedGroup   = 93008 //群被禁言
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

// 长连接网关错误码
const (
	ConnOverMaxNumLimit = 970001
	ConnArgsErr         = 970002
	ConnUpdateErr       = 970003
)
