package middlewares

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
)

// 限流常量
const (
	// 健康检查路径 (限流豁免, 避免负载均衡探活被误伤)
	rateLimitHealthPath = "/api/v1/health"
	// 分片数量: 按 key 哈希路由到不同锁, 降低高并发下的锁竞争
	rateLimitShards = 16
	// 内存保护: 每个分片最多跟踪的 key 数量, 防止伪造 IP 导致内存无限增长
	rateLimitMaxKeysPerShard = 100_000 / rateLimitShards
)

// RateLimiter IP 限流器: 以滑动窗口统计每个 key (IP) 的请求频率
//
// 内部按 key 哈希分片存储 (每片独立 Mutex), 避免单锁在高并发下的热点竞争。
type RateLimiter struct {
	max             int
	window          time.Duration
	buckets         []*rateLimitBucket
	stop            chan struct{}
	maxKeysPerShard int
}

// rateLimitBucket 分片桶: key → 窗口内的时间戳队列 (FIFO)
type rateLimitBucket struct {
	mu    sync.Mutex
	times map[string][]time.Time
}

// NewRateLimiter 创建限流器
//
// max    窗口内允许的最大请求数 (小于 1 时按 1 处理)
// window 滑动窗口时长, 超过该时长的旧记录会被清理
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	return newRateLimiter(max, window, rateLimitShards, rateLimitMaxKeysPerShard)
}

// newRateLimiter 可指定分片数与单分片 key 上限的内部构造 (默认常量见 rateLimitShards/rateLimitMaxKeysPerShard)。
// 分片数与 key 上限参数化主要是便于测试内存保护背压逻辑 (构造小分片/小上限实例)。
func newRateLimiter(max int, window time.Duration, shards, maxKeysPerShard int) *RateLimiter {
	if max < 1 {
		max = 1
	}
	buckets := make([]*rateLimitBucket, shards)
	for i := range buckets {
		buckets[i] = &rateLimitBucket{times: make(map[string][]time.Time)}
	}
	r := &RateLimiter{
		max:             max,
		window:          window,
		buckets:         buckets,
		stop:            make(chan struct{}),
		maxKeysPerShard: maxKeysPerShard,
	}
	// 后台定时清理各分片的过期 key, 与请求路径解耦 (见 janitor)
	go r.janitor()
	return r
}

// janitor 后台定期清理各分片中窗口外已无活动的 key, 防止内存无限增长。
// 与请求路径解耦: Check 只负责窗口内队首清理 (平摊 O(1)),
// 分片级别的整体扫描交给定时任务, 避免攻击者塞满分片后合法请求被每次 O(n) 全量清理拖慢。
func (r *RateLimiter) janitor() {
	ticker := time.NewTicker(r.window)
	defer ticker.Stop()
	for {
		select {
		case <-r.stop:
			return
		case now := <-ticker.C:
			for _, b := range r.buckets {
				b.mu.Lock()
				for k, ts := range b.times {
					if len(ts) == 0 || now.Sub(ts[len(ts)-1]) >= r.window {
						delete(b.times, k)
					}
				}
				b.mu.Unlock()
			}
		}
	}
}

// Close 停止后台清理协程。
// 生产环境由 NewMiddleware 创建单例并随进程存活, 无需调用;
// 主要用于测试/可复用限流器实例的场景, 避免协程泄漏。
func (r *RateLimiter) Close() {
	close(r.stop)
}

// shard 按 key 哈希路由到对应分片
//
// FNV-1a 哈希, 对同一 key 结果确定, 保证同一 key 始终命中同一分片 (计数不会分散)。
// 手写字节循环实现, 避免 fnv.New32a() 的接口分配与 string→[]byte 拷贝 (全局限流, 最高频路径)。
func (r *RateLimiter) shard(key string) *rateLimitBucket {
	h := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}
	return r.buckets[h%uint32(len(r.buckets))]
}

// Check 判断 key 是否允许通过 (会同时清理已过期的记录)
//
// 返回 true 表示未超限 (计数已入窗), false 表示超限。
//
// 性能说明: 时间戳按序追加, 过期的记录必然是队首前缀,
// 因此只清理队首即可 (平摊 O(1)),
// 避免窗口条目多时每次请求全量扫描 (O(n))。
func (r *RateLimiter) Check(key string) bool {
	b := r.shard(key)
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	times := b.times[key]
	// 记录 key 是否已存在 (须在 append 之前判断, 用于内存保护: 只有"新 key"才触发清理)
	_, existed := b.times[key]
	// 移除窗口外的前缀记录 (队首即最旧, times[1:] 仅前移切片头, 不复制)
	for len(times) > 0 && now.Sub(times[0]) >= r.window {
		times = times[1:]
	}
	if len(times) >= r.max {
		b.times[key] = times
		return false
	}
	times = append(times, now)
	b.times[key] = times

	// 内存保护: 分片 key 数超限时仅对新 key 触发清理 (已存在的 key 直接放行),
	// 避免合法请求在分片被塞满后每个请求都执行 O(n) 全量清理;
	// 过期 key 的回收由后台 janitor 定期执行。
	if len(b.times) > r.maxKeysPerShard && !existed {
		for k, ts := range b.times {
			if len(ts) == 0 || now.Sub(ts[len(ts)-1]) >= r.window {
				delete(b.times, k)
			}
		}
		if len(b.times) > r.maxKeysPerShard {
			clear(b.times)
		}
	}
	return true
}

// RateLimit IP 限流中间件:
// 取客户端 IP (优先 X-Forwarded-For, 其次真实连接地址),
// 超过阈值时返回 429 统一 JSON。
func (m *Middleware) RateLimit(c fiber.Ctx) error {
	// 配置关闭限流时直接放行
	if !m.limitConfig.Enable {
		return c.Next()
	}
	// 健康检查路径豁免限流 (负载均衡/容器编排探活不应被 429 拦截)
	if c.Path() == rateLimitHealthPath {
		return c.Next()
	}
	ip := resolveClientIP(c)
	if !m.rateLimiter.Check(ip) {
		m.server.Log().Warnf("接口限流触发 ip=%s", ip)
		// {"code": 429, "message": "请求过于频繁，请稍后再试"}
		return fiber.NewError(fiber.StatusTooManyRequests, "请求过于频繁，请稍后再试")
	}
	return c.Next()
}

// resolveClientIP 解析客户端 IP: 优先代理头 X-Forwarded-For, 其次真实连接地址,
// 兜底使用 "unknown"
//
// 注意: 直接信任 X-Forwarded-For 时客户端可伪造该头绕过限流,
// 生产环境建议仅在有可信反向代理时启用, 或直接去掉该分支只使用真实 IP。
func resolveClientIP(c fiber.Ctx) string {
	if xff := strings.TrimSpace(c.Get(fiber.HeaderXForwardedFor)); xff != "" {
		// 取第一个 (Cut 无需构造切片, 优于 Split)
		first, _, _ := strings.Cut(xff, ",")
		if ip := strings.TrimSpace(first); ip != "" {
			return ip
		}
	}
	if ip := c.IP(); ip != "" {
		return ip
	}
	return "unknown"
}
