FROM alpine:latest

RUN apk update && apk add yarn runuser git

RUN mkdir -p /opt/yarn-contained/workspace

WORKDIR /opt/yarn-contained/workspace
