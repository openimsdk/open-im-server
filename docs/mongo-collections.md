# OpenIM MongoDB 集合说明

本文档基于当前代码实现整理（`pkg/common/storage/database` + `pkg/common/storage/model`），用于说明 OpenIM 使用的 MongoDB 集合、核心字段和用途。

## 说明

- 集合来源：
  - 常量定义：`pkg/common/storage/database/name.go`
  - 硬编码集合：`db.Collection("...")`
- 字段口径：以结构体 `bson` tag 为准。
- 不包含测试代码中的临时集合（如 `msg3`、`version_test`）。

## 集合清单（按名称排序）

### `black`
- 用途：用户拉黑关系。
- 核心结构体：`model.Black`
- 字段：
  - `owner_user_id` (`string`)：拉黑关系所属用户
  - `block_user_id` (`string`)：被拉黑用户
  - `create_time` (`time.Time`)：创建时间
  - `add_source` (`int32`)：添加来源
  - `operator_user_id` (`string`)：操作人
  - `ex` (`string`)：扩展信息

### `conversation`
- 用途：用户会话配置和会话序列范围。
- 核心结构体：`model.Conversation`
- 字段：
  - `owner_user_id` (`string`)：会话所属用户
  - `conversation_id` (`string`)：会话 ID
  - `conversation_type` (`int32`)：会话类型
  - `user_id` (`string`)：单聊对端用户
  - `group_id` (`string`)：群聊群 ID
  - `recv_msg_opt` (`int32`)：消息接收选项
  - `is_pinned` (`bool`)：是否置顶
  - `is_private_chat` (`bool`)：是否私聊
  - `burn_duration` (`int32`)：阅后即焚时长
  - `group_at_type` (`int32`)：群 @ 设置
  - `attached_info` (`string`)：附加信息
  - `ex` (`string`)：扩展信息
  - `max_seq` (`int64`)：会话最大序列号
  - `min_seq` (`int64`)：会话最小序列号
  - `create_time` (`time.Time`)：创建时间
  - `is_msg_destruct` (`bool`)：是否开启消息自毁
  - `msg_destruct_time` (`int64`)：消息自毁时间
  - `latest_msg_destruct_time` (`time.Time`)：最近自毁更新时间

### `conversation_version`
- 用途：会话增量同步版本日志。
- 核心结构体：`model.VersionLogTable`
- 字段：
  - `_id` (`primitive.ObjectID`)：主键
  - `d_id` (`string`)：维度 ID（通常 owner 维度）
  - `logs` (`[]VersionLogElem`)：变更明细
  - `version` (`uint`)：当前版本
  - `deleted` (`uint`)：删除计数/版本标记
  - `last_update` (`time.Time`)：最后更新时间
- `logs` 元素字段（`model.VersionLogElem`）：
  - `e_id` (`string`)：元素 ID
  - `state` (`int32`)：状态（新增/删除/更新）
  - `version` (`uint`)：元素版本
  - `last_update` (`time.Time`)：元素更新时间

### `crypto_device`（硬编码）
- 用途：端到端加密设备信息。
- 核心结构体：`model.CryptoDevice`
- 字段：
  - `device_id` (`string`)：设备 ID
  - `user_id` (`string`)：用户 ID
  - `platform` (`string`)：平台
  - `device_model` (`string`)：设备型号
  - `app_version` (`string`)：应用版本
  - `virgil_identity` (`string`)：加密身份
  - `status` (`string`)：设备状态
  - `last_seen_at` (`time.Time`)：最近活跃时间
  - `create_time` (`time.Time`)：创建时间

### `friend`
- 用途：好友关系与好友属性。
- 核心结构体：`model.Friend`
- 字段：
  - `_id` (`primitive.ObjectID`)：主键
  - `owner_user_id` (`string`)：关系所属用户
  - `friend_user_id` (`string`)：好友用户
  - `remark` (`string`)：备注
  - `create_time` (`time.Time`)：创建时间
  - `add_source` (`int32`)：添加来源
  - `operator_user_id` (`string`)：操作人
  - `ex` (`string`)：扩展信息
  - `is_pinned` (`bool`)：是否置顶
  - `is_muted` (`bool`)：是否免打扰
  - `mute_duration` (`int64`)：免打扰时长
  - `mute_end_time` (`int64`)：免打扰结束时间

### `friend_request`
- 用途：好友申请记录。
- 核心结构体：`model.FriendRequest`
- 字段：
  - `from_user_id` (`string`)：申请人
  - `to_user_id` (`string`)：被申请人
  - `handle_result` (`int32`)：处理结果
  - `req_msg` (`string`)：申请附言
  - `create_time` (`time.Time`)：申请时间
  - `handler_user_id` (`string`)：处理人
  - `handle_msg` (`string`)：处理说明
  - `handle_time` (`time.Time`)：处理时间
  - `ex` (`string`)：扩展信息

### `friend_version`
- 用途：好友列表版本日志（增量同步）。
- 核心结构体：`model.VersionLogTable`
- 字段：同 `conversation_version`

### `group`
- 用途：群资料与群权限配置。
- 核心结构体：`model.Group`
- 字段：
  - `group_id` (`string`)：群 ID
  - `group_name` (`string`)：群名
  - `notification` (`string`)：群公告
  - `introduction` (`string`)：群介绍
  - `face_url` (`string`)：群头像
  - `create_time` (`time.Time`)：创建时间
  - `ex` (`string`)：扩展信息
  - `status` (`int32`)：群状态
  - `creator_user_id` (`string`)：创建者
  - `group_type` (`int32`)：群类型
  - `need_verification` (`int32`)：入群验证策略
  - `look_member_info` (`int32`)：成员信息可见策略
  - `apply_member_friend` (`int32`)：成员加好友策略
  - `notification_update_time` (`time.Time`)：公告更新时间
  - `notification_user_id` (`string`)：公告更新人
  - `allow_send_msg` (`int32`)：发言权限
  - `allow_pin_msg` (`int32`)：置顶消息权限
  - `allow_add_member` (`int32`)：拉人权限
  - `allow_edit_group_info` (`int32`)：编辑群资料权限

### `group_join_version`
- 用途：群成员入群事件版本日志。
- 核心结构体：`model.VersionLogTable`
- 字段：同 `conversation_version`

### `group_key_event`（硬编码）
- 用途：群密钥版本事件。
- 核心结构体：`model.GroupKeyEvent`
- 字段：
  - `event_id` (`string`)：事件 ID
  - `group_id` (`string`)：群 ID
  - `group_key_version` (`int64`)：群密钥版本
  - `event_type` (`string`)：事件类型
  - `operator_user_id` (`string`)：操作人
  - `create_time` (`time.Time`)：事件时间

### `group_key_version`（硬编码）
- 用途：记录群当前密钥版本。
- 核心结构体：`model.GroupKeyVersion`
- 字段：
  - `group_id` (`string`)：群 ID
  - `group_key_version` (`int64`)：当前密钥版本

### `group_member`
- 用途：群成员关系与成员属性。
- 核心结构体：`model.GroupMember`
- 字段：
  - `group_id` (`string`)：群 ID
  - `user_id` (`string`)：用户 ID
  - `nickname` (`string`)：群昵称
  - `face_url` (`string`)：头像
  - `role_level` (`int32`)：角色等级
  - `join_time` (`time.Time`)：入群时间
  - `join_source` (`int32`)：入群来源
  - `inviter_user_id` (`string`)：邀请人
  - `operator_user_id` (`string`)：操作人
  - `mute_end_time` (`time.Time`)：禁言结束时间
  - `ex` (`string`)：扩展信息

### `group_member_version`
- 用途：群成员版本日志。
- 核心结构体：`model.VersionLogTable`
- 字段：同 `conversation_version`

### `group_request`
- 用途：入群申请记录。
- 核心结构体：`model.GroupRequest`
- 字段：
  - `user_id` (`string`)：申请人
  - `group_id` (`string`)：目标群
  - `handle_result` (`int32`)：处理结果
  - `req_msg` (`string`)：申请信息
  - `handled_msg` (`string`)：处理说明
  - `req_time` (`time.Time`)：申请时间
  - `handle_user_id` (`string`)：处理人
  - `handled_time` (`time.Time`)：处理时间
  - `join_source` (`int32`)：入群来源
  - `inviter_user_id` (`string`)：邀请人
  - `ex` (`string`)：扩展信息

### `log`
- 用途：客户端日志上报元数据。
- 核心结构体：`model.Log`
- 字段：
  - `log_id` (`string`)：日志 ID
  - `platform` (`string`)：平台
  - `user_id` (`string`)：用户 ID
  - `create_time` (`time.Time`)：创建时间
  - `url` (`string`)：日志文件 URL
  - `file_name` (`string`)：文件名
  - `system_type` (`string`)：系统类型
  - `app_framework` (`string`)：应用框架
  - `version` (`string`)：版本
  - `ex` (`string`)：扩展信息

### `msg`
- 用途：聊天消息主存储（按会话 + 分片存放）。
- 核心结构体（分层）：
  - `model.MsgDocModel`
  - `model.MsgInfoModel`
  - `model.MsgDataModel`
  - `model.RevokeModel`
  - `model.OfflinePushModel`
- 字段：
  - `MsgDocModel`
    - `doc_id` (`string`)：文档 ID（会话 + 分片序号）
    - `msgs` (`[]*MsgInfoModel`)：消息数组
  - `MsgInfoModel`
    - `msg` (`*MsgDataModel`)：消息体
    - `revoke` (`*RevokeModel`)：撤回信息
    - `del_list` (`[]string`)：逻辑删除用户列表
    - `is_read` (`bool`)：读状态
  - `MsgDataModel`
    - `send_id` (`string`)：发送者
    - `recv_id` (`string`)：接收者（单聊）
    - `group_id` (`string`)：群 ID（群聊）
    - `client_msg_id` (`string`)：客户端消息 ID
    - `server_msg_id` (`string`)：服务端消息 ID
    - `sender_platform_id` (`int32`)：发送平台
    - `sender_nickname` (`string`)：发送者昵称
    - `sender_face_url` (`string`)：发送者头像
    - `session_type` (`int32`)：会话类型
    - `msg_from` (`int32`)：消息来源
    - `content_type` (`int32`)：内容类型
    - `content` (`string`)：消息内容
    - `seq` (`int64`)：消息序号
    - `send_time` (`int64`)：发送时间
    - `create_time` (`int64`)：创建时间
    - `status` (`int32`)：状态
    - `is_read` (`bool`)：已读标记
    - `options` (`map[string]bool`)：扩展选项
    - `offline_push` (`*OfflinePushModel`)：离线推送
    - `at_user_id_list` (`[]string`)：@用户列表
    - `attached_info` (`string`)：附加信息
    - `ex` (`string`)：扩展信息
  - `RevokeModel`
    - `role` (`int32`)：撤回角色
    - `user_id` (`string`)：撤回人 ID
    - `nickname` (`string`)：撤回人昵称
    - `time` (`int64`)：撤回时间
  - `OfflinePushModel`
    - `title` (`string`)：推送标题
    - `desc` (`string`)：推送描述
    - `ex` (`string`)：扩展
    - `ios_push_sound` (`string`)：iOS 声音
    - `ios_badge_count` (`bool`)：iOS 角标策略

### `phone_sn_info`
- 用途：手机号与发送状态信息。
- 核心结构体：`model.PhoneSNInfo`
- 字段：
  - `phone` (`string`)：手机号
  - `user_id` (`int64`)：关联用户 ID
  - `is_snd` (`bool`)：业务发送标识
  - `send_count` (`int64`)：最近 1 分钟发送计数（当前实现新增）
  - `update_time` (`int64`)：更新时间（毫秒）

### `s3`
- 用途：对象存储文件元数据。
- 核心结构体：`model.Object`
- 字段：
  - `name` (`string`)：对象名
  - `user_id` (`string`)：上传用户
  - `hash` (`string`)：文件哈希
  - `engine` (`string`)：存储引擎
  - `key` (`string`)：对象 key
  - `size` (`int64`)：文件大小
  - `content_type` (`string`)：MIME 类型
  - `group` (`string`)：分组
  - `create_time` (`time.Time`)：创建时间

### `seq`
- 用途：会话级消息序列范围。
- 核心结构体：`model.SeqConversation`
- 字段：
  - `conversation_id` (`string`)：会话 ID
  - `max_seq` (`int64`)：最大序列
  - `min_seq` (`int64`)：最小序列

### `seq_user`
- 用途：用户会话维度消息序列与已读位点。
- 核心结构体：`model.SeqUser`
- 字段：
  - `user_id` (`string`)：用户 ID
  - `conversation_id` (`string`)：会话 ID
  - `min_seq` (`int64`)：最小序列
  - `max_seq` (`int64`)：最大序列
  - `read_seq` (`int64`)：已读序列

### `signal_invitation`
- 用途：音视频邀请信令（支持 TTL 过期清理）。
- 核心结构体：`model.SignalInvitation`
- 字段：
  - `room_id` (`string`)：房间 ID
  - `inviter_user_id` (`string`)：邀请人
  - `invitee_user_id_list` (`[]string`)：被邀请人列表
  - `custom_data` (`string`)：自定义数据
  - `group_id` (`string`)：群 ID
  - `timeout` (`int32`)：超时时间
  - `media_type` (`string`)：媒体类型
  - `platform_id` (`int32`)：平台 ID
  - `session_type` (`int32`)：会话类型
  - `initiate_time` (`int64`)：发起时间
  - `busy_line_user_id_list` (`[]string`)：忙线用户
  - `offline_push_title` (`string`)：离线推送标题
  - `offline_push_desc` (`string`)：离线推送描述
  - `offline_push_ex` (`string`)：离线推送扩展
  - `create_time` (`int64`)：创建时间
  - `expire_at` (`time.Time`)：过期时间（TTL）

### `signal_record`
- 用途：音视频通话记录。
- 核心结构体：`model.SignalRecord`
- 字段：
  - `sid` (`string`)：记录 ID
  - `room_id` (`string`)：房间 ID
  - `file_name` (`string`)：文件名
  - `media_type` (`string`)：媒体类型
  - `session_type` (`int32`)：会话类型
  - `inviter_user_id` (`string`)：邀请人
  - `inviter_user_nickname` (`string`)：邀请人昵称
  - `group_id` (`string`)：群 ID
  - `group_name` (`string`)：群名
  - `inviter_user_id_list` (`[]string`)：参与用户列表
  - `send_id` (`string`)：发送者
  - `recv_id` (`string`)：接收者
  - `create_time` (`int64`)：创建时间
  - `end_time` (`int64`)：结束时间
  - `file_size` (`string`)：文件大小
  - `file_url` (`string`)：文件 URL

### `spam_report`
- 用途：垃圾消息/用户举报与处理记录。
- 核心结构体：`model.SpamReport`
- 字段：
  - `_id` (`primitive.ObjectID`)：主键
  - `report_id` (`string`)：举报 ID
  - `reporter_user_id` (`string`)：举报人
  - `reported_user_id` (`string`)：被举报人
  - `conversation_id` (`string`)：会话 ID
  - `client_msg_id` (`string`)：客户端消息 ID
  - `seq` (`int64`)：消息序号
  - `reason_type` (`int32`)：举报类型
  - `reason` (`string`)：举报原因
  - `status` (`int32`)：处理状态
  - `create_time` (`time.Time`)：创建时间
  - `handle_time` (`time.Time`)：处理时间
  - `handler_user_id` (`string`)：处理人
  - `ex` (`string`)：扩展信息

### `user`
- 用途：用户主资料与可见性/通话/消息设置。
- 核心结构体：`model.User`
- 字段：
  - `user_id` (`string`)：用户 ID
  - `nickname` (`string`)：昵称
  - `face_url` (`string`)：头像
  - `ex` (`string`)：扩展
  - `app_manger_level` (`int32`)：应用管理级别
  - `global_recv_msg_opt` (`int32`)：全局消息接收选项
  - `create_time` (`time.Time`)：创建时间
  - `phone` (`string`)：手机号
  - `phone_visibility` (`int32`)：手机号可见性
  - `call_accept_setting` (`int32`)：通话接受策略
  - `msg_receive_setting` (`int32`)：消息接收策略

### `userCommands`（硬编码）
- 用途：用户命令记录。
- 数据来源：`mgo/user.go` 内部通过 `bson.M` 直接操作。
- 字段：
  - `userID` (`string`)：用户 ID
  - `type` (`int32`)：命令类型
  - `uuid` (`string`)：命令唯一 ID
  - `createTime` (`int64`)：创建时间（秒）
  - `value` (`string`)：值
  - `ex` (`string`)：扩展

### `user_global_black_list`
- 用途：全局黑名单。
- 核心结构体：`model.UserGlobalBlack`
- 字段：
  - `user_id` (`string`)：被封禁用户 ID
  - `nickname` (`string`)：昵称
  - `operator_id` (`string`)：操作人
  - `reason` (`string`)：原因
  - `create_time` (`time.Time`)：创建时间

## 常见索引说明（摘要）

- `msg`：`doc_id` 唯一索引。
- `conversation`：`(owner_user_id, conversation_id)` 唯一索引。
- `friend`：`(owner_user_id, friend_user_id)` 唯一索引。
- `group_member`：`(group_id, user_id)` 唯一索引。
- `signal_invitation`：`room_id` 唯一索引，`expire_at` TTL 索引。
- `phone_sn_info`：`phone` 唯一索引。
- `spam_report`：`report_id` 唯一索引。

