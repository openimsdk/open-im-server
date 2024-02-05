# OpenIM V2 至 V3 数据迁移指南

该指南提供了从 OpenIM V2 迁移至 V3 的详细步骤。请确保在开始迁移过程之前，熟悉所有步骤，并按照指南准确执行。

+ [OpenIM Chat](https://github.com/OpenIMSDK/chat)
+ [OpenIM Server](https://github.com/OpenIMSDK/Open-IM-Server)



### 1. 数据备份

在开始数据迁移之前，强烈建议备份所有相关的数据以防止任何可能的数据丢失。

### 2. 迁移 OpenIM MySQL 数据

+ 位置: `open-im-server/tools/data-conversion/openim/cmd/conversion-mysql.go`
+ 配置 `conversion-mysql.go` 文件中的数据库信息。
+ 手动创建 V3 版本的数据库，并确保字符集为 `utf8mb4`。

```bash
// V2 数据库配置
var (
    usernameV2 = "root"
    passwordV2 = "openIM"
    addrV2     = "127.0.0.1:13306"
    databaseV2 = "openIM_v2"
)

// V3 数据库配置
var (
    usernameV3 = "root"
    passwordV3 = "openIM123"
    addrV3     = "127.0.0.1:13306"
    databaseV3 = "openim_v3"
)
```

**执行数据迁移命令：**

```bash
make build BINS="conversion-mysql"
```

启动的二进制在 `_output/bin/tools` 中


### 3. 转换聊天消息（可选）

+ 只支持转换存储在 Kafka 中的消息。
+ 位置: `open-im-server/tools/data-conversion/openim/conversion-msg/conversion-msg.go`
+ 配置 `msg.go` 文件中的消息和服务器信息。

```bash
var (
	topic       = "ws2ms_chat"      // V2 版本 Kafka 主题
	kafkaAddr   = "127.0.0.1:9092"  // V2 版本 Kafka 地址
	rpcAddr     = "127.0.0.1:10130" // V3 版本 RPC 地址
	adminUserID = "openIM123456"    // V3 版本管理员用户ID
	concurrency = 4                 // 并发数量
)
```

**执行数据迁移命令：**

```bash
make build BINS="conversion-msg"
```

### 4. 转换业务服务器数据

+ 只支持转换存储在 Kafka 中的消息。
+ 位置: `open-im-server/tools/data-conversion/chat/cmd/conversion-chat/chat.go`
+ 需要手动创建 V3 版本的数据库，并确保字符集为 `utf8mb4`。
+ 配置 `main.go` 文件中的数据库信息。

```bash
// V2 数据库配置
var (
	usernameV2 = "root"
	passwordV2 = "openIM"
	addrV2     = "127.0.0.1:13306"
	databaseV2 = "admin_chat"
)

// V3 数据库配置
var (
	usernameV3 = "root"
	passwordV3 = "openIM123"
	addrV3     = "127.0.0.1:13306"
	databaseV3 = "openim_enterprise"
)
```

**执行数据迁移命令：**

```bash
make build BINS="conversion-chat"
```
