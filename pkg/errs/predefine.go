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

	ErrRecordNotFound = NewCodeError(RecordNotFoundError, "RecordNotFoundError")

	ErrNotInGroupYet       = NewCodeError(NotInGroupYetError, "NotInGroupYetError")
	ErrDismissedAlready    = NewCodeError(DismissedAlreadyError, "DismissedAlreadyError")
	ErrRegisteredAlready   = NewCodeError(RegisteredAlreadyError, "RegisteredAlreadyError")
	ErrGroupTypeNotSupport = NewCodeError(GroupTypeNotSupport, "")
	ErrGroupRequestHandled = NewCodeError(GroupRequestHandled, "GroupRequestHandled")

	ErrData             = NewCodeError(DataError, "DataError")
	ErrTokenExpired     = NewCodeError(TokenExpiredError, "TokenExpiredError")
	ErrTokenInvalid     = NewCodeError(TokenInvalidError, "TokenInvalidError")         //
	ErrTokenMalformed   = NewCodeError(TokenMalformedError, "TokenMalformedError")     //格式错误
	ErrTokenNotValidYet = NewCodeError(TokenNotValidYetError, "TokenNotValidYetError") //还未生效
	ErrTokenUnknown     = NewCodeError(TokenUnknownError, "TokenUnknownError")         //未知错误
	ErrTokenKicked      = NewCodeError(TokenKickedError, "TokenKickedError")
	ErrTokenNotExist    = NewCodeError(TokenNotExistError, "TokenNotExistError") //在redis中不存在
	ErrDuplicateKey     = NewCodeError(DuplicateKeyError, "DuplicateKeyError")

	ErrMessageHasReadDisable = NewCodeError(MessageHasReadDisable, "MessageHasReadDisable")

	ErrCanNotAddYourself   = NewCodeError(CanNotAddYourselfError, "CanNotAddYourselfError")
	ErrBlockedByPeer       = NewCodeError(BlockedByPeer, "BlockedByPeer")
	ErrNotPeersFriend      = NewCodeError(NotPeersFriend, "NotPeersFriend")
	ErrRelationshipAlready = NewCodeError(RelationshipAlreadyError, "RelationshipAlreadyError")

	ErrMutedInGroup     = NewCodeError(MutedInGroup, "MutedInGroup")
	ErrMutedGroup       = NewCodeError(MutedGroup, "MutedGroup")
	ErrMsgAlreadyRevoke = NewCodeError(MsgAlreadyRevoke, "MsgAlreadyRevoke")

	ErrConnOverMaxNumLimit = NewCodeError(ConnOverMaxNumLimit, "ConnOverMaxNumLimit")

	ErrConnArgsErr = NewCodeError(ConnArgsErr, "args err, need token, sendID, platformID")

	ErrFileUploadedExpired = NewCodeError(FileUploadedExpiredError, "FileUploadedExpiredError")
)
