#!/usr/bin/env bash

#fixme The 10 second delay to start the project is for the docker-compose one-click to start openIM when the infrastructure dependencies are not started
sleep 10

./start_all.sh

sleep 15

#fixme prevents the openIM service exit after execution in the docker container
tail -f /dev/null
