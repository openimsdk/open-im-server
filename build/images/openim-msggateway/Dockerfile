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

ENV BINS=openim-msggateway

# Optionally build your application if needed
RUN mage build ${BINS} check-free-memory seq || true

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


RUN echo -e "serviceBinaries:\n  openim-msggateway: 1 \n" \
    > $SERVER_DIR/start-config.yml && \
    echo -e "toolBinaries:\n  - check-free-memory\n  - seq\n" >> $SERVER_DIR/start-config.yml && \
    echo "maxFileDescriptors: 10000" >> $SERVER_DIR/start-config.yml

RUN go get github.com/openimsdk/openim-gomake@0.0.15-alpha.1

# Set the command to run when the container starts
ENTRYPOINT ["sh", "-c", "mage start && tail -f /dev/null"]
