// Health 健康检查路由 (与 Go/PHP/Rust 版对齐)
// GET /api/v1/health: 探测服务与数据库是否正常
// 数据库不可用时返回 503 + 统一错误格式, 便于负载均衡探活

import { Elysia } from "elysia";
import { prisma } from "../lib/prisma";
import { formatDateTime } from "../util/common";

// 服务版本 (与 Go/PHP/Rust 版 1.0.0 对齐)
const VERSION = "1.0.0";

// 服务启动时间 (进程级常量, 服务器本地时区)
const startedAt = formatDateTime(new Date());

// 数据库探活结果缓存 (5 秒 TTL, 避免探活打爆数据库, 与 Go/PHP/Rust 版一致)
// inFlight: 缓存过期瞬间的并发穿透合并 (single-flight), 避免 N 个并发探活同时执行 SELECT 1
const DB_PING_TTL = 5_000;
let dbCache: { at: number; ok: boolean } | null = null;
let inFlight: Promise<boolean> | null = null;

async function dbHealthy(): Promise<boolean> {
  const now = Date.now();
  if (dbCache && now - dbCache.at < DB_PING_TTL) {
    return dbCache.ok;
  }
  // 复用进行中的探测请求, 只让第一个请求真正执行 (缓存失效瞬间的并发穿透只打一次 DB)
  if (inFlight) {
    return inFlight;
  }
  inFlight = (async () => {
    let ok = false;
    try {
      await prisma.$queryRaw`SELECT 1`;
      ok = true;
    } catch {
      ok = false;
    }
    dbCache = { at: Date.now(), ok };
    return ok;
  })();
  try {
    return await inFlight;
  } finally {
    inFlight = null;
  }
}

// 注意: 用 prefix "/api/v1" + 路径 "/health", 而非 prefix "/api/v1/health" + "/",
// 否则文档路径会变成 "/api/v1/health/" (与其他语言端不一致)
export const healthRoutes = new Elysia({ prefix: "/api/v1" }).get("/health", async ({ set }) => {
  const resp = {
    status: "ok",       // 服务状态 (固定 "ok")
    database: "up",     // 数据库状态 (固定 "up")
    version: VERSION,   // 服务版本
    started_at: startedAt, // 服务启动时间 (服务器本地时区)
    now: formatDateTime(new Date()), // 当前时间 (服务器本地时区)
  };
  if (!(await dbHealthy())) {
    set.status = 503;
    return { code: 503, message: "database unavailable" };
  }
  return resp;
}, {
  detail: {
    tags: ["Health"],
    summary: "健康检查",
    responses: {
      "200": {
        description: "OK",
        content: {
          "application/json": {
            schema: { $ref: "#/components/schemas/HealthResponse" },
          },
        },
      },
      "503": {
        description: "Service Unavailable",
        content: {
          "application/json": { schema: { $ref: "#/components/schemas/Error" } },
        },
      },
    },
  },
});

export default healthRoutes;
