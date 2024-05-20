GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

API_PROTO_FILES=$(shell find profile -name *.proto)

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: api
# generate api proto
api:
	rm -rf ./profile/gen
	protoc --proto_path=./profile \
		--proto_path=./third-party \
		--go_out=./ \
		--go-grpc_out=./ \
		$(API_PROTO_FILES)
	mv ./webapi/profile/gen ./profile && rm -rf ./webapi
	protoc-go-inject-tag -input="./profile/gen/*/*.pb.go"

.PHONY: build
# build cmd dir
build:
	swag init
	go build -ldflags "-X main.Version=$(VERSION)"
