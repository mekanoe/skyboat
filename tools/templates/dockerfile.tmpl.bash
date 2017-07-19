DIRNAME=$(basename $PWD)
FILE="FROM alpine:3.6
# MAINTAINER Katie T. <katie@kat.cafe>

RUN apk add --no-cache su-exec && \
    adduser -S spln

CMD [\"su-exec\", \"spln\", \"$DIRNAME\"]
EXPOSE 2390
COPY ./$DIRNAME /usr/bin/$DIRNAME
"

echo "$FILE" > Dockerfile