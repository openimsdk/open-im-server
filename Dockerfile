FROM golang as build

# go mod Installation source, container environment variable addition will override the default variable value
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# Set up the working directory
WORKDIR /Open-IM-Server
# add all files to the container
COPY . .

WORKDIR /Open-IM-Server/scripts
RUN chmod +x *.sh

RUN /bin/sh -c ./build_all_service.sh

#Blank image Multi-Stage Build
FROM ubuntu

RUN rm -rf /var/lib/apt/lists/*
RUN apt-get update && apt-get install apt-transport-https && apt-get install procps\
&&apt-get install net-tools
#Non-interactive operation
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y vim curl tzdata gawk
#Time zone adjusted to East eighth District
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata


#set directory to map logs,config file,scripts file.
VOLUME ["/Open-IM-Server/logs","/Open-IM-Server/config","/Open-IM-Server/scripts","/Open-IM-Server/db/sdk"]

#Copy scripts files and binary files to the blank image
COPY --from=build /Open-IM-Server/scripts /Open-IM-Server/scripts
COPY --from=build /Open-IM-Server/_output/bin/platforms/linux/amd64 /Open-IM-Server/_output/bin/platforms/linux/amd64

WORKDIR /Open-IM-Server/scripts

CMD ["./docker_start_all.sh"]
