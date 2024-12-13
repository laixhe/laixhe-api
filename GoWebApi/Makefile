GIT_VERSION=$(shell git describe --tags --always)

PROTOCOL_PROTO_FILES=$(shell find protocol/api -name *.proto)

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: protocol
protocol:
	rm -rf ./protocol/gen
	protoc --proto_path=./protocol/api \
		--proto_path=./protocol/third_party \
		--go_out=./ \
		$(PROTOCOL_PROTO_FILES)
	mv ./webapi/protocol/gen ./protocol
	rm -rf ./webapi
	find protocol/gen -type f -name '*.pb.go' -exec sed -i 's/,omitempty"`/\"\`/g' {} +
	protoc-go-inject-tag -input="./protocol/gen/*/*.pb.go"

.PHONY: build
build:
	swag init
	go build -ldflags "-X main.GitVersion=$(GIT_VERSION)"
