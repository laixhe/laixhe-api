GIT_VERSION=$(shell git describe --tags --always)

.PHONY: swag
swag:
	swag init
	rm ./docs/docs.go

.PHONY: build
build:
	swag init
	rm ./docs/docs.go
	go build -ldflags "-X main.GitVersion=$(GIT_VERSION)"
