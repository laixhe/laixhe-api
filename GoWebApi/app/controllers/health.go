package controllers

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"

	"webapi/core"
)

// Version 服务版本: 构建时通过 -ldflags "-X main.GitVersion=xxx" 注入 (见 Makefile/Dockerfile),
// main 启动时同步到此变量 (见 main.go), 健康检查接口返回该版本号; 未注入时默认 "1.0.0"。
var Version = "1.0.0"

// startedAt 服务启动时间 (进程级变量: 启动时赋值一次, 之后只读, 服务器本地时区, 与 entity.User.CreatedAt 格式保持一致)
var startedAt = time.Now().Format(time.DateTime)

// healthPingInterval 健康检查中数据库探测结果的缓存时长。
// 探活请求可能非常频繁, 缓存一段时间可显著降低对数据库的压力;
// 代价是数据库故障后最多延迟该时长才会反映到健康检查结果上。
const healthPingInterval = 5 * time.Second

// Health 健康检查相关
type Health struct {
	server *core.Server

	// RWMutex: 缓存有效期内所有请求走 RLock 并发读, 互不阻塞;
	// 仅当缓存过期时由第一个请求加写锁执行 Ping (double-check 防并发重复探测)。
	mu       sync.RWMutex
	lastPing time.Time
	lastErr  error
}

// dbHealthy 探测数据库连接, 结果缓存 healthPingInterval 时长。
// 读路径 (缓存有效) 无锁竞争, 避免数据库故障时所有健康检查请求被同一把 Mutex 串行卡住。
func (c *Health) dbHealthy() error {
	c.mu.RLock()
	fresh := !c.lastPing.IsZero() && time.Since(c.lastPing) < healthPingInterval
	if fresh {
		err := c.lastErr
		c.mu.RUnlock()
		return err
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	// double-check: 并发第一个抢到写锁的请求已刷新缓存, 后续请求直接复用
	if !c.lastPing.IsZero() && time.Since(c.lastPing) < healthPingInterval {
		return c.lastErr
	}
	c.lastErr = c.server.Orm().Ping()
	c.lastPing = time.Now()
	return c.lastErr
}

func newHealth(server *core.Server) *Health {
	return &Health{server: server}
}

// HealthResponse 健康检查响应体
type HealthResponse struct {
	Status    string `json:"status" validate:"required"`     // 服务状态 (固定 "ok")
	Database  string `json:"database" validate:"required"`   // 数据库状态 (正常时为 "up"; 数据库不可用时直接返回 503 错误体, 不返回本字段)
	Version   string `json:"version" validate:"required"`    // 服务版本
	StartedAt string `json:"started_at" validate:"required"` // 服务启动时间 (服务器本地时区)
	Now       string `json:"now" validate:"required"`        // 当前时间 (服务器本地时区)
}

// Health
// @Summary	健康检查
// @Tags     Health
// @Produce  json
// @Success  200    {object}  controllers.HealthResponse
// @Failure  503    {object}  core.Error
// @Router   /api/v1/health [get]
//
// 通过 Ping 探测数据库连接, 正常返回 200 + 健康信息;
// 数据库不可用时返回 503 + 统一错误格式, 便于负载均衡探活。
func (c *Health) Health(ctx fiber.Ctx) error {
	resp := &HealthResponse{
		Status:    "ok",
		Database:  "up",
		Version:   Version,
		StartedAt: startedAt,
		Now:       time.Now().Format(time.DateTime),
	}
	// 探测数据库连接 (带缓存, 见 dbHealthy)
	if err := c.dbHealthy(); err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "database unavailable")
	}
	return ctx.JSON(resp)
}
