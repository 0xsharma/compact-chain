GO ?= latest
GOBIN = $(CURDIR)/build/bin
GORUN = env GO111MODULE=on go run
GOPATH = $(shell go env GOPATH)

lint:
	@./build/bin/golangci-lint run --config ./.golangci.yml

lintci-deps:
	rm -f ./build/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./build/bin v1.46.0
