package middlewares

import (
	"testing"
	"time"
)

// TestRateLimiterCheck 验证滑动窗口计数: 超限拒绝, 不同 key 互不影响
func TestRateLimiterCheck(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)

	if !rl.Check("ip1") {
		t.Fatal("第 1 次请求应放行")
	}
	if !rl.Check("ip1") {
		t.Fatal("第 2 次请求应放行")
	}
	if rl.Check("ip1") {
		t.Fatal("第 3 次请求应被拒绝")
	}
	// 不同 key 共享同一限流器但计数独立
	if !rl.Check("ip2") {
		t.Fatal("不同 key 的请求应放行")
	}
}

// TestRateLimiterWindowExpire 验证窗口过期后旧记录被清理, 请求恢复放行
func TestRateLimiterWindowExpire(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)

	if !rl.Check("ip1") {
		t.Fatal("第 1 次请求应放行")
	}
	if rl.Check("ip1") {
		t.Fatal("窗口内超限应被拒绝")
	}
	time.Sleep(60 * time.Millisecond)
	if !rl.Check("ip1") {
		t.Fatal("窗口过期后应恢复放行")
	}
}

// TestRateLimiterMaxBelowOne 验证非法 max 被兜底为 1
func TestRateLimiterMaxBelowOne(t *testing.T) {
	rl := NewRateLimiter(0, time.Minute)
	if !rl.Check("ip1") {
		t.Fatal("第 1 次请求应放行")
	}
	if rl.Check("ip1") {
		t.Fatal("max 兜底为 1 后, 第 2 次请求应被拒绝")
	}
}

// TestRateLimiterMemoryProtectionClear 验证内存保护背压:
// 分片 key 数超过上限 (len > maxKeysPerShard) 后, 新 key 触发清理;
// 无过期 key 时仍超限则整体清空 (防内存无限增长), 且已存在的 key 请求不会被拖入 O(n) 清理
func TestRateLimiterMemoryProtectionClear(t *testing.T) {
	// 1 分片, 每分片上限 2 个 key (参数化构造, 见 newRateLimiter);
	// 注意触发条件是 len > maxKeysPerShard, 故第 4 个新 key 才触发清理
	rl := newRateLimiter(100, time.Minute, 1, 2)
	if !rl.Check("ip1") || !rl.Check("ip2") || !rl.Check("ip3") {
		t.Fatal("未超上限前新 key 应正常放行")
	}
	// 第 4 个新 key ip4: 分片超限触发清理 (ip1/ip2/ip3 均在窗口内无过期) → 仍超限 → 整体清空
	if !rl.Check("ip4") {
		t.Fatal("内存保护清空后 ip4 应放行")
	}
	// 整体清空后, 旧 key 计数被重置, 可再次放行
	if !rl.Check("ip1") {
		t.Fatal("清空后 ip1 计数重置, 应可再次放行")
	}
}

// TestRateLimiterMemoryProtectionPruneExpired 验证内存保护优先清理过期 key:
// 分片超限时新 key 触发清理只删除过期项, 保留窗口内活跃 key
func TestRateLimiterMemoryProtectionPruneExpired(t *testing.T) {
	rl := newRateLimiter(100, 50*time.Millisecond, 1, 2)
	if !rl.Check("ip1") {
		t.Fatal("ip1 应放行")
	}
	time.Sleep(60 * time.Millisecond) // ip1 窗口过期
	if !rl.Check("ip2") || !rl.Check("ip3") {
		t.Fatal("ip2/ip3 应放行")
	}
	// 新 key ip4: 清理过期 ip1, 保留活跃 ip2/ip3, 分片未超限 → 放行
	if !rl.Check("ip4") {
		t.Fatal("过期 key 被清理后, ip4 应放行")
	}
	// ip2 仍活跃, 计数未被清空
	if !rl.Check("ip2") {
		t.Fatal("ip2 计数应保留 (未被整体清空)")
	}
}
