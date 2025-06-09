# Go related variables
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_VET=$(GO_CMD) vet
GO_CLEAN=$(GO_CMD) clean
GO_FMT=$(GO_CMD) fmt
GO_TEST=$(GO_CMD) test
LOCALSTACK_CMD=sh ./scripts/localstack.sh

api: vet
	$(GO_BUILD) -C cmd/api -o ../../bin/api

vet: fmt
	$(GO_VET) ./...

fmt:
	$(GO_FMT) ./...

test:
	$(GO_TEST) --cover ./...

clean:
	$(GO_CLEAN) && \
	rm ./bin/* && \

build: collageapi 
all: vet build

.PHONY: build fmt vet clean all

