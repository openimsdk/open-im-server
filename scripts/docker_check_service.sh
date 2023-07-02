#!/usr/bin/env bash
echo "docker-compose ps..........................."
cd ..
docker-compose ps


cd scripts
echo "check OpenIM................................"

sleep 30
./check_all.sh

