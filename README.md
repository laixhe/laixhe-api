# laixhe-api

同一套用户认证 / 用户管理 RESTful API 的 **六种语言实现**（Go / Java / PHP / Python / Rust / TypeScript），接口行为、返回格式、校验规则、错误码全部对齐，适合对比学习不同技术栈的后端开发。

## 仓库结构

| 目录 | 语言 / 框架 | 定位 | 端口 |
| --- | --- | --- | --- |
| [GoWebApi](GoWebApi/README.md) | Go + Fiber v3 + GORM | 原版，结构最清晰 | 6600 |
| [RsWebApi](RsWebApi/README.md) | Rust + axum + sea-orm | Go 版转写，性能优化最深入 | 6600 |
| [TsWebApi](TsWebApi/README.md) | TypeScript + Bun + Elysia + Prisma | 现代 JS 全栈体验 | 6600 |
| [PhpWebApi](PhpWebApi/README.md) | PHP + Laravel 13 | 生态成熟、上手门槛低 | 8000 |
| [JavaWebApi](JavaWebApi/README.md) | Java + Spring Boot 4 + JPA | 企业级生态，可编译 GraalVM 原生镜像 | 6600 |
| [PyWebApi](PyWebApi/README.md) | Python + FastAPI + SQLModel + uv | 现代异步 Python 体验 | 6600 |

公共资源（各端共用）：

- [webapi.sql](webapi.sql) — MySQL 建表脚本（`user` / `user_extend` / `user_third_party` / `config_common`）
- [JWT.md](JWT.md) — JWT 数据结构与使用方式说明

> **文档导航**：六个子目录 README 围绕同一套 API 展开，章节结构大体对齐（快速开始、关键配置、接口列表、统一响应/错误码、校验规则、与 Go 原版的差异、如何新增一个接口、常见问题等），可互相参照。想快速跑起来看「快速开始」，想深入源码看「如何新增一个接口」，想横向对比技术栈看「与 Go 原版的差异」。

## 快速开始

各端详细步骤见对应子目录 README，这里给最小启动路径：

```bash
# 1. 初始化数据库（任选一种语言即可，六端共用同一张表）
mysql -uroot -p < webapi.sql

# 2. 任选一端启动（默认监听 0.0.0.0:6600，PHP 为 8000）
# Go
cd GoWebApi && go run main.go --config=./config.yaml
# Rust
cd RsWebApi && cargo run
# TypeScript（需要 Bun）
cd TsWebApi && bun install && bun run dev
# PHP（需要 Composer + PHP 8.4+）
cd PhpWebApi && composer install && php artisan key:generate && php artisan serve
# Java（需要 JDK 25 / GraalVM；也可用 H2 profile 免装 MySQL）
cd JavaWebApi && ./gradlew bootRun
# Python（需要 uv + Python 3.14+）
cd PyWebApi && uv sync && uv run uvicorn app.main:app --host 0.0.0.0 --port 6600
```

> 各端均需先修改数据库连接配置：Go/Rust 改 `config.yaml` 的 `orm.dsn`，TS/Python 改 `.env` 的 `DATABASE_URL`，PHP 改 `.env` 的 `DB_*`，Java 改 `application.yaml` 的 `spring.datasource`（或用 `./gradlew bootRun --args='--spring.profiles.active=h2'` 切换到内置 H2 内存库，免装 MySQL）。

## API 一览

所有接口前缀 `/api/v1`，鉴权接口需携带请求头 `Authorization: Bearer <token>`。

| 方法 | 路径 | 鉴权 | 说明 |
| --- | --- | --- | --- |
| GET | /api/v1/health | 无（限流豁免） | 健康检查（含数据库探测） |
| POST | /api/v1/auth/register | 无 | 注册，返回 token + 用户信息 |
| POST | /api/v1/auth/login | 无（限流 5 次/分钟） | 登录 |
| POST | /api/v1/auth/refresh | Bearer | 刷新 token |
| GET | /api/v1/user/info?uid=1 | 无 | 查询用户信息 |
| GET | /api/v1/user/list?page=1&page_size=12 | 无 | 用户列表（按 id 倒序分页） |
| POST | /api/v1/user/update | Bearer | 更新昵称/头像（uid 取自 JWT） |
| GET | /api/v1/swagger.json / .yaml | 无 | OpenAPI 文档 |

## 统一约定（六端一致）

- **成功响应**：直接返回业务实体 JSON（如 `{ "token": "...", "user": {...} }`）
- **失败响应**：`{ "code": <HTTP状态码>, "message": "..." }`，HTTP 状态码同步
- **错误码**：400 请求格式错误（如参数非数字，六端一致）/ 422 参数错误 / 401 未授权 / 404 路由不存在 / 408 请求超时 / 429 触发限流 / 500 内部错误（统一固定文案，不泄露内部细节）
- **参数校验**：昵称 2~20 字符；邮箱标准格式；密码 `^[a-zA-Z0-9_@$]{6,64}$`（上限 64 防 bcrypt 72 字节静默截断）；头像地址非空时最长 255 且须以 `http://`/`https://` 开头
- **安全设计**：密码 bcrypt 哈希存储；查询默认排除 `password` 列；封禁账号与密码错误返回同一提示；登录限流防暴力破解
- **限流**：进程内实现（滑动窗口，Java 版为令牌桶、效果等价），单实例部署有效；多实例需换 Redis 等共享存储

## 最短学习路径

想快速理解整个项目，按顺序精读 **Go 端 6 个文件**即可（六端实现完全对应，读透后其他端都能看懂）：

1. [main.go](GoWebApi/main.go) — 程序入口，了解服务启动流程
2. [routers/router.go](GoWebApi/routers/router.go) — 路由注册与中间件链（请求处理顺序）
3. [app/controllers/auth.go](GoWebApi/app/controllers/auth.go) — 控制器层：参数绑定/校验、返回响应
4. [app/services/auth.go](GoWebApi/app/services/auth.go) — 服务层：业务逻辑、数据库操作
5. [app/models/user.go](GoWebApi/app/models/user.go) — 模型层：表结构、数据访问
6. [core/middlewares/jwt.go](GoWebApi/core/middlewares/jwt.go) — 中间件：JWT 鉴权

**对照阅读**：同一接口在六端的对应文件

| 层 | Go | Rust | TS | PHP | Java | Python |
| --- | --- | --- | --- | --- | --- | --- |
| 路由 | [routers/auth.go](GoWebApi/routers/auth.go) | [routes/auth.rs](RsWebApi/src/routes/auth.rs) | [routes/auth.ts](TsWebApi/src/routes/auth.ts) | [routes/api.php](PhpWebApi/routes/api.php) | [AuthController.java](JavaWebApi/src/main/java/com/laixhe/webapi/controller/AuthController.java) | [api/auth.py](PyWebApi/app/api/auth.py) |
| 控制器 | [controllers/auth.go](GoWebApi/app/controllers/auth.go) | [controllers/auth.rs](RsWebApi/src/app/controllers/auth.rs) | [routes/auth.ts](TsWebApi/src/routes/auth.ts) | [AuthController.php](PhpWebApi/app/Http/Controllers/AuthController.php) | [AuthController.java](JavaWebApi/src/main/java/com/laixhe/webapi/controller/AuthController.java) | [api/auth.py](PyWebApi/app/api/auth.py) |
| 服务 | [services/auth.go](GoWebApi/app/services/auth.go) | [services/auth.rs](RsWebApi/src/app/services/auth.rs) | [routes/auth.ts](TsWebApi/src/routes/auth.ts) | [AuthService.php](PhpWebApi/app/Http/Services/AuthService.php) | [AuthService.java](JavaWebApi/src/main/java/com/laixhe/webapi/service/AuthService.java) | [services/auth.py](PyWebApi/app/services/auth.py) |
| 模型 | [models/user.go](GoWebApi/app/models/user.go) | [models/user.rs](RsWebApi/src/app/models/user.rs) | [schema.prisma](TsWebApi/prisma/schema.prisma) | [User.php](PhpWebApi/app/Models/User.php) | [User.java](JavaWebApi/src/main/java/com/laixhe/webapi/entity/User.java) | [models/user.py](PyWebApi/app/models/user.py) |

> 提示：TS 端业务逻辑写在路由 handler 内（Elysia 为薄框架），一个文件对应其他端的 controller + service 两层；Java 用 `@RestController` 注解即路由（无独立路由文件），数据访问独立在 `repository` 层；Python 路由 handler 只做参数绑定/校验，业务逻辑在 `services` 层。各接口上的 Swagger 文档注解（`@openapi`/`#[OA\...]`/`#[utoipa::path]`/`@Summary`）仅供生成接口文档，阅读业务逻辑时可跳过。

## 其他学习路线

- 想学**性能优化**看 Go/Rust/TS 版：Go（bcrypt BcryptPool worker 池可配置、16 分片限流、滑动窗口）；Rust（bcrypt `spawn_blocking` 异步化、16 分片限流、滑动窗口）；TS（原生异步 bcrypt、count 缓存 single-flight）
- 想快速跑起来 / 熟悉 Laravel 生态看 PHP 版
- 想熟悉 Spring Boot 企业级分层 / GraalVM 原生编译（`./gradlew nativeCompile`）看 Java 版，自带 H2 profile 免装 MySQL 即可跑通
- 想体验 Python 异步 Web 开发（FastAPI + SQLModel + uv）看 Python 版
- 六端代码互相对齐，同一接口可以打开六个实现的同名文件对照阅读

## 安全提醒

- **生产环境务必更换 JWT 密钥**：Go/Rust 改 `config.yaml` 的 `jwt.secret_key`，Java 改 `application.yaml` 的 `app.jwt.secret-key`，TS/Python 改 `.env` 的 `JWT_SECRET_KEY`，PHP 改 `.env` 的 `JWT_SECRET`（Java/Python 也支持环境变量 `JWT_SECRET_KEY` 覆盖），默认值任何人都能用来伪造令牌
- `/api/v1/user/info`、`/api/v1/user/list` 是**公开接口**（未挂 JWT），用于教学演示；上生产前请评估是否需要鉴权或字段脱敏
- IP 限流优先信任 `X-Forwarded-For`（Rust 版默认不信任，可配 `trust_proxy`），直接面向公网时客户端可伪造该头绕过限流，建议只在可信反向代理后部署
