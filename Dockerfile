# Build Stage
FROM golang as build

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

RUN apk --no-cache add tzdata

# Set directory to map logs, config files, scripts, and SDK
VOLUME ["/Open-IM-Server/logs", "/Open-IM-Server/config", "/Open-IM-Server/scripts", "/Open-IM-Server/db/sdk"]

# Copy scripts and binary files to the production image
COPY --from=build /Open-IM-Server/scripts /Open-IM-Server/scripts
COPY --from=build /Open-IM-Server/_output/bin/platforms/linux/arm64 /Open-IM-Server/_output/bin/platforms/linux/arm64

WORKDIR /Open-IM-Server/scripts

CMD ["./docker_start_all.sh"]
