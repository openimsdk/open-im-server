# Use Go 1.21 as the base image
FROM golang:1.21 as builder

# Set the working directory
WORKDIR /openim-server

# Copy all files from the current directory to the image
COPY . .

# Execute the script and build command
RUN chmod +x ./bootstrap.sh && \
    ./bootstrap.sh && \
    mage build

# Use scratch as the base image for the minimal production image
FROM scratch

# Copy the compiled binaries from the builder image to the production image
COPY --from=builder /openim-server/_output /openim-server/_output

# Set the working directory
WORKDIR /openim-server

# Set up volume mounts for the config directory and logs directory
VOLUME ["/openim-server/config", "/openim-server/_output/logs"]

