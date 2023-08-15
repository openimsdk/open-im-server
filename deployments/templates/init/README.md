# Systemd Configuration, Installation, and Startup

- [Systemd Configuration, Installation, and Startup](#systemd-configuration-installation-and-startup)
  - [1. Introduction](#1-introduction)
  - [2. Prerequisites (Requires root permissions)](#2-prerequisites-requires-root-permissions)
  - [3. Create `openim-api` systemd unit template file](#3-create-openim-api-systemd-unit-template-file)
  - [4. Copy systemd unit template file to systemd config directory (Requires root permissions)](#4-copy-systemd-unit-template-file-to-systemd-config-directory-requires-root-permissions)
  - [5. Start systemd service](#5-start-systemd-service)

## 1. Introduction

Systemd is the default service management form for the latest Linux distributions, replacing the original init.

Format introduction:

```bash
[Unit] : Unit of the service

Description: Brief description of the service

[Service]: Configuration of the service's runtime behavior

ExecStart: Complete path of the program

Restart: Restart configurations like no, always, on-success, on-failure, on-abnormal, on-abort, on-watchdog

[Install]: Installation configuration

WantedBy: Multi-user, etc.
```

For more details, refer to: [Systemd Service Manual](https://www.freedesktop.org/software/systemd/man/systemd.service.html)

Starting command:

```bash
systemctl daemon-reload && systemctl enable openim-api && systemctl restart openim-api
```

Service status:

```bash
systemctl status openim-api
```

Stop command:

```bash
systemctl stop openim-api
```

More command:
```bash
# 列出正在运行的Unit,可以直接使用systemctl
systemctl list-units

# 列出所有Unit，包括没有找到配置文件的或者启动失败的
systemctl list-units --all

# 列出所有没有运行的 Unit
systemctl list-units --all --state=inactive

# 列出所有加载失败的 Unit
systemctl list-units --failed

# 列出所有正在运行的、类型为service的Unit
systemctl list-units --type=service

# 显示某个 Unit 是否正在运行
systemctl is-active application.service

# 显示某个 Unit 是否处于启动失败状态
systemctl is-failed application.service

# 显示某个 Unit 服务是否建立了启动链接
systemctl is-enabled application.service

# 立即启动一个服务
sudo systemctl start apache.service

# 立即停止一个服务
sudo systemctl stop apache.service

# 重启一个服务
sudo systemctl restart apache.service

# 重新加载一个服务的配置文件
sudo systemctl reload apache.service

# 重载所有修改过的配置文件
sudo systemctl daemon-reload
```

**Why choose systemd?**

**Advanced requirements:**

- Convenient service runtime log recording for problem analysis
- Service management logs
- Option to restart upon abnormal exit

The daemon does not meet these advanced requirements.

`nohup` only logs the service's runtime outputs and errors.

Only systemd can fulfill all of the above requirements.

> The default logs are enhanced with timestamps, usernames, service names, PIDs, etc., making them user-friendly. You can view logs of abnormal service exits. Advanced customization is possible through the configuration files in `/lib/systemd/system/`.

In short, systemd is the current mainstream way to manage backend services on Linux, so I've abandoned `nohup` in my new versions of bash scripts, opting instead for systemd.

## 2. Prerequisites (Requires root permissions)

1. Configure `environment.sh` based on the comments.
2. Create a data directory:

```bash
mkdir -p ${OPENIM_DATA_DIR}/{openim-api,openim-crontask}
```

3. Create a bin directory and copy `openim-api` and `openim-crontask` executable files:

```bash
source ./environment.sh
mkdir -p ${OPENIM_INSTALL_DIR}/bin
cp openim-api openim-crontask ${OPENIM_INSTALL_DIR}/bin
```

4. Copy the configuration files of `openim-api` and `openim-crontask` to the `${OPENIM_CONFIG_DIR}` directory:

```bash
mkdir -p ${OPENIM_CONFIG_DIR}
cp openim-api.yaml openim-crontask.yaml ${OPENIM_CONFIG_DIR}
```

## 3. Create `openim-api` systemd unit template file

For each OpenIM service, we will create a systemd unit template. Follow the steps below for each service:

Run the following shell script to generate the `openim-api.service.template`:

```bash
source ./environment.sh
cat > openim-api.service.template <<EOF
[Unit]
Description=OpenIM Server API
Documentation=https://github.com/marmotedu/iam/blob/master/init/README.md

[Service]
WorkingDirectory=${OPENIM_DATA_DIR}/openim-api
ExecStart=${OPENIM_INSTALL_DIR}/bin/openim-api --config=${OPENIM_CONFIG_DIR}/openim-api.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
```

Following the above style, create the respective template files or generate them in bulk:

First, make sure you've sourced the environment variables:

```bash
source ./environment.sh
```

Use the shell script to generate the systemd unit template for each service:

```bash
declare -a services=(
"openim-api"
... [other services]
)

for service in "${services[@]}"
do
   cat > $service.service.template <<EOF
[Unit]
Description=OpenIM Server - $service
Documentation=https://github.com/marmotedu/iam/blob/master/init/README.md

[Service]
WorkingDirectory=${OPENIM_DATA_DIR}/$service
ExecStart=${OPENIM_INSTALL_DIR}/bin/$service --config=${OPENIM_CONFIG_DIR}/$service.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
done
```

## 4. Copy systemd unit template file to systemd config directory (Requires root permissions)

Ensure you have root permissions to perform this operation:

```bash
for service in "${services[@]}"
do
   sudo cp $service.service.template /etc/systemd/system/$service.service
done
...
```

## 5. Start systemd service

To start the OpenIM services:

```bash
for service in "${services[@]}"
do
   sudo systemctl daemon-reload 
   sudo systemctl enable $service 
   sudo systemctl restart $service
done
```
