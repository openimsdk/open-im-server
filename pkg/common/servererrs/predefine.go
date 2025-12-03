package servererrs

import "github.com/openimsdk/tools/errs"

var (
	ErrSecretNotChanged = errs.NewCodeError(SecretNotChangedError, "secret not changed, please change secret in config/share.yml for security reasons")

	ErrDatabase         = errs.NewCodeError(DatabaseError, "DatabaseError")
	ErrNetwork          = errs.NewCodeError(NetworkError, "NetworkError")
	ErrCallback         = errs.NewCodeError(CallbackError, "CallbackError")
	ErrCallbackContinue = errs.NewCodeError(CallbackError, "ErrCallbackContinue")

	ErrInternalServer = errs.NewCodeError(ServerInternalError, "ServerInternalError")
	ErrArgs           = errs.NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission   = errs.NewCodeError(NoPermissionError, "NoPermissionError")
	ErrDuplicateKey   = errs.NewCodeError(DuplicateKeyError, "DuplicateKeyError")
	ErrRecordNotFound = errs.NewCodeError(RecordNotFoundError, "RecordNotFoundError")

	ErrUserIDNotFound  = errs.NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")
	ErrGroupIDNotFound = errs.NewCodeError(GroupIDNotFoundError, "GroupIDNotFoundError")
	ErrGroupIDExisted  = errs.NewCodeError(GroupIDExisted, "GroupIDExisted")

	ErrNotInGroupYet       = errs.NewCodeError(NotInGroupYetError, "NotInGroupYetError")
	ErrDismissedAlready    = errs.NewCodeError(DismissedAlreadyError, "DismissedAlreadyError")
	ErrRegisteredAlready   = errs.NewCodeError(RegisteredAlreadyError, "RegisteredAlreadyError")
	ErrGroupTypeNotSupport = errs.NewCodeError(GroupTypeNotSupport, "")
	ErrGroupRequestHandled = errs.NewCodeError(GroupRequestHandled, "GroupRequestHandled")

	ErrData             = errs.NewCodeError(DataError, "DataError")
	ErrTokenExpired     = errs.NewCodeError(TokenExpiredError, "TokenExpiredError")
	ErrTokenInvalid     = errs.NewCodeError(TokenInvalidError, "TokenInvalidError")         //
	ErrTokenMalformed   = errs.NewCodeError(TokenMalformedError, "TokenMalformedError")     //
	ErrTokenNotValidYet = errs.NewCodeError(TokenNotValidYetError, "TokenNotValidYetError") //
	ErrTokenUnknown     = errs.NewCodeError(TokenUnknownError, "TokenUnknownError")         //
	ErrTokenKicked      = errs.NewCodeError(TokenKickedError, "TokenKickedError")
	ErrTokenNotExist    = errs.NewCodeError(TokenNotExistError, "TokenNotExistError") //

	ErrMessageHasReadDisable = errs.NewCodeError(MessageHasReadDisable, "MessageHasReadDisable")

	ErrCanNotAddYourself   = errs.NewCodeError(CanNotAddYourselfError, "CanNotAddYourselfError")
	ErrBlockedByPeer       = errs.NewCodeError(BlockedByPeer, "BlockedByPeer")
	ErrNotPeersFriend      = errs.NewCodeError(NotPeersFriend, "NotPeersFriend")
	ErrRelationshipAlready = errs.NewCodeError(RelationshipAlreadyError, "RelationshipAlreadyError")

	ErrMutedInGroup     = errs.NewCodeError(MutedInGroup, "MutedInGroup")
	ErrMutedGroup       = errs.NewCodeError(MutedGroup, "MutedGroup")
	ErrMsgAlreadyRevoke = errs.NewCodeError(MsgAlreadyRevoke, "MsgAlreadyRevoke")

	ErrConnOverMaxNumLimit = errs.NewCodeError(ConnOverMaxNumLimit, "ConnOverMaxNumLimit")

	ErrConnArgsErr          = errs.NewCodeError(ConnArgsErr, "args err, need token, sendID, platformID")
	ErrPushMsgErr           = errs.NewCodeError(PushMsgErr, "push msg err")
	ErrIOSBackgroundPushErr = errs.NewCodeError(IOSBackgroundPushErr, "ios background push err")

	ErrFileUploadedExpired = errs.NewCodeError(FileUploadedExpiredError, "FileUploadedExpiredError")
)
