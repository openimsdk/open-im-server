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


# --------------------------------------------------------------
# OpenIM Protoc Tool v1.0.0
# --------------------------------------------------------------
# OpenIM has released its custom Protoc tool version v1.0.0.
# This tool is customized to meet the specific needs of OpenIM and resides in its separate repository.
# It can be downloaded from the following link:
# https://github.com/OpenIMSDK/Open-IM-Protoc/releases/tag/v1.0.0
# 
# Download link (Windows): https://github.com/OpenIMSDK/Open-IM-Protoc/releases/download/v1.0.0/windows.zip
# Download link (Linux): https://github.com/OpenIMSDK/Open-IM-Protoc/releases/download/v1.0.0/linux.zip
# 
# Installation steps (taking Windows as an example):
# 1. Visit the above link and download the version suitable for Windows.
# 2. Extract the downloaded file.
# 3. Add the extracted tool to your PATH environment variable so that it can be run directly from the command line.
# 
# Note: The specific installation and usage instructions may vary based on the tool's actual implementation. It's advised to refer to official documentation.
# --------------------------------------------------------------

PROTOC_DOWNLOAD_URL="https://github.com/OpenIMSDK/Open-IM-Protoc/releases/download/v1.0.0/linux.zip"
DOWNLOAD_DIR="/tmp/openim-protoc"
INSTALL_DIR="/usr/local/bin"

function help_message {
    echo "Usage: ./install-protobuf.sh [option]"
    echo "Options:"
    echo "-i, --install       Install the OpenIM Protoc tool."
    echo "-u, --uninstall     Uninstall the OpenIM Protoc tool."
    echo "-r, --reinstall     Reinstall the OpenIM Protoc tool."
    echo "-c, --check         Check if the OpenIM Protoc tool is installed."
    echo "-h, --help          Display this help message."
}

function install_protobuf {
    echo "Installing OpenIM Protoc tool..."
    
    # Create temporary directory and download the zip file
    mkdir -p $DOWNLOAD_DIR
    wget $PROTOC_DOWNLOAD_URL -O $DOWNLOAD_DIR/linux.zip

    # Unzip the file
    unzip -o $DOWNLOAD_DIR/linux.zip -d $DOWNLOAD_DIR

    # Move binaries to the install directory and make them executable
    sudo cp $DOWNLOAD_DIR/linux/protoc $INSTALL_DIR/
    sudo cp $DOWNLOAD_DIR/linux/protoc-gen-go $INSTALL_DIR/
    sudo chmod +x $INSTALL_DIR/protoc
    sudo chmod +x $INSTALL_DIR/protoc-gen-go
    
    # Clean up
    rm -rf $DOWNLOAD_DIR

    echo "OpenIM Protoc tool installed successfully!"
}

function uninstall_protobuf {
    echo "Uninstalling OpenIM Protoc tool..."
    
    # Removing binaries from the install directory
    sudo rm -f $INSTALL_DIR/protoc
    sudo rm -f $INSTALL_DIR/protoc-gen-go

    echo "OpenIM Protoc tool uninstalled successfully!"
}

function reinstall_protobuf {
    echo "Reinstalling OpenIM Protoc tool..."
    uninstall_protobuf
    install_protobuf
}

function check_protobuf {
    echo "Checking for OpenIM Protoc tool installation..."
    
    which protoc > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "OpenIM Protoc tool is installed."
    else
        echo "OpenIM Protoc tool is not installed."
    fi
}

while [ "$1" != "" ]; do
    case $1 in
        -i | --install )    install_protobuf
                            ;;
        -u | --uninstall )  uninstall_protobuf
                            ;;
        -r | --reinstall )  reinstall_protobuf
                            ;;
        -c | --check )      check_protobuf
                            ;;
        -h | --help )       help_message
                            exit
                            ;;
        * )                 help_message
                            exit 1
    esac
    shift
done
