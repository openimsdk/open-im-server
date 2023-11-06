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
#
# https://gist.github.com/cubxxw/28f997f2c9aff408630b072f010c1d64
#

set -e
set -o pipefail

############################## OpenIM Github ##############################
# ... rest of the script ...

# TODO
# You can configure this script in three ways. 
# 1. First, set the variables in this column with more comments. 
# 2. The second is to pass an environment variable via a flag such as --help. 
# 3. The third way is to set the variable externally, or pass it in as an environment variable

# Default configuration for OpenIM Repo
# The OpenIM Repo settings can be customized according to your needs.

# OpenIM Repo owner, by default it's set to "OpenIMSDK". If you're using a different owner, replace accordingly.
OWNER="OpenIMSDK" 

# The repository name, by default it's "Open-IM-Server". If you're using a different repository, replace accordingly.
REPO="Open-IM-Server" 

# Version of Go you want to use, make sure it is compatible with your OpenIM-Server requirements.
# Default is 1.18, if you want to use a different version, replace accordingly.
GO_VERSION="1.20"

# Default HTTP_PORT is 80. If you want to use a different port, uncomment and replace the value.
# HTTP_PORT=80

# CPU core number for concurrent execution. By default it's determined automatically.
# Uncomment the next line if you want to set it manually.
# CPU=$(grep -c ^processor /proc/cpuinfo)

# By default, the script uses the latest tag from OpenIM-Server releases.
# If you want to use a specific tag, uncomment and replace "v3.0.0" with the desired tag.
# LATEST_TAG=v3.0.0

# Default OpenIM install directory is /tmp. If you want to use a different directory, uncomment and replace "/test".
# DOWNLOAD_OPENIM_DIR="/test"

# GitHub proxy settings. If you are using a proxy, uncomment and replace the empty field with your proxy URL.
PROXY=

# If you have a GitHub token, replace the empty field with your token.
GITHUB_TOKEN=

# Default user is "root". If you need to modify it, uncomment and replace accordingly.
# USER=root 

# Default password for redis, mysql, mongo, as well as accessSecret in config/config.yaml.
# Remember, it should be a combination of 8 or more numbers and letters. If you want to set a different password, uncomment and replace "openIM123".
# PASSWORD=openIM123

# Default endpoint for minio's external service IP and port. If you want to use a different endpoint, uncomment and replace.
# ENDPOINT=http://127.0.0.1:10005 

# Default API_URL, replace if necessary. 
# API_URL=http://127.0.0.1:10002/object/

# Default data directory. If you want to specify a different directory, uncomment and replace "./".
# DATA_DIR=./

############################## OpenIM Functions ##############################
# Install horizon of the script
#
# Pre-requisites:
#   - git
#   - make
#   - jq
#   - docker
#   - docker-compose
#   - go
#

# Check if the script is run as root
function check_isroot() {
    if [ "$EUID" -ne 0 ]; then
    fatal "Please run the script as root or use sudo."
    fi
}

# check if the current directory is a OpenIM git repository
function check_git_repo() {
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        # Inside a git repository
        for remote in $(git remote); do
            repo_url=$(git remote get-url $remote)
            if [[ $repo_url == "https://github.com/openimsdk/open-im-server.git" || \
                  $repo_url == "https://github.com/openimsdk/open-im-server" || \
                  $repo_url == "git@github.com:openimsdk/open-im-server.git" ]]; then
                # If it's OpenIMSDK repository
                info "Current directory is OpenIMSDK git repository."
                info "Executing installation directly."
                install_openim
                exit 0
            fi
            debug "Remote: $remote, URL: $repo_url"
        done
        # If it's not OpenIMSDK repository
        debug "Current directory is not OpenIMSDK git repository."
    fi
    info "Current directory is not a git repository."
}

# Function to update and install necessary tools
function install_tools() {
    info "Checking and installing necessary tools, about git, make, jq, docker, docker-compose."
    local tools=("git" "make" "jq" "docker" "docker-compose")
    local install_cmd update_cmd os

    if grep -qEi "debian|buntu|mint" /etc/os-release; then
        os="Ubuntu"
        install_cmd="sudo apt install -y"
        update_cmd="sudo apt update"
    elif grep -qEi "fedora|rhel" /etc/os-release; then
        os="CentOS"
        install_cmd="sudo yum install -y"
        update_cmd="sudo yum update"
    else
        fatal "Unsupported OS, please use Ubuntu or CentOS."
    fi

    debug "Detected OS: $os"
    info "Updating system package repositories..."
    $update_cmd

    for tool in "${tools[@]}"; do
        if ! command -v $tool &> /dev/null; then
            warn "$tool is not installed. Installing now..."
            $install_cmd $tool
            success "$tool has been installed successfully."
        else
            info "$tool is already installed."
        fi
    done
}

# Function to check if Docker and Docker Compose are installed
function check_docker() {
    if ! command -v docker &> /dev/null; then
        fatal "Docker is not installed. Please install Docker first."
    fi
    if ! command -v docker-compose &> /dev/null; then
        fatal "Docker Compose is not installed. Please install Docker Compose first."
    fi
}

# Function to download and install Go if it's not already installed
function install_go() {
    command -v go >/dev/null 2>&1
    # Determines if GO_VERSION is defined
    if [ -z "$GO_VERSION" ]; then
        GO_VERSION="1.20"
    fi

    if [[ $? -ne 0 ]]; then
        warn "Go is not installed. Installing now..."
        curl -LO "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"
        if [ $? -ne 0 ]; then
            fatal "Download failed! Please check your network connectivity."
        fi
        sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
        echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
        source ~/.bashrc
        success "Go has been installed successfully."
    else
        info "Go is already installed."
    fi
}

function download_source_code() {

    # If LATEST_TAG was not defined outside the function, get it here example: v3.0.1-beta.1
    if [ -z "$LATEST_TAG" ]; then
        LATEST_TAG=$(curl -s "https://api.github.com/repos/$OWNER/$REPO/tags" | jq -r '.[0].name')
    fi

    # If LATEST_TAG is still empty, set a default value
    local DEFAULT_TAG="v3.0.0"

    LATEST_TAG="${LATEST_TAG:-$DEFAULT_TAG}"

    debug "DEFAULT_TAG: $DEFAULT_TAG"
    info "Use OpenIM Version LATEST_TAG: $LATEST_TAG"

    # If MODIFIED_TAG was not defined outside the function, modify it here,example: 3.0.1-beta.1
    if [ -z "$MODIFIED_TAG" ]; then
        MODIFIED_TAG=$(echo $LATEST_TAG | sed 's/v//')
    fi

    # If MODIFIED_TAG is still empty, set a default value
    local DEFAULT_MODIFIED_TAG="${DEFAULT_TAG#v}" 
    MODIFIED_TAG="${MODIFIED_TAG:-$DEFAULT_MODIFIED_TAG}"
    
    debug "MODIFIED_TAG: $MODIFIED_TAG"

    # Construct the tarball URL
    TARBALL_URL="${PROXY}https://github.com/$OWNER/$REPO/archive/refs/tags/$LATEST_TAG.tar.gz"

    info "Downloaded OpenIM TARBALL_URL: $TARBALL_URL"

    info "Starting the OpenIM automated one-click deployment script."

    # Set the download and extract directory to /tmp
    if [ -z "$DOWNLOAD_OPENIM_DIR" ]; then
        DOWNLOAD_OPENIM_DIR="/tmp"
    fi

    # Check if /tmp directory exists
    if [ ! -d "$DOWNLOAD_OPENIM_DIR" ]; then
        warn "$DOWNLOAD_OPENIM_DIR does not exist. Creating it..."
        mkdir -p "$DOWNLOAD_OPENIM_DIR"
    fi

    info "Downloading OpenIM source code from $TARBALL_URL to $DOWNLOAD_OPENIM_DIR"
    
    curl -L -o "${DOWNLOAD_OPENIM_DIR}/${MODIFIED_TAG}.tar.gz" $TARBALL_URL

    tar -xzvf "${DOWNLOAD_OPENIM_DIR}/${MODIFIED_TAG}.tar.gz" -C "$DOWNLOAD_OPENIM_DIR"
    cd "$DOWNLOAD_OPENIM_DIR/$REPO-$MODIFIED_TAG"
    git init && git add . && git commit -m "init"  --no-verify

    success "Source code downloaded and extracted to $REPO-$MODIFIED_TAG"
}

function set_openim_env() {
    warn "This command can only be executed once. It will modify the component passwords in docker-compose based on the PASSWORD variable in .env, and modify the component passwords in config/config.yaml. If the password in .env changes, you need to first execute docker-compose down; rm components -rf and then execute this command."
    # Set default values for user input
    # If the USER environment variable is not set, it defaults to 'root'
    if [ -z "$USER" ]; then
        USER="root"
        debug "USER is not set. Defaulting to 'root'."
    fi

    # If the PASSWORD environment variable is not set, it defaults to 'openIM123'
    # This password applies to redis, mysql, mongo, as well as accessSecret in config/config.yaml
    if [ -z "$PASSWORD" ]; then
        PASSWORD="openIM123"
        debug "PASSWORD is not set. Defaulting to 'openIM123'."
    fi

    # If the ENDPOINT environment variable is not set, it defaults to 'http://127.0.0.1:10005'
    # This is minio's external service IP and port, or it could be a domain like storage.xx.xx
    # The app must be able to access this IP and port or domain
    if [ -z "$ENDPOINT" ]; then
        ENDPOINT="http://127.0.0.1:10005"
        debug "ENDPOINT is not set. Defaulting to 'http://127.0.0.1:10005'."
    fi

    # If the API_URL environment variable is not set, it defaults to 'http://127.0.0.1:10002/object/'
    # The app must be able to access this IP and port or domain
    if [ -z "$API_URL" ]; then
        API_URL="http://127.0.0.1:10002/object/"
        debug "API_URL is not set. Defaulting to 'http://127.0.0.1:10002/object/'."
    fi

    # If the DATA_DIR environment variable is not set, it defaults to the current directory './'
    # This can be set to a directory with large disk space
    if [ -z "$DATA_DIR" ]; then
        DATA_DIR="./"
        debug "DATA_DIR is not set. Defaulting to './'."
    fi
}

function install_openim() {
    info "Installing OpenIM"
    make -j${CPU} install V=1

    info "Checking installation"
    make check

    success "OpenIM installation completed successfully. Happy chatting!"
}

############################## OpenIM Help ##############################

# Function to display help message
function cmd_help() {
    openim_color
    color_echo ${BRIGHT_GREEN_PREFIX} "Usage: $0 [options]"
    color_echo ${BRIGHT_GREEN_PREFIX} "Options:"
    echo
    color_echo ${BLUE_PREFIX} "-i,  --install       ${CYAN_PREFIX}Execute the installation logic of the script${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-u,  --user          ${CYAN_PREFIX}set user (default: root)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-p,  --password      ${CYAN_PREFIX}set password (default: openIM123)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-e,  --endpoint      ${CYAN_PREFIX}set endpoint (default: http://127.0.0.1:10005)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-a,  --api           ${CYAN_PREFIX}set API URL (default: http://127.0.0.1:10002/object/)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-d,  --directory     ${CYAN_PREFIX}set directory for large disk space (default: ./)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-h,  --help          ${CYAN_PREFIX}display this help message and exit${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-cn, --china         ${CYAN_PREFIX}set to use the Chinese domestic proxy${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-t,  --tag           ${CYAN_PREFIX}specify the tag (default option, set to latest if not specified)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-r,  --release       ${CYAN_PREFIX}specify the release branch (cannot be used with the tag option)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-gt, --github-token  ${CYAN_PREFIX}set the GITHUB_TOKEN (default: not set)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "-g,  --go-version    ${CYAN_PREFIX}set the Go language version (default: GO_VERSION=\"1.20\")${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "--install-dir        ${CYAN_PREFIX}set the OpenIM installation directory (default: /tmp)${COLOR_SUFFIX}"
    color_echo ${BLUE_PREFIX} "--cpu                ${CYAN_PREFIX}set the number of concurrent processes${COLOR_SUFFIX}"
    echo
    color_echo ${RED_PREFIX} "Note: Only one of the -t/--tag or -r/--release options can be used at a time.${COLOR_SUFFIX}"
    color_echo ${RED_PREFIX} "If both are used or none of them are used, the -t/--tag option will be prioritized.${COLOR_SUFFIX}"
    echo
    exit 1
}

function parseinput() {
    # set default values
    # USER=root
    # PASSWORD=openIM123
    # ENDPOINT=http://127.0.0.1:10005
    # API=http://127.0.0.1:10002/object/
    # DIRECTORY=./
    # CHINA=false
    # TAG=latest
    # RELEASE=""
    # GO_VERSION=1.20
    # INSTALL_DIR=/tmp
    # GITHUB_TOKEN=""
    # CPU=$(nproc)

    if [ $# -eq 0 ]; then
        cmd_help
        exit 1
    fi

    while [ $# -gt 0 ]; do
        case $1 in
            -h|--help)
                cmd_help
                exit
                ;;
            -u|--user)
                shift
                USER=$1
                ;;
            -p|--password)
                shift
                PASSWORD=$1
                ;;
            -e|--endpoint)
                shift
                ENDPOINT=$1
                ;;
            -a|--api)
                shift
                API=$1
                ;;
            -d|--directory)
                shift
                DIRECTORY=$1
                ;;
            -cn|--china)
                CHINA=true
                ;;
            -t|--tag)
                shift
                TAG=$1
                ;;
            -r|--release)
                shift
                RELEASE=$1
                ;;
            -g|--go-version)
                shift
                GO_VERSION=$1
                ;;
            --install-dir)
                shift
                INSTALL_DIR=$1
                ;;
            -gt|--github-token)
                shift
                GITHUB_TOKEN=$1
                ;;
            --cpu)
                shift
                CPU=$1
                ;;
            -i|--install)
                openim_main
                exit
                ;;
            *)
                echo "Unknown option: $1"
                cmd_help
                exit 1
                ;;
        esac
        shift
    done
}

############################## OpenIM LOG ##############################
# Set text color to cyan for header and URL
print_with_delay() {
  text="$1"
  delay="$2"

  for i in $(seq 0 $((${#text}-1))); do
    printf "${text:$i:1}"
    sleep $delay
  done
  printf "\n"
}

print_progress() {
  total="$1"
  delay="$2"

  printf "["
  for i in $(seq 1 $total); do
    printf "#"
    sleep $delay
  done
  printf "]\n"
}

# Function for colored echo
color_echo() {
    COLOR=$1
    shift
    echo -e "${COLOR} $* ${COLOR_SUFFIX}"
}

# Color definitions
function openim_color() {
    COLOR_SUFFIX="\033[0m"      # End all colors and special effects

    BLACK_PREFIX="\033[30m"     # Black prefix
    RED_PREFIX="\033[31m"       # Red prefix
    GREEN_PREFIX="\033[32m"     # Green prefix
    YELLOW_PREFIX="\033[33m"    # Yellow prefix
    BLUE_PREFIX="\033[34m"      # Blue prefix
    SKY_BLUE_PREFIX="\033[36m"  # Sky blue prefix
    WHITE_PREFIX="\033[37m"     # White prefix
    BOLD_PREFIX="\033[1m"       # Bold prefix
    UNDERLINE_PREFIX="\033[4m"  # Underline prefix
    ITALIC_PREFIX="\033[3m"     # Italic prefix
    BRIGHT_GREEN_PREFIX='\033[1;32m' # Bright green prefix

    CYAN_PREFIX="\033[0;36m"     # Cyan prefix
}

# --- helper functions for logs ---
info() {
    echo -e "[${GREEN_PREFIX}INFO${COLOR_SUFFIX}] " "$@"
}
warn() {
    echo -e "[${YELLOW_PREFIX}WARN${COLOR_SUFFIX}] " "$@" >&2
}
fatal() {
    echo -e "[${RED_PREFIX}ERROR${COLOR_SUFFIX}] " "$@" >&2
    exit 1
}
debug() {
    echo -e "[${BLUE_PREFIX}DEBUG${COLOR_SUFFIX}]===> " "$@"
}
success() {
    echo -e "${BRIGHT_GREEN_PREFIX}=== [SUCCESS] ===${COLOR_SUFFIX}\n=> " "$@"
}

function openim_logo() {
    # Set text color to cyan for header and URL
    echo -e "\033[0;36m"

    # Display fancy ASCII Art logo
    # look http://patorjk.com/software/taag/#p=display&h=1&v=1&f=Doh&t=OpenIM
    print_with_delay '
                                                                                                                      
                                                                                                                      
     OOOOOOOOO                                                               IIIIIIIIIIMMMMMMMM               MMMMMMMM
   OO:::::::::OO                                                             I::::::::IM:::::::M             M:::::::M
 OO:::::::::::::OO                                                           I::::::::IM::::::::M           M::::::::M
O:::::::OOO:::::::O                                                          II::::::IIM:::::::::M         M:::::::::M
O::::::O   O::::::Oppppp   ppppppppp       eeeeeeeeeeee    nnnn  nnnnnnnn      I::::I  M::::::::::M       M::::::::::M
O:::::O     O:::::Op::::ppp:::::::::p    ee::::::::::::ee  n:::nn::::::::nn    I::::I  M:::::::::::M     M:::::::::::M
O:::::O     O:::::Op:::::::::::::::::p  e::::::eeeee:::::een::::::::::::::nn   I::::I  M:::::::M::::M   M::::M:::::::M
O:::::O     O:::::Opp::::::ppppp::::::pe::::::e     e:::::enn:::::::::::::::n  I::::I  M::::::M M::::M M::::M M::::::M
O:::::O     O:::::O p:::::p     p:::::pe:::::::eeeee::::::e  n:::::nnnn:::::n  I::::I  M::::::M  M::::M::::M  M::::::M
O:::::O     O:::::O p:::::p     p:::::pe:::::::::::::::::e   n::::n    n::::n  I::::I  M::::::M   M:::::::M   M::::::M
O:::::O     O:::::O p:::::p     p:::::pe::::::eeeeeeeeeee    n::::n    n::::n  I::::I  M::::::M    M:::::M    M::::::M
O::::::O   O::::::O p:::::p    p::::::pe:::::::e             n::::n    n::::n  I::::I  M::::::M     MMMMM     M::::::M
O:::::::OOO:::::::O p:::::ppppp:::::::pe::::::::e            n::::n    n::::nII::::::IIM::::::M               M::::::M
 OO:::::::::::::OO  p::::::::::::::::p  e::::::::eeeeeeee    n::::n    n::::nI::::::::IM::::::M               M::::::M
   OO:::::::::OO    p::::::::::::::pp    ee:::::::::::::e    n::::n    n::::nI::::::::IM::::::M               M::::::M
     OOOOOOOOO      p::::::pppppppp        eeeeeeeeeeeeee    nnnnnn    nnnnnnIIIIIIIIIIMMMMMMMM               MMMMMMMM
                    p:::::p                                                                                           
                    p:::::p                                                                                           
                   p:::::::p                                                                                          
                   p:::::::p                                                                                          
                   p:::::::p                                                                                          
                   ppppppppp                                                                                          
                                                                                                                      
    ' 0.0001

    # Display product URL
    print_with_delay "Discover more and contribute at: https://github.com/openimsdk/open-im-server" 0.01

    # Reset text color back to normal
    echo -e "\033[0m"

    # Set text color to green for product description
    echo -e "\033[1;32m"

    print_with_delay "Open-IM-Server: Reinventing Instant Messaging" 0.01
    print_progress 50 0.02

    print_with_delay "Open-IM-Server is not just a product; it's a revolution. It's about bringing the power of seamless, real-time messaging to your fingertips. And it's about joining a global community of developers, dedicated to pushing the boundaries of what's possible." 0.01

    print_progress 50 0.02

    # Reset text color back to normal
    echo -e "\033[0m"

    # Set text color to yellow for the Slack link
    echo -e "\033[1;33m"

    print_with_delay "Join our developer community on Slack: https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q" 0.01

    # Reset text color back to normal
    echo -e "\033[0m"
}

# Main function to run the script
function openim_main() {
    check_git_repo
    check_isroot
    openim_color
    install_tools
    check_docker
    install_go
    download_source_code
    set_openim_env
    install_openim
    openim_logo

}

parseinput "$@"