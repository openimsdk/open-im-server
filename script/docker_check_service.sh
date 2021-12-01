#!/usr/bin/env bash
echo "docker-compose ps....................."
docker-compose ps

sleep 10

echo "check OpenIM.........................."
./check_all.sh

