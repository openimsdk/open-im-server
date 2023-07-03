package errs

var (
	ErrArgs             = NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission     = NewCodeError(NoPermissionError, "NoPermissionError")
	ErrDatabase         = NewCodeError(DatabaseError, "DatabaseError")
	ErrInternalServer   = NewCodeError(ServerInternalError, "ServerInternalError")
	ErrNetwork          = NewCodeError(NetworkError, "NetworkError")
	ErrCallback         = NewCodeError(CallbackError, "CallbackError")
	ErrCallbackContinue = NewCodeError(CallbackError, "ErrCallbackContinue")

	ErrUserIDNotFound  = NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")
	ErrGroupIDNotFound = NewCodeError(GroupIDNotFoundError, "GroupIDNotFoundError")
	ErrGroupIDExisted  = NewCodeError(GroupIDExisted, "GroupIDExisted")
	ErrUserIDExisted   = NewCodeError(UserIDExisted, "UserIDExisted")

	ErrRecordNotFound = NewCodeError(RecordNotFoundError, "RecordNotFoundError")

	ErrOnlyOneOwner        = NewCodeError(OnlyOneOwnerError, "OnlyOneOwnerError")
	ErrInGroupAlready      = NewCodeError(InGroupAlreadyError, "InGroupAlreadyError")
	ErrNotInGroupYet       = NewCodeError(NotInGroupYetError, "NotInGroupYetError")
	ErrDismissedAlready    = NewCodeError(DismissedAlreadyError, "DismissedAlreadyError")
	ErrOwnerNotAllowedQuit = NewCodeError(OwnerNotAllowedQuitError, "OwnerNotAllowedQuitError")
	ErrRegisteredAlready   = NewCodeError(RegisteredAlreadyError, "RegisteredAlreadyError")
	ErrGroupTypeNotSupport = NewCodeError(GroupTypeNotSupport, "")
	ErrGroupNoOwner        = NewCodeError(GroupNoOwner, "ErrGroupNoOwner")

	ErrData                     = NewCodeError(DataError, "DataError")
	ErrTokenExpired             = NewCodeError(TokenExpiredError, "TokenExpiredError")
	ErrTokenInvalid             = NewCodeError(TokenInvalidError, "TokenInvalidError")         //
	ErrTokenMalformed           = NewCodeError(TokenMalformedError, "TokenMalformedError")     //格式错误
	ErrTokenNotValidYet         = NewCodeError(TokenNotValidYetError, "TokenNotValidYetError") //还未生效
	ErrTokenUnknown             = NewCodeError(TokenUnknownError, "TokenUnknownError")         //未知错误
	ErrTokenKicked              = NewCodeError(TokenKickedError, "TokenKickedError")
	ErrTokenNotExist            = NewCodeError(TokenNotExistError, "TokenNotExistError") //在redis中不存在
	ErrTokenDifferentPlatformID = NewCodeError(TokenDifferentPlatformIDError, "TokenDifferentPlatformIDError")
	ErrTokenDifferentUserID     = NewCodeError(TokenDifferentUserIDError, "TokenDifferentUserIDError")
	ErrDuplicateKey             = NewCodeError(DuplicateKeyError, "DuplicateKeyError")

	ErrMessageHasReadDisable = NewCodeError(MessageHasReadDisable, "MessageHasReadDisable")

	ErrRelationshipAlready = NewCodeError(RelationshipAlreadyError, "RelationshipAlreadyError")
	ErrCanNotAddYourself   = NewCodeError(CanNotAddYourselfError, "CanNotAddYourselfError")
	ErrBlockedByPeer       = NewCodeError(BlockedByPeer, "BlockedByPeer")
	ErrNotPeersFriend      = NewCodeError(NotPeersFriend, "NotPeersFriend")

	ErrMutedInGroup     = NewCodeError(MutedInGroup, "MutedInGroup")
	ErrMutedGroup       = NewCodeError(MutedGroup, "MutedGroup")
	ErrUserNotRecvMsg   = NewCodeError(UserNotRecvMsg, "UserNotRecvMsg")
	ErrMsgAlreadyRevoke = NewCodeError(MsgAlreadyRevoke, "MsgAlreadyRevoke")

	ErrConnOverMaxNumLimit = NewCodeError(ConnOverMaxNumLimit, "ConnOverMaxNumLimit")

	ErrConnArgsErr   = NewCodeError(ConnArgsErr, "args err, need token, sendID, platformID")
	ErrConnUpdateErr = NewCodeError(ConnArgsErr, "upgrade http conn err")

	ErrConfig = NewCodeError(ConfigError, "ConfigError")

	ErrFileUploadedComplete = NewCodeError(FileUploadedCompleteError, "FileUploadedComplete")
	ErrFileUploadedExpired  = NewCodeError(FileUploadedExpiredError, "FileUploadedExpiredError")
	ErrGroupRequestHandled  = NewCodeError(GroupRequestHandled, "GroupRequestHandled")
)
