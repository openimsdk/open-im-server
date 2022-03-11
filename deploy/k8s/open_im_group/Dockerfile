FROM golang as build

# go mod Installation source, container environment variable addition will override the default variable value
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV VERSION_TAG="2.0.2"

# Set up the working directory
#WORKDIR /Open-IM-Server
# add all files to the container
#COPY . .

COPY Open-IM-Server.tar.gz /
COPY Open-IM-SDK-Core.tar.gz /

WORKDIR /
RUN tar -zxvf Open-IM-Server.tar.gz && \
    mv Open-IM-Server-$VERSION_TAG Open-IM-Server && \
    rm -rf Open-IM-Server.tar.gz && \
    tar -zxvf Open-IM-SDK-Core.tar.gz && \
    rm -rf Open-IM-Server/cmd/Open-IM-SDK-Core/ && \
    mv Open-IM-SDK-Core-$VERSION_TAG Open-IM-Server/cmd/Open-IM-SDK-Core && \
    rm -rf Open-IM-Server.tar.gz

COPY build-open_im_group.sh /Open-IM-Server/script

WORKDIR /Open-IM-Server/script
RUN chmod +x *.sh && \
    cd /Open-IM-Server/cmd && sed -i "s#GOARCH=amd64##g" `grep "amd64" -rl .` && \
    cd /Open-IM-Server/script && \
    /bin/sh -c ./build-open_im_group.sh

#Blank image Multi-Stage Build
FROM showurl/openim-base:v2.0.2

#set directory to map logs,config file,script file.
VOLUME ["/Open-IM-Server/logs","/Open-IM-Server/script","/Open-IM-Server/db/sdk"]

#Copy scripts files and binary files to the blank image
COPY --from=build /Open-IM-Server/script /Open-IM-Server/script
COPY --from=build /Open-IM-Server/bin /Open-IM-Server/bin
COPY start-open_im_group.sh /Open-IM-Server/script

WORKDIR /Open-IM-Server/script

CMD ["bash", "start-open_im_group.sh"]
