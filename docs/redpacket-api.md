# 红包 API 接口文档

**Base URL:** `/redpacket`  
**协议:** HTTP POST，`Content-Type: application/json`  
**认证:** 请求头携带 `token: <JWT令牌>`（标注需要登录的接口）

> **统一响应结构**
>
> ```json
> {
>   "errCode": 0,
>   "errMsg": "ok",
>   "errDlt": "",
>   "data": { }
> }
> ```
>
> `errCode` 为 `0` 表示成功，非 0 表示错误。

---

## 1. 创建红包订单

**POST** `/redpacket/create_order`  
需要登录。创建一条待上链的红包订单，返回业务 ID（bizID）供后续链上交易关联。

### 请求体

```json
{
  "chainType": "EVM",
  "chainID": 1,
  "contractAddress": "0xAbCd...",
  "creatorWallet": "0x1234...",
  "groupID": "group_001",
  "scopeType": "GROUP",
  "receiverUserID": "",
  "receiverUserIDs": [],
  "packetType": 0,
  "token": "0x0000000000000000000000000000000000000000",
  "totalAmount": "1000000000000000000",
  "totalShares": 10,
  "expiryAt": 1800000000,
  "remark": "新年快乐"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `chainType` | string | **必填** | 链类型，仅支持 `"EVM"` 或 `"TRON"` |
| `chainID` | int64 | 可选 | 链 ID；为 0 时从链客户端自动获取 |
| `contractAddress` | string | 可选 | 合约地址；为空时从链客户端自动获取 |
| `creatorWallet` | string | **必填** | 创建者钱包地址 |
| `scopeType` | string | 可选 | 范围类型：`GROUP`（群组）/`DIRECT`（定向）/`PUBLIC`（公开），默认 `PUBLIC` |
| `groupID` | string | **GROUP 时必填** | 群组 ID（`scopeType=GROUP` 时必须提供） |
| `receiverUserID` | string | **transfer 时必填** | 接收者用户 ID（`packetType=2` 且 `scopeType=DIRECT` 时必须提供） |
| `receiverUserIDs` | []string | **DIRECT + fixed/random 时必填** | 多接收者用户 ID 列表（`scopeType=DIRECT` 且 `packetType=0/1` 时使用） |
| `packetType` | int32 | **必填** | 红包类型：`0`=均分红包，`1`=随机红包，`2`=转账红包 |
| `token` | string | 可选 | ERC20 代币合约地址；为空表示原生代币 |
| `totalAmount` | string | **必填** | 总金额（最小单位整数字符串，如 wei），必须为正整数 |
| `totalShares` | int32 | **必填（packetType 0/1）** | 红包份数；均分/随机红包 >0 且 ≤10000；转账红包固定为 1 |
| `expiryAt` | int64 | 可选 | 过期时间（Unix 时间戳秒）；0 表示不过期；必须为将来时间 |
| `remark` | string | 可选 | 备注 |

**packetType 规则说明：**

- `0`（均分）：`scopeType` 必须为 `GROUP`，`totalAmount` 必须被 `totalShares` 整除
- `1`（随机）：`scopeType` 必须为 `GROUP`，`totalAmount` ≥ `totalShares`
- `2`（转账）：`scopeType` 必须为 `DIRECT`，`totalShares` 必须为 1，`receiverUserID` 必须提供，不能转给自己

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {
    "bizID": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `bizID` | string | 业务订单 ID，用于链上交易完成后提交回调 |

---

## 2. 红包创建回调（链上确认）

**POST** `/redpacket/created_callback`  
需要登录。链上交易广播后，由创建者提交交易哈希以激活红包。仅创建者可调用。

### 请求体

```json
{
  "bizID": "550e8400-e29b-41d4-a716-446655440000",
  "txHash": "0xabc123...",
  "packetID": "",
  "groupID": "",
  "scopeType": "",
  "receiverUserID": "",
  "receiverUserIDs": []
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `bizID` | string | **必填** | 创建订单时返回的业务 ID |
| `txHash` | string | **必填** | 链上交易哈希 |
| `packetID` | string | 条件必填 | 链上红包 ID；链客户端离线时必须手动提供 |
| `groupID` | string | 可选 | 覆盖订单的群组 ID（不填则继承订单值） |
| `scopeType` | string | 可选 | 覆盖订单的范围类型 |
| `receiverUserID` | string | 可选 | 覆盖单一接收者 |
| `receiverUserIDs` | []string | 可选 | 覆盖多接收者列表 |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {}
}
```

---

## 3. 查询红包详情

**POST** `/redpacket/detail`

### 请求体

```json
{
  "packetID": "12345"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `packetID` | string | **必填** | 链上红包 ID |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {
    "record": {
      "bizID": "550e8400-...",
      "chainType": "EVM",
      "packetID": "12345",
      "chainID": 1,
      "contractAddress": "0xAbCd...",
      "creatorUserID": "user_001",
      "creatorWallet": "0x1234...",
      "groupID": "group_001",
      "scopeType": "GROUP",
      "receiverUserID": "",
      "receiverUserIDs": [],
      "packetType": 0,
      "token": "0x0000000000000000000000000000000000000000",
      "totalAmount": "1000000000000000000",
      "totalShares": 10,
      "claimedAmount": "300000000000000000",
      "claimedShares": 3,
      "expiryAt": 1800000000,
      "txHash": "0xabc123...",
      "status": "ACTIVE",
      "createdAt": 1715500000,
      "updatedAt": 1715500100
    },
    "claims": [
      {
        "packetID": "12345",
        "userID": "user_002",
        "claimerWallet": "0x5678...",
        "authNonce": "1715500050000000000",
        "claimTxHash": "0xdef456...",
        "claimedAmount": "100000000000000000",
        "blockNumber": 19000000,
        "status": "CONFIRMED",
        "createdAt": 1715500050,
        "updatedAt": 1715500060
      }
    ]
  }
}
```

**record 字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `status` | string | 红包状态：`PENDING` / `ACTIVE` / `COMPLETED` / `REFUNDED` / `EXPIRED` |
| `totalAmount` | string | 总金额（最小单位） |
| `claimedAmount` | string | 已领金额 |
| `totalShares` / `claimedShares` | int32 | 总份数 / 已领份数 |

---

## 4. 申请领取签名

**POST** `/redpacket/issue_claim_sign`  
需要登录。领取红包前，先获取服务端签名用于链上验证。会校验领取资格（是否群成员/好友、是否已领取、红包状态等）。

### 请求体

```json
{
  "packetID": "12345",
  "claimer": "0x5678...",
  "randomSeed": ""
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `packetID` | string | **必填** | 链上红包 ID |
| `claimer` | string | **必填** | 领取者钱包地址 |
| `randomSeed` | string | 可选 | 随机种子（十进制整数字符串）；为空或 `"0"` 时服务端自动生成 |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {
    "authNonce": "1715500050000000000",
    "deadline": 1715500350,
    "signature": "0xaabbcc...",
    "randomSeed": "8765309000000"
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `authNonce` | string | 认证随机数（传给合约） |
| `deadline` | int64 | 签名过期时间（Unix 时间戳，约 5 分钟后） |
| `signature` | string | 服务端签名（0x 前缀十六进制，65 字节） |
| `randomSeed` | string | 使用的随机种子 |

---

## 5. 提交领取结果

**POST** `/redpacket/claim_result`  
需要登录。链上领取交易广播后提交，服务端解析链上事件并更新领取记录。

### 请求体

```json
{
  "packetID": "12345",
  "claimer": "0x5678...",
  "txHash": "0xdef456..."
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `packetID` | string | **必填** | 链上红包 ID |
| `claimer` | string | **必填** | 领取者钱包地址 |
| `txHash` | string | **必填** | 领取交易哈希 |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {}
}
```

---

## 6. 申请退款

**POST** `/redpacket/request_refund`  
需要登录。仅创建者可调用，红包到期后提交链上退款交易。

### 请求体

```json
{
  "packetID": "12345"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `packetID` | string | **必填** | 链上红包 ID（红包必须已到期） |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {
    "txHash": "0xghi789...",
    "status": "PENDING"
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `txHash` | string | 退款交易哈希 |
| `status` | string | `PENDING`（退款交易已提交）或 `REFUNDED`（已退款） |

---

## 7. 查询退款记录

**POST** `/redpacket/get_refund`

### 请求体

```json
{
  "packetID": "12345"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `packetID` | string | **必填** | 链上红包 ID |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {
    "packetID": "12345",
    "refundTo": "0x1234...",
    "txHash": "0xghi789...",
    "amount": "700000000000000000",
    "createdAt": 1715600000
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `refundTo` | string | 退款目标钱包地址 |
| `amount` | string | 退款金额（最小单位） |
| `createdAt` | int64 | 退款记录创建时间（Unix 时间戳） |

---

## 8. 发起钱包绑定挑战

**POST** `/redpacket/wallet_bind/challenge`  
需要登录。生成一条待用户签名的消息，用于将钱包地址绑定到当前用户账户（EVM 使用 EIP-4361 SIWE，TRON 使用 signMessageV2）。

### 请求体

```json
{
  "chainType": "EVM",
  "chainID": 1,
  "walletAddress": "0x5678...",
  "domain": "myapp.example.com",
  "uri": "https://myapp.example.com/wallet-bind"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `chainType` | string | **必填** | 链类型，`"EVM"` 或 `"TRON"` |
| `walletAddress` | string | **必填** | 待绑定的钱包地址 |
| `chainID` | int64 | 可选 | 链 ID（EVM 时建议提供） |
| `domain` | string | 可选 | 应用域名，默认 `"redpacket"` |
| `uri` | string | 可选 | 应用 URI，默认 `"https://redpacket.local/wallet-bind"` |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {
    "challengeID": "aaaa-bbbb-cccc-dddd",
    "userID": "user_001",
    "chainType": "EVM",
    "chainID": 1,
    "wallet": "0x5678...",
    "protocol": "siwe-eip4361",
    "signMethod": "personal_sign",
    "nonce": "xxxx-yyyy-zzzz",
    "message": "myapp.example.com wants you to sign in with your Ethereum account:\n0x5678...\n\nBind wallet ...",
    "issuedAt": "2026-05-12T08:39:00Z",
    "expiresAt": "2026-05-12T08:49:00Z"
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `challengeID` | string | 挑战 ID，确认绑定时使用 |
| `message` | string | 待签名的完整消息文本 |
| `protocol` | string | EVM: `siwe-eip4361`；TRON: `tron-signmessagev2` |
| `signMethod` | string | EVM: `personal_sign`；TRON: `signMessageV2` |
| `expiresAt` | string | 挑战过期时间（RFC3339，10 分钟有效期） |

---

## 9. 确认钱包绑定

**POST** `/redpacket/wallet_bind/confirm`  
提交用户对挑战消息的签名，服务端验证后完成钱包绑定。挑战有效期 10 分钟。

### 请求体

```json
{
  "challengeID": "aaaa-bbbb-cccc-dddd",
  "signature": "0xaabbccdd..."
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `challengeID` | string | **必填** | 挑战 ID（来自 `/wallet_bind/challenge`） |
| `signature` | string | **必填** | 钱包签名（0x 前缀十六进制，65 字节） |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {
    "userID": "user_001",
    "chainType": "EVM",
    "chainID": 1,
    "walletAddress": "0x5678...",
    "status": "ACTIVE",
    "verifiedAt": "2026-05-12T08:42:00Z"
  }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `status` | string | 绑定状态，成功时为 `ACTIVE` |
| `verifiedAt` | string | 绑定验证时间（RFC3339） |

---

## 10. 查询钱包绑定信息

**POST** `/redpacket/wallet_bind/detail`  
需要登录。查询当前登录用户在指定链上的活跃钱包绑定。

### 请求体

```json
{
  "chainType": "EVM",
  "walletAddress": "0x5678..."
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `chainType` | string | **必填** | 链类型，`"EVM"` 或 `"TRON"` |
| `walletAddress` | string | 可选 | 钱包地址过滤条件 |

### 响应体

```json
{
  "errCode": 0,
  "errMsg": "ok",
  "data": {
    "userID": "user_001",
    "chainType": "EVM",
    "chainID": 1,
    "walletAddress": "0x5678...",
    "status": "ACTIVE",
    "challengeID": "aaaa-bbbb-cccc-dddd",
    "verifiedAt": "2026-05-12T08:42:00Z"
  }
}
```

---

## 附录：公共枚举值

| 枚举 | 可选值 | 说明 |
|------|--------|------|
| `chainType` | `EVM` / `TRON` | 区块链类型 |
| `scopeType` | `GROUP` / `DIRECT` / `PUBLIC` | 红包范围 |
| `packetType` | `0` / `1` / `2` | 均分 / 随机 / 转账 |
| 红包 `status` | `PENDING` / `ACTIVE` / `COMPLETED` / `REFUNDED` / `EXPIRED` | 红包状态 |
| 领取 `status` | `PENDING` / `CONFIRMED` / `FAILED` | 领取记录状态 |
| 钱包绑定 `status` | `ACTIVE` | 绑定状态（查询时只返回活跃绑定） |

---

## 实现参考

- 路由：`internal/api/router.go`（`/redpacket` 分组）
- HTTP 处理：`internal/api/redpacket.go`
- Proto 定义：`protocol/redpacket/redpacket.proto`
- RPC 实现：`internal/rpc/redpacket/service.go`、`internal/rpc/redpacket/wallet.go`
