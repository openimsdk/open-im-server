FROM golang:1.18.0 as build

WORKDIR /openim
COPY . .

RUN make fmt && make tidy
RUN make gateway

FROM ubuntu

WORKDIR /openim
VOLUME ["/openim/logs","/openim/bin"]

#Copy binary files to the blank image
COPY --from=build /openim/bin /openim/bin
COPY --from=build /openim/config /openim/config

EXPOSE 10140
EXPOSE 10001
CMD ["./bin/openim-rpc-msg-gateway","--port", "10140","--ws_port", "10001"]
