#!/bin/bash

echo "Welcome to the Open-IM-Server installation script."
echo "Please select an deploy option:"
echo "1. docker-compose install"
# echo "2. source code install"
# echo "3. source code install with docker-compose dependence"
echo "2. exit"

clear_openimlog() {
    rm -rf ./logs/*
}

is_path() {
  if [ -e "$1" ]; then   
    return 1
  else
    return 0
  fi
}

is_empyt() {
    if [ -z "$1" ]; then
    return 1
    else
    return 0
    fi
}

edit_config() {
    echo "is edit config.yaml?"
    echo "1. vi edit config"
    echo "2. do not edit config"
    read choice
    case $choice in
    1)
      vi config/config.yaml
    ;;
    2)
      echo "do not edit config"    
    ;;
    esac
}

edit_enterprise_config() {
    echo "is edit enterprise config.yaml?"
    echo "1. vi edit enterprise config"
    echo "2. do not edit enterprise config"
    read choice
    case $choice in
    1)
      vi dockker-compose_cfg/config.yaml
    ;;
    2)
      echo "do not edit enterprise config"    
    ;;
    esac
}

install_docker_compose() {
    echo "Please enter the installation path, default is $(pwd)"
    read install_path 
    is_empyt $install_path
    if [ $? -eq 1 ]; then 
        install_path="./"
    fi
    echo "Installing Open-IM-Server to $install_path..."
    is_path $install_path
    mkdir -p $install_path
    cd $install_path
    rm -rf ./Open-IM-Server
    git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive;
    cd ./Open-IM-Server
    git checkout errcode
    echo "git clone success"
    echo "Please enter the data directory, deault is $(pwd), press enter to use default"
    read DATA_DIR
    is_empyt $DATA_DIR
    if [ $? -eq 1 ]; then 
        DATA_DIR="."
    fi 
    
    echo "Please enter the user, deault is root, press enter to use default"
    read USER
    is_empyt $USER
    if [ $? -eq 1 ]; then 
        USER="root"
    fi 

    echo "Please enter the password, default is openIM123, press enter to use default"
    read PASSWORD
    is_empyt $PASSWORD
     if [ $? -eq 1 ]; then 
        PASSWORD="openIM123"
    fi 

    echo "Please enter the minio_endpoint, default will detect auto, press enter to use default:"
    read MINIO_ENDPOINT
    is_empyt $MINIO_ENDPOINT
    if [ $? -eq 1 ]; then
        internet_ip=`curl ifconfig.me -s`
        MINIO_ENDPOINT="http://${internet_ip}:10005"
    fi

    export MINIO_ENDPOINT
    export USER
    export PASSWORD
    export DATA_DIR

    edit_config
    edit_enterprise_config

    cd script;
    chmod +x *.sh;
    ./init_pwd.sh;
    ./env_check.sh;
    cd ..;
    docker-compose up -d;
    cd script;
    ./docker_check_service.sh;
}

read choice

case $choice in
  1)
    install_docker_compose
    ;;
  2)
    
    ;;
  3)
    ;;
  4)
    echo "Exiting installation script..."
    exit 0
    ;;
  *)
    echo "Invalid option, please try again."
    ;;
esac

