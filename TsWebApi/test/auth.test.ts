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

  it("注册失败：重复邮箱返回 422", async () => {
    const res = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({
        nickname: "测试用户",
        email: testEmail,
        password: testPassword,
      }),
    });

    expect(res.status).toBe(422);
    expect((await res.json()).message).toBe("邮箱已存在");
  });

  it("注册失败：参数非法返回 422", async () => {
    const res = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ nickname: "测", email: "bad-email", password: "123" }),
    });

    expect(res.status).toBe(422);
  });

  it("登录成功：返回 token", async () => {
    const res = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: testEmail, password: testPassword }),
    });

    expect(res.status).toBe(200);
    expect((await res.json()).token).toBeTruthy();
  });

  it("登录失败：密码错误返回 422", async () => {
    const res = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: testEmail, password: "wrongpass" }),
    });

    expect(res.status).toBe(422);
    expect((await res.json()).message).toBe("邮箱或密码错误");
  });

  it("鉴权失败：无 token 访问 /user/update 返回 401", async () => {
    const res = await api("/api/v1/user/update", {
      method: "POST",
      body: JSON.stringify({ nickname: "新昵称", avatar_url: "" }),
    });

    expect(res.status).toBe(401);
    expect((await res.json()).message).toBe("Unauthorized");
  });

  it("鉴权失败：refresh 无 token 返回 401", async () => {
    const res = await api("/api/v1/auth/refresh", { method: "POST" });

    expect(res.status).toBe(401);
    expect((await res.json()).message).toBe("Unauthorized");
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

  it("列表分页钳制：page=0 归一为 1, page_size=999 钳制为 100", async () => {
    const res = await api("/api/v1/user/list?page=0&page_size=999");

    expect(res.status).toBe(200);
    const json = await res.json();
    expect(json.page).toBe(1);
    expect(json.page_size).toBe(100);
    expect(Array.isArray(json.list)).toBe(true);
  });

  it("列表分页钳制：page_size=0 回落默认 12 (与 Go/Rust/PHP 端一致)", async () => {
    const res = await api("/api/v1/user/list?page=1&page_size=0");

    expect(res.status).toBe(200);
    const json = await res.json();
    expect(json.page).toBe(1);
    expect(json.page_size).toBe(12);
  });

  it("列表分页钳制：page/page_size 负数归一为 1/12", async () => {
    const res = await api("/api/v1/user/list?page=-3&page_size=-5");

    expect(res.status).toBe(200);
    const json = await res.json();
    expect(json.page).toBe(1);
    expect(json.page_size).toBe(12);
  });

  it("列表非数字参数返回 400 (与 Go/Rust 端绑定层行为一致)", async () => {
    const res = await api("/api/v1/user/list?page=abc&page_size=12");
    expect(res.status).toBe(400);
    expect((await res.json()).code).toBe(400);

    const res2 = await api("/api/v1/user/list?page=1&page_size=xyz");
    expect(res2.status).toBe(400);
  });

  it("info 非数字 uid 返回 400 (与 Go/Rust 端绑定层行为一致)", async () => {
    const res = await api("/api/v1/user/info?uid=abc");

    expect(res.status).toBe(400);
    expect((await res.json()).code).toBe(400);
  });

  it("register/login body 数字字段返回 400 而非 500", async () => {
    // nickname 为数字: 字段类型校验应返回 400 (与 Go/Rust 绑定层行为一致), 而非 Prisma 报错 500
    const res = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ nickname: 123, email: "a@b.com", password: "pass123" }),
    });
    expect(res.status).toBe(400);
    expect((await res.json()).code).toBe(400);

    // password 为数字
    const res2 = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: "a@b.com", password: 123 }),
    });
    expect(res2.status).toBe(400);
  });

  it("update body 数字字段返回 400 而非 500", async () => {
    const res = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify({ nickname: 123 }),
    });

    expect(res.status).toBe(400);
    expect((await res.json()).code).toBe(400);
  });

  it("body 布尔/数组字段返回 400 (与 Go/Rust 端绑定层行为一致)", async () => {
    // 布尔字段
    const boolRes = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify({ nickname: true }),
    });
    expect(boolRes.status).toBe(400);

    // 数组字段
    const arrRes = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify({ nickname: ["a", "b"] }),
    });
    expect(arrRes.status).toBe(400);
  });

  it("register/login body null 字段返回 422 而非 400 (与 Go/PHP 端一致)", async () => {
    // null 视为"无值": 归一化为空串后走业务校验, 返回具体 422 文案
    const res = await api("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ nickname: null, email: "a@b.com", password: "pass123" }),
    });
    expect(res.status).toBe(422);
    expect((await res.json()).message).toBe("昵称长度不能小于2位");

    const res2 = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: null, password: "pass123" }),
    });
    expect(res2.status).toBe(422);
    expect((await res2.json()).message).toBe("邮箱格式错误");
  });

  it("update body null 字段返回 422 而非 400 (与 Go/PHP 端一致)", async () => {
    const res = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify({ nickname: null }),
    });

    expect(res.status).toBe(422);
    expect((await res.json()).message).toBe("昵称长度不能小于2位");
  });

  it("顶层 body 非对象返回 400 (与 Go/Rust 端绑定层行为一致)", async () => {
    // 用 update 接口验证 (无限流; register/login 的顶层检查逻辑与其完全相同)
    const res = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify([1, 2]),
    });
    expect(res.status).toBe(400);

    const res2 = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: "null",
    });
    expect(res2.status).toBe(400);

    const res3 = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify("hello"),
    });
    expect(res3.status).toBe(400);
  });

  it("昵称 emoji 按 Unicode 码点计数 (20 个通过, 21 个拒绝)", async () => {
    // 修复前 string.length 按 UTF-16 单元计数, emoji 会被误判为 2 位 (20 个 emoji 会被误拒)
    const ok = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify({ nickname: "😀".repeat(20) }),
    });
    expect(ok.status).toBe(200);

    const tooLong = await api("/api/v1/user/update", {
      method: "POST",
      headers: { authorization: `Bearer ${registeredToken}` },
      body: JSON.stringify({ nickname: "😀".repeat(21) }),
    });
    expect(tooLong.status).toBe(422);
    expect((await tooLong.json()).message).toBe("昵称长度不能超过20位");
  });

  it("超大整数参数返回 400 (与 Go/Rust 端绑定层溢出行为一致)", async () => {
    // 超出 JS 安全整数范围 (Number 精度丢失), 其他端绑定层溢出直接 400
    const res = await api("/api/v1/user/info?uid=99999999999999999999");
    expect(res.status).toBe(400);

    const res2 = await api("/api/v1/user/list?page=99999999999999999999&page_size=12");
    expect(res2.status).toBe(400);
  });

  it("缺 body 登录返回 422 而非 500", async () => {
    // 不传 body 时 handler 解构 undefined 抛 TypeError, 由全局 onError 兜底为 422
    const res = await api("/api/v1/auth/login", { method: "POST" });

    expect(res.status).toBe(422);
    expect((await res.json()).message).toBe("参数错误");
  });

  it("空 body + content-type 登录返回 400 (JSON 解析失败)", async () => {
    // 带 content-type: application/json 但 body 为空: Elysia 走严格 JSON 解析 → PARSE → 400
    const res = await api("/api/v1/auth/login", {
      method: "POST",
      body: "",
      headers: { "content-type": "application/json" },
    });

    expect(res.status).toBe(400);
  });

  // 登录限流 (5 次/分钟) 测试放在最后, 避免触发后影响前面用例;
  // 注意: 本测试会消耗该 IP 的登录配额, 60s 后自动恢复
  it("登录限流：同一 IP 连发超过阈值返回 429", async () => {
    let lastStatus = 0;
    for (let i = 0; i < 6; i++) {
      const res = await api("/api/v1/auth/login", {
        method: "POST",
        body: JSON.stringify({ email: `nobody_${i}@example.com`, password: "wrong123" }),
      });
      lastStatus = res.status;
    }
    expect(lastStatus).toBe(429);
    const res = await api("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email: "nobody_final@example.com", password: "wrong123" }),
    });
    expect(res.status).toBe(429);
    expect((await res.json()).code).toBe(429);
  });
});
