VERSION=`git --no-pager describe --tags --always`

LDFLAGS+=-s -w
LDFLAGS+=-X main.version=${VERSION}

build:
	go build -ldflags "${LDFLAGS}" -o bin/serve ./cmd/serve

install:
	go install -ldflags "${LDFLAGS}" ./cmd/serve
