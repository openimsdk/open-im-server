# RedPacket 后端接口说明

本文档基于当前后端实现整理，覆盖用户接口与管理员接口，并提供请求/响应示例。

## 基础信息

- Base URL（本地默认）：`http://127.0.0.1:8080`
- 统一响应格式：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

- 错误响应格式：

```json
{
  "code": 400,
  "message": "invalid request body: ..."
}
```

## 健康检查

### GET `/health`

用于服务存活探测。

#### 响应示例

```json
{
  "status": "ok"
}
```

---

## 用户侧接口

## 1) 创建业务订单

### POST `/api/redpacket/create-order`

链上发交易前先创建业务订单，返回 `biz_id`。

#### 请求体

```json
{
  "creator_user_id": "u1001",
  "creator_wallet": "0x1111111111111111111111111111111111111111",
  "packet_type": 1,
  "token": "0x2222222222222222222222222222222222222222",
  "total_amount": "1000000000000000000",
  "total_shares": 10,
  "expiry_at": 0
}
```

#### 字段说明

- `packet_type`: `0` 固定红包，`1` 拼手气红包，`2` 转账红包
- `total_amount`: 链上最小单位的十进制字符串
- `expiry_at`: Unix 秒时间戳，`0` 表示使用合约默认过期时间

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "biz_id": "f8a0f87e-d9cb-4d4a-8350-7bd43ab2e9a4"
  }
}
```

#### 失败响应示例

```json
{
  "code": 400,
  "message": "invalid token address"
}
```

---

## 2) 创建结果回写

### POST `/api/redpacket/created-callback`

前端在链上创建交易确认后，回写 `tx_hash` 和 `packet_id`。

#### 请求体

```json
{
  "biz_id": "f8a0f87e-d9cb-4d4a-8350-7bd43ab2e9a4",
  "tx_hash": "0xabc123...",
  "packet_id": "10001"
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "ok": true
  }
}
```

#### 失败响应示例

```json
{
  "code": 400,
  "message": "biz_id is required"
}
```

---

## 3) 红包详情

### GET `/api/redpacket/detail?packet_id={packetId}`

查询红包业务记录与领取记录。

#### 请求示例

```bash
curl "http://127.0.0.1:8080/api/redpacket/detail?packet_id=10001"
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "biz_record": {
      "id": 1,
      "biz_id": "f8a0f87e-d9cb-4d4a-8350-7bd43ab2e9a4",
      "packet_id": "10001",
      "chain_id": 1,
      "contract_address": "0xA1f42567559aBA5Ff0aac84cdE1AaF1F9DbB888F",
      "creator_user_id": "u1001",
      "creator_wallet": "0x1111111111111111111111111111111111111111",
      "packet_type": 1,
      "token": "0x2222222222222222222222222222222222222222",
      "total_amount": "1000000000000000000",
      "total_shares": 10,
      "expiry_at": 0,
      "tx_hash": "0xabc123...",
      "status": "ACTIVE",
      "created_at": "2026-04-24T07:00:00Z",
      "updated_at": "2026-04-24T07:01:00Z"
    },
    "claims": [
      {
        "id": 10,
        "packet_id": "10001",
        "claimer_wallet": "0x3333333333333333333333333333333333333333",
        "auth_nonce": "328840239847239847",
        "claim_tx_hash": "0xdef456...",
        "claimed_amount": "123456789",
        "block_number": 1234567,
        "status": "CONFIRMED",
        "created_at": "2026-04-24T07:10:00Z",
        "updated_at": "2026-04-24T07:10:00Z"
      }
    ]
  }
}
```

#### 失败响应示例

```json
{
  "code": 404,
  "message": "packet not found: 10001"
}
```

---

## 4) 申请领取签名

### POST `/api/redpacket/claim-sign`

先做业务鉴权，再发放 `claim(...)` 所需签名参数。

#### 鉴权说明

- 该接口不再信任请求体中的 `user_id`
- 当前领取用户从 RPC / 网关注入的登录上下文中获取
- 服务端要求请求上下文里存在 `opUserID`
- 如果缺少登录上下文，接口会直接拒绝

#### 请求头

- `token`: 用户登录 token

> 约定：上游网关或鉴权中间件需要先解析 token，并把当前登录用户写入请求上下文中的 `opUserID`。

#### 请求体

```json
{
  "packet_id": "10001",
  "claimer": "0x3333333333333333333333333333333333333333",
  "random_seed": "0"
}
```

> `random_seed` 可选；传 `0` 或空时后端自动生成。

#### 字段说明

- `packet_id`: 红包链上 ID
- `claimer`: 本次真正发起链上 `claim(...)` 的钱包地址
- `random_seed`: 可选随机种子；空或 `0` 时后端自动生成

#### 服务端处理逻辑

1. 从请求上下文提取当前登录用户 ID
2. 校验红包是否存在、是否过期、是否仍可领取
3. 校验当前登录用户与 `claimer` 钱包地址的绑定关系
4. 校验当前用户在该红包下是否已领取
5. 校验当前钱包在该红包下是否已领取
6. 按红包类型校验群资格 / 指定接收人资格
7. 生成 `auth_nonce`、`deadline`、`random_seed`
8. 调合约 `getSignMessage(packetId, claimer, authNonce, randomSeed, deadline)` 获取摘要
9. 使用后端 `signer` 私钥对摘要裸签名
10. 落库 `red_packet_claim_auth`
11. 返回前端发链所需参数

#### 成功后前端下一步

前端拿到响应后，直接调用链上：

```text
claim(packetId, authNonce, randomSeed, deadline, signature)
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "auth_nonce": "328840239847239847",
    "deadline": 1777012345,
    "signature": "0x7b1e...a2",
    "random_seed": "8888812345"
  }
}
```

#### 常见失败响应

无资格领取：

```json
{
  "code": 403,
  "message": "already claimed"
}
```

同一用户已领取：

```json
{
  "code": 403,
  "message": "user already claimed"
}
```

钱包未绑定：

```json
{
  "code": 403,
  "message": "wallet is not bound to user"
}
```

缺少登录上下文：

```json
{
  "code": 403,
  "message": "op user id missing in context"
}
```

签名服务异常：

```json
{
  "code": 500,
  "message": "failed to issue claim signature: getSignMessage: ..."
}
```

---

## 5) 领取结果回写（可选）

### POST `/api/redpacket/claim-result`

前端在领取交易提交后可调用该接口预写记录。最终状态仍以链监听（indexer）为准。

#### 鉴权说明

- 该接口不再接收可信 `user_id`
- 当前用户从 RPC / 网关注入的登录上下文中获取
- 服务端要求请求上下文里存在 `opUserID`

#### 请求头

- `token`: 用户登录 token

#### 请求体

```json
{
  "packet_id": "10001",
  "claimer": "0x3333333333333333333333333333333333333333",
  "tx_hash": "0xdef456..."
}
```

#### 字段说明

- `packet_id`: 红包链上 ID
- `claimer`: 发起链上领取的钱包地址
- `tx_hash`: 领取交易哈希

#### 服务端处理逻辑

1. 从请求上下文提取当前登录用户 ID
2. 先落一条 `PENDING` 领取记录
3. 如果当前节点能立即解析该交易 receipt，则补全：
   - `auth_nonce`
   - `claimed_amount`
   - `block_number`
   - `status=CONFIRMED`
4. 如果当前节点暂时拿不到 receipt，则保持 `PENDING`
5. 最终仍以链监听器写入结果为准

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "ok": true
  }
}
```

#### 失败响应示例

```json
{
  "code": 403,
  "message": "op user id missing in context"
}
```

---

## 6) 钱包绑定挑战

### POST `/api/redpacket/wallet-bind/challenge`

生成钱包绑定挑战消息，前端拿到消息后调用钱包签名。

#### 鉴权说明

- 该接口不再信任请求体中的 `user_id`
- 当前用户从 RPC / 网关注入的登录上下文中获取

#### 请求头

- `token`: 用户登录 token

#### 请求体

```json
{
  "chain_type": "EVM",
  "chain_id": 1,
  "wallet_address": "0x3333333333333333333333333333333333333333",
  "domain": "redpacket.example.com",
  "uri": "https://redpacket.example.com/wallet-bind"
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "challenge_id": "1f7d9b0d-7b43-4d84-bb11-65f2ecf7e321",
    "user_id": "u2002",
    "chain_type": "EVM",
    "chain_id": 1,
    "wallet": "0x3333333333333333333333333333333333333333",
    "protocol": "siwe-eip4361",
    "sign_method": "personal_sign",
    "nonce": "7b7d8d48-9db6-4e95-9daa-40e9517a2a85",
    "message": "redpacket.example.com wants you to sign in with your Ethereum account:\n0x3333333333333333333333333333333333333333\n\nBind wallet 0x3333333333333333333333333333333333333333 to user u2002.\nURI: https://redpacket.example.com/wallet-bind\nVersion: 1\nChain ID: 1\nNonce: 7b7d8d48-9db6-4e95-9daa-40e9517a2a85\nIssued At: 2026-04-30T03:00:00Z\nExpiration Time: 2026-04-30T03:10:00Z\nRequest ID: 1f7d9b0d-7b43-4d84-bb11-65f2ecf7e321",
    "issued_at": "2026-04-30T03:00:00Z",
    "expires_at": "2026-04-30T03:10:00Z"
  }
}
```

#### 前端下一步

前端收到响应后：

1. 使用 `sign_method` 指定的钱包方法对 `message` 进行签名
2. 把 `challenge_id + signature` 提交给 `/api/redpacket/wallet-bind/confirm`

---

## 7) 钱包绑定确认

### POST `/api/redpacket/wallet-bind/confirm`

提交钱包签名，服务端验签成功后建立钱包绑定关系。

#### 请求体

```json
{
  "challenge_id": "1f7d9b0d-7b43-4d84-bb11-65f2ecf7e321",
  "signature": "0x8f..."
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "user_id": "u2002",
    "chain_type": "EVM",
    "chain_id": 1,
    "wallet_address": "0x3333333333333333333333333333333333333333",
    "status": "ACTIVE",
    "verified_at": "2026-04-30T03:01:00Z"
  }
}
```

---

## 8) 查询钱包绑定

### GET `/api/redpacket/wallet-bind/detail?chain_type={chainType}&wallet_address={walletAddress}`

查询当前登录用户与指定钱包地址的绑定详情。

#### 鉴权说明

- `user_id` 从登录上下文中获取，不需要也不应该由前端传入

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "user_id": "u2002",
    "chain_type": "EVM",
    "chain_id": 1,
    "wallet_address": "0x3333333333333333333333333333333333333333",
    "status": "ACTIVE",
    "challenge_id": "1f7d9b0d-7b43-4d84-bb11-65f2ecf7e321",
    "verified_at": "2026-04-30T03:01:00Z"
  }
}
```

---

## 管理员接口（建议加鉴权）

以下接口属于管理员写链操作，依赖后端配置的 `config_admin_private_key`。

## 6) 设置 signer

### POST `/admin/redpacket/set-signer`

#### 请求体

```json
{
  "new_signer": "0x4444444444444444444444444444444444444444"
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "tx_hash": "0xaaa111..."
  }
}
```

---

## 7) 设置 token 白名单与最小份额

### POST `/admin/redpacket/set-token`

#### 请求体

```json
{
  "token": "0x2222222222222222222222222222222222222222",
  "allowed": true,
  "min_share_amount": "1000000"
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "tx_hash": "0xbbb222..."
  }
}
```

---

## 8) 设置默认过期时间

### POST `/admin/redpacket/set-expiry`

#### 请求体

```json
{
  "duration": "86400"
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "tx_hash": "0xccc333..."
  }
}
```

---

## 9) 设置是否允许所有 token

### POST `/admin/redpacket/set-allow-all-tokens`

#### 请求体

```json
{
  "allow": false
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "tx_hash": "0xddd444..."
  }
}
```

---

## 10) 设置原生币开关

### POST `/admin/redpacket/set-native-token`

#### 请求体

```json
{
  "enabled": true
}
```

#### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "tx_hash": "0xeee555..."
  }
}
```

---

## 11) 按交易哈希解析事件

### POST `/admin/redpacket/parse-tx-events`

支持 ETH/TRON 事件解码。

#### 请求体（ETH）

```json
{
  "chain": "eth",
  "tx_hash": "0xabc123..."
}
```

#### 请求体（TRON）

```json
{
  "chain": "tron",
  "tx_hash": "7d9e...txid"
}
```

#### 成功响应（示例）

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "name": "PacketCreated",
      "data": {
        "packetId": "10001",
        "creator": "0x1111111111111111111111111111111111111111",
        "packetType": 1
      }
    }
  ]
}
```

#### 失败响应示例

TRON 未配置：

```json
{
  "code": 503,
  "message": "TRON client is not configured"
}
```

参数非法：

```json
{
  "code": 400,
  "message": "chain must be \"eth\" or \"tron\""
}
```

---

## 典型调用顺序（前端）

1. `POST /api/redpacket/create-order`
2. 钱包发链上创建交易
3. 解析 `PacketCreated.packetId`
4. `POST /api/redpacket/created-callback`
5. 用户领取前：`POST /api/redpacket/claim-sign`
6. 钱包调用合约 `claim(...)`
7. 可选：`POST /api/redpacket/claim-result`
8. 详情页查询：`GET /api/redpacket/detail?packet_id=...`
