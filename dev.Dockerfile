FROM golang:1.16 as base

FROM base as dev

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct


RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

CMD ["air"]