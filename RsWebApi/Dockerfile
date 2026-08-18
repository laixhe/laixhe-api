# syntax=docker/dockerfile:1
# Rust 多阶段构建

# ---------- 构建阶段 ----------
FROM rust:1.85-alpine AS builder

# 使用国内镜像加速（可按需注释）
ENV RUSTUP_DIST_SERVER=https://rsproxy.cn
ENV RUSTUP_UPDATE_ROOT=https://rsproxy.cn/rustup
ENV CARGO_REGISTRIES_CRATES_IO_PROTOCOL=sparse
ENV TZ=Asia/Shanghai

WORKDIR /webapi

# 先复制依赖清单，利用 Docker 层缓存加速依赖编译
COPY Cargo.toml Cargo.lock ./
# 创建虚拟 src 以预编译依赖
RUN mkdir src && echo "fn main() {}" > src/main.rs \
    && cargo build --release || true \
    && rm -rf src

# 复制真实源码并编译
COPY src ./src
COPY docs ./docs
COPY config.yaml ./config.yaml

# docker build --build-arg GIT_VERSION=xxx ...
ARG GIT_VERSION
RUN cargo build --release \
    && mv target/release/webapi /webapi/webapi

# ---------- 运行阶段 ----------
FROM alpine:3.20

# 镜像加速（中科大），设置时区
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /webapi

COPY --from=builder /webapi/webapi /webapi/webapi
COPY --from=builder /webapi/config.yaml /webapi/config.yaml

EXPOSE 6600

# 优雅停机：SIGTERM 由 Docker stop 发送
STOPSIGNAL SIGTERM

ENTRYPOINT ["/webapi/webapi"]
