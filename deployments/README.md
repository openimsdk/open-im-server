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

假设 OpenIM 项目根目录路径为 `OpenIM_ROOT`

进入 OpenIM 项目根目录

$ cd ${OpenIM_ROOT}



### 容器化安装

具体安装步骤如下：



### Helm安装
