FROM golang:1.26-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct
ENV TZ=Asia/Shanghai

# docker build --build-arg GIT_VERSION=XXX ...
ARG GIT_VERSION
ENV GIT_VERSION=${GIT_VERSION}

WORKDIR /webapi

# 先复制依赖文件，利用 Docker 层缓存加速构建
COPY go.mod go.sum ./
RUN go mod download

COPY . /webapi

RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.GitVersion=$(date +%Y-%m-%dT%H:%M:%S)-${GIT_VERSION}" -o webapi

FROM alpine

# 镜像加速， 添加 alpine apk 源（中科大）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

ENV TZ=Asia/Shanghai
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /webapi

COPY --from=builder /webapi/webapi /webapi/webapi
COPY --from=builder /webapi/config.yaml /webapi/config.yaml

ENTRYPOINT ["/webapi/webapi"]
