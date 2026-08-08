# WebApi (Rust 版)

由 Go 项目 [laixhe-api/GoWebApi](https://github.com/laixhe/laixhe-api/GoWebApi) 转写的 Rust 实现，技术栈：

- [axum](https://github.com/tokio-rs/axum) 0.8 — Web 框架
- [sea-orm](https://www.sea-orm.org/) 2.0 — ORM（支持 MySQL / PostgreSQL / SQLite）
- [jiff](https://docs.rs/jiff) — 时间处理
- [jsonwebtoken](https://docs.rs/jsonwebtoken) — JWT（HS256）
- [bcrypt](https://docs.rs/bcrypt) — 密码哈希
- [tracing](https://docs.rs/tracing) + [rolling-file](https://docs.rs/rolling-file) — 结构化日志（按大小轮转）

## 目录结构

```
RsWebApi/
├── src/
│   ├── main.rs                 # 入口：参数解析 / 优雅停机 / 启动
│   ├── config.rs               # YAML 配置加载（支持 ${ENV} 展开）+ 校验
│   ├── state.rs                # AppState：日志初始化 / 数据库连接池 / 限流器
│   ├── logger.rs               # 统一日志模块（Timer 计时 + log_elapsed! 宏）
│   ├── error.rs                # 统一错误 ApiError（400/401/500…）
│   ├── routes/                 # 路由注册（auth / user / swagger / health / 404 兜底）
│   ├── middleware/
│   │   ├── request_log.rs      # 全局请求日志 + X-Request-ID 关联
│   │   ├── rate_limit.rs       # IP 滑动窗口限流（429 统一 JSON）
│   │   └── jwt.rs              # JWT 鉴权中间件 + JwtAuth 提取器
│   ├── app/
│   │   ├── controllers/        # 控制器（auth / user / health）
│   │   ├── services/           # 业务逻辑（含关键步骤耗时日志）
│   │   ├── entity/             # 接口 DTO（请求/响应）
│   │   ├── models/             # sea-orm 实体（user / user_extend / ...）
│   │   └── util/               # 工具（邮箱/密码正则校验）
│   └── tests.rs                # 集成测试（依赖本机 MySQL）
├── docs/                       # swagger 文档（/api/v1/swagger.json|yaml）
└── config.yaml                 # 配置文件
```

## 快速开始

1. 准备 MySQL 数据库（表结构见 Go 原版或 [docs/schema.sql](./docs/schema.sql)），修改 [config.yaml](./config.yaml) 中的 `orm.dsn`。
2. 运行：

```bash
cargo run -- --config=./config.yaml
# 或省略参数（默认 ./config.yaml）
cargo run
# 或按环境加载 config.{env}.yaml（存在时使用，如 config.prod.yaml）
cargo run -- --env=prod
```

3. 验证：

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

## 配置说明（config.yaml）

| 节点 | 说明 |
| ---- | ---- |
| `http` | 监听地址与端口、请求超时时间 `timeout`（秒，缺省 30） |
| `log` | 日志模式（console/file）、输出格式（`format`：console/json）、级别、按大小轮转（max_size / max_backups） |
| `orm` | 数据库驱动、DSN、连接池大小（max_idle / max_open / max_life_time） |
| `jwt` | 签名密钥、过期时长（秒） |
| `limit` | 接口限流：`enable` 开关、`max` 窗口内最大请求数、`window` 窗口时长（秒）、`trust_proxy` 是否信任代理头 `X-Forwarded-For`（默认 `false`，仅在可信反向代理之后部署时开启，否则客户端可伪造该头绕过限流） |

支持环境变量展开，例如 `dsn: ${MYSQL_DSN}`。`log.format: json` 时输出结构化 JSON 日志，便于 ELK / Loki 采集。

## 响应格式

成功：直接返回业务实体 JSON（裸实体，与 Go 版 `ctx.JSON(resp)` 一致）：

```json
{ "token": "...", "user": { "uid": 1, "nickname": "testuser" } }
```

失败：统一错误格式：

```json
{ "code": 422, "message": "邮箱格式错误" }
```

### 错误码

| code | HTTP 状态 | 含义 |
| ---- | --------- | ---- |
| 422  | 422 | 参数错误（校验失败 / JSON 解析失败 / Query 解析失败，对齐 Go `xfiber.ParamError`） |
| 401  | 401 | 未授权（缺少 / 无效 JWT，或用户被禁用） |
| 404  | 404 | 路由不存在（统一 JSON 兜底） |
| 408  | 408 | 请求超时（超过 `http.timeout` 秒） |
| 429  | 429 | 触发接口限流 |
| 500  | 500 | 服务器内部错误（数据库 / bcrypt / JWT 签发等） |

## API 列表

| 方法 | 路径 | 说明 | 鉴权 |
| ---- | ---- | ---- | ---- |
| GET | `/api/v1/health` | 健康检查（含数据库探测） | 无 |
| GET | `/api/v1/swagger` | Swagger UI 文档页（浏览器访问） | 无 |
| GET | `/api/v1/swagger.json` / `.yaml` | 接口文档（原始 JSON/YAML） | 无 |
| POST | `/api/v1/auth/register` | 注册 | 无 |
| POST | `/api/v1/auth/login` | 登录 | 无 |
| POST | `/api/v1/auth/refresh` | 刷新 JWT | Bearer |
| GET | `/api/v1/user/info` | 获取用户信息（公开视图，不含 email/mobile/account） | 无 |
| GET | `/api/v1/user/list` | 用户列表（分页，公开视图，不含 email/mobile/account） | 无 |
| POST | `/api/v1/user/update` | 更新用户信息 | Bearer |

## 特性

- 成功返回裸实体 JSON、错误统一 `{code, message}` 格式（与 Go 版完全一致）
- 全局请求日志：自动生成 `X-Request-ID` 并贯穿业务日志，便于问题追踪
- IP 滑动窗口限流（16 分片锁，低竞争），超限返回 429 统一 JSON
- CORS 支持（前后端分离部署）
- 请求超时控制（`http.timeout`，返回 408）与响应 gzip 压缩
- 日志按大小轮转 + 优雅停机
- 集成测试：`cargo test`

## 性能优化说明

- **release 构建**（[Cargo.toml](./Cargo.toml)）：`opt-level=3` + `lto` + `codegen-units=1` 链接期优化，`strip` 减小二进制；构建命令：

  ```bash
  cargo build --release
  ```

- **SQL 日志**：默认 debug 级输出（对齐 Go `OrmWriter`），生产以 `info` 级别运行时零 SQL 日志开销
- **限流器**：按 IP 哈希 16 分片锁，同一 IP 串行计数、不同 IP 并行，锁竞争仅为单锁的 1/16
- **静态接口**：`/api/v1/swagger*` 使用 `include_str!` 内嵌常量 + 缓存头，零序列化开销
- 压测建议：`wrk -t4 -c100 -d30s http://127.0.0.1:6600/api/v1/health`

## 请求处理流程（中间件顺序）

```
请求 → 请求日志(Request-ID) → panic恢复 → CORS → gzip 压缩 → 超时控制(408) → IP 限流 → 业务路由
```

## 与 Go 原版的差异

| 项 | Go 原版 | Rust 版 | 说明 |
| -- | ------- | ------- | ---- |
| 成功响应 | 裸实体 JSON | 裸实体 JSON | 已对齐 |
| 参数错误 | HTTP 422（`xfiber.ParamError`） | HTTP 422 | 已对齐 Go 实际行为 |
| uid 类型 | `int` | `u32` | 数据库列为 INT UNSIGNED |
| bcrypt cost | 10（`DefaultCost`） | 10 | 已对齐 |
| 密码哈希 | gonet/crypto | bcrypt crate | 算法一致（bcrypt） |
| ORM | GORM | sea-orm | SQL 日志默认 debug 级（对齐 OrmWriter） |
| Web 框架 | Fiber | axum | 均支持中间件链 |
| 公开接口返回 | 完整 User | 脱敏 UserPublic | `/user/info`、`/user/list` 不再返回 email/mobile/account |
| 额外能力 | - | 健康检查 / 限流 / 优雅停机 / gzip / 超时 / 统一日志 / CI | 增强特性 |

## 测试

依赖本机 MySQL（读取 config.yaml 配置）：

```bash
cargo test
```

## Docker 构建

```bash
docker build -t webapi:latest .
docker run --rm -p 6600:6600 -v /path/to/config.yaml:/webapi/config.yaml webapi:latest
```

## 构建产物

```bash
make build   # 等价于 cargo build --release，可注入 GIT_VERSION
```
