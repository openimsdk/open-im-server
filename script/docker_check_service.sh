#!/usr/bin/env bash
echo "docker-compose ps..........................."
docker-compose ps

echo "check OpenIM, waiting 30s...................."
sleep 30

echo "check OpenIM................................"
./check_all.sh

