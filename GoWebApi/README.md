# WebApi

基于 [Go Fiber v3](https://docs.gofiber.io/) + GORM 的 RESTful API 服务，提供用户注册、登录、JWT 鉴权、用户信息管理、IP 限流、健康检查等基础能力，适合作为新手学习 Go Web 后端的分层项目模板。

## 环境要求

- Go 1.26+（使用了 `clear` 内建函数、`strings.Cut` 等新特性）
- MySQL 5.7+
- 依赖的自研基础库 [laixhe/gonet](https://github.com/laixhe/gonet)（本地 workspace 引用，见下方"核心封装速查"）

## 快速启动

1. 修改 [config.yaml](config.yaml) 中 `orm.dsn` 为你的 MySQL 连接串。
2. 初始化数据库表结构（见"数据库初始化"）。
3. 启动服务：

```bash
go run main.go --config=./config.yaml
```

默认监听 `0.0.0.0:6600`，启动后打印 go 版本、git 版本、配置文件与主机名。

## 目录结构

```
GoWebApi/
├── main.go                  # 程序入口
├── config.yaml              # 配置文件（http/log/orm/jwt/limit）
├── core/                    # 核心基础设施
│   ├── server.go            # Server 聚合：日志、ORM、Fiber 实例
│   ├── config.go            # 配置加载与校验
│   ├── error.go             # 错误处理器（业务错误透传/未知错误统一 500）
│   ├── orm.go               # GORM 日志接入 zap
│   └── middlewares/         # 中间件（JWT 鉴权、IP 限流）
├── routers/                 # 路由注册（auth/user/health + 中间件挂载）
├── app/
│   ├── controllers/         # 控制器层：参数校验、绑定请求、返回响应
│   ├── services/            # 服务层：业务逻辑、数据库操作
│   ├── models/              # 数据模型层：GORM 模型 + 枚举常量
│   ├── entity/              # 请求/响应 DTO（含 swagger 注解）
│   └── util/                # 通用工具（如密码格式校验）
└── docs/                    # swagger 生成文件
```

调用链：`routers → controllers → services → models/entity`，依赖自上而下单向传递。

## 数据库初始化

项目未启用 AutoMigrate，需要手动建表。三张核心表（`user` / `user_extend` / `user_third_party`）的结构见 `app/models/` 下的 GORM tag。

**注意：`user.email` 必须建立唯一索引**。注册接口先查后插，唯一索引兜底并发下的重复注册（重复时返回"邮箱已存在"）：

```sql
ALTER TABLE `user` ADD UNIQUE INDEX idx_user_email (`email`);
```

## API 文档

启动后访问：

```
http://127.0.0.1:6600/api/v1/swagger.yaml
http://127.0.0.1:6600/api/v1/swagger.json
```

接口注解写在 controllers 里，改动后通过 `make swag` 重新生成。

## 核心封装速查

本项目的底层能力大量来自自研库 [laixhe/gonet](https://github.com/laixhe/gonet)，以下是对应关系：

| 代码里看到的 | 底层实现 | 说明 |
| --- | --- | --- |
| `core.NewServer` / `xfiber.New` | gofiber/fiber v3 | Fiber 实例 + 默认错误处理 + requestId |
| `xfiber.UseJwt` | gofiber/contrib jwt + golang-jwt/v5 | JWT 校验中间件 |
| `xfiber.ParamError` / `AuthorizedError` | `*fiber.Error` | 参数错误(422) / 未授权(401) |
| `server.Orm().GetById/FirstByField` | gorm `Take` / `First` | 按 id / 按字段查询 |
| `server.Gorm(ctx)` | gorm `WithContext` | 绑定请求上下文的 *gorm.DB |
| `crypto.BcryptPasswordHash` | bcrypt | 密码哈希 |
| `jwt.GenToken` | golang-jwt/v5 | 签发 JWT |
| `xlog.InitZap` | zap | 日志客户端 |

## 参数校验

`app/entity/` 中各结构体上的 `validate` tag 仅用于 swag 生成 API 文档，项目未注册 validator，tag **不参与请求校验**。真正的校验在 controllers 层手写完成（如昵称长度、邮箱/密码格式、头像地址前缀等），失败时返回 `xfiber.ParamError`（422）。

## 错误约定

- 业务错误（参数校验、鉴权失败、限流、超时）返回 `*fiber.Error`，携带对应状态码与提示信息，原样透传给客户端。
- 未知错误（如数据库异常）由 [core.ErrorHandler](core/error.go) 记录到服务端日志（含 request_id 与请求路径），客户端统一收到固定 500 文案 `internal server error`，不泄露内部细节。
- 服务层使用 `fmt.Errorf("xxx: %w", err)` 包装错误上下文，便于日志排障。

## 安全说明

- **公开接口**：`/api/v1/user/info` 与 `/api/v1/user/list` 未挂载 JWT 中间件，任何人均可查询用户信息与列表。这是有意的设计，但若用于生产请评估是否需要鉴权或脱敏。
- **限流与代理头**：IP 限流优先信任 `X-Forwarded-For`，直接面向公网时客户端可伪造该头绕过限流，建议仅在有可信反向代理时启用（见 [rate_limit.go](core/middlewares/rate_limit.go)）。
- **密码**：bcrypt 哈希存储，注册先查邮箱避免无效计算，`user.email` 唯一索引兜底并发重复注册。
- **JWT**：从 1 开始计的用户 id，`uid <= 0` 视为无效，防御伪造 `{"uid":0}` 的令牌。
- **生产环境务必更换 [config.yaml](config.yaml) 中 `jwt.secret_key` 的默认值**，否则任何人都能用已知密钥伪造令牌。

## 单元测试

```bash
go test ./...
```

当前覆盖纯逻辑部分：密码格式校验、IP 限流器滑动窗口、枚举取值/合法性判断。

## 与 Rust 版本的关系

本仓库是一个 Rust 实现 API 服务的 Go 移植版本（原 Rust 项目中的 `services::auth::register`、`HandleErrorLayer`、`shutdown_signal` 等对应的概念已按 Go 惯例重写）。如需对比两端行为，请以各自仓库为准。
