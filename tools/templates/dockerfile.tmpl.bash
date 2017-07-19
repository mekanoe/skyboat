DIRNAME=$(basename $PWD)
FILE="FROM alpine:3.6
# MAINTAINER Skyboat.io <engineering@skyboat.io>

RUN apk add --no-cache su-exec && \
    adduser -S skyboat

CMD [\"su-exec\", \"skyboat\", \"$DIRNAME\"]
EXPOSE 2390
COPY ./$DIRNAME /usr/bin/$DIRNAME
"

echo "$FILE" > Dockerfile