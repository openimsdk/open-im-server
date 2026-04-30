# RedPacket 前端对接文档

本文档面向前端 / 网关 / App 对接方，说明红包领取和钱包绑定的真实接入方式，重点覆盖：

- 如何把当前登录用户传递给红包服务
- 如何绑定钱包
- 如何申请领取签名
- 前端何时发链、何时回写后端

## 1. 总体原则

红包服务已经切换为 RPC 上下文取当前用户 ID：

- 前端不再把 `user_id` 当作可信业务参数传给红包服务
- 红包服务从请求上下文里的 `opUserID` 获取当前登录用户
- 上下文通常由网关或鉴权中间件根据 `token` 解析后注入

这意味着对接时必须满足一个前提：

- 请求进入红包服务前，网关已经完成 token 解析
- 并且把当前登录用户写入上下文中的 `opUserID`

如果没有这一层，红包服务会返回：

```json
{
  "code": 403,
  "message": "op user id missing in context"
}
```

## 2. 钱包绑定流程

### 2.1 流程图

```text
前端 -> 红包服务: POST /api/redpacket/wallet-bind/challenge
红包服务 -> 前端: challenge_id + message + sign_method
前端 -> 钱包: 对 message 签名
前端 -> 红包服务: POST /api/redpacket/wallet-bind/confirm
红包服务 -> 前端: 绑定成功
```

### 2.2 发起挑战

请求：

```http
POST /api/redpacket/wallet-bind/challenge
token: <user token>
Content-Type: application/json
```

```json
{
  "chain_type": "EVM",
  "chain_id": 1,
  "wallet_address": "0x3333333333333333333333333333333333333333",
  "domain": "redpacket.example.com",
  "uri": "https://redpacket.example.com/wallet-bind"
}
```

返回里最关键的是：

- `challenge_id`
- `message`
- `sign_method`

前端要做的是：

- 按 `sign_method` 调钱包签名
- 当前 EVM 实现使用的是 `personal_sign`

### 2.3 确认绑定

请求：

```http
POST /api/redpacket/wallet-bind/confirm
token: <user token>
Content-Type: application/json
```

```json
{
  "challenge_id": "1f7d9b0d-7b43-4d84-bb11-65f2ecf7e321",
  "signature": "0x8f..."
}
```

成功后代表：

- 当前登录用户
- 当前链类型
- 当前钱包地址

已经在后端建立了有效绑定关系。

## 3. 领取签名流程

### 3.1 流程图

```text
前端 -> 红包服务: POST /api/redpacket/claim-sign
红包服务 -> 红包服务: 校验当前用户、钱包绑定、领取资格
红包服务 -> 合约: getSignMessage(packetId, claimer, authNonce, randomSeed, deadline)
红包服务 -> 前端: auth_nonce + random_seed + deadline + signature
前端 -> 钱包/链上: claim(packetId, authNonce, randomSeed, deadline, signature)
前端 -> 红包服务: POST /api/redpacket/claim-result (可选)
链监听器 -> 红包服务: 最终确认领取结果
```

### 3.2 申请领取签名

请求：

```http
POST /api/redpacket/claim-sign
token: <user token>
Content-Type: application/json
```

```json
{
  "packet_id": "10001",
  "claimer": "0x3333333333333333333333333333333333333333",
  "random_seed": "0"
}
```

说明：

- `claimer` 必须是这次真正发链的地址
- `random_seed` 可省略或传 `0`
- 不需要传 `user_id`

后端会自动完成这些校验：

1. 当前登录用户存在
2. 红包存在且仍可领取
3. 当前登录用户与 `claimer` 已绑定
4. 当前用户在该红包下未领取过
5. 当前钱包在该红包下未领取过
6. 群红包 / 转账红包的附加业务限制通过

成功响应：

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

### 3.3 前端拿到响应后要做什么

前端必须原样把这些参数传给链上：

```text
claim(packetId, authNonce, randomSeed, deadline, signature)
```

对应关系：

- `packetId` -> 前端当前红包 ID
- `authNonce` -> 响应里的 `auth_nonce`
- `randomSeed` -> 响应里的 `random_seed`
- `deadline` -> 响应里的 `deadline`
- `signature` -> 响应里的 `signature`

注意：

- 不要自己改 `auth_nonce`
- 不要重新算摘要
- 不要对摘要再次做 `signMessage`
- 后端返回的 `signature` 已经是最终可上链签名

## 4. 领取结果回写

`claim-result` 是可选的，主要作用是让业务侧尽快看到一条 `PENDING` 领取记录。

请求：

```http
POST /api/redpacket/claim-result
token: <user token>
Content-Type: application/json
```

```json
{
  "packet_id": "10001",
  "claimer": "0x3333333333333333333333333333333333333333",
  "tx_hash": "0xdef456..."
}
```

说明：

- 不需要传 `user_id`
- 当前登录用户仍然从上下文中取
- 如果后端当前能立刻解析 receipt，会把记录补成 `CONFIRMED`
- 如果不能，会先记成 `PENDING`
- 最终仍以链监听器为准

## 5. 前端推荐调用顺序

### 5.1 首次使用钱包领取

1. 用户登录业务系统
2. 前端请求 `/wallet-bind/challenge`
3. 钱包对 `message` 签名
4. 前端请求 `/wallet-bind/confirm`
5. 绑定成功后再进入领取流程

### 5.2 正常领取

1. 前端拿到红包 `packet_id`
2. 用户连接钱包，得到本次 `claimer` 地址
3. 前端请求 `/claim-sign`
4. 拿到 `auth_nonce + random_seed + deadline + signature`
5. 前端调用链上 `claim(...)`
6. 前端可选请求 `/claim-result`
7. 页面轮询详情页或等待业务侧状态同步

## 6. 常见错误和排查

### 6.1 `op user id missing in context`

原因：

- 网关没有解析 token
- 网关没有把 `opUserID` 注入上下文
- 直接绕过网关调用了红包服务

### 6.2 `wallet is not bound to user`

原因：

- 当前钱包还没绑定
- 当前钱包绑定的是别的业务用户
- 链类型不一致

### 6.3 `already claimed`

原因：

- 同一个钱包地址已经领过该红包

### 6.4 `user already claimed`

原因：

- 同一个业务用户已经领取过该红包
- 即使换钱包地址，也会被后端拦截

## 7. 后端接口与代码位置

- 接口契约文档：
  [backend-api.md](/Users/panda/aiCode/red_packet/open-im-server-origin/cmd/openim-rpc/openim-rpc-redpacket/backend-api.md)
- 领取签名核心逻辑：
  [redpacket.go](/Users/panda/aiCode/red_packet/open-im-server-origin/cmd/openim-rpc/openim-rpc-redpacket/internal/service/redpacket.go)
- 用户上下文提取：
  [user.go](/Users/panda/aiCode/red_packet/open-im-server-origin/cmd/openim-rpc/openim-rpc-redpacket/internal/authctx/user.go)
