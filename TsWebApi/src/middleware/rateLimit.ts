// 简易内存速率限制（IP 维度），适用于登录等敏感接口防暴力破解
// 注意：内存实现仅适用于单实例部署，多实例需改用 Redis 等共享存储

import type { Context } from "elysia";
import { warn } from "../util/logger";

interface Entry {
  count: number;
  resetAt: number;
}

const store = new Map<string, Entry>();

// 最大条目数上限，防止攻击者用大量伪造 IP 撑爆内存
const MAX_ENTRIES = 10_000;

// 过期条目清理定时器（每 60s 清理一次）
// unref()：不持有事件循环，避免阻塞进程退出（测试/优雅关闭时友好）
setInterval(() => {
  const now = Date.now();
  for (const [key, entry] of store) {
    if (now > entry.resetAt) {
      store.delete(key);
    }
  }
}, 60_000).unref();

// 获取客户端真实 IP（按可信度从高到低）：
// 1. server.requestIP()：Bun 从 socket 层读取的真实地址，客户端无法伪造
// 2. x-real-ip：通常由可信反向代理设置
// 3. x-forwarded-for：可被客户端直接伪造，仅作为最后的兜底
function getClientIp(request: Request, server: Context["server"]): string {
  const direct = server?.requestIP(request);
  if (direct?.address) {
    return direct.address;
  }

  const realIp = request.headers.get("x-real-ip");
  if (realIp) {
    return realIp;
  }

  const forwarded = request.headers
    .get("x-forwarded-for")
    ?.split(",")[0]
    ?.trim();
  return forwarded || "unknown";
}

// 新条目且达到上限时腾出空间：
// 1. 先清理已过期条目（避免误淘汰仍在窗口内的活跃条目）
// 2. 仍超限则淘汰最旧一条（Map 迭代序即插入序）
// 这样可复用 60s 定时器的清理职责，降低对定时器的依赖
function makeRoomIfNeeded(now: number) {
  if (store.size < MAX_ENTRIES) return;
  for (const [key, entry] of store) {
    if (now > entry.resetAt) {
      store.delete(key);
      if (store.size < MAX_ENTRIES) return;
    }
  }
  const oldestKey = store.keys().next().value;
  if (oldestKey !== undefined) {
    store.delete(oldestKey);
  }
}

/**
 * 速率限制中间件
 * @param maxRequests 时间窗口内最大请求数
 * @param windowMs 时间窗口（毫秒）
 */
export function rateLimit(maxRequests: number, windowMs: number) {
  // status：Elysia 提供的"返回指定状态码响应"的函数（等价于旧版 context.error）
  return ({ request, server, status }: Context) => {
    const ip = getClientIp(request, server);
    const key = `rate:${ip}`;
    const now = Date.now();

    let entry = store.get(key);
    if (!entry || now > entry.resetAt) {
      entry = { count: 1, resetAt: now + windowMs };
      makeRoomIfNeeded(now);
      store.set(key, entry);
      return;
    }

    entry.count++;
    if (entry.count > maxRequests) {
      warn("ratelimit", `IP ${ip} 触发限流`, {
        count: entry.count,
        max: maxRequests,
      });
      return status(429, `请求过于频繁，请 ${Math.ceil((entry.resetAt - now) / 1000)} 秒后再试`);
    }
  };
}
