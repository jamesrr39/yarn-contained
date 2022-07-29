FROM alpine:latest

RUN apk update && apk add yarn runuser

RUN mkdir -p /opt

RUN mkdir -p /opt/yarn-contained/workspace
# feed /opt/start.sh inline so that we don't need to provide a docker build context
# this is a bit painful, main points:
# - end every line with \n\
# - escape variables with backslash, e.g. \${MY_VAR}
RUN echo -e "#!/bin/sh\n\
\n\
# check user exists\n\
\n\
adduser -D -u \${HOST_USER_ID} \${USERNAME}\n\
\n\
chown -R \${HOST_USER_ID} /opt/yarn-contained/workspace\n\
\n\
runuser -l \${USERNAME} -c 'cd /opt/yarn-contained/workspace && yarn \$@'\n\
\n\
\n\
" > /opt/start.sh

RUN chmod +x /opt/start.sh
RUN ls -lh /opt/start.sh
RUN wc -l /opt/start.sh
RUN cat /opt/start.sh

RUN adduser -D user

# COPY start.sh /opt/start.sh

ENTRYPOINT ["/opt/start.sh"]
