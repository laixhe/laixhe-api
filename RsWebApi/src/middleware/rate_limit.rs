//! 接口限流中间件 (基于客户端 IP 的滑动窗口)
//!
//! 超过阈值时返回 429 + 统一 JSON 格式: {"code": 429, "message": "请求过于频繁，请稍后再试"}

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
}

impl RateLimiter {
    /// 创建限流器
    pub fn new(max: usize, window: Duration) -> Self {
        let buckets = (0..SHARDS).map(|_| Mutex::new(HashMap::new())).collect();
        RateLimiter {
            max: max.max(1),
            window,
            buckets,
        }
    }

    /// 按 key 哈希路由到对应分片
    ///
    /// `DefaultHasher` 使用固定种子 (SipHash13)，对同一 key 结果确定，
    /// 保证同一 key 始终命中同一分片（计数不会分散）。
    fn shard(&self, key: &str) -> &Mutex<HashMap<String, VecDeque<Instant>>> {
        let mut hasher = DefaultHasher::new();
        key.hash(&mut hasher);
        let idx = (hasher.finish() as usize) % SHARDS;
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
        let mut buckets = self
            .shard(key)
            .lock()
            .unwrap_or_else(|poisoned| poisoned.into_inner());
        let now = Instant::now();
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
        // 内存保护: key 数量超限时, 先清理窗口内已无活动的 key, 仍超限则整体清空
        if buckets.len() > MAX_KEYS_PER_SHARD {
            buckets.retain(|_, d| {
                d.back()
                    .map_or(false, |t| now.duration_since(*t) < self.window)
            });
            if buckets.len() > MAX_KEYS_PER_SHARD {
                buckets.clear();
            }
        }
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
    let client_ip = resolve_client_ip(&req);

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

/// 解析客户端 IP: 优先代理头 X-Forwarded-For, 其次 ConnectInfo 真实地址,
/// 兜底使用 "unknown" (如 oneshot 测试场景)
///
/// 注意: 直接信任 X-Forwarded-For 时客户端可伪造该头绕过限流,
/// 生产环境建议仅在有可信反向代理时启用, 或直接去掉该分支只使用 ConnectInfo。
fn resolve_client_ip(req: &Request) -> String {
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
    if let Some(ConnectInfo(addr)) = req.extensions().get::<ConnectInfo<SocketAddr>>() {
        return addr.ip().to_string();
    }
    "unknown".to_string()
}
