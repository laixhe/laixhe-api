GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
GIT_VERSION=$(shell git describe --tags --always)

API_PROTO_FILES=$(shell find api -name *.proto)

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: api
api:
	rm -rf ./api/gen
	protoc --proto_path=./api \
		--proto_path=./third_party \
		--go_out=./ \
		$(API_PROTO_FILES)
	mv ./webapi/api/gen ./api
	rm -rf ./webapi
	protoc-go-inject-tag -input="./api/gen/*/*.pb.go"
	#protoc-go-inject-tag -input="./api/gen/*/*/*.pb.go"

.PHONY: build
build:
	swag init
	go build -ldflags "-X main.GitVersion=$(GIT_VERSION)"
