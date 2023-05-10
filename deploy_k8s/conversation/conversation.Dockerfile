FROM ubuntu

# 设置固定的项目路径
ENV WORKDIR /Open-IM-Server
ENV CMDDIR $WORKDIR/cmd
ENV CONFIG_NAME $WORKDIR/config/config.yaml

# 将可执行文件复制到目标目录
ADD ./open_im_conversation $WORKDIR/cmd/main

# 创建用于挂载的几个目录，添加可执行权限
RUN mkdir $WORKDIR/logs $WORKDIR/config $WORKDIR/scripts && \
  chmod +x $WORKDIR/cmd/main

VOLUME ["/Open-IM-Server/logs","/Open-IM-Server/config","/Open-IM-Server/scripts"]

WORKDIR $CMDDIR
CMD ./main