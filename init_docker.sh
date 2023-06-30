#!/usr/bin/env bash

cd scripts ;
chmod +x *.sh ;
./env_check.sh;
cd .. ;
docker-compose up -d;
cd scripts ;
./docker_check_service.sh
