# Use Go 1.21 as the base image for building the application
FROM golang:1.21 as builder

# Define the base directory for the application as an environment variable
ENV OPENIM_SERVER_DIR=/openim-server

# Set the working directory inside the container based on the environment variable
WORKDIR $OPENIM_SERVER_DIR

# Set the Go proxy to improve dependency resolution speed
ENV GOPROXY=https://goproxy.cn,direct

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy all files from the current directory into the container
COPY . .

# Install Mage to use for building the application
RUN go install github.com/magefile/mage@latest

# Execute the build command using Mage
#RUN mage build

# Use Alpine Linux as the final base image due to its small size and included utilities
FROM alpine:latest

# Define the base directory for the application as an environment variable again
ENV OPENIM_SERVER_DIR=/openim-server

# Install necessary packages, such as bash, to ensure compatibility and functionality
RUN apk add --no-cache bash

# Copy the compiled binaries and mage from the builder image to the final image
COPY --from=builder $OPENIM_SERVER_DIR/_output $OPENIM_SERVER_DIR/_output
COPY --from=builder /go/bin/mage /usr/local/bin/mage
COPY --from=builder $OPENIM_SERVER_DIR/magefile_windows.go $OPENIM_SERVER_DIR/
COPY --from=builder $OPENIM_SERVER_DIR/magefile_unix.go $OPENIM_SERVER_DIR/
COPY --from=builder $OPENIM_SERVER_DIR/magefile.go $OPENIM_SERVER_DIR/
COPY --from=builder $OPENIM_SERVER_DIR/start-config.yml $OPENIM_SERVER_DIR/

# Set the working directory to the application directory within the container
WORKDIR $OPENIM_SERVER_DIR

# Set up volume mounts for the configuration directory and logs directory
VOLUME ["$OPENIM_SERVER_DIR/config", "$OPENIM_SERVER_DIR/_output/logs"]

# Set the command to run when the container starts
ENTRYPOINT ["sh", "-c", "mage start && tail -f /dev/null"]
