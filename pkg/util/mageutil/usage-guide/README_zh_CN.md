# mageutil使用指南

**mageutil** 是基于 mage 构建的一个工具，它提供了跨平台和多架构的编译支持，同时也简化了服务的启动、停止、检测流程。

## 使用指南

### 准备工作

1. 将此目录下除README以外的6个文件复制到项目的根目录：
2. 项目根目录下需要包含三个目录：`cmd`、`tools`和`config`。
    - `cmd` 目录专门用于存放那些作为后台服务运行的应用的启动代码。
    - `tools`目录用于存放那些作为工具应用（不以后台服务形式运行）的启动代码。
    - `config`目录用于存放配置文件。
3. `cmd`和`tools`目录可以包含多层多个子目录。对于包含`main`函数的`main package`文件，需以`main.go`命名。例如：
    - `cmd/openim-rpc/openim-rpc-msg/main.go`
    - `tools/check-free-memory/main.go`
    - 所有代码都应属于同一个项目，子目录不应使用独立的`go.mod`和`go.sum`文件。

### 初始化项目

- 对于Linux/Mac系统，先执行`bootstrap.sh`脚本。
- 对于Windows系统，先执行`bootstrap.bat`脚本。

### 编译项目

- 执行`mage`或`mage build`来编译项目。
- 编译完成后，二进制文件将生成在`_output/bin/platforms/<操作系统>/<架构>`目录下，其中二进制文件的命名规则为对应的`main.go`所在的目录名。例如：
    - `_output/bin/platforms/linux/amd64/openim-rpc-msg`
    - `_output/bin/tools/linux/amd64/check-free-memory`
    - **注意：** Windows平台的二进制文件会自动添加`.exe`扩展名。

### 启动工具和服务

1. 首先，编辑`start-config.yml`文件，指定服务和工具相关配置，例如：

    ```yaml
    #cmd服务子目录名: 实例数
    serviceBinaries:
        openim-rpc-msg: 2
    
    #tools工具子目录名
    toolBinaries:
        - check-free-memory
    
    maxFileDescriptors: 10000
    ```

    **注意：**服务和工具名与`cmd`和`tools`下的子目录名保持一致

3. 执行`mage start`来启动服务和工具。
    
    - 工具将以同步方式执行，如果工具执行失败（退出代码非零），则整个启动过程中断。
    - 服务将以异步方式启动。

对于所有工具，将采用以下命令格式启动：`[程序绝对路径] -i 0 -c [配置文件绝对目录]`。

若服务实例数设置为`n`，则服务将启动`n`个实例，每个实例使用的命令格式为：`[程序路径] -i [实例索引] -c [配置文件目录]`，其中实例索引从`0`到`n-1`。

**注意**：本项目仅指定了配置文件的路径，并不负责读取配置文件内容。这样做的目的是为了支持使用多个配置文件的情况。程序和配置文件的路径都自动使用绝对路径。

### 检查和停止服务

- 执行`mage check`来检查服务状态和监听的端口。
- 执行`mage stop`来停止服务，该命令会向服务发送停止信号。

---

