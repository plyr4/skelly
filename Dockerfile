
FROM golang:alpine3.12

EXPOSE 8080

ENV GODEBUG=netdns=go

ADD release/skelly /bin/

ENTRYPOINT [ "/bin/skelly", "server"]
