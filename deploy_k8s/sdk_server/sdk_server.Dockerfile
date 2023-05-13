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

FROM ubuntu

# 设置固定的项目路径
ENV WORKDIR /Open-IM-Server
ENV CMDDIR $WORKDIR/cmd
ENV CONFIG_NAME $WORKDIR/config/config.yaml

# 将可执行文件复制到目标目录
ADD ./open_im_sdk_server $WORKDIR/main

# 创建用于挂载的几个目录，添加可执行权限
RUN mkdir $WORKDIR/logs $WORKDIR/config $WORKDIR/db && \
  chmod +x $WORKDIR/main

VOLUME ["/Open-IM-Server/logs","/Open-IM-Server/config","/Open-IM-Server/script","/Open-IM-Server/db/sdk"]

WORKDIR $CMDDIR
CMD ./main