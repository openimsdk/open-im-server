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
	NoError                  = 0     //无错误
	DatabaseError            = 90002 //redis/mysql等db错误
	NetworkError             = 90004 //网络错误
	IdentityError            = 90008 // 身份错误 非管理员token，且token中userID与请求userID不一致
	GRPCConnIsNil            = 90006 //grpc连接空
	DefaultOtherError        = 90006 //其他错误
	DataError                = 90007 //数据错误
	ConfigError              = 90009
	CallbackError            = 80000
	RelationshipAlreadyError = 92001 //已经是好友关系（或者黑名单）
	NotRelationshipYetError  = 92002 //不是好友关系（或者黑名单）

	//通用错误码
	ServerInternalError = 500  //服务器内部错误
	ArgsError           = 1001 //输入参数错误
	NoPermissionError   = 1002 //权限不足
	DuplicateKeyError   = 1003
	RecordNotFoundError = 1004 //记录不存在

	// 账号错误码
	UserIDNotFoundError    = 1101 //UserID不存在 或未注册
	UserIDExisted          = 1102 //UserID已存在
	RegisteredAlreadyError = 1103 //用户已经注册过了

	// 群组错误码
	GroupIDNotFoundError     = 1201 //GroupID不存在
	GroupIDExisted           = 1202 //GroupID已存在
	OnlyOneOwnerError        = 1203 //只能有一个群主
	InGroupAlreadyError      = 1204 //已在群组中
	NotInGroupYetError       = 1205 //不在群组中
	DismissedAlreadyError    = 1206 //群组已经解散
	OwnerNotAllowedQuitError = 1207 //群主不能退群
	GroupTypeNotSupport      = 1208
	GroupNoOwner             = 1209
	GroupRequestHandled      = 1210

	// 关系链错误码
	CanNotAddYourselfError = 1301 //不能添加自己为好友
	BlockedByPeer          = 1302 //被对方拉黑
	NotPeersFriend         = 1303 //不是对方的好友

	// 消息错误码
	MessageHasReadDisable = 1401
	MutedInGroup          = 1402 //群成员被禁言
	MutedGroup            = 1403 //群被禁言
	UserNotRecvMsg        = 1404 //用户设置了不接收消息
	MsgAlreadyRevoke      = 1405 //消息已撤回

	// token错误码
	TokenExpiredError             = 1501
	TokenInvalidError             = 1502
	TokenMalformedError           = 1503
	TokenNotValidYetError         = 1504
	TokenUnknownError             = 1505
	TokenKickedError              = 1506
	TokenDifferentPlatformIDError = 1507
	TokenDifferentUserIDError     = 1508
	TokenNotExistError            = 1509

	// 长连接网关错误码
	ConnOverMaxNumLimit = 1601
	ConnArgsErr         = 1602
	ConnUpdateErr       = 1603

	// S3错误码
	FileUploadedCompleteError = 1701 // 文件已上传
	FileUploadedExpiredError  = 1702 // 上传过期
)
