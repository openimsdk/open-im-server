#!/usr/bin/env bash
echo "docker-compose ps..........................."
cd ..

if command -v docker-compose &> /dev/null
then
    docker-compose ps
else
    docker compose ps
fi



cd scripts
echo "check OpenIM................................"

sleep 30
./check_all.sh

