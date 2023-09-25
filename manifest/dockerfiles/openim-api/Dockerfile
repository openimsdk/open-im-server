# ******************************************************************************
# 2019 - present Contributed by Apulis Technology (Shenzhen) Co. LTD
#
# This program and the accompanying materials are made available under the
# terms of the MIT License, which is available at
# https://www.opensource.org/licenses/MIT
#
# See the NOTICE file distributed with this work for additional
# information regarding copyright ownership.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.
#
# SPDX-License-Identifier: MIT
#******************************************************************************/

# build container
FROM golang:1.20-alpine3.18 AS builder
ENV GOPROXY https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on

WORKDIR /openim/openim-server
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add git pkgconfig build-base
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .
RUN go build -o cmd/openim-api/openim-api cmd/openim-api/main.go

# archive container
FROM alpine:3.18
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates libdrm
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /openim/openim-server/bin
COPY --from=builder /openim/openim-server/cmd/openim-api/openim-api /openim/openim-server/bin/
ENTRYPOINT ["/openim/openim-server/bin/openim-api"]
