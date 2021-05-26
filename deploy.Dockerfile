FROM golang:1.15 as build

RUN rm -rf /var/lib/apt/lists/*
RUN apt-get install apt-transport-https && apt-get update && apt-get install procps\
&&apt-get install net-tools
#Non-interactive operation
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y vim curl tzdata gawk
#Time zone adjusted to East eighth District
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata
# go mod Installation source, container environment variable addition will override the default variable value
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

#map directory to save logs
VOLUME /home/open_im_server/logs
#map directory to save config
VOLUME /home/open_im_server/config

# Set up the working directory
WORKDIR /home/open_im_server
# add all files to the container
COPY . .

WORKDIR /home/open_im_server/script
RUN chmod +x *.sh
CMD ["./start_all.sh"]




