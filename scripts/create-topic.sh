#!/usr/bin/env bash
# Wait for Kafka to be ready

KAFKA_SERVER=kafka-service:9092

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


topics=("toRedis" "toMongo" "toPush" "toOfflinePush")
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
