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

# OpenIM base image: https://github.com/openim-sigs/openim-base-image

# Set go mod installation source and proxy

FROM golang:1.20 AS builder

ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /openim/openim-server

ENV GO111MODULE=$GO111MODULE
ENV GOPROXY=$GOPROXY

RUN apt-get update && apt-get install -y curl unzip

RUN curl -LO https://app-1302656840.cos.ap-nanjing.myqcloud.com/dist.zip \
   && unzip dist.zip -d ./ \
   && rm dist.zip

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make clean
RUN make build BINS=openim-web

FROM ghcr.io/openim-sigs/openim-ubuntu-image:latest

WORKDIR /openim/openim-server

COPY --from=builder /openim/openim-server/_output/bin/platforms /openim/openim-server/_output/bin/platforms/
COPY --from=builder /app/dist /app/dist

ENV PORT 11001
ENV DISTPATH /app/dist

EXPOSE $PORT

RUN cp -r ${OPENIM_SERVER_BINDIR}/platforms/$(get_os)/$(get_arch)/openim-web /usr/bin/openim-web

ENTRYPOINT ["bash", "-c", "openim-web -port $PORT"]

CMD ["-distPath","/app/dist"]