# Build Stage
FROM golang:1.20 AS builder

LABEL org.opencontainers.image.source=https://github.com/OpenIMSDK/Open-IM-Server
LABEL org.opencontainers.image.description="OpenIM Server image"
LABEL org.opencontainers.image.licenses="Apache 2.0"

# Set go mod installation source and proxy
ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=$GO111MODULE
ENV GOPROXY=$GOPROXY

# Set up the working directory
WORKDIR /Open-IM-Server

# Copy all files to the container
ADD . .

RUN /bin/sh -c "make build"

# Production Stage
FROM alpine

RUN echo "https://mirrors.aliyun.com/alpine/v3.4/main" > /etc/apk/repositories && \
    apk --no-cache add tzdata ca-certificates

# Set directory to map logs, config files, scripts, and SDK
VOLUME ["/Open-IM-Server/logs", "/Open-IM-Server/config", "/Open-IM-Server/scripts", "/Open-IM-Server/db/sdk"]

# Copy scripts and binary files to the production image
COPY --from=builder /Open-IM-Server/scripts /Open-IM-Server/scripts
COPY --from=builder /Open-IM-Server/_output/bin/platforms/linux/amd64 /Open-IM-Server/_output/bin/platforms/linux/amd64

WORKDIR /Open-IM-Server/scripts

CMD ["docker_start_all.sh"]