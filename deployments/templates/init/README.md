# Systemd 配置、安装和启动

- [Systemd 配置、安装和启动](#systemd-配置安装和启动)
  - [0. 介绍](#0-介绍)
  - [1. 前置操作（需要 root 权限）](#1-前置操作需要-root-权限)
  - [2. 创建 openim-api systemd unit 模板文件](#2-创建-openim-api-systemd-unit-模板文件)
  - [3. 创建 openim-crontask systemd unit 模板文件](#3-创建-openim-crontask-systemd-unit-模板文件)
  - [6. 复制 systemd unit 模板文件到 sysmted 配置目录(需要有root权限)](#6-复制-systemd-unit-模板文件到-sysmted-配置目录需要有root权限)
  - [7. 启动 systemd 服务](#7-启动-systemd-服务)

## 0. 介绍

systemd是最新linux发行版管理后台的服务的默认形式，用以取代原有的init。

格式介绍：

```bash
[Unit] :服务单元

Description：对该服务进行简单的描述

[Service]:服务运行时行为配置

ExecStart：程序的完整路径

Restart:重启配置，no、always、on-success、on-failure、on-abnormal、on-abort、on-watchdog

[Install]：安装配置

WantedBy：多用户等
```

更多介绍阅读：https://www.freedesktop.org/software/systemd/man/systemd.service.html

启动命令：

```bash
systemctl daemon-reload && systemctl enable openim-api && systemctl restart openim-api
```

服务状态：

```bash
systemctl status openim-api
```

停止命令：

```bash
systemctl stop openim-api
```

**为什么选择 systemd？**

**高级需求：**

+ 方便分析问题的服务运行日志记录

+ 服务管理的日志

+ 异常退出时可以根据需要重新启动

daemon不能实现上面的高级需求。

nohup 只能记录服务运行时的输出和出错日志。

只有systemd能够实现上述所有需求。

> 默认的日志中增加了时间、用户名、服务名称、PID等，非常人性化。还能看到服务运行异常退出的日志。还能通过/lib/systemd/system/下的配置文件定制各种需求。

总而言之，systemd是目前linux管理后台服务的主流方式，所以我新版本的 bash 抛弃 nohup，改用 systemd 来管理服务。



## 1. 前置操作（需要 root 权限）

1. 根据注释配置 `environment.sh`

2. 创建 data 目录

```
mkdir -p ${OPENIM_DATA_DIR}/{openim-api,openim-crontask}
```

3. 创建 bin 目录，并将 `openim-api` 和 `openim-crontask` 可执行文件复制过去

```bash
source ./environment.sh
mkdir -p ${OPENIM_INSTALL_DIR}/bin
cp openim-api openim-crontask ${OPENIM_INSTALL_DIR}/bin
```

4. 将 `openim-api` 和 `openim-crontask` 配置文件拷贝到 `${OPENIM_CONFIG_DIR}` 目录下

```bash
mkdir -p ${OPENIM_CONFIG_DIR}
cp openim-api.yaml openim-crontask.yaml ${OPENIM_CONFIG_DIR}
```

## 2. 创建 openim-api systemd unit 模板文件

执行如下 shell 脚本生成 `openim-api.service.template`

```bash
source ./environment.sh
cat > openim-api.service.template <<EOF
[Unit]
Description=OpenIM Server API
Documentation=https://github.com/marmotedu/iam/blob/master/init/README.md

[Service]
WorkingDirectory=${OPENIM_DATA_DIR}/openim-api
ExecStart=${OPENIM_INSTALL_DIR}/bin/openim-api --apiconfig=${OPENIM_CONFIG_DIR}/openim-api.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
```

## 3. 创建 openim-crontask systemd unit 模板文件
...


## 6. 复制 systemd unit 模板文件到 sysmted 配置目录(需要有root权限)

```bash
cp openim-api.service.template /etc/systemd/system/openim-api.service
cp openim-crontask.service.template /etc/systemd/system/openim-crontask.service
...
```

## 7. 启动 systemd 服务

```bash
systemctl daemon-reload && systemctl enable openim-api && systemctl restart openim-api
systemctl daemon-reload && systemctl enable openim-crontask && systemctl restart openim-crontask
...
```
