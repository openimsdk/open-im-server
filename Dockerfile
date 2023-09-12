# Build Stage
FROM golang:1.20 AS builder

# Set go mod installation source and proxy
ARG GO111MODULE=on
ARG GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=$GO111MODULE
ENV GOPROXY=$GOPROXY

# Set up the working directory
WORKDIR /openim/openim-server

COPY go.mod go.sum ./
RUN go mod download

# Copy all files to the container
ADD . .

RUN make clean
RUN make build

FROM ghcr.io/openim-sigs/openim-ubuntu-image:latest

WORKDIR ${SERVER_WORKDIR}

# Copy scripts and binary files to the production image
COPY --from=builder ${OPENIM_SERVER_BINDIR} /openim/openim-server/_output/bin
# COPY --from=builder ${OPENIM_SERVER_CMDDIR} /openim/openim-server/scripts
# COPY --from=builder ${SERVER_WORKDIR}/config /openim/openim-server/config

CMD ["/openim/openim-server/scripts/docker-start-all.sh"]
