# OpenIM Application Containerization Deployment Guide

OpenIM supports a variety of cluster deployment methods, including but not limited to `helm`, `sealos`, `kustomize`

Various contributors, as well as previous official releases, have provided some referenceable solutions:

+ [k8s-jenkins Repository](https://github.com/OpenIMSDK/k8s-jenkins)
+ [open-im-server-k8s-deploy Repository](https://github.com/openimsdk/open-im-server-k8s-deploy)
+ [openim-charts Repository](https://github.com/OpenIMSDK/openim-charts)
+ [deploy-openim Repository](https://github.com/showurl/deploy-openim)

### Dependency Check

```bash
Kubernetes: >= 1.16.0-0
Helm: >= 3.0
```

### Minimum Configuration

The recommended minimum configuration for a production environment is as follows:

```yaml
CPU: 4
Memory: 8G
Disk: 100G
```

## Configuration File Generation

We have automated all the files, making the generation of configuration files optional for OpenIM. However, if you desire custom configurations, you can follow the steps below:

```bash
$ make init
# Alternatively, use script:
# ./scripts/init-config.sh
```

At this point, configuration files will be generated under `deployments/openim/config`, which you can modify as per your requirements.

## Cluster Setup

If you already have a `kubernetes` cluster, or if you wish to build a `kubernetes` cluster from scratch, you can skip this step.

For a quick start, I used [sealos](https://github.com/labring/sealos) to rapidly set up the cluster, with sealos also being a wrapper for kubeadm at its core:

```bash
$ SEALOS_VERSION=`curl -s https://api.github.com/repos/labring/sealos/releases/latest | grep -oE '"tag_name": "[^"]+"' | head -n1 | cut -d'"' -f4` && \
  curl -sfL https://raw.githubusercontent.com/labring/sealos/${SEALOS_VERSION}/scripts/install.sh |
  sh -s ${SEALOS_VERSION} labring/sealos
```

**Supported Versions:**

+ docker: `labring/kubernetes-docker`:(v1.24.0~v1.27.0)
+ containerd: `labring/kubernetes`:(v1.24.0~v1.27.0)

#### Cluster Installation:

Cluster details are as follows:

| Hostname | IP Address | System Info                                                  |
| -------- | ---------- | ------------------------------------------------------------ |
| master01 | 10.0.0.9   | `Linux VM-0-9-ubuntu 5.15.0-76-generic #83-Ubuntu SMP Thu Jun 15 19:16:32 UTC 2023 x86_64 x86_64 x86_64 GNU/Linux` |
| node01   | 10.0.0.4   | Similar to master01                                          |
| node02   | 10.0.0.10  | Similar to master01                                          |

```bash
$ export CLUSTER_USERNAME=ubuntu
$ export CLUSTER_PASSWORD=123456
$ sudo sealos run labring/kubernetes:v1.25.0 labring/helm:v3.8.2 labring/calico:v3.24.1 \
    --masters 10.0.0.9 \
    --nodes 10.0.0.4,10.0.0.10 \
    -u "$CLUSTER_USERNAME" \
    -p "$CLUSTER_PASSWORD"
```

> **Node** Uninstallation method: using `kubeadm` for uninstallation does not remove `etcd` and `cni` related configurations. Manual clearance or using `sealos` for uninstallation is needed.
>
> ```bash
> $ sealos reset
> ```

If you are local, you can also use Kind and Minikube to test, for example, using Kind:

```bash
$ sGO111MODULE="on" go get sigs.k8s.io/kind@v0.11.1
$ skind create cluster
```

### Installing helm

Helm simplifies the deployment and management of Kubernetes applications to a large extent by offering version control and release management through packaging.

**Using Script:**

```bash
$ curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

**Adding Repository:**

```bash
$ helm repo add brigade https://openimsdk.github.io/openim-charts
```

### OpenIM Image Strategy

Automated offerings include aliyun, ghcr, docker hub: [Image Documentation](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md)

**Local Test Build Method:**

```bash
$ make image
```

> This command assists in quickly building the required images locally. For a detailed build strategy, refer to the [Build Documentation](https://github.com/openimsdk/open-im-server/blob/main/build/README.md).

## Installation

Explore our Helm-Charts repository and read through: [Helm-Charts Repository](https://github.com/openimsdk/helm-charts)


Using the helm charts repository, you can ignore the following configuration, but if you want to just use the server and scale on top of it, you can go ahead:

**Use Helmfile:**

```bash
GO111MODULE=on go get github.com/roboll/helmfile@latest
```

```bash
export MYSQL_ADDRESS=im-mysql
export MYSQL_PORT=3306
export MONGO_ADDRESS=im-mongo
export MONGO_PORT=27017
export REDIS_ADDRESS=im-redis-master
export REDIS_PORT=6379
export KAFKA_ADDRESS=im-kafka
export KAFKA_PORT=9092
export OBJECT_APIURL="https://openim.server.com/api"
export MINIO_ENDPOINT="http://im-minio:9000"
export MINIO_SIGN_ENDPOINT="https://openim.server.com/im-minio-api"

mkdir ./charts/generated-configs
../scripts/genconfig.sh ../scripts/install/environment.sh ./templates/openim.yaml > ./charts/generated-configs/config.yaml
cp ../config/notification.yaml ./charts/generated-configs/notification.yaml
../scripts/genconfig.sh ../scripts/install/environment.sh ./templates/helm-image.yaml > ./charts/generated-configs/helm-image.yaml
```

```bash
helmfile apply
```