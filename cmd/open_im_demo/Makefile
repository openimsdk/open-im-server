.PHONY: all build run gotool install clean help

BINARY_NAME=open_im_demo
BIN_DIR=../../bin/


all: gotool build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make build
	mv ${BINARY_NAME} ${BIN_DIR}

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

