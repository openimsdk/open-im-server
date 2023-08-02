# Build Stage
FROM golang:1.20 AS builder

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

FROM ghcr.io/openim-sigs/openim-bash-image:latest

# Copy scripts and binary files to the production image
COPY --from=builder /Open-IM-Server/scripts /Open-IM-Server/scripts
COPY --from=builder /Open-IM-Server/_output/bin/platforms/linux/amd64 /Open-IM-Server/_output/bin/platforms/linux/amd64

WORKDIR /Open-IM-Server/scripts

CMD ["./docker_start_all.sh"]