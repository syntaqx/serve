FROM golang:1.19-alpine AS builder

ARG VERSION="0.0.0-docker"

RUN apk add --update --no-cache \
  ca-certificates tzdata openssh git mercurial && update-ca-certificates \
  && rm -rf /var/cache/apk/*

WORKDIR /src

COPY go.mod* go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 go install -ldflags "-X main.version=$VERSION" ./cmd/...

FROM alpine

RUN adduser -S -D -H -h /app appuser
USER appuser

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/* /bin/

ENV PORT=8080
EXPOSE $PORT

VOLUME ["/var/www"]

CMD ["serve", "--dir", "/var/www"]
