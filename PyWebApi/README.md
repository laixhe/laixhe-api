# WebApi (Python 版)

由 GoWebApi 对齐转换的 Python 实现，接口、返回格式、校验规则、错误码与其余各端保持一致。

- 框架：[FastAPI](https://fastapi.tiangolo.com/)（>= 0.141.1）+ [Uvicorn](https://www.uvicorn.org/)
- ORM：[SQLModel](https://sqlmodel.tiangolo.com/) + PyMySQL（支持 SQLite 切换）
- 语言：Python **3.14+**，包管理使用 [uv](https://docs.astral.sh/uv/)
- JWT：[PyJWT](https://pyjwt.readthedocs.io/)（HS256）
- 密码：[bcrypt](https://github.com/pyca/bcrypt) 5.x（cost=10，与 Go 端对齐）
- 配置：pydantic-settings（.env / 环境变量）

## 环境要求

| 依赖 | 版本 |
|---|---|
| Python | >= 3.14 |
| uv | 最新稳定版 |
| MySQL | 5.7+（不装也行：`DATABASE_URL` 切 SQLite） |

## 快速开始

```bash
# 1. 安装依赖
uv sync

# 2. 配置 .env（数据库、JWT 密钥等）
cp .env.example .env

# 3. 启动（默认监听 0.0.0.0:6600；启动时自动建表）
uv run uvicorn app.main:app --host 0.0.0.0 --port 6600
```

验证：

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

## 关键配置（.env）

| 变量 | 说明 | 默认 |
|---|---|---|
| `HOST` / `PORT` | 监听地址与端口 | 0.0.0.0 / 6600 |
| `HTTP_TIMEOUT` | 请求超时（秒），超时返回 408 | 30 |
| `DATABASE_URL` | 数据库连接串（`mysql+pymysql://...`；SQLite 示例 `sqlite:///./pywebapi.db`） | mysql+pymysql://root:123456@127.0.0.1:3306/webapi |
| `JWT_SECRET_KEY` | JWT HMAC 密钥（**生产必改**） | 与 Go 端一致 |
| `JWT_SIGNING_ALGORITHM` | 签名算法 | HS256 |
| `JWT_EXPIRE_SECONDS` | token 过期时长（秒） | 2592000（30 天） |
| `LIMIT_ENABLE` / `LIMIT_MAX` / `LIMIT_WINDOW` | IP 限流开关 / 窗口内最大请求数 / 窗口时长（秒） | true / 1000 / 60 |

配置经 pydantic-settings 加载到 [config.py](app/core/config.py)，支持环境变量与 .env 覆盖。

## 接口列表

Base URL：`http://127.0.0.1:6600`（前缀 `/api/v1`）

| 方法 | 路径 | 鉴权 | 说明 |
|---|---|---|---|
| GET | `/api/v1/health` | 公开（限流豁免） | 健康检查（含数据库探测，DB 异常返回 503） |
| POST | `/api/v1/auth/register` | 公开 | 注册，返回 `{token, user}` |
| POST | `/api/v1/auth/login` | 公开 | 登录，返回 `{token, user}` |
| POST | `/api/v1/auth/refresh` | `Authorization: Bearer` | 刷新 JWT |
| GET | `/api/v1/user/info?uid=` | 公开 | 获取用户信息 |
| GET | `/api/v1/user/list?page=&page_size=` | 公开 | 用户列表（page 默认 1，page_size 默认 12、上限 100） |
| POST | `/api/v1/user/update` | `Authorization: Bearer` | 更新昵称/头像（uid 取自 JWT） |
| GET | `/api/v1/swagger.json` | 公开 | OpenAPI JSON 文档（FastAPI 按代码自动生成） |
| GET | `/api/v1/swagger.yaml` | 公开 | OpenAPI YAML 文档 |
| GET | `/api/v1/swagger` | 公开 | Swagger UI 页面 |

### 用户对象结构

```json
{
  "uid": 1,
  "type_id": 1,
  "account": "全局唯一账号",
  "mobile": "",
  "email": "user@example.com",
  "nickname": "昵称",
  "avatar_url": "",
  "sex": 0,
  "states": 1,
  "created_at": "2026-08-08 17:00:00"
}
```

## 统一响应格式与错误码

- 成功：直接返回业务数据 JSON（`{token, user}` / `User` / `UserListResponse`）
- 失败：`{"code": <HTTP状态码>, "message": "..."}`，HTTP 状态码同步

| code | HTTP | 含义 |
|---|---|---|
| 400 | 400 | 请求格式错误（顶层 body 非对象等，由全局 handler 兜底） |
| 401 | 401 | 未授权（缺少或无效 JWT、用户被禁用） |
| 404 | 404 | 路由不存在（统一 JSON 兜底） |
| 408 | 408 | 请求超时（超过 `HTTP_TIMEOUT` 秒） |
| 422 | 422 | 参数错误（Pydantic 校验失败 / 手写 `validators` 校验 / 业务错误） |
| 429 | 429 | 触发 IP 限流 |
| 500 | 500 | 内部错误（统一固定文案 `internal server error`，不泄露细节） |
| 503 | 503 | 健康检查数据库不可用 |

## 校验规则（与 Go 端一致）

- 昵称：2~20 个字符（`len()` 按字符计数，中文按 1 字符计）
- 邮箱：`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$`
- 密码：`^[a-zA-Z0-9_@$]{6,64}$`（6~64 位；上限 64 防 bcrypt 72 字节静默截断）
- 头像地址：非空时最长 255 且必须以 `http://`/`https://` 开头
- 重复邮箱注册返回 422「邮箱已存在」

> 路由 handler 不做业务校验，统一调用 [validators.py](app/utils/validators.py) 的手写校验函数，失败返回 422 与具体中文文案（与其余各端对齐）；Pydantic 负责缺字段/类型错误的兜底（422）。

## 目录结构

```
app/
├── main.py                 # 应用入口：中间件注册、异常处理器、路由挂载
├── api/                    # 路由（auth / user / health / swagger / deps）
│   └── deps.py             # 依赖：从 Bearer 令牌解析当前 uid（401）
├── core/                   # 配置 / 统一错误 / 日志 / 通用中间件 / IP 限流
├── db/                     # SQLAlchemy 引擎与会话
├── models/                 # SQLModel 模型（user / user_extend / user_third_party / config_common）
├── schemas/                # Pydantic 请求/响应模型
├── security/               # JWT 签发/校验、bcrypt 密码哈希
├── services/               # 业务逻辑（AuthService / UserService）
└── utils/                  # 参数校验
pyproject.toml              # 依赖与元数据（uv 管理）
.env.example                # 环境变量示例
```

调用链：`api(路由) → services(业务) → models(表) → db(引擎)`。路由 handler 只做参数绑定/校验与响应返回，业务逻辑在 `services` 层。

## 中间件顺序

```
请求 → requestId → 访问日志 → CORS → gzip 压缩 → 请求超时(408) → IP 限流(429) → 业务路由
```

（Starlette 后注册的中间件在最外层，因此 [main.py](app/main.py) 按相反顺序注册，与 Go 端中间件链对齐。）

## 技术要点

### 1. IP 滑动窗口限流

[rate_limit.py](app/core/rate_limit.py) 为每个 IP 维护一个 `deque` 时间戳队列（FIFO，队首过期清理），`threading.Lock` 保证线程安全；默认 1000 次/60s。健康检查路径豁免；`max_keys` 内存保护：key 数超限时清理已无活动窗口的 key。

### 2. 请求超时（408）

[timeout_middleware](app/core/middlewares.py) 用 `asyncio.wait_for` 包裹后续处理，超过 `HTTP_TIMEOUT` 秒返回 408 统一 JSON。

### 3. 自动 OpenAPI 文档

FastAPI 根据路由注解自动生成 OpenAPI 文档，[swagger.py](app/api/swagger.py) 将 `/api/v1/swagger.json|yaml` 直接序列化自 `app.openapi()`（OpenAPI 3.1），接口定义保持单一来源，响应带 5 分钟缓存。

### 4. bcrypt cost 对齐

[password.py](app/security/password.py) 显式指定 `bcrypt.gensalt(rounds=10)`——bcrypt 5.x 默认 cost 为 12，显式指定保持与 Go 端 `DefaultCost(10)` 相同的哈希成本。

### 5. 启动自动建表

[init_db](app/db/database.py) 在应用启动（lifespan）时执行 `SQLModel.metadata.create_all` 自动建表，`account`/`email` 为唯一索引；`DATABASE_URL` 切 SQLite 即可免装 MySQL。

### 6. 统一异常处理

[errors.py](app/core/errors.py) 统一返回 `{code, message}`：

- 业务错误（`APIError`，如 422/401/503）状态码与文案原样返回
- Pydantic 校验失败 → 422（`validation error: 字段: 错误`）
- Starlette HTTPException（含 404）→ `{code, message}`
- 其余未知错误：记录服务端日志（含 requestId 与路径）后统一返回固定 500 文案，不泄露内部细节

## 如何新增一个接口

以新增 `GET /api/v1/user/detail`（按 uid 查用户公开信息）为例，走一遍完整链路：

1. **添加路由**：在 [app/api/user.py](app/api/user.py) 追加，模式固定为「FastAPI 依赖注入 → 校验 → `UserService(session)` → 返回 schema」：

   ```python
   @router.get(
       "/detail",
       response_model=UserSchema,
       summary="获取用户详情",
       responses=error_responses,
   )
   def detail(uid: int = Query(description="用户id"), session: Session = Depends(get_session)) -> UserSchema:
       if uid <= 0:
           raise param_error("无效的用户ID")
       return UserService(session).detail(uid)
   ```

   需要 JWT 鉴权时把参数改为 `uid: int = Depends(get_current_uid)`（见 `update`）。

2. **添加业务逻辑**：在 [app/services/user.py](app/services/user.py) 添加方法，复用 `self.session.get(User, uid)`；已知业务错误抛 `param_error(...)`：

   ```python
   def detail(self, uid: int) -> UserSchema:
       user = self.session.get(User, uid)
       if user is None:
           raise param_error("用户不存在")
       return UserSchema.from_model(user)
   ```

3. **注册路由**：无需额外步骤——`user.router` 已在 [app/api/router.py](app/api/router.py) 中 `include_router`。

4. **文档与测试**：FastAPI 根据路由注解自动生成 OpenAPI（`/api/v1/swagger` 页面可见，无需重新生成）；纯逻辑可补 pytest 单测（仓库暂未内置 tests 目录，参考其余各端测试用例）。

## 已知限制

| 项 | 说明 |
|---|---|
| 限流为进程内实现 | 基于 `threading.Lock` 的进程内滑动窗口，多 worker/多实例部署时需换 Redis 等共享存储 |
| 启动自动建表 | `create_all` 为教学简化，表结构变更请使用 [根目录 webapi.sql](../../webapi.sql) 或引入 alembic 迁移 |
| 无内置自动化测试 | 仓库暂未提供 tests 目录；接口行为可参考其余各端测试用例（如 Go `controller_test.go` / Java `ApiIntegrationTests`） |
| 密码长度 6~64 位 | 上限 64 防 bcrypt 72 字节静默截断（六端一致），下限沿用 Go 原版 6 位 |

## 安全说明

- **生产环境务必更换 JWT 密钥**：`.env` 的 `JWT_SECRET_KEY`，默认值任何人都能用来伪造令牌
- `/api/v1/user/info`、`/api/v1/user/list` 是**公开接口**（未挂 JWT），用于教学演示；上生产前请评估是否需要鉴权或字段脱敏
- IP 限流优先信任 `X-Forwarded-For`，直接面向公网时客户端可伪造该头绕过限流，建议仅在可信反向代理后部署

## 常见问题

- **`uv sync` 下载依赖很慢**：国内可配置 uv 使用镜像源（如清华 PyPI 镜像）加速。
- **Python 版本不符**：项目要求 Python 3.14+（见 [pyproject.toml](pyproject.toml)），可用 `uv python install 3.14` 安装。
- **端口 6600 被占用**：修改 `.env` 的 `PORT`。
- **不想装 MySQL**：把 `.env` 的 `DATABASE_URL` 改为 `sqlite:///./pywebapi.db`（见「关键配置」）。
