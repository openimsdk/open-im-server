# Use Go 1.21 as the base image for building the application
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /openim-server

# Set the Go proxy to improve dependency resolution speed
ENV GOPROXY=https://goproxy.cn,direct

# Copy all files from the current directory into the container
COPY . .

RUN go install github.com/magefile/mage@latest

RUN mage -v && mage build

# Use Alpine Linux as the final base image due to its small size and included utilities
FROM alpine:latest

# Install necessary packages, such as bash, to ensure compatibility and functionality
RUN apk add --no-cache bash

# Copy the compiled binaries and mage from the builder image to the final image
COPY --from=builder /openim-server/_output /openim-server/_output
COPY --from=builder /go/bin/mage /usr/local/bin/mage

# Set the working directory to /openim-server within the container
WORKDIR /openim-server

# Set up volume mounts for the configuration directory and logs directory
VOLUME ["/openim-server/config", "/openim-server/_output/logs"]

# Set the command to run when the container starts
ENTRYPOINT ["sh", "-c", "mage start && tail -f /dev/null"]


