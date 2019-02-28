FROM alpine:3.8

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

ADD ./slack-exporter /slack-exporter

ENTRYPOINT ["/slack-exporter"]
