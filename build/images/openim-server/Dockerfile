FROM golang:1.25-alpine AS builder

ARG CMD_PATH
ARG BINARY_NAME
ARG RELEASE=false

ENV SERVER_DIR=/openim-server
WORKDIR $SERVER_DIR

COPY . .
RUN go mod tidy

RUN if [ "$RELEASE" = "true" ]; then \
    go build -trimpath -ldflags "-s -w" -o myapp ./cmd/open-im-api; \
    else \
    RUN go build -o myapp ./cmd/open-im-api

FROM alpine:latest

RUN apk add --no-cache bash

ARG BINARY_NAME
ENV BINARY_NAME=${BINARY_NAME}
ENV SERVER_DIR=/openim-server
WORKDIR $SERVER_DIR

COPY --from=builder $SERVER_DIR/_output $SERVER_DIR/_output

ENTRYPOINT ["sh", "-c", "_output/${BINARY_NAME}"]
