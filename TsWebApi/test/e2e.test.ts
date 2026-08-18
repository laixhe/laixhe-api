// 端到端测试（E2E）：模拟真实用户流程 —— 注册 → 登录 → 携带 token 访问受保护接口
// 与 auth.test.ts（进程内 app.handle）不同，本测试启动真实 HTTP 服务并通过 fetch 访问，
// 覆盖完整网络层（HTTP 解析、JSON 序列化、请求头、状态码、限流等）

import { describe, expect, it, beforeAll, afterAll } from "bun:test";
import app from "../src/index";
import { prisma } from "../src/lib/prisma";

const testEmail = `e2e_${Date.now()}@example.com`;
const testPassword = "pass123";
const baseNickname = "E2E用户";
const updatedNickname = "E2E新昵称";

let token = "";
let uid = 0;
let baseUrl = "";

// 请求封装：自动携带 JSON 头并解析响应
async function api(path: string, options: RequestInit = {}) {
  const res = await fetch(`${baseUrl}${path}`, {
    ...options,
    headers: { "content-type": "application/json", ...options.headers },
  });
  return { status: res.status, body: await res.json() };
}

// 启动真实 HTTP 服务（端口 0 = 由系统分配空闲端口，避免与开发服务器 6600 冲突）
beforeAll(async () => {
  app.listen({ port: 0 });
  const port = app.server?.port;
  if (!port) throw new Error("服务器启动失败：未获取到监听端口");
  baseUrl = `http://127.0.0.1:${port}`;

  // 轮询等待服务就绪（最多 2 秒）
  for (let i = 0; i < 20; i++) {
    try {
      const res = await fetch(`${baseUrl}/`);
      if (res.ok) return;
    } catch {
      await Bun.sleep(100);
    }
  }
  throw new Error("服务器未在预期时间内就绪");
});

// 清理：删除测试用户、关闭服务、断开数据库连接
afterAll(async () => {
  await prisma.user.deleteMany({ where: { email: testEmail } });
  await prisma.$disconnect();
  app.server?.stop();
});

describe("端到端流程：注册 → 登录 → 受保护接口", () => {
  it("健康检查：GET / 返回运行中", async () => {
    const res = await fetch(`${baseUrl}/`);
    expect(res.status).toBe(200);
    expect(await res.text()).toContain("laixhe-api is running");
  });

  it("每个响应带 X-Request-Id 头 (与 Go/PHP/Rust 端对齐)", async () => {
    const res = await fetch(`${baseUrl}/api/v1/health`);
    expect(res.status).toBe(200);
    expect(res.headers.get("x-request-id")).toBeTruthy();
  });

  it("注册：返回 token 与用户信息（不含密码）", async () => {
    const { status, body } = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({
        nickname: baseNickname,
        email: testEmail,
        password: testPassword,
      }),
    });
    expect(status).toBe(200);
    expect(body.token).toBeTruthy();
    expect(body.user.email).toBe(testEmail);
    expect(body.user.password).toBeUndefined();
    token = body.token;
    uid = body.user.uid;
  });

  it("注册：重复邮箱返回 422", async () => {
    const { status, body } = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({
        nickname: baseNickname,
        email: testEmail,
        password: testPassword,
      }),
    });
    expect(status).toBe(422);
    expect(body.message).toBe("邮箱已存在");
  });

  it("登录：正确密码返回新 token", async () => {
    const { status, body } = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: testEmail, password: testPassword }),
    });
    expect(status).toBe(200);
    expect(body.token).toBeTruthy();
    expect(body.user.uid).toBe(uid);
  });

  it("登录：错误密码返回 422", async () => {
    const { status, body } = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: testEmail, password: "wrongpass" }),
    });
    expect(status).toBe(422);
    expect(body.message).toBe("邮箱或密码错误");
  });

  it("受保护接口：无 token 访问 /user/update 返回 401", async () => {
    const { status, body } = await api("/api/v1/user/update", {
      method: "POST",
      body: JSON.stringify({ nickname: "新昵称", avatar_url: "" }),
    });
    expect(status).toBe(401);
    expect(body.message).toBe("Unauthorized");
  });

  it("受保护接口：带 token 更新用户信息", async () => {
    const { status, body } = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${token}` },
      body: JSON.stringify({
        nickname: updatedNickname,
        avatar_url: "https://example.com/avatar.png",
      }),
    });
    expect(status).toBe(200);
    expect(body.nickname).toBe(updatedNickname);
    expect(body.avatar_url).toBe("https://example.com/avatar.png");
  });

  it("受保护接口：带 token 刷新 token", async () => {
    const { status, body } = await api("/api/v1/auth/refresh", {
      method: "POST",
      headers: { authorization: `Bearer ${token}` },
    });
    expect(status).toBe(200);
    expect(body.token).toBeTruthy();
    expect(body.user.uid).toBe(uid);
  });

  it("公开接口：查询用户信息（无需 token）", async () => {
    const { status, body } = await api(`/api/v1/user/info?uid=${uid}`);
    expect(status).toBe(200);
    expect(body.uid).toBe(uid);
  });

  it("公开接口：用户列表分页", async () => {
    const { status, body } = await api("/api/v1/user/list?page=1&page_size=5");
    expect(status).toBe(200);
    expect(Array.isArray(body.list)).toBe(true);
    expect(body.list.length).toBeGreaterThan(0);
    expect(body.total).toBeGreaterThan(0);
  });
});
