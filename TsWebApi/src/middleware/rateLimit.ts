// 简易内存速率限制（IP 维度），适用于登录等敏感接口防暴力破解

import { warn } from "../util/logger";

interface Entry {
  count: number;
  resetAt: number;
}

const store = new Map<string, Entry>();

// 过期条目清理定时器（每 60s 清理一次）
setInterval(() => {
  const now = Date.now();
  for (const [key, entry] of store) {
    if (now > entry.resetAt) {
      store.delete(key);
    }
  }
}, 60_000);

/**
 * 速率限制中间件
 * @param maxRequests 时间窗口内最大请求数
 * @param windowMs 时间窗口（毫秒）
 */
export function rateLimit(maxRequests: number, windowMs: number) {
  return ({ request, error }: any) => {
    const ip =
      request.headers.get("x-forwarded-for")?.split(",")[0]?.trim() ||
      request.headers.get("x-real-ip") ||
      "unknown";
    const key = `rate:${ip}`;
    const now = Date.now();

    let entry = store.get(key);
    if (!entry || now > entry.resetAt) {
      entry = { count: 1, resetAt: now + windowMs };
      store.set(key, entry);
      return;
    }

    entry.count++;
    if (entry.count > maxRequests) {
      warn("ratelimit", `IP ${ip} 触发限流`, {
        count: entry.count,
        max: maxRequests,
      });
      return error(429, `请求过于频繁，请 ${Math.ceil((entry.resetAt - now) / 1000)} 秒后再试`);
    }
  };
}
