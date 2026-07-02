FROM golang:alpine AS builder

ENV SERVER_DIR=/openim-server
WORKDIR $SERVER_DIR

COPY . .
RUN go mod tidy

ARG CMD_PATH
ARG BINARY_NAME
ARG RELEASE=false

RUN if [ "$RELEASE" = "true" ]; then \
    go build -trimpath -ldflags "-s -w" -o _output/${BINARY_NAME} ./${CMD_PATH}; \
    else \
    go build -o _output/${BINARY_NAME} ./${CMD_PATH}; \
    fi

FROM alpine:latest

RUN apk add --no-cache bash

ARG BINARY_NAME
ENV BINARY_NAME=${BINARY_NAME}
ENV SERVER_DIR=/openim-server
WORKDIR $SERVER_DIR

COPY --from=builder $SERVER_DIR/_output $SERVER_DIR/_output

ENTRYPOINT ["sh", "-c", "_output/${BINARY_NAME}"]