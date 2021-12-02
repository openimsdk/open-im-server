#!/usr/bin/env bash
echo "docker-compose ps..........................."
docker-compose ps

sleep 30

echo "check OpenIM................................."
./check_all.sh

