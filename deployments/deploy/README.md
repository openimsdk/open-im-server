# Kubernetes Deployment

## Resource Requests
- CPU: 2 cores
- Memory: 4 GiB
- Disk usage: 20 GiB (on Node)

## Origin Deploy
1. Enter the target dir
`cd ./deployments/deploy/`

2. Deploy configs and dependencies

Apply all config and dependencies
`kubectl apply -f ./openim-config.yml -f `

Run infrasturcture components.

`kubectl apply -f minio-service.yml -f minio-statefulset.yml -f mongo-service.yml -f mongo-statefulset.yml -f redis-service.yml -f redis-statefulset.yml -f kafka-service.yml -f kafka-statefulset.yml`

>Note: Ensure that infrastructure services like MinIO, Redis, and Kafka are running before deploying the main applications.


3. run all deployments and services

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

4. Verification
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

5. clean all

`kubectl delete -f ./`

### Notes:
- If you use a specific namespace for your deployment, be sure to append the -n <namespace> flag to your kubectl commands.