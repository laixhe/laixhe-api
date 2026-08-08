# WebApi (PHP / Laravel 13)

由 GoWebApi 同步转换而来的用户认证与用户管理 API 服务，接口与 Go 端保持一致。

- 框架：Laravel **13.24.0**
- 语言：PHP **8.4+**（实测 8.5.x）
- 数据库：MySQL（表结构与 Go 端共用，含 `user`、`user_extend`、`user_third_party`、`config_common`）

## 环境要求

| 依赖 | 版本 |
|---|---|
| PHP | >= 8.4（Laravel 13 要求 >= 8.3） |
| Composer | 2.x |
| MySQL | 5.7+ / MariaDB 10.3+ |
| 扩展 | ctype、fileinfo、mbstring、openssl、pdo_mysql 等（`composer check-platform-reqs` 可自检） |

## 安装与运行

```bash
# 1. 安装依赖
composer install

# 2. 配置 .env（数据库、JWT 密钥等）
cp .env.example .env   # 项目已内置 .env，可手动修改

# 3. 生成应用密钥
php artisan key:generate

# 4. 启动服务（默认 8000 端口）
php artisan serve
```

生产环境建议额外执行：

```bash
php artisan config:cache
php artisan route:cache
```

## 关键配置（.env）

| 变量 | 说明 | 默认 |
|---|---|---|
| `DB_*` | MySQL 连接 | 127.0.0.1:3306 / webapi |
| `JWT_SECRET` | JWT HMAC 密钥（**生产必改**，为空时启动即报错） | 与 Go 端一致 |
| `JWT_EXPIRE_TIME` | token 过期时长（秒） | 2592000（30 天） |
| `RATE_LIMIT_ENABLE` | 是否启用 IP 限流 | true |
| `RATE_LIMIT_MAX` | 单 IP 窗口内最大请求数 | 1000 |
| `RATE_LIMIT_WINDOW` | 滑动窗口时长（秒） | 60 |
| `CACHE_STORE` | 缓存驱动（生产高并发建议 `redis`） | file |
| `LOG_CHANNEL` / `LOG_LEVEL` | 日志通道与级别 | daily / debug |

## 接口列表

Base URL：`http://127.0.0.1:8000`（前缀 `/api/v1`）

| 方法 | 路径 | 鉴权 | 说明 |
|---|---|---|---|
| GET | `/api/v1/health` | 公开（限流豁免） | 健康检查（含数据库探测，DB 异常返回 503） |
| POST | `/api/v1/auth/register` | 公开 | 注册，返回 `{token, user}` |
| POST | `/api/v1/auth/login` | 公开 | 登录，返回 `{token, user}` |
| POST | `/api/v1/auth/refresh` | `Authorization: Bearer` | 刷新 JWT |
| GET | `/api/v1/user/info?uid=` | 公开 | 获取用户信息 |
| GET | `/api/v1/user/list?page=&page_size=` | 公开 | 用户列表（page 默认 1，page_size 默认 12、上限 100） |
| POST | `/api/v1/user/update` | `Authorization: Bearer` | 更新昵称/头像（uid 取自 JWT） |
| GET | `/api/v1/swagger.json` | 公开 | OpenAPI(Swagger) JSON 文档 |
| GET | `/api/v1/swagger.yaml` | 公开 | OpenAPI(Swagger) YAML 文档 |

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

### 统一响应格式

- 成功：直接返回业务数据 JSON
- 失败：`{"code": <HTTP状态码>, "message": "..."}`（HTTP 状态码同步，如 401/422/429/500）

## 校验规则（与 Go 端一致）

- 昵称：2~20 个字符
- 邮箱：标准邮箱格式
- 密码：`^[a-zA-Z0-9_@$]{6,}$`（字母/数字/`_`/`@`/`$`，至少 6 位）
- 头像地址：非空时最长 255 字符且必须以 `http`/`https` 开头
- 重复邮箱注册返回 422「邮箱已存在」

## 目录结构

```
app/
├── Helpers.php               # 全局响应函数 (response_success/error/exception, format_user)
├── Http/
│   ├── Controllers/          # AuthController / UserController / HealthController / SwaggerController
│   ├── Middleware/           # AssignRequestId / AuthJwt / RateLimit
│   ├── Requests/              # 参数校验 (LoginRequest / RegisterRequest / UserUpdateRequest)
│   └── Services/             # AuthService / UserService (业务逻辑)
├── Models/                   # User / UserExtend / UserThirdParty
├── Result/                   # 统一响应码 ResultCode / 响应体 Result
└── Utils/JwtUtil.php         # JWT 签发与签名校验 (HS256)
config/rate_limit.php         # IP 限流配置
public/swagger/               # OpenAPI 文档文件
routes/api.php                # API 路由
```

## 安全说明

- **JWT 签名校验**：所有受保护接口的 token 均校验 HMAC 签名（`SignedWith`）与有效期（`StrictValidAt`），伪造/篡改/过期 token 一律返回 401
- **密码**：bcrypt 哈希存储；模型 `$hidden` 防止密码被序列化泄露
- **限流**：全局 IP 滑动窗口限流，超阈值返回 429；`/api/v1/health` 豁免，避免探活被误伤
- **统一 500**：未捕获异常统一返回固定文案 `internal server error`，不泄露内部细节
