# openim开源项目k8s部署
## 需要的环境
    MacOS OR Linux
## 需要安装的软件
- wget
- docker
- kubectl
## k8s配置
    编辑 setting.env
```shell
#docker镜像tag前缀
DOCKER_REGISTRY_ADDR="harbor.local/"
#nodePort暴露端口号 api
NODE_PORT_API="10100"
#nodePort暴露端口号 msg_gateway
NODE_PORT_MSG_GATEWAY="10101"
#nodePort暴露端口号 sdk_server
NODE_PORT_SDK_SERVER="10102"
#nodePort暴露端口号 demo
NODE_PORT_DEMO="10103"
#部署在哪个命名空间
K8S_NAMESPACE="openim"
```
## 单个服务部署
```shell
cd $服务路径 && \
bash build.sh
```
## 全部部署
```shell
go run .
```