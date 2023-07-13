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

// UnknownCode 没有解析到code或解析的code=0.
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

// 通用错误码.
const (
	NoError       = 0     // 无错误
	DatabaseError = 90002 // redis/mysql等db错误
	NetworkError  = 90004 // 网络错误
	DataError     = 90007 // 数据错误

	CallbackError = 80000

	//通用错误码.
	ServerInternalError = 500  //服务器内部错误
	ArgsError           = 1001 //输入参数错误
	NoPermissionError   = 1002 //权限不足
	DuplicateKeyError   = 1003
	RecordNotFoundError = 1004 // 记录不存在

	// 账号错误码.
	UserIDNotFoundError    = 1101 // UserID不存在 或未注册
	RegisteredAlreadyError = 1102 // 用户已经注册过了

	// 群组错误码.
	GroupIDNotFoundError  = 1201 // GroupID不存在
	GroupIDExisted        = 1202 // GroupID已存在
	NotInGroupYetError    = 1203 // 不在群组中
	DismissedAlreadyError = 1204 // 群组已经解散
	GroupTypeNotSupport   = 1205
	GroupRequestHandled   = 1206

	// 关系链错误码.
	CanNotAddYourselfError   = 1301 // 不能添加自己为好友
	BlockedByPeer            = 1302 // 被对方拉黑
	NotPeersFriend           = 1303 // 不是对方的好友
	RelationshipAlreadyError = 1304 // 已经是好友关系

	// 消息错误码.
	MessageHasReadDisable = 1401
	MutedInGroup          = 1402 // 群成员被禁言
	MutedGroup            = 1403 // 群被禁言
	MsgAlreadyRevoke      = 1404 // 消息已撤回

	// token错误码.
	TokenExpiredError     = 1501
	TokenInvalidError     = 1502
	TokenMalformedError   = 1503
	TokenNotValidYetError = 1504
	TokenUnknownError     = 1505
	TokenKickedError      = 1506
	TokenNotExistError    = 1507

	// 长连接网关错误码.
	ConnOverMaxNumLimit = 1601
	ConnArgsErr         = 1602

	// S3错误码.
	FileUploadedExpiredError = 1701 // 上传过期
)
