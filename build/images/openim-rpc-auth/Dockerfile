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

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make clean

# RUN make build BINS=openim-rpc

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o ./_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/openim-rpc-auth ./cmd/openim-rpc/openim-rpc-auth

FROM ghcr.io/openim-sigs/openim-bash-image:latest

WORKDIR /openim/openim-server

ENV OS ${OS}
ENV ARCH ${ARCH}

COPY --from=builder /openim/openim-server/_output/bin/platforms /openim/openim-server/_output/bin/platforms/
COPY --from=builder /openim/openim-server/config /openim/openim-server/config

EXPOSE 10160

CMD ["bash", "-c","${OPENIM_SERVER_BINDIR}/platforms/$(get_os)/$(get_arch)/openim-rpc-auth --port 10160 -c ${SERVER_WORKDIR}/config"]