.PHONY: all build run gotool install clean help

BINARY_NAME=open_im_admin_cms
BIN_DIR=../../../bin/

OS:= $(or $(os),linux)
ARCH:=$(or $(arch),amd64)
all: gotool build

build:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags="-w -s"

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
