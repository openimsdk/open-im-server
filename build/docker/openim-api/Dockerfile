FROM golang:1.20.0 as build

WORKDIR /openim

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/miniblog

COPY miniblog /openim

EXPOSE 10002
CMD ["./bin/openim-api","--port", "10002"]