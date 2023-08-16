#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# Print the necessary information after installation
function openim::redis::info() {
cat << EOF
Redis Login: redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a '${REDIS_PASSWORD}'
EOF
}

# 安装
function openim::redis::install()
{
  # 1. 安装 Redis
  openim::common::sudo "apt-get -y install redis-server"

  # 2. 配置 Redis
  # 2.1 修改 `/etc/redis/redis.conf` 文件，将 daemonize 由 no 改成 yes，表示允许 Redis 在后台启动
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^daemonize/{s/no/yes/}' /etc/redis/redis.conf

  # 2.2 在 `bind 127.0.0.1` 前面添加 `#` 将其注释掉，默认情况下只允许本地连接，注释掉后外网可以连接 Redis
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^# bind 127.0.0.1/{s/# //}' /etc/redis/redis.conf

  # 2.3 修改 requirepass 配置，设置 Redis 密码
  echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^# requirepass.*$/requirepass '"${REDIS_PASSWORD}"'/' /etc/redis/redis.conf

  # 2.4 因为我们上面配置了密码登录，需要将 protected-mode 设置为 no，关闭保护模式
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^protected-mode/{s/yes/no/}' /etc/redis/redis.conf

  # 3. 为了能够远程连上 Redis，需要执行以下命令关闭防火墙，并禁止防火墙开机启动（如果不需要远程连接，可忽略此步骤）
  openim::common::sudo "sudo ufw disable"
  openim::common::sudo "sudo ufw status"

  # 4. 启动 Redis
  openim::common::sudo "redis-server /etc/redis/redis.conf"

  openim::redis::status || return 1
  openim::redis::info
  openim::log::info "install Redis successfully"
}

# 卸载
function openim::redis::uninstall()
{
  set +o errexit
  openim::common::sudo "/etc/init.d/redis-server stop"
  openim::common::sudo "apt-get -y remove redis-server"
  openim::common::sudo "rm -rf /var/lib/redis"
  set -o errexit
  openim::log::info "uninstall Redis successfully"
}

# 状态检查
function openim::redis::status()
{
  if [[ -z "`pgrep redis-server`" ]];then
    openim::log::error_exit "Redis not running, maybe not installed properly"
    return 1
  fi

  redis-cli --no-auth-warning -h ${REDIS_HOST} -p ${REDIS_PORT} -a "${REDIS_PASSWORD}" --hotkeys || {
    openim::log::error "can not login with ${REDIS_USERNAME}, redis maybe not initialized properly"
    return 1
  }

  openim::log::info "redis-server status active"
}

#eval $*
if [[ "$*" =~ openim::redis:: ]];then
  eval $*
fi
