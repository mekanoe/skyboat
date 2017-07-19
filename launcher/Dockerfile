FROM alpine:3.6
# MAINTAINER Katie T. <katie@kat.cafe>

RUN apk add --no-cache su-exec && \
    adduser -S skyboat

CMD ["su-exec", "skyboat", "launcher"]
EXPOSE 2390
COPY ./launcher /usr/bin/launcher
