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

# 4. 初始化数据库（建表，account / email 均为唯一索引）
php artisan migrate

# 5. 启动服务（默认 8000 端口）
php artisan serve
```

验证：

```bash
# 健康检查
curl http://127.0.0.1:8000/api/v1/health

# Swagger UI 文档页（浏览器打开）
http://127.0.0.1:8000/api/v1/swagger

# 注册
curl -X POST http://127.0.0.1:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"nickname":"test","email":"test@example.com","password":"abc123"}'
```

生产环境建议额外执行：

```bash
php artisan config:cache
php artisan route:cache
```

> **注意**：执行 `config:cache` 后 `env()` 在配置文件之外一律返回 `null`。
> 因此 JWT 密钥必须通过 `config('jwt.*')` 读取（见 [config/jwt.php](config/jwt.php) 顶部说明），
> 本项目已按要求实现，`config:cache` 不会导致鉴权失效。

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
| GET | `/api/v1/swagger` | 公开 | Swagger UI 可视化页面（浏览器访问） |

> 文档由 swagger-php 根据代码注解生成：修改 `#[OA\...]` 注解后执行 `composer swagger`（或 `php scripts/generate-swagger.php`）重新生成 `public/swagger/swagger.json|yaml`，对应 Go 端 `make swag`。

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
- 密码：`^[a-zA-Z0-9_@$]{6,64}$`（字母/数字/`_`/`@`/`$`，6~64 位；上限 64 防 bcrypt 72 字节静默截断）
- 头像地址：非空时最长 255 字符且必须以 `http://`/`https://` 开头
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
config/jwt.php                # JWT 配置 (经 config() 读取, 兼容 config:cache)
config/rate_limit.php         # IP 限流配置
database/migrations/          # 建表迁移 (user / user_extend / user_third_party)
public/swagger/               # OpenAPI 文档文件
routes/api.php                # API 路由
tests/                        # PHPUnit 单元测试 (JWT 签验等纯逻辑)
phpunit.xml                   # PHPUnit 配置
```

## 测试

```bash
composer test    # 等价于 vendor/bin/phpunit
```

- **单元测试**（`tests/Unit`）：JWT 签发/解析/验签/伪造拒绝、uid<=0 拒绝，无需数据库
- **集成测试**（`tests/Feature`）：注册/登录/参数校验/限流 429，使用 sqlite 内存库（`RefreshDatabase` 自动跑 migrations），不依赖外部 MySQL

## 教学取舍说明（重要）

本项目作为**跨语言对照教学仓库**，部分写法刻意偏离了 Laravel 生产实践，以与其余五端逐行对齐。**请勿照抄到生产项目**：

- **未使用 FormRequest 自动校验**：控制器里用 `new RegisterRequest()` + 手动调用 `$request->validator($req)`（见 [AuthController.php](app/Http/Controllers/AuthController.php)），等价于其余各端的"绑定层 + 手写校验"流程。生产项目应使用 Laravel 的依赖注入（在控制器方法参数声明 `RegisterRequest $request`），由框架自动完成校验与 422 响应。
- **Service 用 `new` 实例化而非容器 DI**：如 `new AuthService()`（见 [AuthController.php](app/Http/Controllers/AuthController.php)）。生产项目应通过构造函数注入或 `app(AuthService::class)` 解析，以利用容器的依赖解析与测试替换能力。
- **封禁账号/密码错误统一 422**：故意不用 401，为与六端错误码对齐；语义上这类场景用 422 是教学取舍，生产可评估使用 401/403。

对应地，本仓库真正的生产级做法（如 bcrypt 时序侧信道防护、Snowflake 单例、Redis 限流）都在各"技术要点"中保留了，取舍时注意区分"教学简化"与"生产必要"。

## 技术要点

### 1. 缓存容器别名机制

> 排查记录（2026-08-08）：`app(\Illuminate\Contracts\Cache\Repository::class)` 在 Laravel 13 中依然可正常解析，**并非**从 `CacheServiceProvider` 移除。

- Laravel 7.x~13.x 的 `Illuminate\Cache\CacheServiceProvider` 历来只绑定字符串键 `cache` / `cache.store` / `cache.psr6` / `memcached.connector` / `RateLimiter`，**从未注册过**缓存接口类名
- 缓存接口类名能通过 `app()` 解析，靠的是 `Application::registerCoreContainerAliases()`（`Illuminate\Foundation\Application`）注册的**容器类名别名**：

```php
'cache'       => [CacheManager::class, Contracts\Cache\Factory::class],
'cache.store' => [Repository::class, Contracts\Cache\Repository::class, Psr\SimpleCache\CacheInterface::class],
```

- 以下三种写法等价，均返回默认 store 的 `Illuminate\Cache\Repository` 单例：
  - `$app->make('cache')->store()`
  - `app('cache.store')`
  - `app(\Illuminate\Contracts\Cache\Repository::class)`

### 2. 缓存配置与 store 选择

Snowflake 序列与 IP 限流共用 `CACHE_STORE`（`.env` 变量 → `config/cache.php` 的 `default`）。

**切换 Redis 的方式**：

```bash
# .env
CACHE_STORE=redis
REDIS_CLIENT=phpredis          # 或 predis（需 composer require predis/predis）
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=null
```

启用后注意：`config/cache.php` 的 redis store 默认使用 `REDIS_CACHE_CONNECTION`（即 `config/database.php` 中名为 `cache` 的 redis 连接），与 `REDIS_*` 主连接解耦，可按需独立配置。

**各 store 特性**：

| store | 原子 `add/increment` | 说明 |
|---|---|---|
| redis / memcached / database | 支持 | 生产推荐（尤其 redis） |
| file | 不支持 | 单进程 FPM 开发可用；但 `snowflake:<毫秒>` 这类时间戳 key 过期是惰性判断、**物理文件永不删除**，长期运行会在 `storage/framework/cache/data` 累积大量小文件，仅 `cache:clear` 清理 |
| array | 支持（仅进程内） | 仅测试用，重启即丢 |

**缓存 key 一览**：

| key | 用途 | 生命周期 |
|---|---|---|
| `snowflake:<毫秒时间戳>` | Snowflake 同毫秒序列计数 | TTL 10 秒 |
| `rate_limit:<IP>` | IP 滑动窗口时间戳数组 | TTL = 窗口时长（默认 60s） |
| `health_started_at` / `health_db` | 健康检查缓存 | 长期 / 5 秒 |

### 3. Snowflake 单例绑定说明

绑定位于 [app/Providers/AppServiceProvider.php](app/Providers/AppServiceProvider.php)（配置见 `config/snowflake.php`）：

```php
$this->app->singleton(Snowflake::class, function (Application $app) {
    return (new Snowflake(
        (int) config('snowflake.datacenter'),
        (int) config('snowflake.worker'),
    ))->setSequenceResolver(
        (new LaravelSequenceResolver($app->make('cache')->store()))
            ->setCachePrefix('snowflake:')
    );
});
```

**为什么必须单例**：Snowflake ID 由 `时间戳(41bit) | datacenter(5bit) | worker(5bit) | sequence(12bit)` 组成。`sequence` 的状态（同一毫秒内递增）保存在解析器实例内，**必须跨请求保留**；若每次 `new` 都会重建状态，同一毫秒内的多次调用可能生成相同 ID。

**LaravelSequenceResolver 如何保证同毫秒唯一**：

1. 同毫秒内第一个请求：`cache->add("snowflake:<ms>", 1, 10)` 成功 → 序列返回 0
2. 后续请求：`add` 失败 → `cache->increment("snowflake:<ms>")` 返回 1、2、3…
3. 由于自增由缓存（redis 等）原子完成，**跨进程/跨实例**也能保证同毫秒唯一

**datacenter / worker 的作用**：把 ID 空间划分为 `31×31` 个互不重叠的分区（各 5bit）。单实例默认 `-1`（构造时随机分配即可）；多实例部署时建议为各实例分配不同的组合——这是"双保险"（即使不区分，共享原子缓存下的序列也能保证唯一），同时可降低对缓存原子性的依赖。

**验证绑定是否生效**：

```bash
php artisan tinker --execute='$s = app(\Godruoyi\Snowflake\Snowflake::class); echo $s->id();'
```

连续执行两次应得到递增的 ID；若改用 redis，可执行 `redis-cli keys "*snowflake:*"` 观察序列计数 key（若配置了 `CACHE_PREFIX`，key 会带该前缀）。

## 与 Go 原版的差异

| 项 | Go 原版 | PHP 版 | 说明 |
|---|---|---|---|
| Web 框架 | Fiber v3 | Laravel 13 | 路由集中在 routes/api.php，控制器方法对应各端点 |
| ORM | GORM | Eloquent | `User::query()->select(User::noPassword())` 排除 password |
| 密码哈希 | gonet/crypto（bcrypt） | PHP 原生 `password_hash`/`password_verify`（PASSWORD_BCRYPT） | 含时序侧信道防护（登录走假哈希） |
| JWT | golang-jwt/v5 | lcobucci/jwt（[JwtUtil](app/Utils/JwtUtil.php) 单例封装） | HS256 对齐 |
| 限流 | 进程内滑动窗口（16 分片锁） | 缓存驱动滑动窗口（file/redis，key `rate_limit:<IP>`） | 多实例可用 Redis 共享计数 |
| 端口 | 6600 | 8000 | `php artisan serve` 默认端口 |
| OpenAPI | swag 生成静态文件 | swagger-php 注解生成（`composer swagger`） | 端点一致 |
| 500 文案 | `internal server error` | `internal server error`（未捕获异常统一固定文案） | 已对齐 |
| 教学取舍 | - | 手动校验 / `new` 实例化 Service | 见「教学取舍说明」 |

## 如何新增一个接口

以新增 `GET /api/v1/user/detail`（按 uid 查用户公开信息）为例，走一遍完整链路：

1. **添加控制器方法**：在 [app/Http/Controllers/UserController.php](app/Http/Controllers/UserController.php) 添加，模式固定为「取参/校验 → `new UserService()` → `response_success(...)`」：

   ```php
   #[OA\Get(
       path: '/api/v1/user/detail',
       summary: '获取用户详情',
       tags: ['User'],
       parameters: [
           new OA\QueryParameter(name: 'uid', description: '用户id', required: true, schema: new OA\Schema(type: 'integer')),
       ],
       responses: [
           new OA\Response(response: 200, description: 'OK', content: new OA\JsonContent(ref: '#/components/schemas/User')),
           new OA\Response(response: 400, description: '请求格式错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
           new OA\Response(response: 422, description: '参数错误', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
           new OA\Response(response: 500, description: 'Internal Server Error', content: new OA\JsonContent(ref: '#/components/schemas/Error')),
       ],
   )]
   public function detail(Request $request): JsonResponse
   {
       $uidRaw = $request->input('uid', 0);
       if ($uidRaw !== 0 && filter_var($uidRaw, FILTER_VALIDATE_INT) === false) {
           return response_error(ResultCode::BadRequest, '无效的用户ID');
       }
       $uid = (int)$uidRaw;
       if ($uid <= 0) {
           return response_error(ResultCode::Param, '无效的用户ID');
       }
       $user = (new UserService())->info($uid);
       if (empty($user)) {
           return response_error(ResultCode::Param, '用户不存在');
       }
       return response_success(format_user($user));
   }
   ```

   需要 JWT 鉴权时给路由挂 `AuthJwt` 中间件（见下方第 3 步），控制器里从 `$request->attributes->get('uid')` 取当前用户（见 `update`）。

2. **添加业务逻辑**：若逻辑与已有方法不同，在 [app/Http/Services/UserService.php](app/Http/Services/UserService.php) 添加方法，复用 `User::query()->select(User::noPassword())`（排除 password）；已知业务错误抛 `RuntimeException($msg, ResultCode::Param->value)`（见 `update`）。

3. **注册路由**：在 [routes/api.php](routes/api.php) 追加一行：

   ```php
   Route::get('user/detail', [UserController::class, 'detail']);
   ```

   需要 JWT 时改为 `Route::get('user/detail', [UserController::class, 'detail'])->middleware(AuthJwt::class)`（见 `auth/refresh`）。

4. **文档与测试**：方法上写好 `#[OA\Get(...)]` 注解，改动后 `composer swagger` 重新生成文档；可在 [tests/Feature](tests/Feature) 补一条 PHPUnit 集成测试（sqlite 内存库自动跑 migrations，见 AuthTest）。

## 安全说明

- **JWT 签名校验**：所有受保护接口的 token 均校验 HMAC 签名（`SignedWith`）与有效期（`StrictValidAt`），伪造/篡改/过期 token 一律返回 401
- **密码**：bcrypt 哈希存储；模型 `$hidden` 防止密码被序列化泄露
- **限流**：全局 IP 滑动窗口限流，超阈值返回 429；`/api/v1/health` 豁免，避免探活被误伤
- **统一 500**：未捕获异常统一返回固定文案 `internal server error`，不泄露内部细节

## 常见问题

- **`composer install` 很慢**：国内可配置 Composer 镜像源（如阿里云 packagist 镜像）。
- **启动报错提示缺少 APP_KEY**：先执行 `php artisan key:generate`（见「安装与运行」）。
- **端口 8000 被占用**：改用 `php artisan serve --port=8080` 启动。
- **执行 `config:cache` 后担心鉴权失效**：JWT 密钥必须经 `config('jwt.*')` 读取，项目已按此实现不会失效（见「安装与运行」注意事项）。
