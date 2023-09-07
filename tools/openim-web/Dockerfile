# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# 使用官方Go镜像作为基础镜像
FROM golang:1.21 AS build-env
ENV CGO_ENABLED=0
# 设置工作目录
WORKDIR /app

# 安装curl和unzip工具
#RUN apt-get update && apt-get install -y curl unzip

# 从GitHub下载并解压dist.zip
#RUN curl -LO https://github.com/OpenIMSDK/dist.zip \
 #   && unzip dist.zip -d ./ \
 #   && rm dist.zip

# 复制Go代码到容器
COPY . .

# 编译Go代码
RUN go build -o openim-web

# 使用轻量级的基础镜像
FROM debian:buster-slim

# 将编译好的二进制文件和dist资源复制到新的容器
WORKDIR /app
COPY --from=build-env /app/openim-web /app/openim-web
COPY --from=build-env /app/dist /app/dist

# 开放容器的20001端口
EXPOSE 20001

# 指定容器启动命令
ENTRYPOINT ["/app/openim-web"]