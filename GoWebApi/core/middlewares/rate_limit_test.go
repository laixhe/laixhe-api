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
