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

#Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${SCRIPTS_ROOT}")/..

#Include shell font styles and some basic information
source $SCRIPTS_ROOT/lib/init.sh
source $SCRIPTS_ROOT/path_info.sh

cd $SCRIPTS_ROOT

source $OPENIM_ROOT/.env

# Check if PASSWORD only contains letters and numbers
if [[ "$PASSWORD" =~ ^[a-zA-Z0-9]+$ ]]
then
    echo "PASSWORD is valid."
else
    echo "ERR: PASSWORD should only contain letters and numbers. " $PASSWORD
    exit
fi

echo ""
echo -e "===> ${PURPLE_PREFIX} you user is:$USER ${COLOR_SUFFIX}"
echo -e "===> ${PURPLE_PREFIX} you password is:$PASSWORD ${COLOR_SUFFIX}" 
echo -e "===> ${PURPLE_PREFIX} you minio endpoint is:$MINIO_ENDPOINT ${COLOR_SUFFIX}"
echo -e "===> ${PURPLE_PREFIX} you api url is:$API_URL ${COLOR_SUFFIX}"
echo ""

# Specify the config file
config_file="${OPENIM_ROOT}"/config/config.yaml

# Load variables from .env file
source "${OPENIM_ROOT}"/.env

# Replace the password and username field for mysql
sed -i "/mysql:/,/database:/ s/password:.*/password: $PASSWORD/" $config_file
sed -i "/mysql:/,/database:/ s/username:.*/username: $USER/" $config_file

# Replace the password and username field for mongo
sed -i "/mongo:/,/maxPoolSize:/ s/password:.*/password: $PASSWORD/" $config_file
sed -i "/mongo:/,/maxPoolSize:/ s/username:.*/username: $USER/" $config_file

# Replace the password field for redis
sed -i '/redis:/,/password:/s/password: .*/password: '${PASSWORD}'/' $config_file

# Replace accessKeyID and secretAccessKey for minio
sed -i "/minio:/,/isDistributedMod:/ s/accessKeyID:.*/accessKeyID: $USER/" $config_file
sed -i "/minio:/,/isDistributedMod:/ s/secretAccessKey:.*/secretAccessKey: $PASSWORD/" $config_file
sed -i '/minio:/,/endpoint:/s|endpoint: .*|endpoint: '${MINIO_ENDPOINT}'|' $config_file
sed -i '/object:/,/apiURL:/s|apiURL: .*|apiURL: '${API_URL}'|' $config_file


# Replace secret for token
sed -i "s/secret: .*/secret: $PASSWORD/" $config_file
