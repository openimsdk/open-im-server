# Kubernetes Deployment

## Resource Requests

- CPU: 2 cores
- Memory: 4 GiB
- Disk usage: 20 GiB (on Node)

## Preconditions

ensure that you have already deployed the following components:

- Redis
- MongoDB
- Kafka
- MinIO

## Origin Deploy

### Enter the target dir

`cd ./deployments/deploy/`

### Deploy configs and dependencies

Upate your configMap `openim-config.yml`. **You can check the official docs for more details.**

In `openim-config.yml`, you need modify the following configurations:

**discovery.yml**

- `kubernetes.namespace`: default is `default`, you can change it to your namespace.

**mongodb.yml**

- `address`: set to your already mongodb address or mongo Service name and port in your deployed.
- `database`: set to your mongodb database name.
- `authSource`: set to your mongodb authSource. (authSource is specify the database name associated with the user's credentials, user need create in this database.)

**share.yml**

- `secret`: same to **OpenIM Chat** secret.
- `imAdminUserID`: default is `imAdmin`.

**kafka.yml**

- `address`: set to your already kafka address or kafka Service name and port in your deployed.

**redis.yml**

- `address`: set to your already redis address or redis Service name and port in your deployed.

**minio.yml**

- `internalAddress`: set to your already minio internal address or minio Service name and port in your deployed.
- `externalAddress`: set to your already expose minio external address.

### Set the secret

A Secret is an object that contains a small amount of sensitive data. Such as password and secret. Secret is similar to ConfigMaps.

#### Example:

create a secret for redis password. You can update `redis-secret.yml`.

you need update `redis-password` value to your redis password in base64.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: openim-redis-secret
type: Opaque
data:
  redis-password: b3BlbklNMTIz # you need update to your redis password in base64
```

#### Usage:

use secret in deployment file. If you apply the secret to IM Server, you need adapt the Env Name to config file and all toupper.

OpenIM Server use prefix `IMENV_`, OpenIM Chat use prefix `CHATENV_`. Next adapt is the config file name. Like `redis.yml`. Such as `IMENV_REDIS_PASSWORD` is mapped to `redis.yml` password filed in OpenIM Server.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rpc-user-server
spec:
  template:
    spec:
      containers:
        - name: rpc-user-server
          env:
            - name: IMENV_REDIS_PASSWORD # adapt to redis.yml password field in OpenIM Server config, Don't modify it.
              valueFrom:
                secretKeyRef:
                  name: openim-redis-secret # You deployed secret name
                  key: redis-password # You deployed secret key name
```

So, you need following configurations to set secret:

- `MONGODB_USERNAME`
- `MONGODB_PASSWORD`
- `REDIS_PASSWORD`
- `MINIO_ACCESSKEYID`
- `MINIO_SECRETACCESSKEY`

Apply the secret.

```shell
kubectl apply -f redis-secret.yml -f minio-secret.yml -f mongo-secret.yml
```

### Apply all config and dependencies

`kubectl apply -f ./openim-config.yml`

> Attation: If you use `default` namespace, you can excute `clusterRile.yml` to create a cluster role binding for default service account.
>
> Namespace is modify to `discovery.yml` in `openim-config.yml`, you can change `kubernetes.namespace` to your namespace.

**Excute `clusterRole.yml`**

`kubectl apply -f ./clusterRole.yml`

**If you have already deployed the storage component, you need to update corresponding config and secret. And pass corresponding deployments and services build.**

> Note: Ensure that infrastructure services like MinIO, Redis, and Kafka are running before deploying the main applications.

### run all deployments and services

```bash
kubectl apply \
  -f openim-api-deployment.yml \
  -f openim-api-service.yml \
  -f openim-crontask-deployment.yml \
  -f openim-rpc-user-deployment.yml \
  -f openim-rpc-user-service.yml \
  -f openim-msggateway-deployment.yml \
  -f openim-msggateway-service.yml \
  -f openim-push-deployment.yml \
  -f openim-push-service.yml \
  -f openim-msgtransfer-service.yml \
  -f openim-msgtransfer-deployment.yml \
  -f openim-rpc-conversation-deployment.yml \
  -f openim-rpc-conversation-service.yml \
  -f openim-rpc-auth-deployment.yml \
  -f openim-rpc-auth-service.yml \
  -f openim-rpc-group-deployment.yml \
  -f openim-rpc-group-service.yml \
  -f openim-rpc-friend-deployment.yml \
  -f openim-rpc-friend-service.yml \
  -f openim-rpc-msg-deployment.yml \
  -f openim-rpc-msg-service.yml \
  -f openim-rpc-third-deployment.yml \
  -f openim-rpc-third-service.yml
```

### Verification

After deploying the services, verify that everything is running smoothly:

```bash
# Check the status of all pods
kubectl get pods

# Check the status of services
kubectl get svc

# Check the status of deployments
kubectl get deployments

# View all resources
kubectl get all
```

### clean all

`kubectl delete -f ./`

### Notes:

- If you use a specific namespace for your deployment, be sure to append the -n <namespace> flag to your kubectl commands.
