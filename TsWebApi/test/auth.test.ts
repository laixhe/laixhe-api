// Auth 接口冒烟测试（bun test）
// 覆盖注册、登录的正常与异常路径
// 注意：测试会真实写入数据库，测试结束自动清理测试用户

import { describe, expect, it, afterAll } from "bun:test";
import app from "../src/index";
import { prisma } from "../src/lib/prisma";

// 使用时间戳生成唯一邮箱，避免与已有数据冲突
const testEmail = `test_${Date.now()}@example.com`;
const testPassword = "pass123";

// 注册成功后保存的 token，供鉴权用例使用
let registeredToken = "";
let registeredUid = 0;

// 构造请求的公共方法
function api(path: string, options: RequestInit = {}) {
  return app.handle(
    new Request(`http://localhost${path}`, {
      ...options,
      headers: { "content-type": "application/json", ...options.headers },
    })
  );
}

describe("auth API 冒烟测试", () => {
  // 测试结束后清理数据并断开数据库连接
  afterAll(async () => {
    await prisma.user.deleteMany({ where: { email: testEmail } });
    await prisma.$disconnect();
  });

  it("注册成功：返回 token 和用户信息", async () => {
    const res = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({
        nickname: "测试用户",
        email: testEmail,
        password: testPassword,
      }),
    });

    expect(res.status).toBe(200);
    const json = await res.json();
    expect(json.token).toBeTruthy();
    expect(json.user.email).toBe(testEmail);
    // 响应中不得包含密码
    expect(json.user.password).toBeUndefined();
    registeredToken = json.token;
    registeredUid = json.user.uid;
  });

  it("注册失败：重复邮箱返回 400", async () => {
    const res = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({
        nickname: "测试用户",
        email: testEmail,
        password: testPassword,
      }),
    });

    expect(res.status).toBe(400);
    expect((await res.json()).message).toBe("邮箱已存在");
  });

  it("注册失败：参数非法返回 400", async () => {
    const res = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ nickname: "测", email: "bad-email", password: "123" }),
    });

    expect(res.status).toBe(400);
  });

  it("登录成功：返回 token", async () => {
    const res = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: testEmail, password: testPassword }),
    });

    expect(res.status).toBe(200);
    expect((await res.json()).token).toBeTruthy();
  });

  it("登录失败：密码错误返回 400", async () => {
    const res = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: testEmail, password: "wrongpass" }),
    });

    expect(res.status).toBe(400);
    expect((await res.json()).message).toBe("邮箱或密码错误");
  });

  it("鉴权失败：无 token 访问 /user/update 返回 401", async () => {
    const res = await api("/api/v1/user/update", {
      method: "POST",
      body: JSON.stringify({ nickname: "新昵称", avatar_url: "" }),
    });

    expect(res.status).toBe(401);
    expect((await res.json()).message).toBe("未授权");
  });

  it("鉴权失败：refresh 无 token 返回 401", async () => {
    const res = await api("/api/v1/auth/refresh", { method: "POST" });

    expect(res.status).toBe(401);
    expect((await res.json()).message).toBe("未授权");
  });

  it("鉴权成功：带 token 更新用户信息", async () => {
    const res = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify({
        nickname: "新昵称",
        avatar_url: "https://example.com/avatar.png",
      }),
    });

    expect(res.status).toBe(200);
    const json = await res.json();
    expect(json.nickname).toBe("新昵称");
    expect(json.avatar_url).toBe("https://example.com/avatar.png");
  });

  it("鉴权成功：带 token 刷新 token", async () => {
    const res = await api("/api/v1/auth/refresh", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
    });

    expect(res.status).toBe(200);
    expect((await res.json()).token).toBeTruthy();
  });

  it("公开接口：无需 token 即可查询用户信息", async () => {
    const res = await api(`/api/v1/user/info?uid=${registeredUid}`);

    expect(res.status).toBe(200);
    expect((await res.json()).uid).toBe(registeredUid);
  });

  it("公开接口：无需 token 即可访问根路径", async () => {
    const res = await api("/");

    expect(res.status).toBe(200);
    expect(await res.text()).toContain("laixhe-api is running");
  });
});
