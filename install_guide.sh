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

is_directory_exists() {
  if [ -d "$1" ]; then
    return 1
  else
    return 0 
  fi
}

edit_config() {
    echo "Is edit config.yaml?"
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
    echo "Is edit enterprise config.yaml?"
    echo "1. vi edit enterprise config"
    echo "2. do not edit enterprise config"
    read choice
    case $choice in
    1)
      vi docker-compose_cfg/config.yaml
    ;;
    2)
      echo "Do not edit enterprise config"    
    ;;
    esac
}

install_docker_compose() {
    echo "Please input the installation path, default is $(pwd)/Open-IM-Server, press enter to use default"
    read install_path 
    is_empyt $install_path
    if [ $? -eq 1 ]; then 
        install_path="."
    fi
    echo "Installing Open-IM-Server to ${install_path}/Open-IM-Server..."
    is_path $install_path
    mkdir -p $install_path
    cd $install_path
    is_directory_exists "${install_path}/Open-IM-Server"
    if [ $? -eq 1 ]; then
        echo "Error: Directory $install_path/Open-IM-Server exist, please ensure your path"
        echo "1. delete the directory and install"
        echo "2. exit"
        read choice
        case $choice in
        1)
          rm -rf "${install_path}/Open-IM-Server"
        ;;
        2)
          exit 1
        ;;
        esac
    fi
    rm -rf ./Open-IM-Server
    git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive;
    cd ./Open-IM-Server
    git checkout errcode
    echo "======== git clone success ========"
    source .env
    echo "Please input the components data directory, deault is ${DATA_DIR} press enter to use default"
    read NEW_DATA_DIR
    is_empyt $NEW_DATA_DIR
    if [ $? -eq 0 ]; then 
        DATA_DIR=$NEW_DATA_DIR
    fi 
    echo "Please input the user, deault is root, press enter to use default"
    read NEW_USER
    is_empyt $NEW_USER
    if [ $? -eq 0 ]; then 
        USER=$NEW_USER
    fi 

    echo "Please input the password, default is openIM123, press enter to use default"
    read NEW_PASSWORD
    is_empyt $NEW_PASSWORD
     if [ $? -eq 0 ]; then 
        PASSWORD=$NEW_PASSWORD
    fi 

    echo "Please input the minio_endpoint, default will detect auto, press enter to use default:"
    read NEW_MINIO_ENDPOINT
    is_empyt $NEW_MINIO_ENDPOINT
    if [ $? -eq 1 ]; then
        internet_ip=`curl ifconfig.me -s`
        MINIO_ENDPOINT="http://${internet_ip}:10005"
    else 
        MINIO_ENDPOINT=$NEW_MINIO_ENDPOINT  
    fi

    export MINIO_ENDPOINT
    export USER
    export PASSWORD
    export DATA_DIR
    sed -i "s/^MINIO_ENDPOINT=.*/MINIO_ENDPOINT=$MINIO_ENDPOINT/" .env
    sed -i "s/^USER=.*/USER=$USER/" .env
    sed -i "s/^PASSWORD=.*/PASSWORD=$PASSWORD/" .env
    sed -i "s/^DATA_DIR=.*/DATA_DIR=$DATA_DIR/" .env


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

