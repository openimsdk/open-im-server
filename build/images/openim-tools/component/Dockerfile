# # Copyright © 2023 OpenIM. All rights reserved.
# #
# # Licensed under the Apache License, Version 2.0 (the "License");
# # you may not use this file except in compliance with the License.
# # You may obtain a copy of the License at
# #
# #     http://www.apache.org/licenses/LICENSE-2.0
# #
# # Unless required by applicable law or agreed to in writing, software
# # distributed under the License is distributed on an "AS IS" BASIS,
# # WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# # See the License for the specific language governing permissions and
# # limitations under the License.

# # OpenIM base image: https://github.com/openim-sigs/openim-base-image

# # Set go mod installation source and proxy

# FROM golang:1.20 AS builder

# ARG GO111MODULE=on

# WORKDIR /openim/openim-server

# ENV GO111MODULE=$GO111MODULE
# ENV GOPROXY=$GOPROXY

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN make clean
# RUN make build BINS=component

# # FROM ghcr.io/openim-sigs/openim-bash-image:latest
# FROM ghcr.io/openim-sigs/openim-bash-image:latest

# WORKDIR /openim/openim-server

# COPY --from=builder /openim/openim-server/_output/bin/tools /openim/openim-server/_output/bin/tools/
# COPY --from=builder /openim/openim-server/config /openim/openim-server/config

# ENV OPENIM_SERVER_CONFIG_NAME=/openim/openim-server/config

# RUN mv ${OPENIM_SERVER_BINDIR}/platforms/$(get_os)/$(get_arch)/component /usr/bin/component

# ENTRYPOINT ["bash", "-c", "component -c $OPENIM_SERVER_CONFIG_NAME"]


# Use Go 1.22 Alpine as the base image for building the application
FROM golang:1.22-alpine AS builder

ARG GO111MODULE=on

# Define the base directory for the application as an environment variable
ENV SERVER_DIR=/openim-server

# Set the working directory inside the container based on the environment variable
WORKDIR $SERVER_DIR

# Set the Go proxy to improve dependency resolution speed
ENV GO111MODULE=$GO111MODULE
ENV GOPROXY=https://goproxy.io,direct

# Copy all files from the current directory into the container
COPY . .

RUN go mod download

# Install Mage to use for building the application
RUN go install github.com/magefile/mage@v1.15.0

# ENV BINS=openim-rpc-user

# Optionally build your application if needed
# RUN mage build ${BINS} check-free-memory seq || true
RUN mage build check-free-memory seq || true

# Using Alpine Linux with Go environment for the final image
FROM golang:1.22-alpine

# Install necessary packages, such as bash
RUN apk add bash

# Set the environment and work directory
ENV SERVER_DIR=/openim-server
WORKDIR $SERVER_DIR


# Copy the compiled binaries and mage from the builder image to the final image
COPY --from=builder $SERVER_DIR/_output $SERVER_DIR/_output
COPY --from=builder $SERVER_DIR/config $SERVER_DIR/config
COPY --from=builder /go/bin/mage /usr/local/bin/mage
COPY --from=builder $SERVER_DIR/magefile_windows.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/magefile_unix.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/magefile.go $SERVER_DIR/
# COPY --from=builder $SERVER_DIR/start-config.yml $SERVER_DIR/
COPY --from=builder $SERVER_DIR/go.mod $SERVER_DIR/
COPY --from=builder $SERVER_DIR/go.sum $SERVER_DIR/


RUN echo -e "serviceBinaries:\n  \n" \
    > $SERVER_DIR/start-config.yml && \
    echo -e "toolBinaries:\n  - check-free-memory\n  - seq\n" >> $SERVER_DIR/start-config.yml && \
    echo "maxFileDescriptors: 10000" >> $SERVER_DIR/start-config.yml

RUN go get github.com/openimsdk/openim-gomake@0.0.15-alpha.1

# Set the command to run when the container starts
ENTRYPOINT ["sh", "-c", "mage start && tail -f /dev/null"]
