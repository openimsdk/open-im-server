FROM golang:alpine AS builder

ARG RELEASE=false
ARG COMPRESS=false
WORKDIR /openim-server

RUN apk add --no-cache upx

RUN go install github.com/magefile/mage@latest

COPY . .
RUN go mod download
RUN RELEASE=${RELEASE} COMPRESS=${COMPRESS} mage build
RUN mage -compile ./mage -ldflags "-s -w"

FROM alpine:latest

WORKDIR /openim-server

COPY --from=builder /openim-server/_output ./_output
COPY --from=builder /openim-server/config ./config
COPY --from=builder /openim-server/start-config.yml ./start-config.yml
COPY --from=builder /openim-server/mage ./mage

ENTRYPOINT ["sh", "-c", "./mage start && sleep infinity"]