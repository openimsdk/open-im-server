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

package errs

var (
	ErrArgs             = NewCodeError(ArgsError, "ArgsError")
	ErrDatabase         = NewCodeError(DatabaseError, "DatabaseError")
	ErrInternalServer   = NewCodeError(ServerInternalError, "ServerInternalError")
	ErrNetwork          = NewCodeError(NetworkError, "NetworkError")
	ErrNoPermission     = NewCodeError(NoPermissionError, "NoPermissionError")
	ErrIdentity         = NewCodeError(IdentityError, "IdentityError")
	ErrCallback         = NewCodeError(CallbackError, "CallbackError")
	ErrCallbackContinue = NewCodeError(CallbackError, "ErrCallbackContinue")

	ErrUserIDNotFound  = NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")
	ErrGroupIDNotFound = NewCodeError(GroupIDNotFoundError, "GroupIDNotFoundError")
	ErrGroupIDExisted  = NewCodeError(GroupIDExisted, "GroupIDExisted")
	ErrUserIDExisted   = NewCodeError(UserIDExisted, "UserIDExisted")

	ErrRecordNotFound = NewCodeError(RecordNotFoundError, "RecordNotFoundError")

	ErrRelationshipAlready = NewCodeError(RelationshipAlreadyError, "RelationshipAlreadyError")
	ErrNotRelationshipYet  = NewCodeError(NotRelationshipYetError, "NotRelationshipYetError")
	ErrCanNotAddYourself   = NewCodeError(CanNotAddYourselfError, "CanNotAddYourselfError")

	ErrOnlyOneOwner        = NewCodeError(OnlyOneOwnerError, "OnlyOneOwnerError")
	ErrInGroupAlready      = NewCodeError(InGroupAlreadyError, "InGroupAlreadyError")
	ErrNotInGroupYet       = NewCodeError(NotInGroupYetError, "NotInGroupYetError")
	ErrDismissedAlready    = NewCodeError(DismissedAlreadyError, "DismissedAlreadyError")
	ErrOwnerNotAllowedQuit = NewCodeError(OwnerNotAllowedQuitError, "OwnerNotAllowedQuitError")
	ErrRegisteredAlready   = NewCodeError(RegisteredAlreadyError, "RegisteredAlreadyError")
	ErrGroupTypeNotSupport = NewCodeError(GroupTypeNotSupport, "")
	ErrGroupNoOwner        = NewCodeError(GroupNoOwner, "ErrGroupNoOwner")

	ErrDefaultOther             = NewCodeError(DefaultOtherError, "DefaultOtherError")
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

	ErrBlockedByPeer = NewCodeError(BlockedByPeer, "BlockedByPeer")
	//不是对方的好友
	ErrNotPeersFriend = NewCodeError(NotPeersFriend, "NotPeersFriend")

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
