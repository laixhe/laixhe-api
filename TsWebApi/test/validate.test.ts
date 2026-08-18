// validate.ts 纯逻辑单元测试（bun test）
// 不依赖数据库/网络，直接验证参数校验函数的边界行为，
// 与 auth.test.ts（HTTP 冒烟）互补，覆盖 bodyError 联合类型契约等底层逻辑。

import { describe, expect, it } from "bun:test";
import {
  bodyError,
  isEmail,
  isNicknameTooLong,
  isNicknameTooShort,
  isPasswordInvalid,
  isPasswordTooLong,
  isPasswordTooShort,
  nonStringField,
  normalizeNulls,
} from "../src/util/validate";

describe("isEmail", () => {
  it("合法邮箱通过", () => {
    expect(isEmail("user@example.com")).toBe(true);
    expect(isEmail("a.b+c@sub.domain.cn")).toBe(true);
  });

  it("非法邮箱拒绝", () => {
    expect(isEmail("not-an-email")).toBe(false);
    expect(isEmail("")).toBe(false);
    expect(isEmail("a@b")).toBe(false);
    expect(isEmail("a@b.")).toBe(false);
  });
});

describe("密码校验 (长度 6~64, 仅字母数字 _ @ $)", () => {
  it("长度边界: 5/6/64/65", () => {
    expect(isPasswordTooShort("12345")).toBe(true);
    expect(isPasswordTooShort("123456")).toBe(false);
    expect(isPasswordTooLong("a".repeat(64))).toBe(false);
    expect(isPasswordTooLong("a".repeat(65))).toBe(true);
  });

  it("缺字段 (undefined) 视为空串, 不抛 TypeError", () => {
    expect(isPasswordTooShort(undefined)).toBe(true);
    expect(isPasswordTooLong(undefined)).toBe(false);
    expect(isPasswordInvalid(undefined)).toBe(true);
  });

  it("字符集非法拒绝, 合法字符通过", () => {
    expect(isPasswordInvalid("abc def")).toBe(true);
    expect(isPasswordInvalid("abc-def")).toBe(true);
    expect(isPasswordInvalid("密码abc123")).toBe(true);
    expect(isPasswordInvalid("abc_123@$")).toBe(false);
  });
});

describe("昵称校验 (按 Unicode 码点计数 2~20)", () => {
  it("emoji 按码点而非 UTF-16 单元计数", () => {
    // string.length 会把 emoji 计为 2, 此处应按 1 个字符
    expect(isNicknameTooShort("😀")).toBe(true); // 1 个码点 → 过短
    expect(isNicknameTooShort("😀😀")).toBe(false); // 2 个码点 → 合法
    expect(isNicknameTooLong("😀".repeat(20))).toBe(false);
    expect(isNicknameTooLong("😀".repeat(21))).toBe(true);
  });

  it("中文与英文边界", () => {
    expect(isNicknameTooShort("好")).toBe(true);
    expect(isNicknameTooShort("你好")).toBe(false);
    expect(isNicknameTooLong("a".repeat(20))).toBe(false);
    expect(isNicknameTooLong("a".repeat(21))).toBe(true);
  });
});

describe("nonStringField", () => {
  it("数字/布尔/数组字段被识别为类型错误", () => {
    expect(nonStringField({ nickname: 123 }, ["nickname"])).toBe("nickname");
    expect(nonStringField({ nickname: true }, ["nickname"])).toBe("nickname");
    expect(nonStringField({ nickname: ["a"] }, ["nickname"])).toBe("nickname");
  });

  it("字符串/null/缺字段不算类型错误", () => {
    expect(nonStringField({ nickname: "ok" }, ["nickname"])).toBeNull();
    expect(nonStringField({ nickname: null }, ["nickname"])).toBeNull();
    expect(nonStringField({}, ["nickname"])).toBeNull();
  });
});

describe("normalizeNulls", () => {
  it("null 归一化为空串, 其余值不变", () => {
    const body = { nickname: null, avatar_url: "https://x.com/a.png" };
    normalizeNulls(body, ["nickname", "avatar_url"]);
    expect(body.nickname).toBe("");
    expect(body.avatar_url).toBe("https://x.com/a.png");
  });
});

describe("bodyError (自描述联合类型契约)", () => {
  it("缺 body (undefined) → { ok:false, reason:'missing' }", () => {
    expect(bodyError(undefined, ["nickname"])).toEqual({ ok: false, reason: "missing" });
  });

  it("顶层非对象 (null/数组/标量) → 'top-level'", () => {
    expect(bodyError(null, ["nickname"])).toEqual({ ok: false, reason: "top-level" });
    expect(bodyError([1, 2], ["nickname"])).toEqual({ ok: false, reason: "top-level" });
    expect(bodyError("hello", ["nickname"])).toEqual({ ok: false, reason: "top-level" });
    expect(bodyError(42, ["nickname"])).toEqual({ ok: false, reason: "top-level" });
  });

  it("字段类型非字符串 → 返回该字段名", () => {
    expect(bodyError({ nickname: 123, email: "a@b.com" }, ["nickname", "email"])).toEqual({
      ok: false,
      reason: "nickname",
    });
  });

  it("合法 body → { ok:true }, 且 null 已就地归一化", () => {
    const body = { nickname: null, email: "a@b.com" } as Record<string, unknown>;
    const result = bodyError(body, ["nickname", "email"]);
    expect(result).toEqual({ ok: true });
    // 归一化副作用: null 字段已变为空串
    expect(body.nickname).toBe("");
  });
});
