FROM golang:1.20.0 as build

WORKDIR /openim

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/miniblog

COPY miniblog /openim

EXPOSE 10002
CMD ["./bin/openim-api","--port", "10002"]


# OpenIM base image: https://github.com/openim-sigs/openim-base-image

# Set go mod installation source and proxy

ARG GOARCH
ARG GOOS

FROM golang:1.20 AS builder

ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /openim/openim-chat

ENV GO111MODULE=$GO111MODULE
ENV GOPROXY=$GOPROXY

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o /openim/openim-chat/bin/open_im_admin_api ./cmd/api/admin_api

FROM ghcr.io/openim-sigs/openim-bash-image:latest

WORKDIR ${CHAT_WORKDIR}

COPY --from=builder ${OPENIM_CHAT_BINDIR}/open_im_admin_api /openim/openim-chat/bin/open_im_admin_api
COPY --from=builder ${OPENIM_CHAT_CONFIG_NAME} /openim/openim-server/config/config.yaml

env PORT=10002

EXPOSE ${PORT}
CMD ["sh","-c","${OPENIM_CHAT_BINDIR}/open_im_admin_api","--port", "${PORT}","--config_folder_path","${OPENIM_CHAT_CONFIG_NAME}"]
