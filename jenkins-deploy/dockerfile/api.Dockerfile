FROM golang:1.19.0 as build

WORKDIR /openim
COPY . .

RUN make fmt  \
    && make tidy
RUN make api

FROM ubuntu

WORKDIR /openim
VOLUME ["/openim/logs","/openim/bin"]

COPY --from=build /openim/bin /openim/bin
COPY --from=build /openim/config /openim/config

EXPOSE 10002
CMD ["./bin/openim-api","--port", "10002"]
