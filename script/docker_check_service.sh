#!/usr/bin/env bash
echo "docker-compose ps..........................."
docker-compose ps

sleep 20

echo "check OpenIM................................."
./check_all.sh

