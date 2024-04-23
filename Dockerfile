# Use Go 1.21 as the base image for building the application
FROM golang:1.21 as builder

# Define the base directory for the application as an environment variable
ENV SERVER_DIR=/openim-chat

# Set the working directory inside the container based on the environment variable
WORKDIR $SERVER_DIR

# Set the Go proxy to improve dependency resolution speed
ENV GOPROXY=https://goproxy.io,direct

# Copy all files from the current directory into the container
COPY . .

RUN go mod download


# Install Mage to use for building the application
RUN go install github.com/magefile/mage@latest

# Uncomment and ensure your build command is correctly specified
#RUN mage build


# Use Alpine Linux as the final base image due to its small size and included utilities
FROM alpine:latest

# Install necessary packages, such as bash, to ensure compatibility and functionality
RUN apk add --no-cache bash

ENV SERVER_DIR=/openim-chat

# Set the working directory inside the container based on the environment variable
WORKDIR $SERVER_DIR



# Copy the compiled binaries and mage from the builder image to the final image

COPY --from=builder $SERVER_DIR/_output $SERVER_DIR/
COPY --from=builder $SERVER_DIR/config $SERVER_DIR/
COPY --from=builder /go/bin/mage /usr/local/bin/mage
COPY --from=builder $SERVER_DIR/magefile_windows.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/magefile_unix.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/magefile.go $SERVER_DIR/
COPY --from=builder $SERVER_DIR/start-config.yml $SERVER_DIR/


# Set up volume mounts for the configuration directory and logs directory
VOLUME ["$SERVER_DIR/config", "$SERVER_DIR/_output/logs"]

# Set the command to run when the container starts
ENTRYPOINT ["sh", "-c", "mage start && tail -f /dev/null"]
