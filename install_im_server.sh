#!/usr/bin/env bash
internet_ip=`curl ifconfig.me -s`
echo $internet_ip

source .env
echo $MINIO_ENDPOINT
if [ $MINIO_ENDPOINT == "http://127.0.0.1:10005" ]; then
	sed -i "s/127.0.0.1/${internet_ip}/" .env 

fi
	
cd scripts ;
chmod +x *.sh ;
./init_pwd.sh
./env_check.sh;
cd .. ;

if command -v docker-compose &> /dev/null
then
    docker-compose up -d ;
else
    docker compose up -d ;
fi


cd scripts ;
./docker_check_service.sh
