// 简易内存速率限制（IP 维度），适用于登录等敏感接口防暴力破解
// 注意：内存实现仅适用于单实例部署，多实例需改用 Redis 等共享存储
// 实现为滑动窗口（时间戳队列），与 Go/PHP/Rust 端限流算法保持一致

import type { Context } from "elysia";
import { warn } from "../util/logger";

interface Entry {
  // 窗口内的时间戳队列（毫秒），升序排列；窗口过期后整体删除
  timestamps: number[];
  // 该条目所属限流器的窗口时长（每个 key 只被一个限流器使用, 故窗口固定）
  windowMs: number;
}

const store = new Map<string, Entry>();

// 最大条目数上限，防止攻击者用大量伪造 IP 撑爆内存
const MAX_ENTRIES = 10_000;

// 过期条目清理定时器（每 60s 清理一次）
// unref()：不持有事件循环，避免阻塞进程退出（测试/优雅关闭时友好）
setInterval(() => {
  const now = Date.now();
  for (const [key, entry] of store) {
    const last = entry.timestamps[entry.timestamps.length - 1];
    if (!last || now - last > entry.windowMs) {
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
// 1. 先清理窗口已完全过期的条目（仅扫描有界预算, 避免 store 被伪造 IP 塞满时
//    每个新 key 都触发对全部条目的 O(n) 扫描——过期条目的整体回收主要依赖上方 60s 定时器,
//    这里的扫描只是新插入路径上的补充)
// 2. 仍超限则淘汰最旧一条（Map 迭代序即插入序）
function makeRoomIfNeeded(now: number) {
  if (store.size < MAX_ENTRIES) return;
  // 有界扫描预算: 单次插入最多付出固定工作量, 最坏情况退化为 O(预算)
  const SCAN_BUDGET = 200;
  let scanned = 0;
  for (const [key, entry] of store) {
    if (++scanned > SCAN_BUDGET) break;
    const last = entry.timestamps[entry.timestamps.length - 1];
    if (!last || now - last > entry.windowMs) {
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
 * @param name 限流器名称, 用于隔离不同限流器的计数
 *              (全局限流与各单端点限流使用独立计数, 互不干扰)
 */
export function rateLimit(maxRequests: number, windowMs: number, name: string) {
  // status：Elysia 提供的"返回指定状态码响应"的函数（等价于旧版 context.error）
  return ({ request, server, status }: Context) => {
    const ip = getClientIp(request, server);
    const key = `rate:${name}:${ip}`;
    const now = Date.now();

    let entry = store.get(key);
    if (!entry) {
      entry = { timestamps: [now], windowMs };
      makeRoomIfNeeded(now);
      store.set(key, entry);
      return;
    }

    // 清理窗口外的旧时间戳：时间戳升序追加, 过期记录必为队首前缀, 故只需弹队首
    // 注意: JS 数组 shift() 会移动剩余元素, 实际为 O(n) 而非 O(1)
    // (Go 端 times[1:] 切片头不复制才是 O(1)); 窗口内条数受 maxRequests 限制,
    // 量级很小, 此处成本可接受, 与 Go/Rust 端滑动窗口的"只清理队首"语义一致
    const cutoff = now - windowMs;
    while (entry.timestamps.length > 0 && entry.timestamps[0] <= cutoff) {
      entry.timestamps.shift();
    }

    if (entry.timestamps.length >= maxRequests) {
      warn("ratelimit", `IP ${ip} 触发限流`, {
        count: entry.timestamps.length,
        max: maxRequests,
        name,
      });
      // 统一 JSON 格式与文案 (对齐 Go/PHP/Rust 版): {"code":429,"message":"请求过于频繁，请稍后再试"}
      return status(429, { code: 429, message: "请求过于频繁，请稍后再试" });
    }

    entry.timestamps.push(now);
  };
}

/**
 * 全局限流中间件 (对齐 Go/PHP/Rust 端: 单个 IP 在 60s 滑动窗口内最多 1000 次请求)
 * 健康检查路径豁免限流, 避免负载均衡探活被误伤 (与 Go/PHP/Rust 端一致)
 */
export function globalRateLimit(ctx: Context) {
  // 健康检查路径豁免限流: 用前缀匹配兼容带/不带尾斜杠两种探活 URL
  // (Elysia 的 loose-path 使 /api/v1/health 与 /api/v1/health/ 都可访问)
  if (ctx.path.startsWith("/api/v1/health")) return;
  return rateLimit(1000, 60_000, "global")(ctx);
}
