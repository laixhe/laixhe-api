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

   默认监听 `http://0.0.0.0:6600`，可用 `curl http://localhost:6600/` 验证服务已启动。

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

## 参数校验分工

- **结构校验**：路由定义处的 Elysia `t.Object` schema，负责字段是否存在、类型是否正确，不合法自动返回 400；
- **业务规则校验**：`src/util/validate.ts`（如昵称长度、密码字符集），负责给出具体的中文错误提示。

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
