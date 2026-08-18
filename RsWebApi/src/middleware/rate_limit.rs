//! 接口限流中间件 (基于客户端 IP 的滑动窗口)
//!
//! 超过阈值时返回 429 + 统一 JSON 格式: {"code": 429, "message": "请求过于频繁，请稍后再试"}
//!
//! # 已知限制
//!
//! 限流器为进程内内存实现 (与 Go 版一致): 多实例部署时每个实例独立计数,
//! 实际限流上限会随实例数放大。如需多实例全局限流, 需引入 Redis 等共享存储。

use std::collections::{HashMap, VecDeque};
use std::hash::{DefaultHasher, Hash, Hasher};
use std::net::SocketAddr;
use std::sync::Mutex;
use std::time::{Duration, Instant};

use axum::extract::{ConnectInfo, Request, State};
use axum::http::StatusCode;
use axum::middleware::Next;
use axum::response::{IntoResponse, Response};
use axum::Json;

use crate::state::AppState;

/// 健康检查路径 (限流豁免, 避免负载均衡探活被误伤)
const HEALTH_PATH: &str = "/api/v1/health";

/// 分片数量: 按 key 哈希路由到不同锁, 降低高并发下的锁竞争
const SHARDS: usize = 16;

/// 内存保护: 每个分片最多跟踪的 key 数量, 防止伪造 IP 导致内存无限增长
const MAX_KEYS_PER_SHARD: usize = 100_000 / SHARDS;

/// IP 限流器: 以滑动窗口统计每个 key (IP) 的请求频率
///
/// 内部按 key 哈希分片存储 (每片独立 Mutex), 避免单锁在高并发下的热点竞争。
#[derive(Debug)]
pub struct RateLimiter {
    /// 窗口内允许的最大请求数
    max: usize,
    /// 滑动窗口时长
    window: Duration,
    /// 分片桶: key → 窗口内的时间戳队列 (FIFO)
    buckets: Vec<Mutex<HashMap<String, VecDeque<Instant>>>>,
    /// 每个分片最多跟踪的 key 数量 (内存保护)
    max_keys_per_shard: usize,
}

impl RateLimiter {
    /// 创建限流器 (使用默认分片数与单分片 key 上限, 见 SHARDS / MAX_KEYS_PER_SHARD)
    pub fn new(max: usize, window: Duration) -> Self {
        Self::with_capacity(max, window, SHARDS, MAX_KEYS_PER_SHARD)
    }

    /// 可指定分片数与单分片 key 上限的内部构造。
    /// 参数化主要是便于测试内存保护背压逻辑 (构造小分片/小上限实例)。
    fn with_capacity(max: usize, window: Duration, shards: usize, max_keys_per_shard: usize) -> Self {
        let buckets = (0..shards).map(|_| Mutex::new(HashMap::new())).collect();
        RateLimiter {
            max: max.max(1),
            window,
            buckets,
            max_keys_per_shard,
        }
    }

    /// 按 key 哈希路由到对应分片
    ///
    /// `DefaultHasher` 使用固定种子 (SipHash13)，对同一 key 结果确定，
    /// 保证同一 key 始终命中同一分片（计数不会分散）。
    fn shard(&self, key: &str) -> &Mutex<HashMap<String, VecDeque<Instant>>> {
        let mut hasher = DefaultHasher::new();
        key.hash(&mut hasher);
        let idx = (hasher.finish() as usize) % self.buckets.len();
        &self.buckets[idx]
    }

    /// 判断 key 是否允许通过 (会同时清理已过期的记录)
    ///
    /// # 示例
    ///
    /// ```ignore
    /// let limiter = RateLimiter::new(100, Duration::from_secs(60));
    /// if limiter.check(&client_ip) {
    ///     // 未超限, 放行
    /// } else {
    ///     // 超限, 返回 429
    /// }
    /// ```
    pub fn check(&self, key: &str) -> bool {
        // 分片锁: 若临界区内曾 panic 导致锁中毒, 取回数据继续使用,
        // 避免后续每次请求都 panic (限流计数场景数据一致性可接受)
        let mut buckets = crate::sync::lock_unpoison(self.shard(key));
        let now = Instant::now();
        // 内存保护: 新 key 且分片已满时, 先清理窗口内已无活动的 key, 仍满则拒绝本次请求。
        // 注意不能整体 clear(): 攻击者可批量伪造 IP 触发全量清零, 从而绕过自己已触发的
        // 限流计数 (旧 IP 的计数被一并重置); 拒绝新 key 才是安全的背压策略。
        if !buckets.contains_key(key) && buckets.len() >= self.max_keys_per_shard {
            buckets.retain(|_, d| {
                d.back()
                    .map_or(false, |t| now.duration_since(*t) < self.window)
            });
            if buckets.len() >= self.max_keys_per_shard {
                return false;
            }
        }
        // 优先复用已有 key 的队列, 仅在首次出现时分配 String (全局限流中间件, 最高频路径)
        let deque = match buckets.get_mut(key) {
            Some(d) => d,
            None => buckets.entry(key.to_string()).or_default(),
        };
        // 移除窗口外的旧记录: 时间戳按序追加, 过期记录必为队首前缀,
        // 只弹队首即平摊 O(1), 避免每次请求全量扫描 (与 Go 版 times[1:] 语义一致)
        while let Some(&front) = deque.front() {
            if now.duration_since(front) >= self.window {
                deque.pop_front();
            } else {
                break;
            }
        }
        if deque.len() >= self.max {
            return false;
        }
        deque.push_back(now);
        true
    }
}

/// 全局限流中间件:
/// 取客户端 IP (优先 X-Forwarded-For, 其次 ConnectInfo 真实地址),
/// 超过阈值时返回 429 统一 JSON
pub async fn rate_limit_middleware(
    State(state): State<AppState>,
    req: Request,
    next: Next,
) -> Response {
    // 配置关闭限流时直接放行
    if !state.config.limit.enable {
        return next.run(req).await;
    }
    // 健康检查路径豁免限流 (负载均衡/容器编排探活不应被 429 拦截)
    if req.uri().path() == HEALTH_PATH {
        return next.run(req).await;
    }
    let client_ip = resolve_client_ip(&req, state.config.limit.trust_proxy);

    if !state.limiter.check(&client_ip) {
        tracing::warn!(ip = %client_ip, "接口限流触发");
        return (
            StatusCode::TOO_MANY_REQUESTS,
            Json(serde_json::json!({
                "code": StatusCode::TOO_MANY_REQUESTS.as_u16(),
                "message": "请求过于频繁，请稍后再试",
            })),
        )
            .into_response();
    }
    next.run(req).await
}

/// 解析客户端 IP
///
/// - `trust_proxy = true`: 优先取代理头 X-Forwarded-For 的第一个 IP (仅应在可信反向代理
///   之后部署时开启, 此时代理会覆写该头, 客户端无法伪造);
/// - `trust_proxy = false` (默认): 直接使用 ConnectInfo 中的真实对端地址, 客户端无法伪造;
/// - 两者均取不到时兜底 "unknown" (如 oneshot 测试场景)。
fn resolve_client_ip(req: &Request, trust_proxy: bool) -> String {
    if trust_proxy {
        if let Some(ip) = req
            .headers()
            .get("X-Forwarded-For")
            .and_then(|v| v.to_str().ok())
            .and_then(|s| s.split(',').next())
            .map(|s| s.trim().to_string())
            .filter(|s| !s.is_empty())
        {
            return ip;
        }
    }
    if let Some(ConnectInfo(addr)) = req.extensions().get::<ConnectInfo<SocketAddr>>() {
        return addr.ip().to_string();
    }
    "unknown".to_string()
}

#[cfg(test)]
mod tests {
    use super::*;

    /// 窗口滑动边界: 窗口内超限拒绝, 窗口过期后恢复放行
    #[test]
    fn window_sliding_boundary() {
        let rl = RateLimiter::new(2, Duration::from_millis(50));
        assert!(rl.check("ip1"), "第 1 次应放行");
        assert!(rl.check("ip1"), "第 2 次应放行");
        assert!(!rl.check("ip1"), "第 3 次窗口内应拒绝");
        std::thread::sleep(Duration::from_millis(60));
        assert!(rl.check("ip1"), "窗口过期后应恢复放行");
    }

    /// 过期清理: 窗口过期后旧时间戳应被弹出 (队列不保留陈旧记录)
    #[test]
    fn stale_records_pruned() {
        let rl = RateLimiter::new(3, Duration::from_millis(30));
        for _ in 0..3 {
            assert!(rl.check("ip1"));
        }
        assert!(!rl.check("ip1"), "未过期前超限应拒绝");
        std::thread::sleep(Duration::from_millis(100));
        // 旧记录全部过期弹出后, 新请求视为从空窗口开始
        assert!(rl.check("ip1"), "过期记录被清理后应放行");
        assert!(rl.check("ip1"), "再次放行 (已清空旧计数)");
    }

    /// 不同 key 计数独立 (同一限流器, 互不影响)
    #[test]
    fn different_keys_independent() {
        let rl = RateLimiter::new(1, Duration::from_secs(60));
        assert!(rl.check("ip1"));
        assert!(!rl.check("ip1"), "ip1 超限");
        assert!(rl.check("ip2"), "ip2 与 ip1 计数独立, 应放行");
    }

    /// 同 IP 命中同一分片: 哈希路由确定性保证计数不分散
    #[test]
    fn same_key_same_shard() {
        let rl = RateLimiter::new(5, Duration::from_secs(60));
        let key = "192.168.1.100";
        let shard1 = rl.shard(key) as *const _;
        for _ in 0..5 {
            assert!(rl.check(key), "同一 key 前 5 次应放行 (计数未被分散)");
        }
        assert!(!rl.check(key), "第 6 次应拒绝");
        let shard2 = rl.shard(key) as *const _;
        assert_eq!(shard1, shard2, "同一 key 必须始终路由到同一分片");
    }

    /// 内存保护背压: 分片已满时新 key 被拒绝, 已存在的 key 不受影响 (防整体清空被绕过)
    #[test]
    fn memory_protection_backpressure() {
        let rl = RateLimiter::with_capacity(100, Duration::from_secs(60), 1, 2);
        assert!(rl.check("ip1"));
        assert!(rl.check("ip2"));
        // 分片满 (2/2), 新 key ip3 触发清理 (无过期) → 仍满 → 拒绝
        assert!(!rl.check("ip3"), "分片已满时新 key 应被拒绝 (背压)");
        // 已存在的 key 正常计数, 不被拒绝
        assert!(rl.check("ip1"), "已存在的 key 不受背压影响");
    }

    /// 内存保护优先清理过期 key: 分片满但存在过期 key 时, 新 key 清理后放行
    #[test]
    fn memory_protection_prunes_expired() {
        let rl = RateLimiter::with_capacity(100, Duration::from_millis(30), 1, 2);
        assert!(rl.check("ip1"));
        std::thread::sleep(Duration::from_millis(100)); // ip1 窗口过期
        assert!(rl.check("ip2"));
        // 新 key ip3: 清理过期 ip1 → 剩 ip2 → 未满 → 放行
        assert!(rl.check("ip3"), "过期 key 被清理后, 新 key 应放行");
    }

    /// max 非法值兜底为 1
    #[test]
    fn max_below_one_fallback() {
        let rl = RateLimiter::new(0, Duration::from_secs(60));
        assert!(rl.check("ip1"));
        assert!(!rl.check("ip1"), "max 兜底为 1 后第 2 次应拒绝");
    }
}
