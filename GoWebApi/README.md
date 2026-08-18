# WebApi

基于 [Go Fiber v3](https://docs.gofiber.io/) + GORM 的 RESTful API 服务，提供用户注册、登录、JWT 鉴权、用户信息管理、IP 限流、健康检查等基础能力，适合作为新手学习 Go Web 后端的分层项目模板。

## 环境要求

- Go 1.26+（使用了 `clear` 内建函数、`strings.Cut` 等新特性）
- MySQL 5.7+
- 依赖的自研基础库 [laixhe/gonet](https://github.com/laixhe/gonet)：版本已在 [go.mod](go.mod) 锁定，构建时由 Go 模块代理自动拉取，无需手动 clone。若你想在本地直接调试 gonet 源码（比如改限流/错误处理逻辑），可以在该目录下创建 `go.work` 把 GoWebApi 与 gonet 加入同一工作区：

```bash
# 在 gonet 仓库根目录执行, 将用到的 gonet 子模块与 GoWebApi 加入同一工作区
# (/绝对路径/laixhe-api/GoWebApi 换成实际路径):
go work init
go work use ./config ./db/gorm/mysql ./db/gorm/orm ./jwt ./utils ./xfiber ./xlog /绝对路径/laixhe-api/GoWebApi
```

## 快速启动

1. 修改 [config.yaml](config.yaml) 中 `orm.dsn` 为你的 MySQL 连接串。
2. 初始化数据库表结构（见"数据库初始化"）。
3. 启动服务：

```bash
go run main.go --config=./config.yaml
```

默认监听 `0.0.0.0:6600`，启动后打印 go 版本、git 版本、配置文件与主机名。验证：

```bash
# 健康检查
curl http://127.0.0.1:6600/api/v1/health

# Swagger UI 文档页（浏览器打开）
http://127.0.0.1:6600/api/v1/swagger

# 注册
curl -X POST http://127.0.0.1:6600/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"nickname":"test","email":"test@example.com","password":"abc123"}'
```

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

项目未启用 AutoMigrate，需要手动建表。可直接导入仓库根目录的 [webapi.sql](../../webapi.sql)，或参考 `app/models/` 下的 GORM tag 建表。

**注意：`user.email` 为唯一索引（与 webapi.sql 一致）**。注册接口先查后插 + 数据库唯一约束双重防重。

## API 文档

启动后访问：

```
http://127.0.0.1:6600/api/v1/swagger        # Swagger UI 可视化页面（浏览器打开）
http://127.0.0.1:6600/api/v1/swagger.json   # OpenAPI JSON 文档
http://127.0.0.1:6600/api/v1/swagger.yaml   # OpenAPI YAML 文档
```

接口注解写在 controllers 里，改动后通过 `make swag` 重新生成（生成 docs/docs.go 后立即删除，实际生效的是 [docs/swagger.go](docs/swagger.go) 用 `//go:embed` 嵌入的 swagger.json/yaml）。

## 核心封装速查

本项目的底层能力大量来自自研库 [laixhe/gonet](https://github.com/laixhe/gonet)，以下是对应关系：

| 代码里看到的 | 底层实现 | 说明 |
| --- | --- | --- |
| `core.NewServer` / `xfiber.New` | gofiber/fiber v3 | Fiber 实例 + 默认错误处理 + requestId |
| `xfiber.UseJwt` | gofiber/contrib jwt + golang-jwt/v5 | JWT 校验中间件 |
| `xfiber.ParamError` / `AuthorizedError` | `*fiber.Error` | 参数错误(422) / 未授权(401) |
| `server.Orm().GetById/FirstByField` | gorm `Take` / `First` | 按主键 / 按唯一字段查全行; 需要裁剪敏感列或复杂条件时改用下方 `server.Gorm(ctx).Select(...)` |
| `ctx.Bind().WithAutoHandling().JSON(req)` | gonet xfiber | 请求体绑定到结构体; 绑定失败时自动返回统一 400 错误响应, 调用方直接 `return err` 透传 |
| `server.Gorm(ctx)` | gorm `WithContext` | 绑定请求上下文的 *gorm.DB |
| `server.Bcrypt().Hash/Check` | bcrypt (BcryptPool) | 密码哈希/校验，CPU 密集计算在独立 worker 池执行，不阻塞请求 goroutine；worker 数由 `config.yaml` 的 `bcrypt.workers` 配置（见 [bcrypt_pool.go](core/bcrypt_pool.go)） |
| `jwt.GenToken` | golang-jwt/v5 | 签发 JWT |
| `xlog.InitZap` | zap | 日志客户端 |

**查询风格**：项目里三种查库写法是刻意并存、各有用途（新手可先统一用 `Orm().GetById/FirstByField`，需要列裁剪时再切换）：
- 查重/只取个别列：`Gorm(ctx).Select("id").First(user, "email = ?", req.Email)`（见 [services/auth.go](app/services/auth.go) 注册）
- 按唯一字段取全行（登录需要 password）：`Orm().FirstByField(ctx, user, "email", req.Email)`
- 按主键取行且排除敏感列：`Gorm(ctx).Select(UserColumnsNoPassword).Where("id", uid).Take(user)`

## 参数校验

`app/entity/` 中各结构体上的 `validate` tag 仅用于 swag 生成 API 文档，项目未注册 validator，tag **不参与请求校验**。真正的校验在 controllers 层手写完成（如昵称长度、邮箱/密码格式、头像地址前缀等），失败时返回 `xfiber.ParamError`（422）。

其中昵称长度按**字符**统计（`utf8.RuneCountInString`），中文等多字节字符不会被按字节数误判。

## 错误约定

- 业务错误（参数校验、鉴权失败、限流、超时）返回 `*fiber.Error`，携带对应状态码与提示信息，原样透传给客户端。
- 未知错误（如数据库异常）由 [core.ErrorHandler](core/error.go) 记录到服务端日志（含 request_id 与请求路径），客户端统一收到固定 500 文案 `internal server error`，不泄露内部细节。
- 服务层使用 `fmt.Errorf("xxx: %w", err)` 包装错误上下文，便于日志排障。

## 安全说明

- **公开接口**：`/api/v1/user/info` 与 `/api/v1/user/list` 未挂载 JWT 中间件，任何人均可查询用户信息与列表。这是有意的设计，但若用于生产请评估是否需要鉴权或脱敏。
- **限流与代理头**：IP 限流优先信任 `X-Forwarded-For`，直接面向公网时客户端可伪造该头绕过限流，建议仅在有可信反向代理时启用（见 [rate_limit.go](core/middlewares/rate_limit.go)）。
- **密码**：bcrypt 哈希存储（成本 cost=10，单次约 50-100ms 的 CPU 密集计算），提交到 [BcryptPool](core/bcrypt_pool.go) worker 池执行，避免阻塞请求 goroutine（对齐 Rust 版 `spawn_blocking` 思路）；worker 数由 `config.yaml` 的 `bcrypt.workers` 配置（缺省 0=CPU 核数）。注册先查邮箱避免无效计算，`user.email` 唯一索引兜底并发重复注册（先查后插 + 数据库唯一约束双重防重）。
- **JWT**：从 1 开始计的用户 id，`uid <= 0` 视为无效，防御伪造 `{"uid":0}` 的令牌。
- **生产环境务必更换 [config.yaml](config.yaml) 中 `jwt.secret_key` 的默认值**，否则任何人都能用已知密钥伪造令牌。推荐通过环境变量 `JWT_SECRET_KEY` 注入（优先于配置文件），避免密钥进入版本库。注意 `jwt.signing_method` 是必填项（仅支持 HS256/HS384/HS512），缺失会导致启动失败。

## 优雅停机与请求超时

这两个能力由 fiber 中间件 + 系统信号实现，配置入口统一在 [config.yaml](config.yaml)：

```yaml
http:
  # 请求超时时间(单位秒), 缺省 30 秒, 用于请求超时中间件
  timeout: 30
```

**请求超时 (408)**：超过 `http.timeout` 秒未完成处理的请求返回统一 JSON 408，避免慢接口长期占用连接。核心代码在 [routers/router.go](routers/router.go)：

```go
// 请求超时中间件 (超过 http.timeout 秒未完成返回 408 统一 JSON)
r.server.Fiber().App().Use(timeout.New(func(c fiber.Ctx) error {
    return c.Next()
}, timeout.Config{
    Timeout: time.Duration(r.server.Config().Http.Timeout) * time.Second,
    // 超时响应统一 JSON (避免纯文本)
    OnTimeout: func(c fiber.Ctx) error {
        return c.Status(fiber.StatusRequestTimeout).
            JSON(fiber.NewError(fiber.StatusRequestTimeout, "Request Timeout"))
    },
}))
```

**优雅停机**：收到 Ctrl+C / SIGTERM 时先停止接收新连接，等已处理中的请求完成后退出，避免强杀导致请求中断。核心代码在 [routers/router.go](routers/router.go) 的 `HttpStart`：

```go
// Ctrl+C / SIGTERM → 优雅停机 (等待进行中的请求完成后再退出)
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
return r.server.Fiber().App().Listen(r.server.Config().Http.Address(), fiber.ListenConfig{
    GracefulContext: ctx,
})
```

**验证**：

```bash
# 优雅停机: 启动服务后按 Ctrl+C, 观察服务先停止接收新连接再退出
go run main.go --config=./config.yaml
# 408: 将 config.yaml 的 http.timeout 临时调小 (如 1), 再访问会阻塞的接口即可看到 408
# (项目无内置慢接口, 可自行加一个 time.Sleep 的路由验证; Rust 端有对应集成测试可参考)
```

## 如何新增一个接口

以新增 `GET /api/v1/user/detail`（按 uid 查用户公开信息）为例，走一遍完整链路：

1. **定义请求 DTO**：在 [app/entity/user.go](app/entity/user.go) 添加

   ```go
   // UserDetailRequest 请求-获取用户详情
   type UserDetailRequest struct {
       Uid int `query:"uid" validate:"required"` // 用户id
   }
   ```

   注意：结构体上的 `validate` tag 只用于 swag 生成文档，项目未注册 validator，请求到达时不会自动校验（见该文件顶部注释）。

2. **添加控制器**：在 [app/controllers/user.go](app/controllers/user.go) 添加方法，模式固定为「绑定 → 校验 → 调 service → `ctx.JSON(resp)`」：

   ```go
   // Detail 获取用户详情
   // @Summary 获取用户详情
   // @Router   /api/v1/user/detail [get]
   func (c *User) Detail(ctx fiber.Ctx) error {
       req := &entity.UserDetailRequest{}
       if err := ctx.Bind().WithAutoHandling().Query(req); err != nil {
           return err
       }
       if req.Uid <= 0 {
           return xfiber.ParamError("无效的用户ID")
       }
       resp, err := c.service.User.Detail(ctx, req)
       if err != nil {
           return err
       }
       return ctx.JSON(resp)
   }
   ```

   需要鉴权时在方法开头调用 `middlewares.GetJwtClaims(ctx)` 取当前用户（见同文件 `Update`）。

3. **添加业务逻辑**：在 [app/services/user.go](app/services/user.go) 添加对应方法，复用模型查询，错误用 `fmt.Errorf("...: %w", err)` 包装后透传：

   ```go
   func (s *User) Detail(ctx fiber.Ctx, req *entity.UserDetailRequest) (*entity.User, error) {
       user := &models.User{}
       err := s.server.Gorm(ctx.Context()).
           Select(models.UserColumnsNoPassword).
           Where("id", req.Uid).
           Take(user).Error
       if err != nil {
           if errors.Is(err, gorm.ErrRecordNotFound) {
               return nil, xfiber.ParamError("用户不存在")
           }
           return nil, fmt.Errorf("user detail: query user by id: %w", err)
       }
       return entity.NewUserFromModel(user, "", ""), nil
   }
   ```

   注意查询统一 `Select(models.UserColumnsNoPassword)`，避免把 password 哈希拉进内存。

4. **注册路由**：在 [routers/user.go](routers/user.go) 的公开分组下追加一行：

   ```go
   groupRouter.Get("detail", r.app.Controller.User.Detail) // 获取用户详情
   ```

   需要 JWT 鉴权时移到 `groupRouter.Use(xfiber.UseJwt(...))` 之后注册（见 `update`）。

5. **文档与测试**：在控制器方法上写好 swag 注解（`@Summary`/`@Param`/`@Success`/`@Router`），改动后 `make swag` 重新生成文档；纯逻辑部分可补 `_test.go` 单测（见 [controller_test.go](app/controllers/controller_test.go)）。

## 单元测试

```bash
go test ./...
```

当前覆盖纯逻辑部分：昵称/密码格式校验、IP 限流器滑动窗口、枚举取值/合法性判断。

## 与 Rust 版本的关系

本仓库是整套 API 的参考原版，Rust 版（RsWebApi）由本仓库转写（差异对照见 [RsWebApi/README.md](../RsWebApi/README.md) 的「与 Go 原版的差异」），TS/PHP/Java/Python 版同样以本版为基准对齐实现。如需对比各端行为，请以各子目录 README 与仓库根 README 的「对照阅读」为准。

## 常见问题

- **启动报数据库连接失败**：检查 [config.yaml](config.yaml) 的 `orm.dsn` 用户名/密码/库名，并确认 MySQL 已启动、已导入根目录 [webapi.sql](../webapi.sql)。
- **端口 6600 被占用**：修改 `config.yaml` 的 `http.port`，或先用 `netstat -ano | findstr 6600` 定位占用进程。
- **Windows 下没有 `make swag`**：手动执行 `swag init` 后删除 `docs\docs.go`（见 [Makefile](Makefile) 顶部说明）。
- **`go run` 拉取依赖很慢**：国内可配置 `go env -w GOPROXY=https://goproxy.cn,direct`。
