#!/bin/bash

# Check if the script is run as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run the script as root or use sudo."
  exit
fi

set -e
set -o pipefail
set -o noglob

# Color definitions
openim_color() {
    BLACK_PREFIX="\033[30m"  # Black prefix
    RED_PREFIX="\033[31m"  # Red prefix
    GREEN_PREFIX="\033[32m"  # Green prefix
    YELLOW_PREFIX="\033[33m"  # Yellow prefix
    BLUE_PREFIX="\033[34m"  # Blue prefix
    SKY_BLUE_PREFIX="\033[36m"  # Sky blue prefix
    WHITE_PREFIX="\033[37m"  # White prefix
    BOLD_PREFIX="\033[1m"  # Bold prefix
    UNDERLINE_PREFIX="\033[4m"  # Underline prefix
    ITALIC_PREFIX="\033[3m"  # Italic prefix
}


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

############### OpenIM Github ###############
# ... rest of the script ...

# OpenKF Repo
OWNER="OpenIMSDK"
REPO="Open-IM-Server"

# Update your Go version here
GO_VERSION="1.18"

# --- helper functions for logs ---
info()
{
    echo -e "[${GREEN_PREFIX}INFO${COLOR_SUFFIX}] " "$@"
}
warn()
{
    echo -e "[${YELLOW_PREFIX}WARN${COLOR_SUFFIX}] " "$@" >&2
}
fatal()
{
    echo -e "[${RED_PREFIX}ERROR${COLOR_SUFFIX}] " "$@" >&2
    exit 1
}

# Function to download and install Go if it's not already installed
install_go() {
    command -v go >/dev/null 2>&1
    if [[ $? -ne 0 ]]; then
        warn "Go is not installed. Installing now..."
        curl -LO "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"
        if [ $? -ne 0 ]; then
            fatal "Download failed! Please check your network connectivity."
        fi
        sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
        echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
        source ~/.bashrc
        info "Go has been installed successfully."
    else
        info "Go is already installed."
    fi
}

# Function for colored echo
color_echo() {
    COLOR=$1
    shift
    echo -e "${COLOR}===========> $* ${COLOR_SUFFIX}"
}

# Function to update and install necessary tools
install_tools() {
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

    info "Detected OS: $os"
    info "Updating system package repositories..."
    $update_cmd

    for tool in "${tools[@]}"; do
        if ! command -v $tool &> /dev/null; then
            warn "$tool is not installed. Installing now..."
            $install_cmd $tool
            info "$tool has been installed successfully."
        else
            info "$tool is already installed."
        fi
    done
}

############### OpenIM LOGO ###############
# Set text color to cyan for header and URL

openim_logo() {
    echo -e "\033[0;36m"

    # Display fancy ASCII Art logo
    print_with_delay '
##########################################################################
    ____                   _ _
    / __ \                 (_) |
    | |  | |_ __   ___  _ __ _| |_ _   _ _ __ ___  _ __
    | |  | | '"'"'_ \ / _ \| '"'"'__| | __| | | | '"'"'_ ` _ \| '"'"'_ \
    | |__| | |_) | (_) | |  | | |_| |_| | | | | | | |_) |
    \____/| .__/ \___/|_|  |_|\__|\__,_|_| |_| |_| .__/
        | |                                    | |
        |_|                                    |_|
##########################################################################
    ' 0.01

    # Display product URL
    print_with_delay "Discover more and contribute at: https://github.com/OpenIMSDK/Open-IM-Server" 0.01

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

    print_with_delay "Join our developer community on Slack: https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg" 0.01

    # Reset text color back to normal
    echo -e "\033[0m"   
}

# Main function to run the script
main() {
    openim_color

    install_tools
    install_go

    LATEST_TAG=$(curl -s "https://api.github.com/repos/$OWNER/$REPO/tags" | jq -r '.[0].name')
    MODIFIED_TAG=$(echo $LATEST_TAG | sed -r 's/(v3\.0\.)[1-9][0-9]*$/\10/g')
    TARBALL_URL="https://github.com/$OWNER/$REPO/archive/refs/tags/$MODIFIED_TAG.tar.gz"

    color_echo ${GREEN_PREFIX} "Starting the OpenIM automated one-click deployment script."

    color_echo ${GREEN_PREFIX} "Downloading OpenIM source code from $TARBALL_URL"
    curl -L -o "${MODIFIED_TAG}.tar.gz" $TARBALL_URL
    tar -xzvf "${MODIFIED_TAG}.tar.gz"
    cd "$REPO-$MODIFIED_TAG"
    
    openim_logo

    info "Source code downloaded and extracted to $REPO-$MODIFIED_TAG"

    # Add the logic to modify .env based on user input here

    warn "This command can only be executed once. It will modify the component passwords in docker-compose based on the PASSWORD variable in .env, and modify the component passwords in config/config.yaml. If the password in .env changes, you need to first execute docker-compose down; rm components -rf and then execute this command."

    color_echo ${GREEN_PREFIX} "Installing OpenIM"
    make --debug -j install V=1

    color_echo ${GREEN_PREFIX} "Checking installation"
    make --debug check

    color_echo ${GREEN_PREFIX} "OpenIM installation completed successfully. Happy chatting!"
}

main "$@"