FROM alpine:3.18

RUN apk update
RUN apk add --no-cache yarn nodejs runuser git

RUN mkdir -p /opt/yarn-contained/workspace

WORKDIR /opt/yarn-contained/workspace
