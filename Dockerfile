FROM golang:1.16-alpine AS builder
WORKDIR /go/src/github.com/syntaqx/serve

ARG VERSION="0.0.0-docker"

RUN apk add --no-cache git ca-certificates
ENV CGO_ENABLED=0 GO111MODULE=on

COPY go.* ./
RUN go mod download

COPY . /go/src/github.com/syntaqx/serve
RUN go install -ldflags "-X main.version=$VERSION" ./cmd/...

FROM alpine:3
LABEL maintainer="Chase Pierce <syntaqx@gmail.com>"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/syntaqx/serve/static /var/www
COPY --from=builder /go/bin/serve /usr/bin/

RUN addgroup -S serve \
  && adduser -D -S -s /sbin/nologin -G serve serve
USER serve

VOLUME ["/var/www"]

EXPOSE 8080
CMD ["serve", "-dir", "/var/www"]
