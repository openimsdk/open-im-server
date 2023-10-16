# v2数据迁移工具

### <font color=red>转换前请做好数据备份！！！</font>

### 转换OPENIM MYSQL数据
 - open-im-server/v3/tools/data-conversion/openim/mysql.go
 - 配置mysql.go数据库信息
 - 需要手动创建v3版本数据库,字符集`utf8mb4`

```go
var (
    usernameV2 = "root"            // v2版本mysql用户名
    passwordV2 = "openIM"          // v2版本mysql密码
    addrV2     = "127.0.0.1:13306" // v2版本mysql地址
    databaseV2 = "openIM_v2"       // v2版本mysql数据库名字
)

var (
    usernameV3 = "root"            // v3版本mysql用户名
    passwordV3 = "openIM123"       // v3版本mysql密码
    addrV3     = "127.0.0.1:13306" // v3版本mysql地址
    databaseV3 = "openIM_v3"       // v3版本mysql数据库名字
)
```
```shell
go run mysql.go
```

### 转换聊天消息(可选)
- 目前只支持转换kafka中的消息
- open-im-server/v3/tools/data-conversion/openim/msg.go
- 配置msg.go数据库信息
```go
var (
	topic       = "ws2ms_chat"      // v2版本配置文件kafka.topic.ws2ms_chat
	kafkaAddr   = "127.0.0.1:9092"  // v2版本配置文件kafka.topic.addr
	rpcAddr     = "127.0.0.1:10130" // v3版本配置文件rpcPort.openImMessagePort
	adminUserID = "openIM123456"    // v3版本管理员userID
	concurrency = 4                 // 并发数量
)
```
```shell
go run msg.go
```

### 转换业务服务器数据(使用官方业务服务器需要转换)
- 目前只支持转换kafka中的消息
- open-im-server/v3/tools/data-conversion/chat/chat.go
- 需要手动创建v3版本数据库,字符集`utf8mb4`
- main.go数据库信息
```go
var (
	usernameV2 = "root"            // v2版本mysql用户名
	passwordV2 = "openIM"          // v2版本mysql密码
	addrV2     = "127.0.0.1:13306" // v2版本mysql地址
	databaseV2 = "admin_chat"      // v2版本mysql数据库名字
)

var (
	usernameV3 = "root"              // v3版本mysql用户名
	passwordV3 = "openIM123"         // v3版本mysql密码
	addrV3     = "127.0.0.1:13306"   // v3版本mysql地址
	databaseV3 = "openim_enterprise" // v3版本mysql数据库名字
)
```
```shell
go run chat.go
```