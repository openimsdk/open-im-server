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
OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

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

do_sth() {
    #Main program to run
    echo "++++++++++++++++++++++++"
    sleep 5
    echo "++++++++++++++++++++++++"

    sleep 10
}

#Import environment variables
source .env

#Get the public IP address of the local machine
internet_ip=$(curl ifconfig.me -s)
echo -e "\033[1;34mCurrent public IP address: ${internet_ip}\033[0m\n"

#If MINIO_ENDPOINT is "http://127.0.0.1:10005", replace it with the current public IP address
if [[ $MINIO_ENDPOINT == "http://127.0.0.1:10005" ]]; then
    sed -i "s/127.0.0.1/${internet_ip}/" .env
fi

do_progress_bar() {
    local duration=${1}
    local max_progress=20
    local current_progress=0

    while true; do
        ((current_progress++))
        if [[ $current_progress -gt $max_progress ]]; then
            break
        fi
        sleep "$duration"
        echo "=====> Progress: [${current_progress}/${max_progress}]"
    done
}

#Start Docker containers
start_docker_containers() {
    if command -v docker-compose >/dev/null 2>&1; then
        echo -e "\033[1;34mFound docker-compose command, starting docker containers...\033[0m\n"
        docker-compose -f ${OPENIM_ROOT}/${docker_compose_file_name} up -d
    else
        if command -v docker >/dev/null 2>&1; then
            echo -e "\033[1;34mFound docker command, starting docker containers...\033[0m\n"
            docker compose -f ${OPENIM_ROOT}/${docker_compose_file_name} up -d
        else
            echo -e "\033[1;31mFailed to find docker-compose or docker command, please make sure they are installed and configured correctly.\033[0m"
            return 1
        fi
    fi
}

#Execute scripts
setup_script() {
    chmod +x ${SCRIPTS_ROOT}/*.sh
    echo -e "\033[1;34m============>Executing init_pwd.sh script...\033[0m\n"
    ${SCRIPTS_ROOT}/init_pwd.sh
    echo -e "\033[1;34m============>Executing env_check.sh script...\033[0m\n"
    ${SCRIPTS_ROOT}/env_check.sh
}

setup_script &

#Start Docker containers (timeout 10 seconds)
start_docker_containers

docker_pid=$!
timeout 10s tail --pid=${docker_pid} -f /dev/null
docker_exit_code=$?

if [ $docker_exit_code -eq 0 ]; then
    echo -e "\033[1;32m============>Docker containers started successfully!\033[0m\n"
else
    echo -e "\033[1;31m============>Failed to start Docker containers, please check the environment configuration and dependencies.\033[0m\n"
    exit 1
fi

echo -e "\033[1;34m============>Executing docker_check_service.sh script...\033[0m\n"

#View running Docker containers
echo -e "\033[1;34m============>Viewing running Docker containers...\033[0m\n"
echo ""
docker ps

#Replace the progress bar section with the pv command
echo -e "\033[1;34m============>Starting progress bar...\033[0m\n"
do_progress_bar 0.5 | pv -l -s 20 > /dev/null
echo -e "\033[1;34m============>Progress bar completed.\033[0m\n"

#Execute the main program
do_sth
