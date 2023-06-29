
### 以docker-compose 形式单独部署
```sh
# 查看 ./Makefile ，先编译各个需要的源码到 ../bin 
# win-* 表示在win平台编译位linux二进制，其实就是处理了 go env -w GOOS=linux 
make win-build-all

# 得到各个二进制程序之后，打包为镜像
# 目前没有处理 Open-IM-SDK-Core ，需要的话可以自己单独处理这个模块
make image-all

# docker-compose.yaml 分成了两部分，一部分是openIM的镜像容器 openim.yaml，一部分是依赖的环境 env.yaml
# 两部分使用一个外部的网络来联通，所以首先创建用到的 network
docker network create openim --attachable=true -d bridge

# 处理openim组件需要的挂载目录，主要是处理config目录
mkdir ./config
cp ./config.example.yaml ./config/config.yaml # 修改 ./config/config.yaml 内容，比如各个依赖组件的 host

# 然后拉起env.yaml
docker-compose -f ./env.yaml up -d

# 等env 容器全部拉起成功之后，拉起openim.yaml
docker-compose -f ./openim.yaml up -d

# 查看容器运行，推荐使用下 portainer ，web查看容器情况，查看日志等等
docker container ps -a | grep openim

# 正常应该是查看api,demo等的容器日志，看到gin打印的路由日志才算是成功
```