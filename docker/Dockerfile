FROM alpine:latest

RUN apk update && apk upgrade
RUN apk add --update --no-cache yarn nodejs runuser git

RUN mkdir -p /opt/yarn-contained/workspace

WORKDIR /opt/yarn-contained/workspace
