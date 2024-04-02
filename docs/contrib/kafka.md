# OpenIM Kafka Guide

This document aims to provide a set of concise guidelines to help you quickly install and use Kafka through Docker Compose.

## Installing Kafka

With the Docker Compose script provided by OpenIM, you can easily install Kafka. Use the following command to start Kafka:

```bash
docker compose up -d
```

After executing this command, Kafka will be installed and started. You can confirm the Kafka container is running with the following command:

```bash
docker ps | grep kafka
```

The output of this command, as shown below, displays the status information of the Kafka container:

```
be416b5a0851   bitnami/kafka:3.5.1 "/opt/bitnami/script…"   3 days ago   Up 2 days   9092/tcp, 0.0.0.0:19094->9094/tcp, :::19094->9094/tcp   kafka
```

### References

- Official Docker installation documentation: [Click here](http://events.jianshu.io/p/b60afa35303a)
- Detailed installation guide: [Tutorial on Towards Data Science](https://towardsdatascience.com/how-to-install-apache-kafka-using-docker-the-easy-way-4ceb00817d8b)

## Using Kafka

### Entering the Kafka Container

To execute Kafka commands, you first need to enter the Kafka container. Use the following command:

```bash
docker exec -it kafka bash
```

### Kafka Command Tools

Inside the Kafka container, you can use various command-line tools to manage Kafka. These tools include but are not limited to:

- `kafka-topics.sh`: For creating, deleting, listing, or altering topics.
- `kafka-console-producer.sh`: Allows sending messages to a specified topic from the command line.
- `kafka-console-consumer.sh`: Allows reading messages from the command line, with the ability to specify topics.
- `kafka-consumer-groups.sh`: For managing consumer group information.

### Kafka Client Tool Installation

For easier Kafka management, you can install Kafka client tools. If you installed Kafka through OpenIM's Docker Compose, you can install the Kafka client tools with the following command:

```bash
make install.kafkactl
```

### Automatic Topic Creation

When installing Kafka through OpenIM's Docker Compose method, OpenIM automatically creates the following topics:

- `latestMsgToRedis`
- `msgToPush`
- `offlineMsgToMongoMysql`

These topics are created using the `scripts/create-topic.sh` script. The script waits for Kafka to be ready before executing the commands to create topics:

```bash
# Wait for Kafka to be ready
until /opt/bitnami/kafka/bin/kafka-topics.sh --list --bootstrap-server localhost:9092; do
  echo "Waiting for Kafka to be ready..."
  sleep 2
done

# Create topics
/opt/bitnami/kafka/bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 8 --topic latestMsgToRedis
/opt/bitnami/kafka/bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 8 --topic msgToPush
/opt/bitnami/kafka/bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 8 --topic offlineMsgToMongoMysql

echo "Topics created."
```

The optimized and expanded documentation further details some basic commands for operations inside the Kafka container, as well as basic commands for managing Kafka using `kafkactl`. Here is a more detailed guide.


## Basic Commands in the Kafka Container

### Listing Topics

To list all existing topics, you can use the following command:

```bash
kafka-topics.sh --list --bootstrap-server localhost:9092
```

### Creating a New Topic

When creating a new topic, you can specify the number of partitions and the replication factor. Here is the command to create a new topic:

```bash
kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic your_topic_name
```

### Producing Messages

To send messages to a specific topic, you can use the producer command. The following command prompts you to enter messages, which are sent to the specified topic with each press of the Enter key:

```bash
kafka-console-producer.sh --broker-list localhost:9092 --topic your_topic_name
```

### Consuming Messages

To read messages from a specific topic, you can use the consumer command. The following command reads new messages from the specified topic and outputs them on the console:

```bash
kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic your_topic_name --from-beginning
```

The `

--from-beginning` parameter reads messages from the beginning of the topic. If this parameter is omitted, only new messages will be read.


## Basic Commands Using `kafkactl`

`kafkactl` is a command-line tool for managing and operating Kafka clusters. It offers a more modern way to interact with Kafka.

### Listing Topics

To list all topics, you can use:

```bash
kafkactl get topics
```

### Creating a New Topic

To create a new topic with `kafkactl`, use:

```bash
kafkactl create topic your_topic_name --partitions 1 --replication-factor 1
```

### Producing Messages

To send messages to a topic, you can use:

```bash
kafkactl produce your_topic_name --value "your message"
```

Here, `"your message"` is the content of the message you want to send.

### Consuming Messages

To consume messages from a topic, use:

```bash
kafkactl consume your_topic_name --from-beginning
```

Again, the `--from-beginning` parameter will start consuming messages from the beginning of the topic. If you do not wish to start from the beginning, you can omit this parameter.