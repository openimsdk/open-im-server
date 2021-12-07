FROM alpine:3.13

# 设置固定的项目路径
ENV WORKDIR /app
ENV CONFIG_NAME $WORKDIR/config/config.yaml

# 将可执行文件复制到目标目录
ADD ./open_im_rpc_user $WORKDIR/main

# 创建用于挂载的几个目录，添加可执行权限
RUN mkdir $WORKDIR/logs $WORKDIR/config $WORKDIR/db && \
  chmod +x $WORKDIR/main


WORKDIR $WORKDIR
CMD ./main