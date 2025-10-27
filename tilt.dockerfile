FROM alpine:3.22.2
WORKDIR /
COPY ./bin/manager /manager

ENTRYPOINT ["/manager"]
