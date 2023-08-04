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
RUN make build BINS=openim-msggateway

FROM ghcr.io/openim-sigs/openim-bash-image:latest

WORKDIR /openim/openim-server

COPY --from=builder /openim/openim-server/_output/bin/platforms /openim/openim-server/_output/bin/platforms
COPY --from=builder /openim/openim-server/config /openim/openim-server/config

ENV OS ${OS}
ENV ARCH ${ARCH}

EXPOSE 10140
EXPOSE 10001

CMD ${OPENIM_SERVER_BINDIR}/platforms/${OS}/${ARCH}/openim-msggateway --port 10140 --ws_port 10001