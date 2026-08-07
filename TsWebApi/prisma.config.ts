// Prisma 7 迁移配置
// Prisma 7 不在 schema.prisma 中定义 datasource url，改为在此文件统一配置
// 仅用于 `prisma migrate`/`prisma db push` 等 CLI 命令，运行时由 PrismaClient adapter 接管

import { defineConfig } from "prisma/config";

export default defineConfig({
  datasources: {
    db: {
      url: process.env.DATABASE_URL!,
    },
  },
});
