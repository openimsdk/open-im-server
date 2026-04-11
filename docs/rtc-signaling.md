# OpenIM 音视频（RTC）信令与媒体 — 技术说明

## 1. 职责边界

| 维度 | 说明 |
|------|------|
| OpenIM | 呼叫信令编排、邀请状态（Mongo）、LiveKit 房间创建/删除、进房 JWT、通过消息链路把信令投递到对端（在线 WebSocket + 离线推送）。 |
| LiveKit | WebRTC 媒体面（客户端持 Token 连接 `externalAddress`）。 |
| 协议 | `protocol/rtc/rtc.proto`；通知类 `ContentType` 见 `protocol/constant/rtc.go`。 |

---

## 2. 服务与配置

| 项 | 位置 |
|----|------|
| RTC 进程入口 | `pkg/common/cmd/rpc_rtc.go` → `internal/rpc/rtc.Start` |
| 实现 | `internal/rpc/rtc/server.go`、`internal/rpc/rtc/signal.go` |
| 注册名 | `config.Share.RpcRegisterName.Rtc` |
| LiveKit 配置 | `pkg/common/config/config.go`：`LiveKit`（`internalAddress`、`externalAddress`、`apiKey`、`apiSecret`、`tokenExpiry`） |

---

## 3. 接入方式（概览）

### 3.1 WebSocket

1. `internal/msggateway/ws_server.go`：注入 `RtcServiceClient`。
2. `internal/msggateway/client.go`：`ReqIdentifier == WSSendSignalMsg`（**1004**）。
3. `internal/msggateway/message_handler.go`：`SendSignalMessage` → `SignalMessageAssemble`（体为 `SignalMessageAssembleReq` 或裸 `SignalReq`）。
4. 返回：`SignalResp` 的 protobuf 二进制放在 `Resp.Data`。

### 3.2 HTTP（`/rtc`）

路由：`internal/api/router.go`；封装：`internal/api/rtc.go`（`a2r.Call` 将 HTTP 体映射到同名 gRPC）。Prometheus 发现：`GET /prometheus_discovery/rtc`。

**详细接口与链路见第 4 节。**

---

## 4. 接口清单与各接口调用链路

对外暴露形态包括：**gRPC 方法名**（服务 `openim.rtc.RtcService`）、**HTTP**（OpenIM API 网关）、以及 **WebSocket**（仅映射到 `SignalMessageAssemble`）。以下链路按代码真实调用顺序描述。

### 4.1 接口总表

| # | gRPC 方法 | HTTP 路径 | WebSocket |
|---|-----------|-----------|-----------|
| 1 | `SignalMessageAssemble` | `POST /rtc/signal_message_assemble` | `ReqIdentifier=1004`（`WSSendSignalMsg`） |
| 2 | `SignalGetRoomByGroupID` | `POST /rtc/signal_get_room_by_group_id` | — |
| 3 | `SignalGetTokenByRoomID` | `POST /rtc/signal_get_token_by_room_id` | — |
| 4 | `SignalGetRooms` | `POST /rtc/signal_get_rooms` | — |
| 5 | `GetSignalInvitationInfo` | `POST /rtc/get_signal_invitation_info` | — |
| 6 | `GetSignalInvitationInfoStartApp` | `POST /rtc/get_signal_invitation_info_start_app` | — |
| 7 | `SignalSendCustomSignal` | `POST /rtc/signal_send_custom_signal` | — |
| 8 | `GetSignalInvitationRecords` | `POST /rtc/get_signal_invitation_records` | — |
| 9 | `DeleteSignalRecords` | `POST /rtc/delete_signal_records` | — |

说明：`SignalReq` 内嵌的 **`getTokenByRoomID`** 与独立 RPC **`SignalGetTokenByRoomID`** 在服务端均落到 `genToken` + 返回 `liveURL`，前者经 `SignalMessageAssemble` 分发，后者经 HTTP 直达同名 gRPC。

---

### 4.2 `SignalMessageAssemble`

**作用**：处理一路信令请求，返回 `SignalResp`；部分分支会写 Mongo、调 LiveKit、并通过 Msg 服务发 1601 通知。

**入口 A — WebSocket**

1. 客户端发送二进制帧 → `internal/msggateway/client.go` 按 `ReqIdentifier==1004` 分支。
2. `LongConnServer.SendSignalMessage` → `GrpcHandler.SendSignalMessage`（`internal/msggateway/message_handler.go`）。
3. `proto.Unmarshal`：`SignalMessageAssembleReq`；若失败则解 `SignalReq` 并填入 `assembleReq.SignalReq`。
4. `RtcServiceClient.SignalMessageAssemble(ctx, assembleReq)`（gRPC 至 **rtc 进程**）。
5. `internal/rpc/rtc/signal.go`：`rtcServer.SignalMessageAssemble` → `switch req.SignalReq.Payload` → `handleInvite` / `handleInviteInGroup` / `handleCancel` / `handleAccept` / `handleHungUp` / `handleReject` / `handleGetTokenByRoomID`。
6. 返回 `SignalMessageAssembleResp` → 网关将 `SignalResp` `proto.Marshal` → `Resp.Data` 回客户端。

**入口 B — HTTP**

1. `POST /rtc/signal_message_assemble` → `internal/api/rtc.go`：`RtcApi.SignalMessageAssemble`。
2. `github.com/openimsdk/tools/a2r.Call`：解析 Gin 请求体 → 调用 `RtcServiceClient.SignalMessageAssemble`。
3. 后续与步骤 5–6 相同（响应经 HTTP 返回，而非 WS `Resp`）。

**分支内典型下游（仅当对应 payload 触发时）**

| 子逻辑 | LiveKit | Mongo（`controller.RtcDatabase` → `mgo/signal`） | Msg（`rpcli.MsgClient.SendMsg`） |
|--------|---------|---------------------------------------------------|----------------------------------|
| `handleInvite` | `CreateRoom` | `CreateInvitation` | 对每个被叫 `sendSignalingNotification`（1601） |
| `handleInviteInGroup` | `CreateRoom` | `CreateInvitation` | 同上（`SessionType` 为群） |
| `handleAccept` | — | — | 通知主叫 1601 |
| `handleReject` | — | `DeleteInvitation` / `RemoveInvitee` | 通知主叫 1601 |
| `handleCancel` | — | `DeleteInvitation` | 通知被叫 1601 |
| `handleHungUp` | `DeleteRoom` | `DeleteInvitation` | 通知对端 1601 |
| `handleGetTokenByRoomID` | — | — | — |

**若发生 `SendMsg`（1601）**，后续链路见 **第 7 节**（Kafka → msg_transfer → push → 网关 `WSPushMsg` 2001）。

---

### 4.3 `SignalGetRoomByGroupID`

**作用**：按群 ID 查当前（或最近）邀请信息，返回 `InvitationInfo` 与 `roomID`。

**HTTP 链路**

1. `POST /rtc/signal_get_room_by_group_id` → `RtcApi.SignalGetRoomByGroupID` → `a2r.Call` → gRPC。
2. `internal/rpc/rtc/signal.go`：`SignalGetRoomByGroupID` → `db.GetInvitationByGroupID` → Mongo `signal_invitation`（`mgo/signal.go`）。
3. `modelToInvitationInfo` 填响应返回。

**不经过**：LiveKit、Msg、Kafka。

---

### 4.4 `SignalGetTokenByRoomID`（独立 RPC）

**作用**：已有房间时，仅为指定用户签发 LiveKit JWT，并返回 `liveURL`（`ExternalAddress`）。

**HTTP 链路**

1. `POST /rtc/signal_get_token_by_room_id` → `RtcApi.SignalGetTokenByRoomID` → `a2r.Call` → gRPC。
2. `internal/rpc/rtc/signal.go`：`SignalGetTokenByRoomID`（与 `handleGetTokenByRoomID` 同源逻辑）→ `genToken(roomID, userID)`。
3. 返回 `SignalGetTokenByRoomIDResp`。

**不经过**：Mongo（不校验邀请是否存在）、Msg、Kafka。

---

### 4.5 `SignalGetRooms`

**作用**：批量 `roomID` 查询邀请信息列表。

**HTTP 链路**

1. `POST /rtc/signal_get_rooms` → `RtcApi.SignalGetRooms` → `a2r.Call` → gRPC。
2. `SignalGetRooms` → `db.GetInvitationsByRoomIDs` → Mongo。
3. 组装 `[]*SignalGetRoomByGroupIDResp` 返回。

**不经过**：LiveKit、Msg、Kafka。

---

### 4.6 `GetSignalInvitationInfo`

**作用**：按 **roomID** 查邀请详情及离线推送字段。

**HTTP 链路**

1. `POST /rtc/get_signal_invitation_info` → `RtcApi.GetSignalInvitationInfo` → `a2r.Call` → gRPC。
2. `GetSignalInvitationInfo` → `db.GetInvitationByRoomID` → Mongo。
3. 填充 `InvitationInfo`、`OfflinePushInfo` 返回。

**不经过**：LiveKit、Msg、Kafka。

---

### 4.7 `GetSignalInvitationInfoStartApp`

**作用**：按 **被叫 userID** 查其相关待处理邀请（冷启动拉铃场景）。

**HTTP 链路**

1. `POST /rtc/get_signal_invitation_info_start_app` → `RtcApi.GetSignalInvitationInfoStartApp` → `a2r.Call` → gRPC。
2. `GetSignalInvitationInfoStartApp` → `db.GetInvitationByInviteeUserID` → Mongo（`invitee_user_id_list` 查询）。
3. 返回邀请与 `OfflinePushInfo`。

**不经过**：LiveKit、Msg、Kafka。

---

### 4.8 `SignalSendCustomSignal`

**作用**：向房间内除操作者外的参与者广播 **自定义信令**（系统消息 **1605**）。

**HTTP 链路**

1. `POST /rtc/signal_send_custom_signal` → `RtcApi.SignalSendCustomSignal` → `a2r.Call` → gRPC。
2. `SignalSendCustomSignal` → `db.GetInvitationByRoomID`（取邀请内 `InviteeUserIDList` + `InviterUserID`）。
3. `mcontext.GetOpUserID(ctx)` 排除发送者自己。
4. 对每个目标用户 `sendCustomSignalNotification` → `MsgClient.SendMsg`（`ContentType=CustomSignalNotification`，JSON body）。
5. 若第 2 步查无邀请：打日志后返回空成功（不报错）。

**若发生 `SendMsg`**：后续同第 7 节（1605 走消息总线与推送）。

---

### 4.9 `GetSignalInvitationRecords`

**作用**：分页查询通话/信令话单（`signal_record`）。

**HTTP 链路**

1. `POST /rtc/get_signal_invitation_records` → `RtcApi.GetSignalInvitationRecords` → `a2r.Call` → gRPC。
2. `GetSignalInvitationRecords` → `db.SearchRecords`（`sendID` / `recvID` / `sessionType` / 时间范围 / 分页）→ Mongo `signal_record`。
3. 映射为 `[]*rtc.SignalRecord` 返回。

**不经过**：LiveKit、Msg、Kafka。

---

### 4.10 `DeleteSignalRecords`

**作用**：按话单主键 `SID` 列表删除记录。

**HTTP 链路**

1. `POST /rtc/delete_signal_records` → `RtcApi.DeleteSignalRecords` → `a2r.Call` → gRPC。
2. `DeleteSignalRecords` → `db.DeleteRecords(sIDs)` → Mongo。

**不经过**：LiveKit、Msg、Kafka。

---

## 5. `SignalMessageAssemble` 行为摘要（payload 与副作用）

（实现：`internal/rpc/rtc/signal.go`）

| 动作 | LiveKit | Mongo | 通知 |
|------|---------|-------|------|
| Invite | CreateRoom | CreateInvitation | 向被叫发 1601 |
| InviteInGroup | CreateRoom | CreateInvitation | 向被叫发 1601（群 SessionType） |
| Accept | — | — | 通知主叫 1601 |
| Reject | — | DeleteInvitation / RemoveInvitee | 通知主叫 |
| Cancel | — | DeleteInvitation | 通知被叫 |
| HungUp | **DeleteRoom** | DeleteInvitation | 通知对端 |
| GetTokenByRoomID（嵌在 SignalReq） | — | — | — |

Token：`github.com/livekit/protocol/auth`，`VideoGrant`（`RoomJoin` + `Room` + `Identity`），有效期由配置决定。

---

## 6. Mongo

- 集合：`signal_invitation`、`signal_record`（`pkg/common/storage/database/name.go`）。
- 模型：`pkg/common/storage/model/signal.go`。
- DAO：`pkg/common/storage/database/mgo/signal.go`。
- 控制器：`pkg/common/storage/controller/rtc.go`。

话单 `SignalRecord` 的写入需结合业务；`GetSignalInvitationRecords` 依赖该集合已有数据。

---

## 7. 信令进消息链路（`SendMsg` 之后）

适用于：`sendSignalingNotification`（1601）、`sendCustomSignalNotification`（1605）。

1. `MsgClient.SendMsg` → `internal/rpc/msg/send.go`（按 `SessionType` 走单聊/群聊分支）。
2. `MsgToMQ` → Kafka **`toRedisTopic`**（key：单聊为 `GenConversationUniqueKeyForSingle`；群为 `GroupID`）。
3. `msg_transfer`（`internal/msgtransfer/online_history_msg_handler.go`）消费 → Redis seq → **`toMongoTopic`** → **`toPushTopic`**。
4. `push`（`internal/push/push_handler.go`）消费 `toPushTopic` → `Push2User` / `Push2Group` → 网关 RPC。
5. 网关 `Client.PushMessage`（`internal/msggateway/client.go`）：**`ReqIdentifier = WSPushMsg`（2001）**，`Data` 为 `sdkws.PushMessages` 的 protobuf。

离线推送：`SignalingNotification` 可走离线；`RoomParticipantsConnected/Disconnected`（1602/1603）在 push 逻辑中默认不触发离线推。

---

## 8. `MsgData.Options` 与会话 ID（缺省行为）

`pkg/msgprocessor/options.go`：`Options.Is(key)` 在 **key 未设置时视为 true**。

RTC 侧 `sendSignalingNotification` 使用 `make(map[string]bool)` 空 map，故 `IsHistory` / `IsNotNotification` 等表现为 true，信令在 transfer 中多走**落库 + 带 seq 后推送**路径。

网关下行使用 `GetConversationIDByMsg`：单聊信令默认挂在 **`si_*`** 的 `PushMessages.Msgs` 中（`IsNotification` 为 false 的前缀规则）。

---

## 9. 已知风险与排查

- **群通话信令**：当前构造通知时若未设置 `MsgData.GroupID`，`sendMsgGroupChat` 的 Kafka key 与 `Push2Group(ctx, groupID, ...)` 可能异常；建议在发群信令时写入与 `InvitationInfo.group_id` 一致的 `GroupID`。
- **常量 1602–1604**：协议与 push 有特殊分支，但 `internal/rpc/rtc` 主路径主要发 1601/1605；若产品需要房间成员/流状态通知，需在扩展路径发送。

---

## 10. 常量（节选）

| 值 | 含义 |
|----|------|
| 1601 | `SignalingNotification` |
| 1605 | `CustomSignalNotification` |
| 1602–1604 | 房间参与者/流变更等（push 对 1602/1603 限制离线推） |

---

## 11. 端到端链路（简图）

```text
客户端 → [WS 1004 或 HTTP /rtc] → rtc RPC → LiveKit + Mongo
                    → msg SendMsg → Kafka(toRedis) → msg_transfer → Kafka(toPush)
                    → push → msg_gateway → WS 2001 (PushMessages)
客户端 ← LiveKit(media) + OpenIM(信令推送)
```
