#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# 本脚本功能：根据 scripts/environment.sh 配置，生成 OPENIM 组件 YAML 配置文件。
# 示例：genconfig.sh scripts/environment.sh configs/openim-api.yaml

# Path to the original script file
env_file="$1"
# Path to the generated config file
template_file="$2"

. $(dirname ${BASH_SOURCE})/lib/init.sh

if [ $# -ne 2 ];then
    openim::log::error "Usage: genconfig.sh scripts/environment.sh configs/openim-api.yaml"
    exit 1
fi

source "${env_file}"

declare -A envs

set +u
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file})
do
    if [ -z "$(eval echo \$${env})" ];then
        openim::log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ "${missing}" ];then
    openim::log::error 'You may run `source scripts/environment.sh` to set these environment'
    exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"
