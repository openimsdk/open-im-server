#!/usr/bin/env bash
internet_ip=`curl ifconfig.me -s`
echo $internet_ip

source .env
echo $MINIO_ENDPOINT
if [ $MINIO_ENDPOINT == "http://127.0.0.1:10005" ]; then
	sed -i "s/127.0.0.1/${internet_ip}/" .env

fi

cd script ;
chmod +x *.sh ;
./init_pwd.sh
./env_check.sh;
cd .. ;
docker-compose -f im-compose.yaml up -d
docker ps



