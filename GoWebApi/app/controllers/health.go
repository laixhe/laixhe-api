package controllers

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"

	"webapi/core"
)

// version 服务版本 (可通过 -ldflags "-X webapi/app/controllers.version=xxx" 注入)
var version = "1.0.0"

// startedAt 服务启动时间 (进程级常量, 服务器本地时区, 与 entity.User.CreatedAt 格式保持一致)
var startedAt = time.Now().Format(time.DateTime)

// healthPingInterval 健康检查中数据库探测结果的缓存时长。
// 探活请求可能非常频繁, 缓存一段时间可显著降低对数据库的压力;
// 代价是数据库故障后最多延迟该时长才会反映到健康检查结果上。
const healthPingInterval = 5 * time.Second

// Health 健康检查相关
type Health struct {
	server *core.Server

	mu       sync.Mutex
	lastPing time.Time
	lastErr  error
}

// dbHealthy 探测数据库连接, 结果缓存 healthPingInterval 时长。
func (c *Health) dbHealthy() error {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	Status    string `json:"status"`     // 服务状态 (固定 "ok")
	Database  string `json:"database"`   // 数据库状态 (固定 "up")
	Version   string `json:"version"`    // 服务版本
	StartedAt string `json:"started_at"` // 服务启动时间 (服务器本地时区)
	Now       string `json:"now"`        // 当前时间 (服务器本地时区)
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
		Version:   version,
		StartedAt: startedAt,
		Now:       time.Now().Format(time.DateTime),
	}
	// 探测数据库连接 (带缓存, 见 dbHealthy)
	if err := c.dbHealthy(); err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "database unavailable")
	}
	return ctx.JSON(resp)
}
