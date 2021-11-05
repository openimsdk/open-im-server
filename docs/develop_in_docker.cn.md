
## 在 Docker 容器中开发

### 运行服务容器

```
# 这个会在运行容器的时候，会占用大量 CPU，导致容器启动失败
docker-compose -f docker-compose.dev.yaml up -d
```

### 检查容器运行状态

```
docker-compose -f docker-compose.dev.yaml ps
```

应该能看到以下结果

```
NAME                   COMMAND                  SERVICE                STATUS              PORTS
etcd                   "/usr/local/bin/etcd…"   etcd                   running             0.0.0.0:2379-2380->2379-2380/tcp
kafka                  "start-kafka.sh"         kafka                  running             0.0.0.0:9092->9092/tcp
mongo                  "docker-entrypoint.s…"   mongodb                running             0.0.0.0:27017->27017/tcp
mysql                  "docker-entrypoint.s…"   mysql                  running             0.0.0.0:3306->3306/tcp
open_im_api            "air"                    open_im_api            running             0.0.0.0:10000->10000/tcp
open_im_auth           "air"                    open_im_auth           running             0.0.0.0:10600->10600/tcp
open_im_friend         "air"                    open_im_friend         running             0.0.0.0:10200->10200/tcp
open_im_group          "air"                    open_im_group          running             0.0.0.0:10500->10500/tcp
open_im_msg_gateway    "air"                    open_im_msg_gateway    running
open_im_msg_transfer   "air"                    open_im_msg_transfer   running
open_im_push           "air"                    open_im_push           running             0.0.0.0:10700->10700/tcp
open_im_timed_task     "air"                    open_im_timed_task     running
open_im_user           "air"                    open_im_user           running             0.0.0.0:10100->10100/tcp
redis                  "docker-entrypoint.s…"   redis                  running             0.0.0.0:6379->6379/tcp
zookeeper              "/bin/sh -c '/usr/sb…"   zookeeper              running             0.0.0.0:2181->2181/tcp
```

### 检查容器日志

```

docker-compose -f docker-compose.dev.yaml logs -f
```

若要检查指定容器的日志，比如 `open_im_api`，则

```
docker-compose -f docker-compose.dev.yaml logs -f open_im_api
```