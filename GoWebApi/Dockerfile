FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct
ENV TZ=Asia/Shanghai

WORKDIR /webapi
COPY . /webapi

RUN go mod download

RUN go build -o webapi

FROM alpine

ENV TZ=Asia/Shanghai

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /webapi

COPY --from=builder /webapi/webapi /webapi/webapi
COPY --from=builder /webapi/config.yaml /webapi/config.yaml

ENTRYPOINT ["/webapi/webapi"]
