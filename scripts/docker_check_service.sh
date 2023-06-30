#!/usr/bin/env bash
echo "docker-compose ps..........................."
docker-compose ps

echo "check OpenIM, waiting 30s...................."
sleep 60

echo "check OpenIM................................"
./check_all.sh
# chmod +x ./enterprise/*.sh
# ./enterprise/check_all.sh

