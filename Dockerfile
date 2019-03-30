ARG VERSION="0.0.0-docker"

ARG GO_VERSION=1.12
ARG ALPINE_VERSION=3.9

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
WORKDIR /go/src/github.com/syntaqx/serve

RUN apk add --no-cache git ca-certificates
ENV CGO_ENABLED=0 GO111MODULE=on

ADD go.mod go.sum ./
RUN go mod download

COPY . /go/src/github.com/syntaqx/serve
RUN go build -installsuffix cgo -ldflags '-s -w -X main.version=$VERSION' -o ./bin/serve ./cmd/serve

FROM alpine:${ALPINE_VERSION}
LABEL maintainer="Chase Pierce <syntaqx@gmail.com>"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/syntaqx/serve/bin/serve /usr/bin/
COPY --from=builder /go/src/github.com/syntaqx/serve/static /var/www

RUN addgroup -S serve \
  && adduser -D -S -s /sbin/nologin -G serve serve
USER serve

VOLUME ["/var/www"]

CMD ["serve", "-dir", "/var/www"]
EXPOSE 8080
