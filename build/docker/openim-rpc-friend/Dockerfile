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

RUN make build BINS=openim-rpc
# RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o ./_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/openim-rpc-friend ./cmd/openim-rpc/openim-rpc-friend

RUN rm -rf ./_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/openim-rpc-group && \
    rm -rf ./_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/openim-rpc-msg && \
    rm -rf ./_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/openim-rpc-third && \
    rm -rf ./_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/openim-rpc-user && \
    rm -rf ./_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/openim-rpc-conversation && \
    rm -rf ./_output/bin/platforms/$(go env GOOS)/$(go env GOARCH)/openim-rpc-auth

FROM ghcr.io/openim-sigs/openim-bash-image:latest

WORKDIR /openim/openim-server

ENV OS ${OS}
ENV ARCH ${ARCH}

COPY --from=builder /openim/openim-server/_output/bin/platforms /openim/openim-server/_output/bin/platforms/
COPY --from=builder /openim/openim-server/config /openim/openim-server/config

EXPOSE 10120
EXPOSE 20120

CMD ["sh","-c","${OPENIM_SERVER_BINDIR}/platforms/${OS}/${ARCH}/openim-rpc-friend --port 10120 --prometheus_port 20120 -c ${SERVER_WORKDIR}/config"]