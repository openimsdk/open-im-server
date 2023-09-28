
# build container
FROM golang:1.20-alpine3.18 AS builder
ENV GOPROXY https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on

WORKDIR /app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add git pkgconfig build-base
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .
RUN go build -o cmd/openim-rpc/openim-rpc-msg/openim-rpc-msg cmd/openim-rpc/openim-rpc-msg/main.go

# archive container
FROM alpine:3.18
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates libdrm
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /app/bin
COPY --from=builder /app/cmd/openim-rpc/openim-rpc-msg/openim-rpc-msg /app/bin/
ENTRYPOINT ["/app/bin/openim-rpc-msg"]
