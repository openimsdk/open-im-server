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

RUN go mod download
RUN make clean
RUN make build BINS=openim-cmdutils

FROM ghcr.io/openim-sigs/openim-bash-image:latest

WORKDIR /openim/openim-server

COPY --from=builder ${SERVER_WORKDIR}/_output/bin/platforms /openim/openim-server/_output/bin/platforms
COPY --from=builder ${SERVER_WORKDIR}/config /openim/openim-server/config

CMD ["sh","-c","${OPENIM_SERVER_BINDIR}/platforms/${OS}/${ARCH}/openim-cmdutils"]
