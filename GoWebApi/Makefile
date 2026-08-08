# 注意: 本 Makefile 面向 Unix/macOS (依赖 swag 命令)。
# Windows 用户可手动执行:
#   swag init
#   del docs\docs.go
#   go build -ldflags "-X main.GitVersion=xxx"
#
# 说明: swag init 会生成 docs/docs.go (依赖 github.com/swaggo/swag, 项目不引入该依赖),
# 因此生成后立即删除; main.go 导入的 webapi/docs 实际用的是手写的 docs/swagger.go
# (通过 //go:embed 嵌入 swagger.json/yaml)。若误删 swagger.go, 构建会报 "no Go files"。
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

.PHONY: test
test:
	go test ./...
