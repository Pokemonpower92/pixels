# Go related variables
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_VET=$(GO_CMD) vet
GO_CLEAN=$(GO_CMD) clean
GO_FMT=$(GO_CMD) fmt
GO_TEST=$(GO_CMD) test
LOCALSTACK_CMD=sh ./scripts/localstack.sh

authapi: vet
	$(GO_BUILD) -C cmd/authapi -o ../../bin/authapi

collageapi: vet
	$(GO_BUILD) -C cmd/collageapi -o ../../bin/collageapi

thumbnail-worker: vet
	$(GO_BUILD) -C cmd/thumbnail-worker -o ../../bin/thumbnail-worker

metadata-worker: vet
	$(GO_BUILD) -C cmd/metadata-worker -o ../../bin/metadata-worker

filestore: vet
	$(GO_BUILD) -C cmd/filestore -o ../../bin/filestore
	
vet: fmt
	$(GO_VET) ./...

fmt:
	$(GO_FMT) ./...

test:
	$(GO_TEST) --cover ./...

clean:
	$(GO_CLEAN) && \
	rm ./bin/* && \
	find ./resources -type f ! -name '.keep' -delete 

build: collageapi authapi thumbnail-worker metadata-worker filestore
all: vet build

.PHONY: build fmt vet clean all

