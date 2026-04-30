# RedPacket Go 后端对接说明（ETH + TRON）

本文档基于当前 OpenIM 版红包服务实现整理，重点说明 Go 后端如何接入 EVM / TRON 链能力、如何签发 claim 授权、如何解析交易事件，以及当前实现中哪些能力是完整实现、哪些仍是 mock 或待补齐。

相关代码位置：

- RPC 入口：`cmd/openim-rpc/openim-rpc-redpacket/main.go`
- 服务启动：`pkg/common/cmd/rpc_redpacket.go`
- 业务逻辑：`internal/rpc/redpacket/service.go`
- 管理接口：`internal/rpc/redpacket/admin.go`
- 钱包绑定：`internal/rpc/redpacket/wallet.go`
- 链客户端：`internal/rpc/redpacket/chain`
- 合约 ABI：`internal/rpc/redpacket/chain/abi/RedPacket.json`
- 配置文件：`config/openim-rpc-redpacket.yml`

## 1. 当前架构

`openim-rpc-redpacket` 已经不再是独立 Gin + GORM 服务，而是标准 OpenIM RPC 服务：

```text
openim-api
  -> /redpacket/* HTTP API
  -> pbredpacket.RedPacketClient
  -> openim-rpc-redpacket
  -> MongoDB + EVM/TRON clients
```

服务启动时会初始化：

- MongoDB DAO：`controller.NewRedPacketDatabase(...)`
- EVM client：当 `chain.rpcURL` 与 `chain.contractAddress` 配置完整时启用
- TRON client：当 `tron.fullNodeURL` 与 `tron.contractBase58` 配置完整时启用
- signer 私钥：当 `chain.signerPrivateKey` 配置完整时用于 claim 裸签名

链客户端初始化失败不会阻止服务启动，但会导致链上确认、事件解析或签名 digest 获取降级。

## 2. 配置

`config/openim-rpc-redpacket.yml` 示例：

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
  rpcURL: "https://eth-mainnet.g.alchemy.com/v2/xxx"
  contractAddress: "0x..."
  chainID: 1
  signerPrivateKey: "0x..."
  configAdminPrivateKey: "0x..."

tron:
  fullNodeURL: "https://api.trongrid.io"
  contractBase58: "T..."
  ownerBase58: "T..."
  privateKeyHex: "..."
  feeLimit: 100000000

indexer:
  pollInterval: 5
```

配置含义：

- `chain.rpcURL`: EVM JSON-RPC 地址
- `chain.contractAddress`: EVM RedPacket 合约地址
- `chain.chainID`: EVM 链 ID；用于记录业务单与构造交易
- `chain.signerPrivateKey`: claim 授权签名私钥，应对应合约 `signer`
- `chain.configAdminPrivateKey`: 管理写链私钥，当前 EVM admin 仍是 mock
- `tron.fullNodeURL`: TRON FullNode / TronGrid 地址
- `tron.contractBase58`: TRON 合约 Base58 地址
- `tron.ownerBase58`: TRON 管理交易发送地址
- `tron.privateKeyHex`: TRON 管理交易私钥
- `tron.feeLimit`: TRON 交易 fee limit

安全建议：

- `signerPrivateKey` 与 `configAdminPrivateKey` 必须分离
- 生产不要把管理私钥明文放在普通配置文件中，建议接入 KMS/HSM 或密钥托管服务
- `signerPrivateKey` 是高频签名密钥，权限只能用于 claim 授权，不应拥有合约配置权限

## 3. Claim 签名

### 3.1 合约签名事实

当前后端签名逻辑对应合约的：

```text
getSignMessage(packetId, claimer, authNonce, randomSeed, deadline)
claim(packetId, authNonce, randomSeed, deadline, signature)
```

后端流程：

1. 业务鉴权：登录用户、钱包绑定、红包状态、重复领取、群/转账资格
2. 生成 `authNonce`、`randomSeed`、`deadline`
3. EVM client 可用时调用链上 `getSignMessage(...)` 获取 digest
4. 用 `signerPrivateKey` 对 digest 做裸签名
5. 如果 `v` 是 0/1，转换为 27/28
6. 保存 `red_packet_claim_auth`
7. 返回前端调用 `claim(...)` 所需参数

注意：不要使用 `personal_sign` 对 claim digest 签名。claim 授权使用的是裸 ECDSA 签名，不带 Ethereum Signed Message 前缀。

### 3.2 Go 裸签名示例

```go
func signClaimDigest(priv *ecdsa.PrivateKey, digest [32]byte) (string, error) {
    sig, err := crypto.Sign(digest[:], priv)
    if err != nil {
        return "", err
    }
    if len(sig) == 65 && sig[64] < 27 {
        sig[64] += 27
    }
    return "0x" + hex.EncodeToString(sig), nil
}
```

### 3.3 当前降级行为

当前代码有两个降级点：

- EVM client 不可用时，后端会用本地 `keccak256(packetID:claimer:nonce:randomSeed:deadline)` 生成 digest；该 digest 不保证与合约一致，仅适合调试。
- signer 私钥未配置时，后端会返回 placeholder 签名；该签名不能通过链上验签。

生产环境必须配置可用的 EVM client 和 signer 私钥。

## 4. ETH 接入

### 4.1 创建红包

推荐调用顺序：

1. 后端 `CreateOrder` 生成 `biz_id`
2. 前端或托管钱包发起链上创建交易
3. 从 `PacketCreated` 事件解析 `packetId`
4. 调用 `CreatedCallback` 回写 `biz_id + tx_hash + packet_id`
5. 后端使用 EVM client 解析 receipt 并校验事件字段
6. 校验通过后业务单变为 `ACTIVE`

当前代码中的校验点：

- `tx_hash` 必填
- receipt 中必须有可识别的 `PacketCreated`
- event 解析出的 creator / packetType / token / amount / shares / expiry 要与业务单一致
- 如果链客户端不可用，允许请求体提供 `packet_id` fallback

### 4.2 领取红包

推荐调用顺序：

1. 前端确认用户已经绑定当前 EVM 钱包
2. 调用 `IssueClaimSign`
3. 前端使用返回参数调用合约 `claim(...)`
4. 交易提交后调用 `ClaimResult`
5. 后端解析 `PacketClaimed`，补全 amount、authNonce、blockNumber

`ClaimResult` 当前行为：

- 先落 `PENDING` 领取记录
- 能解析 receipt 时更新为 `CONFIRMED`
- 解析到 `PacketClaimed` 后更新红包领取进度
- 已领取份数达到 `total_shares` 时状态更新为 `COMPLETED`

### 4.3 事件解析

EVM 事件解析由 `internal/rpc/redpacket/chain/parser.go` 负责。管理接口也提供手动解析入口：

```http
POST /redpacket/admin/parse_tx_events
```

请求：

```json
{
  "chain": "eth",
  "tx_hash": "0xabc123..."
}
```

响应示例：

```json
{
  "chain": "eth",
  "tx_hash": "0xabc123...",
  "events": [
    {
      "name": "PacketCreated",
      "data": {
        "packetId": "10001",
        "creator": "0x1111111111111111111111111111111111111111",
        "packetType": "1"
      }
    }
  ]
}
```

核心事件：

- `PacketCreated`: 创建成功，提供唯一可信 `packetId`
- `PacketClaimed`: 领取成功，提供实际领取金额
- `PacketRefunded`: 退款成功，提供退款金额与接收方

### 4.4 ETH 管理接口现状

当前 `internal/rpc/redpacket/admin.go` 中 EVM 管理接口是 mock：

- `SetSigner`
- `SetToken`
- `SetExpiry`
- `SetAllowAllTokens`
- `SetNativeTokenEnabled`

这些接口在 EVM client 可用时只记录日志并返回成功 message，不会真正发链上交易。上线前如需后端托管管理交易，需要补充 EVM admin transaction 实现。

## 5. TRON 接入

### 5.1 TRON 创建与领取

TRON 合约兼容 EVM ABI 的 topic/data 事件模型，但地址、签名与交易广播流程和 EVM 不同。

当前后端支持：

- 创建业务单时 `chain_type=TRON`
- `contract_address` 可从 `tron.contractBase58` 自动填充
- TRON 钱包绑定 challenge 生成
- TRON admin 写交易通过 `SendAdminTransaction(...)` 尝试调用 FullNode

当前后端尚未完整支持：

- TRON 钱包绑定签名验签
- TRON claim digest 获取与 claim 签名链上闭环
- TRON receipt 事件完整解析与索引

### 5.2 TRON 管理交易

当前 TRON admin 使用 FullNode HTTP 流程：

```text
triggersmartcontract
  -> gettransactionsign
  -> broadcasttransaction
```

配置依赖：

- `tron.fullNodeURL`
- `tron.contractBase58`
- `tron.ownerBase58`
- `tron.privateKeyHex`
- `tron.feeLimit`

管理接口会把方法映射到合约调用：

- `SetSigner` -> `setSigner`
- `SetToken` -> `setAllowedToken`
- `SetExpiry` -> `setDefaultExpiryDuration`
- `SetAllowAllTokens` -> `setAllowAllTokens`
- `SetNativeTokenEnabled` -> `setNativeTokenEnabled`

### 5.3 TRON 事件解析现状

`ParseTxEvents(chain=tron)` 当前返回：

```json
{
  "chain": "tron",
  "tx_hash": "7d9e...txid",
  "note": "TRON event parsing not fully implemented in this version"
}
```

后续如果要补齐，应实现：

1. 调用 `/wallet/gettransactioninfobyid`
2. 从 `log` 读取 topics/data
3. 将 TRON 地址字段规范化为 Base58 或 hex
4. 使用 `RedPacket.json` ABI 解码事件
5. 复用 EVM 的 `PacketCreated` / `PacketClaimed` / `PacketRefunded` 业务回写逻辑

## 6. 钱包绑定

### 6.1 EVM 绑定

EVM 绑定采用 SIWE 风格消息：

- protocol: `siwe-eip4361`
- sign method: `personal_sign`
- challenge 有效期: 10 分钟

确认绑定时，后端会：

1. 读取 `wallet_binding_challenge`
2. 检查状态为 `PENDING`
3. 检查未过期
4. 用 `personalSignMessage(message)` 计算 hash
5. `SigToPub` recover 地址
6. 比对 recover 地址与 challenge wallet
7. challenge 更新为 `VERIFIED`
8. upsert `wallet_binding`

### 6.2 TRON 绑定

TRON challenge 会生成：

- protocol: `tron-signmessagev2`
- sign method: `signMessageV2`

但确认绑定当前未实现，会返回：

```text
TRON wallet binding verification is not implemented yet
```

## 7. MongoDB 数据

当前使用 6 个 collection：

- `red_packet`: 红包主记录
- `red_packet_claim`: 领取记录
- `red_packet_claim_auth`: claim 签名授权记录
- `red_packet_refund`: 退款记录
- `wallet_binding_challenge`: 钱包绑定 challenge
- `wallet_binding`: 钱包绑定关系

关键幂等约束：

- `red_packet.biz_id` 唯一
- `red_packet_claim.claim_tx_hash` 唯一
- `red_packet_claim_auth.auth_nonce` 唯一
- `wallet_binding_challenge.challenge_id` 唯一
- `wallet_binding.user_id + chain_type + wallet_address` 唯一

## 8. 部署检查清单

上线前至少确认：

- `share.yml` 中存在 `rpcRegisterName.redPacket: redPacket`
- `openim-rpc-redpacket.yml` 已加入配置目录
- `openim-api` watch service list 包含 `redPacket`
- MongoDB 可用且服务启动时能创建索引
- EVM 环境配置了有效 `rpcURL`、`contractAddress`、`signerPrivateKey`
- 生产关闭 placeholder signer 降级路径
- 管理接口补充管理员鉴权与操作审计
- 如需 ETH admin 写链，补齐当前 mock 实现
- 如需 TRON 完整闭环，补齐绑定验签、事件解析、claim 签名链路
