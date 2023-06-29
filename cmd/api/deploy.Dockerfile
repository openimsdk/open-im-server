FROM ubuntu

WORKDIR /Open-IM-Server/bin

RUN apt-get update && apt-get install apt-transport-https && apt-get install procps\
&&apt-get install net-tools
#Non-interactive operation
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y vim curl tzdata gawk
#Time zone adjusted to East eighth District
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata
RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl
COPY ./open_im_api ./

VOLUME ["/Open-IM-Server/logs","/Open-IM-Server/config"]

CMD ["./open_im_api","--port", "10002"]
