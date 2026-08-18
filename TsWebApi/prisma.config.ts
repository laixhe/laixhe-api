// Prisma 7 迁移配置
// Prisma 7 不在 schema.prisma 中定义 datasource url，改为在此文件统一配置
// 注意：Prisma 7 CLI 不再自动加载 .env（c12 dotenv 已关闭），需在此显式加载
// 仅用于 `prisma migrate`/`prisma db push` 等 CLI 命令，运行时由 PrismaClient adapter 接管

import { existsSync } from "node:fs";
import { loadEnvFile } from "node:process";
import { defineConfig, env } from "prisma/config";

// 加载 .env（文件不存在时跳过，便于 CI 等场景直接注入环境变量）
if (existsSync(".env")) {
  loadEnvFile(".env");
}

export default defineConfig({
  // 注意：Prisma 7 配置为 datasource.url（单数），不是旧版的 datasources.db.url
  datasource: {
    url: env("DATABASE_URL"),
  },
});
