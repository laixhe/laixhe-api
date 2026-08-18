package routers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

// TestNotFoundUnifiedJSON 验证 404 兜底模式: 注册在业务路由之后的 Use 对未匹配路由返回统一 JSON
// (与 router.go init() 末尾的 404 兜底逻辑等价, 不依赖数据库)
func TestNotFoundUnifiedJSON(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/health", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})
	// 与 init() 中一致: 404 兜底注册在业务路由之后
	app.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.NewError(fiber.StatusNotFound, "Not Found"))
	})

	// 已匹配路由正常返回 200
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("已匹配路由应返回 200, got %d", resp.StatusCode)
	}

	// 未匹配路由 → 404 统一 JSON (与 Rust/TS/PHP 端格式一致)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/no-such-route", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("未匹配路由应返回 404, got %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("读取响应失败: %v", err)
	}
	// JSON 字段顺序与框架序列化有关, 按字段内容断言 (与 Rust/TS/PHP 端格式一致)
	var got struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("404 响应应为 JSON, got %s: %v", string(body), err)
	}
	if got.Code != http.StatusNotFound || got.Message != "Not Found" {
		t.Fatalf("404 响应应为 {code:404,message:\"Not Found\"}, got %s", string(body))
	}
}
