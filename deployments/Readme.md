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
- `database`: set to your mongodb database name.(Need have a created database.)
- `authSource`: set to your mongodb authSource. (authSource is specify the database name associated with the user's credentials, user need create in this database.)

**kafka.yml**

- `address`: set to your already kafka address or kafka Service name and port in your deployed.

**redis.yml**

- `address`: set to your already redis address or redis Service name and port in your deployed.

**minio.yml**

- `internalAddress`: set to your minio Service name and port in your deployed.
- `externalAddress`: set to your already expose minio external address.

### Set the secret

A Secret is an object that contains a small amount of sensitive data. Such as password and secret. Secret is similar to ConfigMaps.

#### Redis:

Update the `redis-password` value in `redis-secret.yml` to your Redis password encoded in base64.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: openim-redis-secret
type: Opaque
data:
  redis-password: b3BlbklNMTIz # update to your redis password encoded in base64, if need empty, you can set to ""
```

#### Mongo:

Update the `mongo_openim_username`, `mongo_openim_password` value in `mongo-secret.yml` to your Mongo username and password encoded in base64.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: openim-mongo-secret
type: Opaque
data:
  mongo_openim_username: b3BlbklN # update to your mongo username encoded in base64, if need empty, you can set to "" (this user credentials need in authSource database).
  mongo_openim_password: b3BlbklNMTIz # update to your mongo password encoded in base64, if need empty, you can set to ""
```

#### Minio:

Update the `minio-root-user` and `minio-root-password` value in `minio-secret.yml` to your MinIO accessKeyID and secretAccessKey encoded in base64.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: openim-minio-secret
type: Opaque
data:
  minio-root-user: cm9vdA== # update to your minio accessKeyID encoded in base64, if need empty, you can set to ""
  minio-root-password: b3BlbklNMTIz # update to your minio secretAccessKey encoded in base64, if need empty, you can set to ""
```

#### Kafka:

Update the `kafka-password` value in `kafka-secret.yml` to your Kafka password encoded in base64.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: openim-kafka-secret
type: Opaque
data:
  kafka-password: b3BlbklNMTIz # update to your kafka password encoded in base64, if need empty, you can set to ""
```

### Apply the secret.

```shell
kubectl apply -f redis-secret.yml -f minio-secret.yml -f mongo-secret.yml -f kafka-secret.yml
```

### Apply all config

`kubectl apply -f ./openim-config.yml`

> Attation: If you use `default` namespace, you can excute `clusterRile.yml` to create a cluster role binding for default service account.
>
> Namespace is modify to `discovery.yml` in `openim-config.yml`, you can change `kubernetes.namespace` to your namespace.

**Excute `clusterRole.yml`**

`kubectl apply -f ./clusterRole.yml`

### run all deployments and services

> Note: Ensure that infrastructure services like MinIO, Redis, and Kafka are running before deploying the main applications.

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
