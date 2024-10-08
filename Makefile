# Makefile for a Go project with two applications: app1 and app2

# Go related variables
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_VET=$(GO_CMD) vet
GO_CLEAN=$(GO_CMD) clean
GO_FMT=$(GO_CMD) fmt
GO_TEST=$(GO_CMD) test -cover

db: vet
	$(GO_CMD) run ./scripts/db/create_db.go

localstack: 
	sh ./build/localstack.sh

imagesetparser: vet
	$(GO_BUILD) -C cmd/imagesetparser -o ../../bin/imagesetparser

collageapi: vet
	$(GO_BUILD) -C cmd/collageapi -o ../../bin/collageapi

vet: fmt
	$(GO_VET) ./...

fmt:
	$(GO_FMT) ./...

test:
	$(GO_TEST) ./...

clean:
	$(GO_CLEAN) && \
	rm ./bin/*

build: imagesetparser collageapi
all: vet build

.PHONY: build fmt vet clean all

