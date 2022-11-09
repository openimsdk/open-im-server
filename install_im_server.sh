#!/usr/bin/env bash

cd script ;
chmod +x *.sh ;
./env_check.sh;
cd .. ;
docker-compose up -d;
echo "upppp...."
cd script ;
./docker_check_service.sh
