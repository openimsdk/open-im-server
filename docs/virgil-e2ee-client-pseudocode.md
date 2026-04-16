# OpenIM + Virgil E2EE 客户端伪代码模板

本文给出可直接改造成业务代码的四段伪代码模板：

1. 登录初始化
2. 单聊加密发送与接收解密
3. 群聊加密发送与接收解密
4. 群密钥版本追平（增量同步）

> 约定：
> - OpenIM SDK 负责业务消息收发、会话管理。
> - OpenIM CryptoService 负责设备注册、JWT、群版本事件。
> - Virgil E3Kit 负责本地加解密。
> - 服务端不解密消息明文。

---

## 0. 公共结构与工具函数

```typescript
type DeviceInfo = {
  deviceID: string
  platform: string
  model: string
  appVersion: string
}

type CipherEnvelope = {
  alg: "virgil-e3kit-v1"
  senderUserID: string
  senderDeviceID: string
  // base64 ciphertext
  payload: string
  // 扩展字段：会话ID、消息版本等
  meta?: Record<string, string>
}

type GroupCipherEnvelope = {
  alg: "virgil-group-v1"
  groupID: string
  groupKeyVersion: number
  senderUserID: string
  senderDeviceID: string
  payload: string
  meta?: Record<string, string>
}

function toBase64(input: Uint8Array): string { /* ... */ }
function fromBase64(input: string): Uint8Array { /* ... */ }
function nowMs(): number { return Date.now() }

function buildDeviceInfo(): DeviceInfo {
  return {
    deviceID: getStableDeviceID(), // 本地持久化的设备唯一ID
    platform: getPlatformName(),   // iOS/Android/Web/Desktop
    model: getDeviceModel(),
    appVersion: getAppVersion(),
  }
}
```

---

## 1) 登录初始化（设备注册 + Virgil 初始化）

```typescript
async function loginAndInitE2EE(userID: string, token: string) {
  // Step 1: 初始化 OpenIM SDK 登录态
  await openim.login({ userID, token })

  // Step 2: 注册设备（幂等）
  const device = buildDeviceInfo()
  await openim.crypto.registerDevice({
    deviceID: device.deviceID,
    platform: device.platform,
    deviceModel: device.model,
    appVersion: device.appVersion,
  })

  // Step 3: 获取 Virgil JWT（短期）
  const jwtResp = await openim.crypto.getVirgilJWT({
    deviceID: device.deviceID,
  })
  const virgilJWT = jwtResp.virgilJWT
  const virgilIdentity = jwtResp.virgilIdentity // 一般是 userID:deviceID

  // Step 4: 初始化 E3Kit（客户端本地）
  const e3 = await E3Kit.initialize(async () => virgilJWT)

  // Step 5: 首次设备场景（仅示例）
  // - 若本地无私钥：生成并 publishCard
  // - 若有 Keyknox 备份：restorePrivateKey
  if (!(await e3.hasLocalPrivateKey())) {
    await e3.register() // 内部通常会 publish card
  }

  // Step 6: 安全预检（可选但建议）
  const precheck = await openim.crypto.securityPrecheck({
    deviceID: device.deviceID,
    action: "login_init_e2ee",
  })
  if (!precheck.allowed) {
    throw new Error(`security precheck denied: ${precheck.reason}`)
  }

  // Step 7: 持有上下文
  return {
    userID,
    deviceID: device.deviceID,
    virgilIdentity,
    e3,
  }
}
```

---

## 2) 单聊发收（发送加密 / 接收解密）

### 2.1 发送单聊加密消息

```typescript
async function sendPrivateEncryptedText(
  ctx: { userID: string; deviceID: string; e3: E3Kit },
  toUserID: string,
  plainText: string
) {
  // 1) 拉取接收方 card（Virgil）
  const recipientCard = await ctx.e3.findUsers(toUserID)
  if (!recipientCard) throw new Error("recipient card not found")

  // 2) 本地加密
  const encrypted = await ctx.e3.encrypt(plainText, recipientCard)

  // 3) 封装消息体（发给 OpenIM 的 content）
  const envelope: CipherEnvelope = {
    alg: "virgil-e3kit-v1",
    senderUserID: ctx.userID,
    senderDeviceID: ctx.deviceID,
    payload: toBase64(encrypted),
  }

  // 4) 通过 OpenIM 普通消息通道发送（服务端仅转发存储密文）
  await openim.sendMessage({
    recvID: toUserID,
    conversationType: "single",
    contentType: "custom_e2ee_text", // 你项目定义的 contentType
    content: JSON.stringify(envelope),
  })
}
```

### 2.2 接收单聊加密消息并解密

```typescript
async function onPrivateMessageReceived(
  ctx: { e3: E3Kit },
  msg: { contentType: string; content: string; sendID: string }
) {
  if (msg.contentType !== "custom_e2ee_text") return

  const envelope = JSON.parse(msg.content) as CipherEnvelope
  if (envelope.alg !== "virgil-e3kit-v1") return

  // 1) 拉取发送方 card
  const senderCard = await ctx.e3.findUsers(msg.sendID)
  if (!senderCard) {
    markMessageDecryptFailed(msg, "sender card not found")
    return
  }

  // 2) 本地解密
  try {
    const plainText = await ctx.e3.decrypt(fromBase64(envelope.payload), senderCard)
    renderMessage(msg, plainText)
  } catch (err) {
    markMessageDecryptFailed(msg, `decrypt failed: ${String(err)}`)
  }
}
```

---

## 3) 群聊发收（发送加密 / 接收解密）

### 3.1 发送群聊加密消息

```typescript
async function sendGroupEncryptedText(
  ctx: { userID: string; deviceID: string; e3: E3Kit },
  groupID: string,
  plainText: string
) {
  // 1) 先确保本地群密钥版本已追平（见第4段）
  const version = await ensureGroupKeySynced(ctx, groupID)

  // 2) 获取/创建本地 group session（示意）
  const groupSession = await getOrCreateGroupSession(ctx.e3, groupID, version)

  // 3) 本地加密
  const encrypted = await groupSession.encrypt(plainText)

  // 4) 封装消息体
  const envelope: GroupCipherEnvelope = {
    alg: "virgil-group-v1",
    groupID,
    groupKeyVersion: version,
    senderUserID: ctx.userID,
    senderDeviceID: ctx.deviceID,
    payload: toBase64(encrypted),
  }

  // 5) 发送到 OpenIM
  await openim.sendMessage({
    recvID: groupID,
    conversationType: "group",
    contentType: "custom_e2ee_group_text",
    content: JSON.stringify(envelope),
  })
}
```

### 3.2 接收群聊加密消息并解密

```typescript
async function onGroupMessageReceived(
  ctx: { e3: E3Kit },
  msg: { groupID: string; contentType: string; content: string }
) {
  if (msg.contentType !== "custom_e2ee_group_text") return

  const envelope = JSON.parse(msg.content) as GroupCipherEnvelope
  if (envelope.alg !== "virgil-group-v1") return

  // 1) 若消息携带版本比本地新，先追平
  const localVersion = await loadLocalGroupKeyVersion(msg.groupID)
  if (envelope.groupKeyVersion > localVersion) {
    await ensureGroupKeySynced({ e3: ctx.e3 } as any, msg.groupID)
  }

  // 2) 取本地会话并解密
  try {
    const groupSession = await getOrCreateGroupSession(ctx.e3, msg.groupID, envelope.groupKeyVersion)
    const plainText = await groupSession.decrypt(fromBase64(envelope.payload))
    renderMessage(msg, plainText)
  } catch (err) {
    markMessageDecryptFailed(msg, `group decrypt failed: ${String(err)}`)
  }
}
```

---

## 4) 群密钥版本追平（增量同步模板）

```typescript
async function ensureGroupKeySynced(
  ctx: { userID?: string; deviceID?: string; e3: E3Kit },
  groupID: string
): Promise<number> {
  // A) 服务端当前版本
  const latestResp = await openim.crypto.getGroupKeyVersion({ groupID })
  const latestVersion = Number(latestResp.groupKeyVersion || 0)

  // B) 本地版本
  let localVersion = await loadLocalGroupKeyVersion(groupID) // 默认 0
  if (localVersion >= latestVersion) return localVersion

  // C) 拉取增量事件
  const eventsResp = await openim.crypto.getGroupKeyEvents({
    groupID,
    sinceVersion: localVersion,
  })
  const events = eventsResp.events || []

  // D) 按版本顺序应用事件（关键）
  events.sort((a, b) => Number(a.groupKeyVersion) - Number(b.groupKeyVersion))
  for (const ev of events) {
    const targetVersion = Number(ev.groupKeyVersion)
    if (targetVersion <= localVersion) continue

    // 示例：根据事件刷新 group ticket / 重新分发会话材料
    // 具体实现取决于你对 Virgil Group Tickets 的封装
    await applyGroupKeyEvent(ctx.e3, groupID, {
      eventType: ev.eventType,
      operatorUserID: ev.operatorUserID,
      targetVersion,
    })

    localVersion = targetVersion
    await saveLocalGroupKeyVersion(groupID, localVersion)
  }

  // E) 防御性校验
  if (localVersion < latestVersion) {
    // 说明事件缺失，可做一次全量恢复
    await rebuildGroupSessionFromSource(ctx.e3, groupID, latestVersion)
    localVersion = latestVersion
    await saveLocalGroupKeyVersion(groupID, localVersion)
  }

  return localVersion
}
```

---

## 5) 建议落地策略（简版）

- **密钥刷新时机**：应用启动、会话进入、收到群消息且版本落后、网络重连后。
- **失败重试**：`get_virgil_jwt`、`get_group_key_events` 使用指数退避重试。
- **本地缓存**：按 `groupID -> groupKeyVersion` 持久化，避免重复全量同步。
- **消息兼容**：保留 `alg` 字段，支持后续算法版本升级。
- **安全日志**：不要记录明文与私钥，仅记录 `groupID/userID/version/eventType`。

---

## 6) 与当前服务端接口对应

当前服务端公开路由（`internal/api/router.go`）：

- `/crypto/register_device`
- `/crypto/get_devices`
- `/crypto/revoke_device`
- `/crypto/get_virgil_jwt`
- `/crypto/get_group_key_version`
- `/crypto/get_group_key_events`
- `/crypto/security_precheck`
- `/crypto/integrity_report`

> `bump_group_key_version` 为内部服务调用，不提供给客户端公开使用。

