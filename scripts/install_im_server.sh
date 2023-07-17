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


# Common utilities, variables and checks for all build scripts.
set -o errexit
set -o nounset
set -o pipefail

#Include shell font styles and some basic information
SCRIPTS_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source $SCRIPTS_ROOT/style_info.sh

# docker-compose.yaml file name
docker_compose_file_name="docker-compose.yaml"

trap 'onCtrlC' INT
function onCtrlC () {
    #Capture CTRL+C, terminate the background process of the program when the script is terminated in the form of ctrl+c
    kill -9 ${do_sth_pid} ${progress_pid}
    echo
    echo 'Ctrl+C is captured'
    exit 1
}

# Get the public internet IP address
internet_ip=$(curl ifconfig.me -s)
echo -e "${PURPLE_PREFIX}=========> Your public internet IP address is ${internet_ip} ${COLOR_SUFFIX} \n"

# Load environment variables from .env file
source ${OPENIM_ROOT}/.env

echo -e "${PURPLE_PREFIX}=========> Your minio endpoint is ${MINIO_ENDPOINT} ${COLOR_SUFFIX} \n"

# Change directory to scripts folder

chmod +x ${SCRIPTS_ROOT}/*.sh

# Execute necessary scripts
echo -e "${PURPLE_PREFIX}=========> init_pwd.sh ${COLOR_SUFFIX} \n"

${SCRIPTS_ROOT}/init_pwd.sh

echo -e "${PURPLE_PREFIX}=========> env_check.sh ${COLOR_SUFFIX} \n"

${SCRIPTS_ROOT}/env_check.sh

# Replace local IP address with the public IP address in .env file
if [ $API_URL == "http://127.0.0.1:10002/object/" ]; then
    sed -i "s/127.0.0.1/${internet_ip}/" ${OPENIM_ROOT}/.env
fi

if [ $MINIO_ENDPOINT == "http://127.0.0.1:10005" ]; then
    sed -i "s/127.0.0.1/${internet_ip}/" ${OPENIM_ROOT}/.env
fi 

# Go back to the previous directory
cd ${OPENIM_ROOT}

# Check if docker-compose command is available
if command -v docker-compose &> /dev/null
then
    docker-compose up -d
else
    docker compose up -d
fi

${SCRIPTS_ROOT}/docker_check_service.sh