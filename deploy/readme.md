
### 以docker-compose 形式单独部署
```sh
# 查看 ./Makefile ，先编译各个需要的源码到 ../bin 
# win-* 表示在win平台编译位linux二进制，其实就是处理了 go env -w GOOS=linux 
make win-build-all

# 得到各个二进制程序之后，打包为镜像
#
make image-all

# docker-compose.yaml 分成了两部分，一部分是openIM的镜像容器 openim.yaml，一部分是依赖的环境 env.yaml
# 两部分使用一个外部的网络来联通，所以首先创建用到的 network
docker network create openim --attachable=true -d bridge

# 然后拉起env.yaml
docker-compose -f ./env.yaml up -d

# 等env 容器全部拉起成功之后，拉起openim.yaml
docker-compose -f ./openim.yaml up -d
```