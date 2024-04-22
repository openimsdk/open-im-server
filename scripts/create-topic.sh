#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Wait for Kafka to be ready

KAFKA_SERVER=localhost:9092

MAX_ATTEMPTS=300
attempt_num=1

echo "Waiting for Kafka to be ready..."

until /opt/bitnami/kafka/bin/kafka-topics.sh --list --bootstrap-server $KAFKA_SERVER; do
  echo "Attempt $attempt_num of $MAX_ATTEMPTS: Kafka not ready yet..."
  if [ $attempt_num -eq $MAX_ATTEMPTS ]; then
    echo "Kafka not ready after $MAX_ATTEMPTS attempts, exiting"
    exit 1
  fi
  attempt_num=$((attempt_num+1))
  sleep 1
done

echo "Kafka is ready. Creating topics..."


topics=("toRedis" "toMongo" "toPush")
partitions=8
replicationFactor=1

for topic in "${topics[@]}"; do
  if /opt/bitnami/kafka/bin/kafka-topics.sh --create \
    --bootstrap-server $KAFKA_SERVER \
    --replication-factor $replicationFactor \
    --partitions $partitions \
    --topic $topic
  then
    echo "Topic $topic created."
  else
    echo "Failed to create topic $topic."
  fi
done

echo "All topics created."
