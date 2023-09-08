# OpenIM enviroment


## How to change the configuration


**Modify the configuration files:**

Three ways to modify the configuration:

#### **1. Recommended using environment variables:**

```bash
export PASSWORD="openIM123" # Set password
export USER="root" # Set username
# Choose chat version and server version https://github.com/openimsdk/open-im-server/blob/main/docs/conversions/images.md, eg: main, release-v*.*
export CHAT_BRANCH="main"
export SERVER_BRANCH="main"
#... Other environment variables
# MONGO_USERNAME: This sets the MongoDB username
# MONGO_PASSWORD: Set the MongoDB password
# MONGO_DATABASE: Sets the MongoDB database name
# MINIO_ENDPOINT: set the MinIO service address
# API_URL: under network environment, set OpenIM Server API address
export API_URL="http://127.0.0.1:10002"
```

Next, update the configuration using `make init`:

```bash
make init
```

#### **2. Modify the automation script:**

```bash
scripts/install/environment.sh
```

Next, update the configuration using `make init`:

```bash
make init
```

#### 3. Modify `config.yaml` and `.env` files (but will be overwritten when using `make init` again).

The `config/config.yaml` file has detailed configuration instructions for the storage components.


The config file is available via [environment.sh](https://github.com/openimsdk/open-im-server/blob/main/scripts/install/environment.sh) configuration [openim.yaml](https://github.com/openimsdk/open-im-server/blob/main/deployments/templates/openim.yaml) template, and then through the `make init` to automatically generate a new configuration.


## Environment variable

By setting the environment variable below, You can then refresh the configuration using `make init` or `./scripts/init-config.sh`

##### MINIO

+ [MINIO DOCS](https://min.io/docs/minio/kubernetes/upstream/index.html)

apiURL is the address of the api, the access address of the app, use s3 must be configured

#### Overview

MinIO is an object storage server that is API compatible with Amazon S3. It's best suited for storing unstructured data such as photos, videos, log files, backups, and container/VM images. In this guide, we'll walk through the process of configuring MinIO with custom settings.

#### Default Configuration

Configuration can be achieved by modifying the default variables in the `./scripts/install/environment.sh` file. However, for more flexibility and dynamic adjustments, setting environment variables is recommended.

#### Setting Up the Environment Variables

##### IP Configuration

By default, the system generates the public IP of the machine. To manually set a public or local IP address, use:

```bash
export IP=127.0.0.1
```

##### API URL

This is the address your application uses to communicate with MinIO. By default, it uses the public IP. However, you can adjust it to a public domain or another IP.

```bash
export API_URL=127.0.0.1:10002
```

##### MinIO Endpoint Configuration

This is the primary address MinIO uses for communications:

```bash
export MINIO_ENDPOINT="127.0.0.1"
```

##### MinIO Sign Endpoint

For direct external access to stored content:

```bash
export MINIO_SIGN_ENDPOINT=127.0.0.1:10005
```

##### Modifying MinIO's Port

If you need to adjust MinIO's port from the default:

```bash
export MINIO_PORT="10005"
```

#### Applying the Configuration

After setting your desired environment variables, restart the MinIO server to apply the changes.

#### Verification

It's crucial to verify the configurations by checking the connectivity between your application and MinIO using the set API URL and ensuring that the data can be directly accessed using the `signEndpoint`.


## Configuration Details

###### Zookeeper

- **Purpose**: Used for RPC service discovery and registration, cluster support.
  
    ```bash
    zookeeper:
      schema: openim                          # Not recommended to modify
      address: [ 127.0.0.1:2181 ]            # Address
      username:                               # Username
      password:                               # Password
    ```

###### MySQL

- **Purpose**: Used for storing users, relationships, and groups. Supports master-slave database.

    ```bash
    mysql:
      address: [ 127.0.0.1:13306 ]            # Address
      username: root                          # Username
      password: openIM123                     # Password
      database: openIM_v2                     # Not recommended to modify
      maxOpenConn: 1000                       # Maximum connection
      maxIdleConn: 100                        # Maximum idle connection
      maxLifeTime: 60                         # Max time a connection can be reused (seconds)
      logLevel: 4                             # Log level (1=silent, 2=error, 3=warn, 4=info)
      slowThreshold: 500                      # Slow statement threshold (milliseconds)
    ```

###### Mongo

- **Purpose**: Used for storing offline messages. Supports mongo sharded clusters.

    ```bash
    mongo:
      uri:                                    # Use this value directly if not empty
      address: [ 127.0.0.1:37017 ]            # Address
      database: openIM                        # Default mongo db
      username: root                          # Username
      password: openIM123                     # Password
      maxPoolSize: 100                        # Maximum connections
    ```

###### Redis

- **Purpose**: Used for storing message sequence numbers, latest messages, user tokens, and MySQL cache. Supports cluster deployment.

    ```bash
    redis:
      address: [ 127.0.0.1:16379 ]            # Address
      username:                               # Username
      password: openIM123                     # Password
    ```

###### Kafka

- **Purpose**: Used for message queues for decoupling. Supports cluster deployment.

    ```bash
    kafka:
      username:                               # Username
      password:                               # Password
      addr: [ 127.0.0.1:9092 ]                # Address
      topics:
        latestMsgToRedis: "latestMsgToRedis"
        offlineMsgToMongo: "offlineMsgToMongoMysql"
        msgToPush: "msgToPush"
        msgToModify: "msgToModify"
      consumerGroupID:
        msgToRedis: redis
        msgToMongo: mongo
        msgToMySql: mysql
        msgToPush: push
        msgToModify: modify
    ```





## Config options

...

## Use the default values

A method to revert to the default value:

```bash
export IP=127.0.0.1
```
