FROM golang:1.18.0 as build

WORKDIR /openim
COPY . .

RUN make fmt && make tidy
RUN make transfer

FROM ubuntu

WORKDIR /openim
VOLUME ["/openim/logs","/openim/bin"]

#Copy binary files to the blank image
COPY --from=build /openim/bin /openim/bin
COPY --from=build /openim/config /openim/config

CMD ["./bin/openim-msgtransfer"]
