GO ?= latest
GOBIN = $(CURDIR)/build/bin
GORUN = env GO111MODULE=on go run
GOPATH = $(shell go env GOPATH)
GO_FLAGS += -buildvcs=false
GOTEST = GODEBUG=cgocheck=0 go test $(GO_FLAGS) -p 1

lint:
	@./build/bin/golangci-lint run --config ./.golangci.yml

lintci-deps:
	rm -f ./build/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./build/bin v1.53.2

test:
	$(GOTEST) --timeout 5m -shuffle=on ./...

devtools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@type "protoc" 2> /dev/null || echo '\n\nPlease install protoc\n Linux : apt install -y protobuf-compiler\n Mac : brew install protobuf\n Windows : choco install protoc\n'

protoc:
	protoc --go_out=. --go-grpc_out=. ./protos/*.proto
