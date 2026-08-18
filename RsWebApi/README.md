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

1. 准备 MySQL 数据库，导入表结构（与 Go 原版共用，email 为唯一索引，注册先查后插 + 数据库唯一约束双重防重）：

```bash
mysql -uroot -p < docs/schema.sql
```

2. 修改 [config.yaml](./config.yaml) 中的 `orm.dsn`，然后运行：

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

支持环境变量展开：配置值可写成 `${VAR}`（或 `$VAR`），加载时自动替换为环境变量值，未设置时替换为空字符串。常用注入示例：

| 环境变量 | 用途 |
| ------- | ---- |
| `MYSQL_DSN` | 覆盖 `orm.dsn`，避免把真实数据库连接串写进版本库 |
| `JWT_SECRET_KEY` | 覆盖 `jwt.secret_key`（生产必改，理由同 Go 版 README） |

`log.format: json` 时输出结构化 JSON 日志，便于 ELK / Loki 采集。

## make 命令

| 命令 | 说明 |
| ---- | ---- |
| `make check` | 快速编译检查 |
| `make test` | 运行集成测试（依赖本机 MySQL） |
| `make fmt` | 代码格式化（`cargo fmt --all`） |
| `make lint` | 静态检查（`cargo clippy --all-targets -- -D warnings`） |
| `make docs` | 重新生成 swagger 文档（utoipa 注解 → `docs/swagger.json\|yaml`，对齐 Go 端 `make swag`） |
| `make build` | 发布构建，可注入版本号 `make build GIT_VERSION=v1.0.0` |
| `make run` | 开发运行（默认加载 `./config.yaml`） |
| `make clean` | 清理构建产物 |
| `make all` | 依次执行 check → test → build |

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
| GET | `/api/v1/user/info` | 获取用户信息（公开接口，返回完整用户实体） | 无 |
| GET | `/api/v1/user/list` | 用户列表（分页，公开接口，返回完整用户实体） | 无 |
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

- **SQL 日志**：默认 warn 级（仅慢查询/错误），生产运行时零 SQL 日志开销；排查 SQL 问题时临时将 `log.level` 改为 `debug`
- **限流器**：按 IP 哈希 16 分片锁，同一 IP 串行计数、不同 IP 并行，锁竞争仅为单锁的 1/16；分片满时清理过期 key、拒绝新 key（安全背压），不做全表清零
- **Swagger 文档**：`/api/v1/swagger*` 由 utoipa 注解在运行期生成并缓存 (`LazyLock`), 零序列化开销; 修改注解后运行 `make docs` 重新生成 `docs/swagger.json|yaml` (对应 Go 端 `make swag`)
- 压测建议：`wrk -t4 -c100 -d30s http://127.0.0.1:6600/api/v1/health`

## 优雅停机与请求超时

**请求超时 (408)**：超过 `http.timeout` 秒未完成处理的请求返回统一 JSON 408（等价 Go 端 fiber `timeout.OnTimeout`）。配置在 [config.yaml](./config.yaml)：

```yaml
http:
  # 请求超时时间(单位秒), 缺省 30 秒, 用于请求超时中间件 (与 Go config.yaml 显式写 timeout: 30 等价)
  timeout: 30
```

中间件实现位于 [src/routes/mod.rs](./src/routes/mod.rs)，用 `tokio::time::timeout` 包裹后续处理：

```rust
/// 请求超时中间件: 超过 http.timeout 秒未完成返回 408 统一 JSON (对应 Go fiber timeout.OnTimeout)
async fn timeout_middleware(State(state): State<AppState>, req: Request, next: Next) -> Response {
    let timeout = Duration::from_secs(
        state.config.http.timeout.filter(|&t| t > 0).unwrap_or(30) as u64
    );
    match tokio::time::timeout(timeout, next.run(req)).await {
        Ok(resp) => resp,
        Err(_) => (
            StatusCode::REQUEST_TIMEOUT,
            Json(serde_json::json!({"code": 408, "message": "Request Timeout"})),
        ).into_response(),
    }
}
```

**优雅停机**：`axum::serve(...).with_graceful_shutdown(shutdown_signal())` 收到 Ctrl+C / SIGTERM 后停止接收新连接并等待进行中的请求完成。核心代码在 [src/main.rs](./src/main.rs)：

```rust
if let Err(e) = axum::serve(
    listener,
    app.into_make_service_with_connect_info::<std::net::SocketAddr>(),
)
.with_graceful_shutdown(shutdown_signal()) // Ctrl+C / SIGTERM → 优雅停机
.await
{
    panic!("server error: {e}");
}

/// 等待退出信号: Ctrl+C / SIGTERM, 触发优雅停机
async fn shutdown_signal() {
    #[cfg(unix)]
    {
        let mut term = tokio::signal::unix::signal(tokio::signal::unix::SignalKind::terminate())
            .expect("install SIGTERM handler failed");
        tokio::select! {
            _ = tokio::signal::ctrl_c() => tracing::info!("收到 Ctrl+C, 开始优雅停机"),
            _ = term.recv() => tracing::info!("收到 SIGTERM, 开始优雅停机"),
        }
    }
    #[cfg(not(unix))]
    {
        tokio::signal::ctrl_c().await.expect("install ctrl_c handler failed");
        tracing::info!("收到 Ctrl+C, 开始优雅停机");
    }
}
```

**验证**：

```bash
# 408: 集成测试已覆盖 (request_timeout_returns_408), 直接运行
cargo test request_timeout_returns_408
# 优雅停机: cargo run 启动后按 Ctrl+C, 观察日志输出 "收到 Ctrl+C, 开始优雅停机"
```

> 说明：这两项是 Go/Rust/Java/Python 应用层中间件能力（Java 见 `TimeoutFilter`、Python 见 `timeout_middleware`）；TS(Bun)/PHP(FPM) 端受框架限制未实现——Bun 无内置请求超时中间件，PHP 应由 Nginx 等服务器层配置超时。

## 已知限制

| 项 | 说明 |
| -- | ---- |
| 限流器为进程内实现 | 多实例部署时每实例独立计数，实际限流上限会随实例数放大；如需全局限流需引入 Redis 等共享存储 |
| `bcrypt` cost 固定 10 | 对齐 Go 原版 `DefaultCost`；OWASP 现建议 12+，升级会拖慢注册/登录接口（约 4 倍耗时） |
| `/user/list` 的 `count(*)` | 表数据量大时全表 count 较慢（InnoDB 全扫），与 Go 版语义一致暂不缓存；可用计数表 / 游标分页优化 |
| 密码长度 6~64 位 | 上限 64 防 bcrypt 72 字节静默截断（六端一致），下限沿用 Go 原版 6 位 |

## 请求处理流程（中间件顺序）

```
请求 → 请求日志(Request-ID) → panic恢复 → CORS → gzip 压缩 → 超时控制(408) → IP 限流 → 业务路由
```

## 与 Go 原版的差异

| 项 | Go 原版 | Rust 版 | 说明 |
| -- | ------- | ------- | ---- |
| 成功响应 | 裸实体 JSON | 裸实体 JSON | 已对齐 |
| 参数错误 | HTTP 422（`xfiber.ParamError`） | HTTP 422 | 已对齐 Go 实际行为 |
| uid 类型 | `int` | `i32` | 与 Go 版一致, 兼容已存在的有符号 INT 旧库 (无需改表) |
| bcrypt cost | 10（`DefaultCost`） | 10 | 已对齐 |
| 密码哈希 | gonet/crypto | bcrypt crate | 算法一致（bcrypt） |
| ORM | GORM | sea-orm | SQL 日志默认 warn 级（仅错误；Go 侧 gorm 可配 log_level） |
| Web 框架 | Fiber | axum | 均支持中间件链 |
| 公开接口返回 | 完整 User | 完整 User | 已对齐（六端 `/user/info`、`/user/list` 均返回完整用户实体，不返回 password） |
| 额外能力 | - | 健康检查 / 限流 / 优雅停机 / gzip / 超时 / 统一日志 / CI | 增强特性 |

## 如何新增一个接口

以新增 `GET /api/v1/user/detail`（按 uid 查用户公开信息）为例，走一遍完整链路：

1. **定义请求 DTO**：在 [src/app/entity/user.rs](./src/app/entity/user.rs) 添加
   ```rust
   #[derive(Debug, Deserialize)]
   pub struct UserDetailRequest {
       /// 用户id
       #[serde(default)]
       pub uid: u32,
   }
   ```
   注意 `#[serde(default)]`：Query/JSON 缺省字段不报错（对齐 Go 的零值语义）。
2. **添加控制器**：在 [src/app/controllers/user.rs](./src/app/controllers/user.rs) 添加 handler，模式固定为
   ```rust
   pub async fn detail(
       State(state): State<AppState>,
       QueryParams(req): QueryParams<UserDetailRequest>,
   ) -> Result<Json<User>, ApiError> {
       // 参数校验 → 调 service → 返回裸实体 JSON
       let resp = services::user::detail(&state, &req).await?;
       Ok(Json(resp))
   }
   ```
   需要鉴权时追加 `JwtAuth(claims): JwtAuth` 参数即可。
3. **添加业务逻辑**：在 [src/app/services/user.rs](./src/app/services/user.rs) 添加对应函数，复用 `user_model::find_by_id`，错误用 `?` 传播（`DbErr` 自动转 500 统一文案）。
4. **添加数据访问**：表查询函数写在 [src/app/models/user.rs](./src/app/models/user.rs)；若已有等价查询（如 `find_by_id`）直接复用。
5. **注册路由**：在 [src/routes/user.rs](./src/routes/user.rs) 的 `Router::new()` 链上追加
   ```rust
   .route("/user/detail", get(ctrl::detail))
   ```
   需要 JWT 鉴权时改用 `.route("/user/detail", get(ctrl::detail)).route_layer(from_fn_with_state(state, jwt_middleware))`。
6. **文档与测试**：用 `///` 在 handler 上写接口说明，`cargo doc --open` 可生成参考文档（本项目注释质量较高，`cargo doc` 产物很适合入门阅读）；接口涉及数据库时在 [src/tests.rs](./src/tests.rs) 补一条集成测试。

## 测试

依赖本机 MySQL（读取 config.yaml 配置）：

```bash
cargo test
```

代码文档（基于 `///` 注释生成，含模块结构图与接口说明）：

```bash
cargo doc --open
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

## 常见问题

- **`cargo build` 下载依赖很慢**：国内可配置 crates 镜像（如 rsproxy.cn）加速。
- **`cargo test` 失败**：集成测试依赖本机 MySQL（读取 config.yaml 的 `orm.dsn`），先确认数据库可用并已导入表结构（`docs/schema.sql`）。
- **端口 6600 被占用**：修改 `config.yaml` 的 `http.port`。
- **限流看起来不生效**：确认 `config.yaml` 的 `limit.enable: true`；多实例部署时进程内限流上限会被放大（见「已知限制」）。
