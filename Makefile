.PHONY: run build test clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/user-service cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...
