GIT_VERSION=$(shell git describe --tags --always)
DATE_TIME=$(shell date "+%Y-%m-%dT%H:%M:%S")

.PHONY: swag
swag:
	swag init
	rm ./docs/docs.go

.PHONY: build
build:
	swag init
	rm ./docs/docs.go
	go build -ldflags "-X main.GitVersion=$(DATE_TIME)-$(GIT_VERSION)"
