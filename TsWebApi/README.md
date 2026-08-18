# laixhe-api (TsWebApi)

基于 [Bun](https://bun.sh) + [Elysia](https://elysiajs.com) + [Prisma](https://www.prisma.io)（MariaDB/MySQL）+ JWT(jose) 的 TypeScript Web API 服务，提供注册、登录、Token 刷新、用户信息查询/列表/更新等接口。

## 技术栈

- 运行时：Bun
- Web 框架：Elysia
- ORM：Prisma（数据库驱动为 MariaDB，兼容 MySQL）
- 鉴权：JWT（jose）签发/校验，bcrypt 密码加密
- 测试：bun test

## 环境要求

- Bun >= 1.3
- MariaDB / MySQL 5.7+

## 快速开始

1. 安装依赖

   ```bash
   bun install
   ```

2. 配置环境变量

   ```bash
   cp .env.example .env
   ```

   编辑 `.env`，至少需要填写 `DATABASE_URL`（格式见 `.env.example`）。

3. 初始化数据库（首次）

   ```bash
   bun prisma migrate dev --name init
   ```

   或跳过迁移直接同步表结构：

   ```bash
   bun run db:push
   ```

4. 启动开发服务器

   ```bash
   bun run dev
   ```

   默认监听 `http://0.0.0.0:6600`，可用 `curl http://localhost:6600/` 验证服务已启动。完整验证：

   ```bash
   # 健康检查
   curl http://127.0.0.1:6600/api/v1/health

   # Swagger UI 文档页（浏览器打开）
   # http://127.0.0.1:6600/api/v1/swagger

   # 注册
   curl -X POST http://127.0.0.1:6600/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"nickname":"test","email":"test@example.com","password":"abc123"}'
   ```

## 常用命令

| 命令 | 说明 |
| --- | --- |
| `bun run dev` | 开发模式（热重载） |
| `bun run build` | 构建到 dist/ |
| `bun run start` | 运行构建产物 |
| `bun run typecheck` | TypeScript 类型检查 |
| `bun test` | 运行接口冒烟测试（需可用的数据库） |
| `bun run db:generate` | 重新生成 Prisma Client（schema 变更后执行） |
| `bun run db:push` | 直接同步数据库表结构 |
| `bun run db:studio` | 打开 Prisma Studio 可视化查看数据 |

## 接口一览

所有接口前缀 `/api/v1`，鉴权接口需在请求头携带 `Authorization: Bearer <token>`。

| 方法 | 路径 | 鉴权 | 说明 |
| --- | --- | --- | --- |
| POST | /api/v1/auth/register | 否 | 注册，返回 token + 用户信息 |
| POST | /api/v1/auth/login | 否 | 登录（IP 维度限流，默认 5 次/分钟） |
| POST | /api/v1/auth/refresh | 是 | 刷新 token |
| GET | /api/v1/user/info?uid=1 | 否 | 查询用户信息（限流 60 次/分钟） |
| GET | /api/v1/user/list?page=1&page_size=12 | 否 | 用户列表（分页，按 id 倒序，限流 60 次/分钟） |
| POST | /api/v1/user/update | 是 | 更新昵称/头像 |
| GET | /api/v1/swagger | 否 | Swagger UI 可视化页面（浏览器打开） |
| GET | /api/v1/swagger.json | 否 | OpenAPI JSON 文档 |
| GET | /api/v1/swagger.yaml | 否 | OpenAPI YAML 文档 |

## 统一响应格式与错误码

- 成功：直接返回业务数据 JSON（`{token, user}` / `User` / `UserListResponse`）
- 失败：`{"code": <HTTP状态码>, "message": "..."}`，HTTP 状态码同步（由 [response.ts](src/util/response.ts) 的 `fail()` 与 [index.ts](src/index.ts) 的全局 `onError` 统一生成）

| code | HTTP | 含义 |
| --- | --- | --- |
| 400 | 400 | 请求格式错误（顶层 body 非对象 / 字段类型非字符串 / JSON 解析失败 / 非数字 query 参数） |
| 401 | 401 | 未授权（缺少或无效 JWT，由 `requireAuth` 插件返回） |
| 404 | 404 | 路由不存在 |
| 422 | 422 | 参数错误（缺 body / 业务校验失败，返回具体中文文案） |
| 429 | 429 | 触发限流（文案「请求过于频繁，请稍后再试」） |
| 500 | 500 | 内部错误（固定文案「服务器内部错误」，不泄露细节） |

> 注意：TS 版 500 文案为「服务器内部错误」，与其余各端的英文 `internal server error` 略有不同（见 [index.ts](src/index.ts)）。

## 校验规则（与 Go 端一致）

- 昵称：2~20 字符（`[...nickname]` 按 Unicode 码点计数，emoji 按 1 字符计）
- 邮箱：标准邮箱正则
- 密码：6~64 位，仅含字母 数字 _ @ $（`^[a-zA-Z0-9_@$]{6,64}$`；上限 64 防 bcrypt 72 字节静默截断）
- 头像地址：非空时最长 255 且必须以 `http://`/`https://` 开头
- 重复邮箱注册返回 422「邮箱已存在」

具体实现见 [validate.ts](src/util/validate.ts)，由各 handler 手写调用（与 Go/PHP/Rust 端一致）。

## API 文档

文档由 [@elysia/openapi](https://github.com/elysiajs/elysia-openapi) 插件（官方推荐，替代已弃用的 `@elysiajs/swagger`）在运行时根据路由 detail 注解动态生成，无需手动重新生成；公共响应 Schema 维护在 [src/docs/schemas.ts](src/docs/schemas.ts)。`/api/v1/swagger.yaml` 由 [src/index.ts](src/index.ts) 请求 JSON spec 后动态转换为 YAML，保持端点与其余各端一致。

## 参数校验分工

项目采用**全手动参数校验**（与 Go/PHP/Rust 端一致），未使用 Elysia 的 `t.Object` 结构 schema：

- **业务规则校验**：`src/util/validate.ts`（昵称长度、邮箱格式、密码字符集），在 handler 内手写调用，失败返回 422 + 具体中文提示；
- **缺字段/空 body**：请求体缺失字段导致解构 `undefined` 时，由全局 `onError`（`src/index.ts`）兜底返回 422「参数错误」，不会落到 500。

### bodyError 为什么这么写（400 / 422 / 500 的判定流程）

Go/Rust 端的框架绑定层（`ctx.Bind().JSON()` / `JsonBody<T>`）对"顶层 body 是数组/标量/null"和"字段类型错误"会自动返回 400；TS 端没有这一层，若不加处理，这类输入会一路漏到 Prisma/业务层变成 500。`bodyError`（[validate.ts](src/util/validate.ts)）就是手动补上这层"绑定层类型校验"：

```
请求体 body 到达 handler
        │
        ├─ body === undefined（缺 body）──────────► 返回 null，交给 handler
        │                                            │ 解构 {email,password} 抛 TypeError
        │                                            ▼
        │                                     全局 onError 兜底 ──► 422「参数错误」
        │
        ├─ body 非纯对象（数组/标量/null）──────► 返回 "top-level" ──► handler 返回 400
        │
        ├─ 字段类型非 string（数字/布尔/数组）──► 返回字段名 ────────► handler 返回 400
        │
        └─ 全部通过 → 返回 null，handler 继续走业务校验 ────────────► 422 + 具体文案
```

设计要点：

1. **null 视为"无值"而非类型错误**（与 Go/PHP 端一致），先 `normalizeNulls` 归一化为空串，再走业务校验返回具体的 422 文案（如「昵称长度不能小于2位」），而不是 400。
2. **缺 body 与缺字段走不同路径**：缺 body 依赖全局 `onError` 的 TypeError 兜底返回 422（区分不出具体字段）；缺字段由各校验函数对 `undefined` 做空串兜底（`?? ""`），返回**具体的** 422 文案，更友好。
3. **fields 白名单必须与 handler 解构的字段一致**：新增字段忘记加入白名单，该字段会绕过类型校验，错误类型输入可能漏成 500——这是该模式唯一需要手动维护的地方。

## 目录结构

```
src/
  index.ts        应用入口（路由注册、全局错误处理、优雅关闭）
  config.ts       环境变量配置
  lib/prisma.ts   Prisma 客户端单例（连接池）
  routes/         路由（auth / user）
  middleware/     鉴权插件(requireAuth)、JWT 工具、速率限制
  entity/         响应类型 / 枚举（auth、user）
  util/           日志、校验、统一响应、通用工具
  generated/      Prisma 生成的代码（勿手改）
test/             接口冒烟测试
prisma/schema.prisma  数据库模型定义
```

## 注意事项

- 登录限流基于内存 Map，仅适用于单实例部署；多实例请改用 Redis 等共享存储（见 `src/middleware/rateLimit.ts`）。
- 生产环境务必修改 `.env` 中的 `JWT_SECRET_KEY` 为强随机字符串。
- 修改 `prisma/schema.prisma` 后需重新生成客户端（`bun run db:generate`）并执行迁移（`bun prisma migrate dev`）。

## 与 Go 原版的差异

| 项 | Go 原版 | TS 版 | 说明 |
| --- | --- | --- | --- |
| Web 框架 | Fiber v3 | Elysia | 业务逻辑写在路由 handler 内（薄框架），无独立 controller/service 文件 |
| ORM | GORM | Prisma | 表结构定义见 prisma/schema.prisma |
| 密码哈希 | gonet/crypto（bcrypt） | `Bun.password`（bcrypt cost=10） | 原生异步、独立线程计算，不阻塞 JS 主线程 |
| JWT | golang-jwt/v5 | jose | 载荷均含 `uid` + `iat`/`exp`；TS 版不设置 `nbf`（校验端不要求） |
| 限流 | 全局滑动窗口（16 分片锁） | 内存 Map 滑动窗口 + 定时清理 | 另加登录/注册 5 次/分钟、info/list 60 次/分钟单端点限流 |
| 500 文案 | `internal server error` | `服务器内部错误` | 未完全对齐，见「统一响应格式与错误码」 |
| OpenAPI | swag 生成静态文件 | @elysia/openapi 运行时生成 | 端点一致（swagger / swagger.json / swagger.yaml） |
| count 优化 | 无缓存 | total 5s 缓存 + single-flight | 避免高频翻页反复全表 count |

## 如何新增一个接口

以新增 `GET /api/v1/user/detail`（按 uid 查用户公开信息）为例，走一遍完整链路：

1. **添加路由 handler**：在 [src/routes/user.ts](src/routes/user.ts) 的公开分组追加 `.get("/detail", ...)`，模式固定为「类型校验 → 业务校验 → Prisma 查询 → `toUserInfo` / `fail`」：

   ```ts
   // GET /api/v1/user/detail?uid=xxx（公开接口）
   .get("/detail", async ({ query, set }) => {
     let uid: number;
     if (query.uid === undefined || query.uid === "") {
       uid = 0;
     } else {
       uid = Number(query.uid);
       if (!Number.isSafeInteger(uid)) {
         return fail(set, 400, "Bad Request");
       }
     }
     if (uid <= 0) {
       return fail(set, 422, "无效的用户ID");
     }
     const user = await prisma.user.findUnique({
       where: { id: uid },
       select: userPublicSelect,
     });
     if (!user) {
       return fail(set, 422, "用户不存在");
     }
     return toUserInfo(user);
   }, {
     // 公开接口防刷：宽松限流（60 次/分钟）
     beforeHandle: rateLimit(60, 60_000, "user:detail"),
     detail: {
       tags: ["User"],
       summary: "获取用户详情",
       parameters: [
         { name: "uid", in: "query", required: true, description: "用户id", schema: { type: "integer" } },
       ],
       responses: {
         "200": { description: "OK", content: { "application/json": { schema: { $ref: "#/components/schemas/User" } } } },
         "400": { description: "请求格式错误", content: { "application/json": { schema: { $ref: "#/components/schemas/Error" } } } },
         "422": { description: "参数错误", content: { "application/json": { schema: { $ref: "#/components/schemas/Error" } } } },
         "500": { description: "Internal Server Error", content: { "application/json": { schema: { $ref: "#/components/schemas/Error" } } } },
       },
     },
   })
   ```

   需要 JWT 鉴权时把 handler 移到 `.use(requireAuth)` 之后注册（handler 直接拿到 `user`，见 `update`）。

2. **响应复用既有 Schema**：`User` / `Error` 已定义在 [src/docs/schemas.ts](src/docs/schemas.ts)，无需新增；只有引入新实体时才需要在那里补 schema。

3. **注册路由**：无需额外步骤——`userRoutes` 已在 [src/index.ts](src/index.ts) 中 `.use(userRoutes)`，新端点自动生效。

4. **文档与测试**：`detail` 注解由 @elysia/openapi 自动进文档（`/api/v1/swagger` 页面可见）；可在 [test/e2e.test.ts](test/e2e.test.ts) 补一条冒烟测试。

## 排错记录

### 1. Prisma CLI 报 `The datasource.url property is required`（2026-08-08）

**现象**：执行 `bun run db:push` / `bun prisma migrate dev` 报错 `The datasource.url property is required in your Prisma config file`，即使 `.env` 中已配置 `DATABASE_URL`。

**根因**（Prisma 7 的两处行为变化）：
1. Prisma 7 配置项为 `datasource.url`（**单数**），旧版的 `datasources.db.url` 写法会被 CLI 静默忽略；
2. Prisma 7 CLI 不再自动加载 `.env`（底层 c12 的 `dotenv` 已关闭），配置里 `process.env.DATABASE_URL` 始终为 `undefined`。

**修复**（`prisma.config.ts`）：

```ts
import { existsSync } from "node:fs";
import { loadEnvFile } from "node:process";
import { defineConfig, env } from "prisma/config";

// 显式加载 .env（文件不存在时跳过，便于 CI 直接注入环境变量）
if (existsSync(".env")) {
  loadEnvFile(".env");
}

export default defineConfig({
  datasource: {
    url: env("DATABASE_URL"),
  },
});
```

### 2. 给 `user` 表加唯一约束时的历史数据清理（2026-08-08）

**背景**：为 `account`、`email` 增加 `@unique` 约束（防止并发注册产生重复数据）时，`db push` 分阶段报错。

**阶段一：`P2002 Unique constraint failed on user_account_key`**

原因：历史测试数据中有 4 条记录的 `account` 为空字符串，重复导致唯一约束失败。邮箱无重复。

处理：将空账号按应用注册逻辑回填为随机 UUID：

```ts
await prisma.user.update({
  where: { id: u.id },
  data: { account: crypto.randomUUID() },
});
```

**阶段二：外键约束失败 `user_extend_uid_fkey`**

原因：`user_extend` / `user_third_party` 中存在 18 条孤儿行（`uid` 对应的 `user` 记录已删除），重建表添加外键时校验失败。

处理：删除孤儿行：

```ts
await prisma.userExtend.delete({ where: { id: o.id } });
await prisma.userThirdParty.delete({ where: { id: o.id } });
```

**收尾**：以上数据修复完成后执行（注意历史库存在 unsigned→signed 类型漂移，需带数据损失确认标志）：

```bash
bun x prisma db push --accept-data-loss
```

> 提示：若你的库中 `account` 为空或存在重复数据，应先查询并修复再执行上述命令；新写入的数据（注册接口）均使用随机 UUID 作为 `account`，不会重复。
>
> 说明：本文档为排错历史记录。当前 [schema.prisma](./prisma/schema.prisma) 中 `account`、`email` 均保持 `@unique` 唯一约束，与本文结论一致；改动 schema 后需执行 `bun x prisma db push` 同步数据库。

## 常见问题

- **端口 6600 被占用**：修改 `src/config.ts` 的 `http.port`。
- **数据库连接失败**：检查 `.env` 的 `DATABASE_URL`，并确认 MySQL 已启动。
- **Prisma 相关报错（`datasource.url`、`db push` 失败等）**：见上文「排错记录」。
- **修改 `schema.prisma` 后接口字段报错**：先执行 `bun run db:generate` 重新生成 Prisma Client。
