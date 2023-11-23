# OpenIM V3.4.0 至 V3.5.0 数据迁移指南

---
从3.5.0起，我们从MySQL切换到了MongoDB，这意味着您需要将数据从MySQL迁移到MongoDB。我们提供了一个工具来帮助您完成这项工作。本次迁移完成后完全兼容v3之前的数据。

### 1. 数据备份

在开始数据迁移之前，强烈建议备份所有相关的数据以防止任何可能的数据丢失。

### 2. 迁移数据

+ 位置: `open-im-server/tools/mysql2mongo/main.go`

```bash
// 数据库配置
var (
	mysqlUsername = "root"            // mysql用户名
	mysqlPassword = "openIM123"       // mysql密码
	mysqlAddr     = "127.0.0.1:13306" // mysql地址
	mysqlDatabase = "openIM_v3"       // mysql数据库名字
)

var s3 = "minio" // 文件储存方式 minio, cos, oss

var (
	mongoUsername = "root"            // mongodb用户名
	mongoPassword = "openIM123"       // mongodb密码
	mongoHosts    = "127.0.0.1:13306" // mongodb地址
	mongoDatabase = "openIM_v3"       // mongodb数据库名字
)
```

**执行数据迁移命令：**

```bash
go run main.go
```