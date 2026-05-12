# RedPacket 后端接口说明

本文档按当前 `internal/api/redpacket.go` 与 `internal/rpc/redpacket/*` 实现整理。红包服务已经从独立 Gin 服务迁移为 OpenIM 标准 RPC 服务：

- HTTP 入口在 `internal/api` 网关，路由前缀为 `/redpacket`
- 网关通过 `pbredpacket.RedPacketClient` 调用 `internal/rpc/redpacket`
- RPC 服务使用 MongoDB 存储，通过 `pkg/common/storage/controller.RedPacketDatabase` 聚合 DAO
- 服务注册名为 `redPacket`，配置文件为 `config/openim-rpc-redpacket.yml`

## 1. 基础约定

### 1.1 Base URL

网关地址由 `openim-api` 部署决定，例如：

```text
http://127.0.0.1:10002
```

红包接口统一挂在：

```text
/redpacket
```

### 1.2 鉴权

当前 `internal/api/router.go` 的 `Whitelist` 未包含 `/redpacket/*`，因此所有红包 HTTP 接口默认都需要登录 token。

请求头：

```http
token: <OpenIM user token>
operationID: <request id>
```

RPC 层不信任请求体中的 `user_id`。当前登录用户统一从 `mcontext.GetOpUserID(ctx)` 读取。

### 1.3 请求字段命名

HTTP 请求建议使用 snake_case。网关使用 `a2r.ParseRequestNotCheck` 解析到 protobuf 请求对象。

示例：

- HTTP: `packet_id`
- protobuf Go 字段: `PacketID`

### 1.4 响应格式

网关使用 `apiresp.GinSuccess` / `apiresp.GinError` 包装响应。不同 OpenIM 版本的外层字段可能略有差异，下面示例重点展示 `data` 内容。

成功示意：

```json
{
  "errCode": 0,
  "errMsg": "",
  "data": {}
}
```

失败示意：

```json
{
  "errCode": 1001,
  "errMsg": "packet_id is required"
}
```

## 2. 接口总览

用户侧接口：

- `POST /redpacket/create_order`
- `POST /redpacket/created_callback`
- `POST /redpacket/detail`
- `POST /redpacket/issue_claim_sign`
- `POST /redpacket/claim_result`
- `POST /redpacket/wallet_bind/challenge`
- `POST /redpacket/wallet_bind/confirm`
- `POST /redpacket/wallet_bind/detail`

管理员接口：

- `POST /redpacket/admin/set_signer`
- `POST /redpacket/admin/set_token`
- `POST /redpacket/admin/set_expiry`
- `POST /redpacket/admin/set_allow_all_tokens`
- `POST /redpacket/admin/set_native_token_enabled`
- `POST /redpacket/admin/parse_tx_events`

## 3. 用户侧接口

### 3.1 创建红包业务单

```text
POST /redpacket/create_order
gRPC: CreateOrder(CreateOrderReq) returns (CreateOrderResp)
```

链上创建红包前调用，服务端创建一条 `PENDING` 业务记录并返回 `biz_id`。

请求示例：

```json
{
  "chain_type": "EVM",
  "chain_id": 1,
  "contract_address": "0xA1f42567559aBA5Ff0aac84cdE1AaF1F9DbB888F",
  "creator_wallet": "0x1111111111111111111111111111111111111111",
  "group_id": "g001",
  "scope_type": "GROUP",
  "receiver_user_id": "",
  "receiver_user_ids": [],
  "packet_type": 1,
  "token": "0x2222222222222222222222222222222222222222",
  "total_amount": "1000000000000000000",
  "total_shares": 10,
  "expiry_at": 0,
  "remark": "happy new year"
}
```

字段说明：

- `chain_type`: 必填，当前支持 `EVM`、`TRON`
- `chain_id`: 可选；EVM client 可用时为空会使用配置的 chainID
- `contract_address`: 可选；EVM/TRON client 可用时为空会使用配置地址
- `creator_wallet`: 必填，发红包钱包地址
- `scope_type`: `GROUP`、`DIRECT`、`PUBLIC`；空值默认 `PUBLIC`
- `group_id`: `scope_type=GROUP` 时必填
- `receiver_user_id` / `receiver_user_ids`: `scope_type=DIRECT` 时至少一个非空
- `packet_type`: `0` 固定红包，`1` 拼手气红包，`2` 转账
- `total_amount`: 链上最小单位十进制字符串
- `total_shares`: 总份数
- `expiry_at`: Unix 秒；`0` 表示使用合约默认过期

成功响应 `data`：

```json
{
  "biz_id": "f8a0f87e-d9cb-4d4a-8350-7bd43ab2e9a4"
}
```

服务端写入：

- collection: `red_packet`
- status: `PENDING`
- creatorUserID: 来自登录上下文，不来自请求体

### 3.2 创建交易回写

```text
POST /redpacket/created_callback
gRPC: CreatedCallback(CreatedCallbackReq) returns (CreatedCallbackResp)
```

链上创建交易确认后调用，用于把 `biz_id` 与链上 `packet_id` / `tx_hash` 绑定。

请求示例：

```json
{
  "biz_id": "f8a0f87e-d9cb-4d4a-8350-7bd43ab2e9a4",
  "tx_hash": "0xabc123...",
  "packet_id": "10001",
  "group_id": "g001",
  "scope_type": "GROUP",
  "receiver_user_id": "",
  "receiver_user_ids": []
}
```

成功响应 `data`：

```json
{}
```

服务端逻辑：

- `biz_id` 与 `tx_hash` 必填
- 如果链客户端可用，会解析交易 receipt 中的 `PacketCreated`
- 解析成功后校验 creator、packetType、token、amount、shares、expiry 是否与业务单一致
- 如果链客户端不可用或解析失败，但请求提供了 `packet_id`，会使用 fallback
- 成功后更新 `red_packet.status=ACTIVE`

### 3.3 查询红包详情

```text
POST /redpacket/detail
gRPC: GetDetail(GetDetailReq) returns (GetDetailResp)
```

请求示例：

```json
{
  "packet_id": "10001"
}
```

成功响应 `data`：

```json
{
  "record": {
    "biz_id": "f8a0f87e-d9cb-4d4a-8350-7bd43ab2e9a4",
    "chain_type": "EVM",
    "packet_id": "10001",
    "chain_id": 1,
    "contract_address": "0xA1f42567559aBA5Ff0aac84cdE1AaF1F9DbB888F",
    "creator_user_id": "u1001",
    "creator_wallet": "0x1111111111111111111111111111111111111111",
    "group_id": "g001",
    "scope_type": "GROUP",
    "receiver_user_id": "",
    "receiver_user_ids": [],
    "packet_type": 1,
    "token": "0x2222222222222222222222222222222222222222",
    "total_amount": "1000000000000000000",
    "total_shares": 10,
    "claimed_amount": "123456789",
    "claimed_shares": 1,
    "expiry_at": 0,
    "tx_hash": "0xabc123...",
    "status": "ACTIVE",
    "created_at": 1777000000,
    "updated_at": 1777000060
  },
  "claims": [
    {
      "packet_id": "10001",
      "user_id": "u2002",
      "claimer_wallet": "0x3333333333333333333333333333333333333333",
      "auth_nonce": "328840239847239847",
      "claim_tx_hash": "0xdef456...",
      "claimed_amount": "123456789",
      "block_number": 1234567,
      "status": "CONFIRMED",
      "created_at": 1777000100,
      "updated_at": 1777000100
    }
  ]
}
```

说明：

- `created_at` / `updated_at` 为 Unix 秒
- `claims` 按 Mongo 查询返回，DAO 层按 `created_at desc` 排序

### 3.4 申请领取签名

```text
POST /redpacket/issue_claim_sign
gRPC: IssueClaimSign(IssueClaimSignReq) returns (IssueClaimSignResp)
```

请求示例：

```json
{
  "packet_id": "10001",
  "claimer": "0x3333333333333333333333333333333333333333",
  "random_seed": "0"
}
```

成功响应 `data`：

```json
{
  "auth_nonce": "328840239847239847",
  "deadline": 1777012345,
  "signature": "0x7b1e...a2",
  "random_seed": "8888812345"
}
```

校验逻辑：

1. 当前用户必须存在：`mcontext.GetOpUserID(ctx) != ""`
2. `packet_id` 与 `claimer` 必填
3. 红包必须存在且 `status=ACTIVE`
4. 未过期、未退款
5. 当前用户与 `claimer` 钱包必须有 `ACTIVE` 绑定
6. 同一用户 / 同一钱包不能重复领取
7. 固定红包和拼手气红包要求 `group_id` 存在
8. 转账红包要求当前用户为 `receiver_user_id`

签名逻辑：

- EVM client 可用时调用 `getSignMessage(packetId, claimer, authNonce, randomSeed, deadline)` 获取 digest
- 使用 `chain.signerPrivateKey` 裸签 digest
- `v` 从 0/1 调整为 27/28
- 如果 signer 私钥未配置，当前代码会返回 placeholder 签名，仅适合本地调试

### 3.5 领取结果回写

```text
POST /redpacket/claim_result
gRPC: ClaimResult(ClaimResultReq) returns (ClaimResultResp)
```

请求示例：

```json
{
  "packet_id": "10001",
  "claimer": "0x3333333333333333333333333333333333333333",
  "tx_hash": "0xdef456..."
}
```

成功响应 `data`：

```json
{}
```

服务端逻辑：

- 先保存一条 `PENDING` claim
- 若能立即解析 `PacketClaimed` 事件，则更新为 `CONFIRMED`
- 成功解析后会累计 `claimed_amount` / `claimed_shares`
- 红包领完时状态变为 `COMPLETED`
- 如果 receipt 暂不可用，保持 `PENDING`，等待 indexer 补偿

### 3.6 发起钱包绑定挑战

```text
POST /redpacket/wallet_bind/challenge
gRPC: IssueWalletBindChallenge(IssueWalletBindChallengeReq)
```

请求示例：

```json
{
  "chain_type": "EVM",
  "chain_id": 1,
  "wallet_address": "0x3333333333333333333333333333333333333333",
  "domain": "redpacket.example.com",
  "uri": "https://redpacket.example.com/wallet-bind"
}
```

成功响应 `data`：

```json
{
  "challenge_id": "1f7d9b0d-7b43-4d84-bb11-65f2ecf7e321",
  "user_id": "u2002",
  "chain_type": "EVM",
  "chain_id": 1,
  "wallet": "0x3333333333333333333333333333333333333333",
  "protocol": "siwe-eip4361",
  "sign_method": "personal_sign",
  "nonce": "7b7d8d48-9db6-4e95-9daa-40e9517a2a85",
  "message": "redpacket.example.com wants you to sign in with your Ethereum account:\n...",
  "issued_at": "2026-04-30T03:00:00Z",
  "expires_at": "2026-04-30T03:10:00Z"
}
```

说明：

- EVM 使用 `siwe-eip4361` + `personal_sign`
- TRON 使用 `tron-signmessagev2` + `signMessageV2`
- challenge 有效期为 10 分钟

### 3.7 确认钱包绑定

```text
POST /redpacket/wallet_bind/confirm
gRPC: ConfirmWalletBind(ConfirmWalletBindReq)
```

请求示例：

```json
{
  "challenge_id": "1f7d9b0d-7b43-4d84-bb11-65f2ecf7e321",
  "signature": "0x8f..."
}
```

成功响应 `data`：

```json
{
  "user_id": "u2002",
  "chain_type": "EVM",
  "chain_id": 1,
  "wallet_address": "0x3333333333333333333333333333333333333333",
  "status": "ACTIVE",
  "verified_at": "2026-04-30T03:01:00Z"
}
```

当前限制：

- EVM 验签已实现
- TRON 验签当前返回 `TRON wallet binding verification is not implemented yet`

### 3.8 查询钱包绑定

```text
POST /redpacket/wallet_bind/detail
gRPC: GetWalletBinding(GetWalletBindingReq)
```

请求示例：

```json
{
  "chain_type": "EVM",
  "wallet_address": "0x3333333333333333333333333333333333333333"
}
```

成功响应 `data`：

```json
{
  "user_id": "u2002",
  "chain_type": "EVM",
  "chain_id": 1,
  "wallet_address": "0x3333333333333333333333333333333333333333",
  "status": "ACTIVE",
  "challenge_id": "1f7d9b0d-7b43-4d84-bb11-65f2ecf7e321",
  "verified_at": "2026-04-30T03:01:00Z"
}
```

## 4. 管理员接口

### 4.1 设置 signer

```text
POST /redpacket/admin/set_signer
gRPC: SetSigner(SetSignerReq)
```

请求：

```json
{
  "signer_address": "0x4444444444444444444444444444444444444444"
}
```

响应：

```json
{
  "message": "signer address updated successfully"
}
```

### 4.2 设置 token 白名单

```text
POST /redpacket/admin/set_token
gRPC: SetToken(SetTokenReq)
```

请求：

```json
{
  "token_address": "0x2222222222222222222222222222222222222222",
  "allowed": true,
  "min_amount": "1000000"
}
```

响应：

```json
{
  "message": "token configuration updated"
}
```

### 4.3 设置默认过期时间

```text
POST /redpacket/admin/set_expiry
gRPC: SetExpiry(SetExpiryReq)
```

请求：

```json
{
  "expiry_seconds": 86400
}
```

### 4.4 设置是否允许所有 token

```text
POST /redpacket/admin/set_allow_all_tokens
gRPC: SetAllowAllTokens(SetAllowAllTokensReq)
```

请求：

```json
{
  "allow_all": false
}
```

### 4.5 设置原生币开关

```text
POST /redpacket/admin/set_native_token_enabled
gRPC: SetNativeTokenEnabled(SetNativeTokenEnabledReq)
```

请求：

```json
{
  "enabled": true
}
```

### 4.6 解析交易事件

```text
POST /redpacket/admin/parse_tx_events
gRPC: ParseTxEvents(ParseTxEventsReq)
```

请求：

```json
{
  "tx_hash": "0xabc123...",
  "chain": "eth"
}
```

EVM 响应：

```json
{
  "chain": "eth",
  "tx_hash": "0xabc123...",
  "events": [
    {
      "name": "PacketCreated",
      "data": {
        "packetId": "10001",
        "creator": "0x1111111111111111111111111111111111111111"
      }
    }
  ]
}
```

TRON 当前响应：

```json
{
  "chain": "tron",
  "tx_hash": "7d9e...txid",
  "note": "TRON event parsing not fully implemented in this version"
}
```

### 4.7 管理接口当前行为边界

- EVM admin 接口当前为 mock，仅记录日志并返回 message，不发链上交易。
- TRON admin 接口会调用 `SendAdminTransaction(...)` 尝试发链上交易。
- 管理接口目前没有单独管理员校验，默认只依赖 API 网关 token。生产建议补管理员鉴权与审计。

## 5. 业务状态

红包状态：

- `PENDING`: 已创建业务单，尚未确认链上创建
- `ACTIVE`: 链上创建已确认，可领取
- `COMPLETED`: 已领取完成
- `REFUNDED`: 已退款

领取状态：

- `PENDING`: 已提交领取 txHash，receipt 尚未解析或未确认
- `CONFIRMED`: 已解析 `PacketClaimed`
- `FAILED`: 预留失败状态，当前逻辑仅用于重复领取判断时放行失败记录

钱包绑定 challenge 状态：

- `PENDING`
- `VERIFIED`
- `FAILED`
- `EXPIRED`

钱包绑定状态：

- `ACTIVE`

## 6. 常见错误

- `op user id is empty`: 缺少 token 或 token 未正确注入上下文
- `unsupported chain_type`: `chain_type` 不是 `EVM` 或 `TRON`
- `packet_id is required`: 缺少红包链上 ID
- `wallet is not bound to user`: 当前用户未绑定该领取钱包
- `user already claimed`: 当前用户已领取
- `already claimed`: 当前钱包已领取
- `packet is not active`: 红包尚未激活或已经完成/退款
- `packet is expired`: 红包已过期
- `TRON wallet binding verification is not implemented yet`: 当前未实现 TRON 绑定验签

## 7. 前端推荐调用顺序

创建红包：

1. `POST /redpacket/create_order`
2. 钱包发起 `createFixedPacket/createRandomPacket/createTransfer`
3. 从 `PacketCreated` 解析 `packetId`
4. `POST /redpacket/created_callback`
5. `POST /redpacket/detail` 刷新状态

绑定钱包：

1. `POST /redpacket/wallet_bind/challenge`
2. 钱包按 `sign_method` 签名 `message`
3. `POST /redpacket/wallet_bind/confirm`
4. `POST /redpacket/wallet_bind/detail`

领取红包：

1. `POST /redpacket/detail`
2. `POST /redpacket/issue_claim_sign`
3. 钱包调用链上 `claim(packetId, authNonce, randomSeed, deadline, signature)`
4. 可选：`POST /redpacket/claim_result`
5. `POST /redpacket/detail` 刷新状态

## 8. 存储与索引

Mongo collections：

- `red_packet`
- `red_packet_claim`
- `red_packet_claim_auth`
- `red_packet_refund`
- `wallet_binding_challenge`
- `wallet_binding`

主要索引：

- `red_packet.biz_id` 唯一
- `red_packet.packet_id`
- `red_packet.group_id`
- `red_packet_claim.claim_tx_hash` 唯一
- `red_packet_claim.packet_id + user_id`
- `red_packet_claim.packet_id + claimer_wallet`
- `red_packet_claim_auth.auth_nonce` 唯一
- `wallet_binding_challenge.challenge_id` 唯一
- `wallet_binding.user_id + chain_type + wallet_address` 唯一

## 9. 配置文件

`config/openim-rpc-redpacket.yml`：

```yaml
rpc:
  registerIP: ""
  listenIP: 0.0.0.0
  autoSetPorts: false
  ports: [10560]

prometheus:
  enable: false
  ports: [12560]

chain:
  rpcURL: ""
  contractAddress: ""
  chainID: 0
  signerPrivateKey: ""
  configAdminPrivateKey: ""

tron:
  fullNodeURL: ""
  contractBase58: ""
  ownerBase58: ""
  privateKeyHex: ""
  feeLimit: 100000000

indexer:
  pollInterval: 5
```

`chain.rpcURL` 为空时 EVM client 初始化会失败并降级；`tron.fullNodeURL` 为空时 TRON client 不启用。服务会继续启动。
