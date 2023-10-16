# OpenIM 应用容器化部署指南

OpenIM 支持很多种集群化部署方式，包括但不限于 helm, sealos, kubeam, kubesphere, kubeflow, kuboard, kubespray, k3s, k3d, k3c, k3sup, k3v, k3x

**目前还在开发这个模块，预计 v3.2.0 之前会有一个集群方案。**

目前各个贡献者，以及之前的官方有出过一些可以参考的方案：

- https://github.com/OpenIMSDK/k8s-jenkins
- https://github.com/openimsdk/open-im-server-k8s-deploy
- https://github.com/OpenIMSDK/openim-charts
- https://github.com/showurl/deploy-openim


### 依赖检查

```bash
Kubernetes: >= 1.16.0-0
Helm: >= 3.0
```


### 最低配置

建议生产环境的最低配置如下：

```bash
CPU: 4
Memory: 8G
Disk: 100G
```

## 生成配置文件

我们将自动文件全部自动化处理了，所以生成配置文件对于 openim 来说是可选的，但是如果你想要自定义配置，可以参考下面的步骤：

```bash
$ make init
# 或者是使用脚本:
# ./scripts/init-config.sh
```
此时会帮你在 `deployments/openim/config` 目录下生成配置文件，你可以根据自己的需求进行修改。


## 集群搭建

如果你已经有了一个 `kubernetes` 集群，或者是你希望自己从头开始搭建一个 `kubernetes` 那么你可以直接跳过这一步。

为了快速开始，我使用 [sealos](https://github.com/labring/sealos) 来快速搭建集群，sealos 底层也是对 kubeadm 的封装:

```bash
$ SEALOS_VERSION=`curl -s https://api.github.com/repos/labring/sealos/releases/latest | grep -oE '"tag_name": "[^"]+"' | head -n1 | cut -d'"' -f4` && \
  curl -sfL https://raw.githubusercontent.com/labring/sealos/${SEALOS_VERSION}/scripts/install.sh |
  sh -s ${SEALOS_VERSION} labring/sealos
```

**支持的版本：**

+ docker: `labring/kubernetes-docker`:(v1.24.0~v1.27.0)
+ containerd: `labring/kubernetes`:(v1.24.0~v1.27.0)


#### 安装集群：

集群的信息如下：

| 机器名   | IP地址          | 系统信息                                                                                                    |
|---------|-----------------|------------------------------------------------------------------------------------------------------------|
| master01| 10.0.0.9   | Linux VM-0-9-ubuntu 5.15.0-76-generic #83-Ubuntu SMP Thu Jun 15 19:16:32 UTC 2023 x86_64 x86_64 x86_64 GNU/Linux |
| node01  | 10.0.0.4   | Linux VM-0-9-ubuntu 5.15.0-76-generic #83-Ubuntu SMP Thu Jun 15 19:16:32 UTC 2023 x86_64 x86_64 x86_64 GNU/Linux |
| node02  | 10.0.0.10  | Linux VM-0-9-ubuntu 5.15.0-76-generic #83-Ubuntu SMP Thu Jun 15 19:16:32 UTC 2023 x86_64 x86_64 x86_64 GNU/Linux |

```bash
$ export CLUSTER_USERNAME=ubuntu
$ export CLUSTER_PASSWORD=123456
$ sudo sealos run labring/kubernetes:v1.25.0 labring/helm:v3.8.2 labring/calico:v3.24.1 \
    --masters 10.0.0.9 \
    --nodes 10.0.0.4,10.0.0.10 \
    -u "$CLUSTER_USERNAME" \
    -p "$CLUSTER_PASSWORD"
```

### 安装 helm

helm通过打包的方式，支持发布的版本管理和控制，很大程度上简化了Kubernetes应用的部署和管理。


**使用脚本：**

```bash
$ curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

**添加仓库：**

```bash
$ helm repo add brigade https://openimsdk.github.io/openim-charts
```


### 容器化安装

具体安装步骤如下：



### Helm安装
