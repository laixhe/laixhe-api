# WebApi (Java 版)

由 GoWebApi 对齐转换的 Java 实现，接口、返回格式、校验规则、错误码与其余各端保持一致。

- 框架：Spring Boot **4.1.0**（Spring MVC + Spring Security + Spring Data JPA）
- 语言：Java **25**（GraalVM 25 工具链）
- JWT：[jjwt](https://github.com/jwtk/jjwt) 0.13.0（HS256）
- 限流：[Bucket4j](https://github.com/bucket4j/bucket4j) 8.19.0（令牌桶，等效滑动窗口）
- OpenAPI：[springdoc-openapi](https://springdoc.org/) 3.1.0（代码注解动态生成）
- 数据库：MySQL（生产）/ H2（本地与测试免安装，MySQL 兼容模式）

## 环境要求

| 依赖 | 版本 |
|---|---|
| JDK | 25（GraalVM 25；构建工具链按 [build.gradle.kts](build.gradle.kts) 的 `JavaLanguageVersion.of(25)` 解析） |
| Gradle | 9.7（项目内置 wrapper，无需手动安装） |
| MySQL | 5.7+（不装也行：`h2` profile 使用内置内存库） |

## 快速开始

```bash
# 1. 修改配置（application.yaml 的 spring.datasource 指向你的 MySQL）

# 2. 启动（默认监听 0.0.0.0:6600）
./gradlew bootRun

# 免装 MySQL：切 h2 profile，使用 H2 内存库（MySQL 兼容模式，启动自动建表）
./gradlew bootRun --args='--spring.profiles.active=h2'
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

## 关键配置（application.yaml）

| 配置 | 说明 | 默认 |
|---|---|---|
| `server.port` | 监听端口 | 6600 |
| `spring.datasource` | MySQL 连接（`jdbc:mysql://127.0.0.1:3306/webapi`，账号 root / 123456） | - |
| `spring.jpa.hibernate.ddl-auto` | 演示项目自动建表（`account`/`email` 唯一索引由实体生成） | update |
| `app.http-timeout` | 请求超时（秒），超时返回 408 | 30 |
| `app.jwt.secret-key` | JWT 密钥（**生产必改**，支持环境变量 `JWT_SECRET_KEY` 覆盖） | 与 Go 端一致 |
| `app.jwt.expire-seconds` | token 过期时长（秒） | 2592000（30 天） |
| `app.limit.enable` / `max` / `window-seconds` | IP 限流开关 / 窗口内最大请求数 / 窗口时长（秒） | true / 1000 / 60 |
| `springdoc.swagger-ui.enabled` | 关闭框架自带 UI，使用自研 `/api/v1/swagger` 页面 | false |

配置经 `@ConfigurationProperties(prefix = "app")` 绑定到 [AppProperties.java](src/main/java/com/laixhe/webapi/config/AppProperties.java)。

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
| GET | `/api/v1/swagger.json` | 公开 | OpenAPI JSON 文档（springdoc 动态生成，原端点 `/v3/api-docs`） |
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
| 400 | 400 | 请求体 JSON 解析失败 / 缺少必填 query 参数 / query 参数类型错误 |
| 401 | 401 | 未授权（缺少或无效 JWT、用户被禁用） |
| 404 | 404 | 路由不存在（[NotFoundController](src/main/java/com/laixhe/webapi/controller/NotFoundController.java) 兜底） |
| 408 | 408 | 请求超时（超过 `app.http-timeout` 秒） |
| 422 | 422 | 参数错误（`@Valid` 校验失败 / `Validators` 手写校验 / 业务错误） |
| 429 | 429 | 触发 IP 限流 |
| 500 | 500 | 内部错误（统一固定文案 `internal server error`，不泄露细节） |
| 503 | 503 | 健康检查数据库不可用 |

## 校验规则（与 Go 端一致）

- 昵称：2~20 个字符（`Validators` 按 Unicode 码点统计，中文按 1 字符计）
- 邮箱：`@NotBlank` + `@Email`，错误文案「邮箱格式错误」
- 密码：`^[a-zA-Z0-9_@$]{6,64}$`（6~64 位；上限 64 防 bcrypt 72 字节静默截断）
- 头像地址：非空时最长 255 且必须以 `http://`/`https://` 开头
- 重复邮箱注册返回 422「邮箱已存在」

> DTO 上的 Bean Validation 注解（`@NotBlank`/`@Email`/`@Pattern`）由框架自动校验，失败统一 422 并取第一个字段错误信息；昵称/头像因需要按 Unicode 码点与自定义文案，由控制器手写调用 [Validators.java](src/main/java/com/laixhe/webapi/common/Validators.java)。

## 目录结构

```
src/
├── main/
│   ├── java/com/laixhe/webapi/
│   │   ├── WebApiApplication.java    # 启动入口
│   │   ├── common/                   # ApiException / Error / GlobalExceptionHandler / Validators
│   │   ├── config/                   # AppProperties / SecurityConfig / OpenApiConfig
│   │   ├── controller/               # Auth / User / Health / NotFound / SwaggerDoc
│   │   ├── dto/                      # 请求/响应 DTO（校验注解 + swagger 注解）
│   │   ├── entity/                   # JPA 实体（User / UserExtend / UserThirdParty / ...）
│   │   ├── middleware/               # RateLimitFilter / TimeoutFilter / 过滤器注册
│   │   ├── repository/               # Spring Data JPA 数据访问
│   │   ├── security/                 # JwtService / JwtAuthenticationFilter / ClaimsHolder
│   │   └── service/                  # AuthService / UserService（业务逻辑）
│   └── resources/                    # application.yaml / application-h2.yaml
└── test/java/com/laixhe/webapi/      # 集成测试（H2 内存库，事务回滚）
build.gradle.kts                      # 构建配置（含 GraalVM Native 插件）
```

调用链：`controller → service → repository → entity`，JWT 校验由 Spring Security 过滤器链完成。

## 测试

```bash
./gradlew test
```

- [ApiIntegrationTests](src/test/java/com/laixhe/webapi/ApiIntegrationTests.java)：注册/登录/刷新/更新/参数校验/404/Swagger 全链路，H2 内存库 + 每用例事务回滚，不依赖外部 MySQL
- [RateLimitTests](src/test/java/com/laixhe/webapi/RateLimitTests.java)：`app.limit.max=3` 独立上下文验证超限 429 与健康检查豁免

## GraalVM 原生镜像

```bash
./gradlew nativeCompile   # 产物: build/native/nativeCompile/webapi
```

需要本机安装 GraalVM 25 及 native-image；原生编译可显著降低启动时间与内存占用（插件配置见 [build.gradle.kts](build.gradle.kts)）。

## 技术要点

### 1. 无状态 JWT 鉴权（Spring Security）

- 受保护路径只有 `/api/v1/auth/refresh` 与 `/api/v1/user/update`，其余全部公开（[SecurityConfig.java](src/main/java/com/laixhe/webapi/config/SecurityConfig.java)）
- [JwtAuthenticationFilter](src/main/java/com/laixhe/webapi/security/JwtAuthenticationFilter.java) 解析 `Authorization: Bearer` 令牌写入安全上下文；令牌缺失/无效时不直接报错，由安全链对受保护路径统一返回 401，公开路径即使带无效令牌也能访问
- `uid` 从 1 开始，`uid <= 0` 视为无效（防御伪造 `{"uid":0}` 的令牌）；`ClaimsHolder.uid()` 从安全上下文取当前用户

### 2. Bucket4j 令牌桶限流

[RateLimitFilter](src/main/java/com/laixhe/webapi/middleware/RateLimitFilter.java) 为每个 IP 维护一个令牌桶：容量 = `limit.max`，按窗口时长匀速补充（`Refill.greedy(max, window)`），效果等价于滑动窗口（默认 1000 次/60s）。健康检查路径豁免；后台 janitor 线程周期性清理空闲超 2 个窗口的桶，防止伪造 IP 导致内存无限增长。过滤器注册顺序 `-200`，早于 Spring Security 过滤器链（`-100`），与 Go 版「限流 → JWT」的中间件顺序对齐（见 [RateLimitConfig.java](src/main/java/com/laixhe/webapi/middleware/RateLimitConfig.java)）。

### 3. 请求超时（408）

[TimeoutFilter](src/main/java/com/laixhe/webapi/middleware/TimeoutFilter.java) 将过滤链提交到独立线程执行，超过 `app.http-timeout` 秒未完成返回 408 统一 JSON。注意：超时后下游请求线程无法被强制中断，会继续运行到结束（教学规模下仅用于兜底异常慢请求）。

### 4. 全局异常处理

[GlobalExceptionHandler](src/main/java/com/laixhe/webapi/common/GlobalExceptionHandler.java) 统一返回 `{code, message}`：

- 业务异常（`ApiException`）状态码与文案原样返回
- `@Valid` 校验失败 → 422，取第一个字段错误信息
- JSON 解析失败 / 缺少必填参数 / 参数类型错误 → 400
- 未匹配路由 → 404
- 其余未知错误：记录服务端日志后统一返回固定 500 文案，不泄露内部细节

### 5. H2 profile 与自动建表

[application-h2.yaml](src/main/resources/application-h2.yaml) 使用 MySQL 兼容模式（`MODE=MySQL;NON_KEYWORDS=USER`）的内存库，`ddl-auto: update` 自动建表（`account`/`email` 唯一索引由实体 `@UniqueConstraint` 生成），本地开发与测试无需安装 MySQL。

### 6. 健康检查

[HealthController](src/main/java/com/laixhe/webapi/controller/HealthController.java) 执行 `SELECT 1` 探测数据库，结果缓存 5 秒避免频繁探活压垮数据库（与 Go 端 `healthPingInterval` 对齐），数据库不可用时返回 503 便于负载均衡探活。

## 与 Go 原版的差异

| 项 | Go 原版 | Java 版 | 说明 |
|---|---|---|---|
| Web 框架 | Fiber v3 | Spring MVC | 注解式路由（无独立路由文件） |
| ORM | GORM | Spring Data JPA | 数据访问独立在 `repository` 层 |
| 密码哈希 | gonet/crypto（bcrypt） | Spring Security `BCryptPasswordEncoder` | cost=10 对齐 |
| JWT | golang-jwt/v5 | jjwt | HS256 对齐 |
| 限流 | 滑动窗口 | Bucket4j 令牌桶 | 等效滑动窗口 |
| OpenAPI | swag 生成静态文件 | springdoc 动态生成 | `/api/v1/swagger*` 由 [SwaggerDocController](src/main/java/com/laixhe/webapi/controller/SwaggerDocController.java) 转发到 `/v3/api-docs` |
| 建表 | 手动导入 webapi.sql | `ddl-auto: update` 自动建表 | 也可按根目录 webapi.sql 建表 |
| 额外能力 | - | H2 profile、GraalVM 原生镜像、MockMvc 集成测试 | 增强 |

## 如何新增一个接口

以新增 `GET /api/v1/user/detail`（按 uid 查用户公开信息）为例，走一遍完整链路：

1. **添加控制器方法**：在 [UserController.java](src/main/java/com/laixhe/webapi/controller/UserController.java) 添加，模式固定为「参数绑定 → 调 service → 返回 DTO」：

   ```java
   @Operation(summary = "获取用户详情")
   @ApiResponses({
           @ApiResponse(responseCode = "200", description = "OK", content = @Content(schema = @Schema(implementation = UserResponse.class))),
           @ApiResponse(responseCode = "400", description = "Bad Request", content = @Content(schema = @Schema(implementation = Error.class))),
           @ApiResponse(responseCode = "422", description = "Unprocessable Entity", content = @Content(schema = @Schema(implementation = Error.class))),
           @ApiResponse(responseCode = "500", description = "Internal Server Error", content = @Content(schema = @Schema(implementation = Error.class))),
   })
   @GetMapping("/detail")
   public UserResponse detail(@Parameter(description = "用户id", required = true)
                              @RequestParam("uid") int uid) {
       return userService.detail(uid);
   }
   ```

   需要 JWT 鉴权时把路径加进 [SecurityConfig.java](src/main/java/com/laixhe/webapi/config/SecurityConfig.java) 的 `authenticated()` 列表，控制器用 `ClaimsHolder.uid()` 取当前用户（见 `update`）。

2. **添加业务逻辑**：在 [UserService.java](src/main/java/com/laixhe/webapi/service/UserService.java) 添加方法，复用 `userRepository.findById`；已知业务错误抛 `ApiException.paramError(...)`：

   ```java
   public UserResponse detail(int uid) {
       if (uid <= 0) {
           throw ApiException.paramError("无效的用户ID");
       }
       User user = userRepository.findById(uid)
               .orElseThrow(() -> ApiException.paramError("用户不存在"));
       return UserResponse.from(user);
   }
   ```

   复杂查询在 [repository](src/main/java/com/laixhe/webapi/repository) 层定义方法（Spring Data JPA 方法名推导或 `@Query`），在 service 注入调用。

3. **注册路由**：无需额外步骤——`@RestController` + `@GetMapping` 注解即路由，启动即生效（区别于 Go/Rust/PHP 的独立路由文件）。

4. **文档与测试**：`@Operation`/`@ApiResponses` 注解由 springdoc 自动进文档（`/api/v1/swagger` 页面可见，无需重新生成）；可在 [ApiIntegrationTests.java](src/test/java/com/laixhe/webapi/ApiIntegrationTests.java) 补一条 MockMvc 用例（H2 内存库，事务回滚）。

## 安全说明

- **生产环境务必更换 JWT 密钥**：`application.yaml` 的 `app.jwt.secret-key`，或用环境变量 `JWT_SECRET_KEY` 覆盖，默认值任何人都能用来伪造令牌
- `/api/v1/user/info`、`/api/v1/user/list` 是**公开接口**（未挂 JWT），用于教学演示；上生产前请评估是否需要鉴权或字段脱敏
- IP 限流优先信任 `X-Forwarded-For`，直接面向公网时客户端可伪造该头绕过限流，建议仅在可信反向代理后部署

## 常见问题

- **首次 `./gradlew bootRun` 下载依赖很慢**：国内可在 [build.gradle.kts](build.gradle.kts) 的 `repositories` 增加阿里云 Maven 镜像加速。
- **提示找不到 JDK 25 / GraalVM 工具链**：Gradle toolchain 会自动解析或下载对应 JDK；仍失败请安装 GraalVM 25 并检查 `JAVA_HOME`。
- **端口 6600 被占用**：修改 `application.yaml` 的 `server.port`。
- **不想装 MySQL**：用 `./gradlew bootRun --args='--spring.profiles.active=h2'` 启动（见「快速开始」）。
