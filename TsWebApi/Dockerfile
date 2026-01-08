FROM oven/bun:1.3-alpine AS builder

ENV TZ=Asia/Shanghai

WORKDIR /tswebapi
COPY . /tswebapi
#RUN rm -rf /tswebapi/node_modules

RUN bun install
RUN bun prisma generate
RUN bun run build

FROM oven/bun:1.3-alpine

# 镜像加速， 添加alpine apk源（中科大）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

ENV TZ=Asia/Shanghai
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /tswebapi

COPY --from=builder /tswebapi/node_modules/.prisma /tswebapi/node_modules/.prisma
COPY --from=builder /tswebapi/dist/ /tswebapi/
#COPY --from=builder /tswebapi/.env /tswebapi/.env

EXPOSE 80

ENTRYPOINT ["bun", "/tswebapi/index.js"]
