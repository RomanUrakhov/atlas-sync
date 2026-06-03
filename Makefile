BINARY := atlas-sync
CMD := ./cmd/atlas-sync

.PHONY: build run lint test test-integration

build:
	go build -o ${BINARY} ${CMD}

run: build
	./${BINARY} sync

lint:
	go vet ./...

test:
	go test -v ./...

test-integration:
	go test -tags=integration -v ./...
