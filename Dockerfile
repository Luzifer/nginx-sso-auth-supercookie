FROM golang:alpine as builder

ENV CGO_ENABLED=1

# Add required packages
RUN set -ex \
 && apk add --update \
      build-base \
      curl \
      git \
 && curl -sSfLo /tmp/dumb-init "https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64" \
 && install -Dm 0755 -t /build/usr/local/bin \
      /tmp/dumb-init

# Install base application
RUN set -ex \
 && git clone https://github.com/Luzifer/nginx-sso.git /go/src/github.com/Luzifer/nginx-sso \
 && cd /go/src/github.com/Luzifer/nginx-sso \
 && git reset --hard $(git describe --tags --abbrev=0) \
 && go build \
      -ldflags "-X main.version=$(git describe --tags || git rev-parse --short HEAD || echo dev)" \
      -mod=readonly \
 && install -Dm 0755 -t /build/usr/local/bin \
      docker-start.sh \
      nginx-sso \
 && install -Dm 0644 -t /build/usr/local/share/nginx-sso \
      config.yaml \
 && install -Dm 0644 -t /build/usr/local/share/nginx-sso/frontend \
      frontend/index.html

# Install local source-code as plugin
ADD . /go/src/github.com/Luzifer/nginx-sso-auth-supercookie
WORKDIR /go/src/github.com/Luzifer/nginx-sso-auth-supercookie

RUN set -ex \
 && go build \
      -buildmode=plugin \
      -mod=readonly \
      -o auth-supercookie.so \
 && install -Dm 0644 -t /build/usr/local/share/nginx-sso/plugins \
      auth-supercookie.so


FROM alpine

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      bash \
      ca-certificates

COPY --from=builder /build/ /

EXPOSE 8082
VOLUME ["/data"]

ENTRYPOINT ["/usr/local/bin/docker-start.sh"]
CMD ["--"]

# vim: set ft=Dockerfile:
