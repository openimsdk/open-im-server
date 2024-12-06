# Use Go 1.22 Alpine as the base image for building the application
FROM golang:1.22-alpine AS builder
# Define the base directory for the application as an environment variable
ENV SERVER_DIR=/openim-server

# Set the working directory inside the container based on the environment variable
WORKDIR $SERVER_DIR

# Set the Go proxy to improve dependency resolution speed

#ENV GOPROXY=https://goproxy.io,direct

# Copy all files from the current directory into the container
COPY . .

RUN go mod tidy



RUN go build -o _output/openim-push ./cmd/openim-push


# Using Alpine Linux for the final image
FROM alpine:latest

# Install necessary packages, such as bash
RUN apk add --no-cache bash

# Set the environment and work directory
ENV SERVER_DIR=/openim-server
WORKDIR $SERVER_DIR


# Copy the compiled binaries and mage from the builder image to the final image
COPY --from=builder $SERVER_DIR/_output $SERVER_DIR/_output
# COPY --from=builder $SERVER_DIR/config $SERVER_DIR/config

# Set the command to run when the container starts
ENTRYPOINT ["sh", "-c", "_output/openim-push"]