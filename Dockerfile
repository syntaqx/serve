FROM golang:1.26-alpine AS builder

ARG VERSION="0.0.0-docker"

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

WORKDIR /src

COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o /out/serve ./cmd/serve

FROM alpine:3.24

RUN adduser -S -D -H -h /app appuser
USER appuser

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/serve /bin/serve

ENV PORT=8080
EXPOSE 8080

VOLUME ["/var/www"]

# Liveness probe against the auth-free health endpoint (works with TLS off).
HEALTHCHECK --interval=30s --timeout=3s --start-period=2s --retries=3 \
  CMD wget -qO- "http://127.0.0.1:${PORT}/healthz" >/dev/null 2>&1 || exit 1

CMD ["serve", "--dir", "/var/www"]
