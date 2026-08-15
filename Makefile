VERSION=`git --no-pager describe --tags --always`

LDFLAGS+=
LDFLAGS+=-X main.version=${VERSION}

.PHONY: build install run test cover lint vet fmt tidy docker clean

build:
	go build -ldflags "${LDFLAGS}" -o bin/serve ./cmd/serve

install:
	go install -ldflags "${LDFLAGS}" ./cmd/serve

run:
	go run ./cmd/serve

test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

cover: test
	go tool cover -func=coverage.out

lint:
	golangci-lint run

vet:
	go vet ./...

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

docker:
	docker build -t syntaqx/serve --build-arg VERSION=${VERSION} .

clean:
	rm -rf bin dist coverage.out
