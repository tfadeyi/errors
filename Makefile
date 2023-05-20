ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# include files with the `// +build mock` annotation
TEST_TAGS:=-tags mock -coverprofile cover.out

.PHONY: build generate test build-all-platforms clean install-ci-tools install-local-tools licenses

install-ci-tools:
	go install github.com/matryer/moq@v0.2.7
	go install github.com/google/go-licenses@c781b427440f8ea100841eefdd308e660d26d121

install-local-tools: install-ci-tools
	go install github.com/goreleaser/goreleaser@v1.17.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

build:
	cd $(ROOT_DIR) && $(GO_BUILD) -o builds/$(BIN_NAME) .

# used to generate struct mocks
generate:
	cd $(ROOT_DIR) && go generate ./...

test: build
	cd $(ROOT_DIR) &&  go test $(GOFLAGS) $(TEST_TAGS) ./...

coverage: test
	go tool cover -html=cover.out -o cover.html
	@echo "open ./cover.html to see coverage"

lint:
	golangci-lint run

# generates the licenses used by the tool
licenses:
	rm -rf licenses
	go-licenses save . --save_path="./licenses"
