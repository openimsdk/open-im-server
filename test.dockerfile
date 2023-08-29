# FROM ghcr.io/openim-sigs/openim-bash-image:latest
FROM ghcr.io/openim-sigs/openim-bash-image:latest

WORKDIR /openim/openim-server

COPY ./_output/bin/platforms /openim/openim-server/_output/bin/platforms
COPY ./config /openim/openim-server/config

ENV PORT 10002

EXPOSE 10002

RUN cp -r ${OPENIM_SERVER_BINDIR}/platforms/$(get_os)/$(get_arch)/openim-api /usr/bin/openim-api

ENTRYPOINT ["/usr/bin/openim-api","-c","${SERVER_WORKDIR}/config"]

CMD ["--port 10002"]