// Prisma 数据库客户端
// 使用 MariaDB 适配器连接 MySQL，全局单例模式防止热重载时重复创建

import { PrismaClient } from "../generated/prisma/client";
import { PrismaMariaDb } from "@prisma/adapter-mariadb";

// 将 PrismaClient 挂载到 globalThis，Bun --watch 热重载时复用已有实例
const globalForPrisma = globalThis as unknown as {
  prisma: PrismaClient | undefined;
};

const connectionString = process.env.DATABASE_URL!;

// 解析连接字符串构建 PoolConfig（显式配置 pool size，默认 10）
const url = new URL(connectionString);
const poolConfig = {
  host: url.hostname,
  port: parseInt(url.port || "3306", 10),
  user: decodeURIComponent(url.username),
  password: decodeURIComponent(url.password),
  database: url.pathname.replace(/^\//, ""),
  connectionLimit: parseInt(process.env.POOL_SIZE || "10", 10),
};

export const prisma =
  globalForPrisma.prisma ??
  new PrismaClient({
    adapter: new PrismaMariaDb(poolConfig), // MariaDB 驱动适配器 + 连接池配置
  });

// 非生产环境缓存实例，避免 hmr 时创建过多连接
if (process.env.NODE_ENV !== "production") {
  globalForPrisma.prisma = prisma;
}
