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

#### 请求体

```json
{
  "packet_id": "10001",
  "claimer": "0x3333333333333333333333333333333333333333",
  "user_id": "u2002",
  "random_seed": "0"
}
```

> `random_seed` 可选；传 `0` 或空时后端自动生成。

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

#### 请求体

```json
{
  "packet_id": "10001",
  "claimer_wallet": "0x3333333333333333333333333333333333333333",
  "tx_hash": "0xdef456...",
  "auth_nonce": "328840239847239847"
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
  "message": "packet_id and tx_hash are required"
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

