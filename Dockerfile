FROM golang:alpine

MAINTAINER Knut Ahlers <knut@ahlers.me>

ADD . /go/src/github.com/Luzifer/locationmaps
WORKDIR /go/src/github.com/Luzifer/locationmaps

RUN set -ex \
 && apk add --update git \
 && go install -ldflags "-X main.version=$(git describe --tags || git rev-parse --short HEAD || echo dev)" \
 && apk del --purge git

EXPOSE 3000

ENTRYPOINT ["/go/bin/locationmaps"]
CMD ["--"]
