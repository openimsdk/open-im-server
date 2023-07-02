#!/usr/bin/env bash
echo "docker-compose ps..........................."
cd ..
docker-compose ps


cd scripts
echo "check OpenIM................................"
./check_all.sh

