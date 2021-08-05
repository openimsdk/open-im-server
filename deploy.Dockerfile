FROM golang:1.15 as build

# go mod Installation source, container environment variable addition will override the default variable value
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# Set up the working directory
WORKDIR /home/Open-IM-Server
# add all files to the container
COPY . .

WORKDIR /home/Open-IM-Server/script
RUN chmod +x *.sh

RUN /bin/sh -c ./build_all_service.sh

#Blank image Multi-Stage Build
FROM ubuntu

RUN rm -rf /var/lib/apt/lists/*
RUN apt-get install apt-transport-https && apt-get update && apt-get install procps\
&&apt-get install net-tools
#Non-interactive operation
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y vim curl tzdata gawk
#Time zone adjusted to East eighth District
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata


#set directory to map logs,config file,script file.
VOLUME ["/home/Open-IM-Server/logs","/home/Open-IM-Server/config","/home/Open-IM-Server/script"]

#Copy scripts files and binary files to the blank image
COPY --from=build /home/Open-IM-Server/script /home/Open-IM-Server/script
COPY --from=build /home/Open-IM-Server/bin /home/Open-IM-Server/bin

WORKDIR /home/Open-IM-Server/script

# "&& tail -f /dev/null " Prevent the container exit after the command is executed
CMD ["./start_all.sh && tail -f /dev/null"]