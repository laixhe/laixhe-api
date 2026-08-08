GIT_VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo dev)

# 默认目标
.PHONY: all
all: check test build

# 检查代码能否编译（快速）
.PHONY: check
check:
	cargo check

# 运行集成测试（依赖本机 MySQL，见 config.yaml）
.PHONY: test
test:
	cargo test

# 格式化
.PHONY: fmt
fmt:
	cargo fmt --all

# 静态检查
.PHONY: lint
lint:
	cargo clippy --all-targets -- -D warnings

# 发布构建，可注入版本号: make build GIT_VERSION=v1.0.0
.PHONY: build
build:
	cargo build --release
	@echo "GIT_VERSION=$(GIT_VERSION)" > .version

# 直接运行（开发调试）
.PHONY: run
run:
	cargo run -- --config=./config.yaml

# 清理构建产物
.PHONY: clean
clean:
	cargo clean
	rm -f .version
