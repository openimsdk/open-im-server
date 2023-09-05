# OpenIM 应用容器化部署指南

OpenIM 支持很多种集群化部署方式，包括但不限于 helm, sealos, kubeam, kubesphere, kubeflow, kuboard, kubespray, k3s, k3d, k3c, k3sup, k3v, k3x

**目前还在开发这个模块，预计 v3.2.0 之前会有一个集群方案。**

目前各个贡献者，以及之前的官方有出过一些可以参考的方案：

- https://github.com/OpenIMSDK/k8s-jenkins
- https://github.com/OpenIMSDK/Open-IM-Server-k8s-deploy
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

### 容器化安装

具体安装步骤如下：



### Helm安装
