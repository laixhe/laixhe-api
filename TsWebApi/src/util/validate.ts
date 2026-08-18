// 邮箱格式正则
const emailRe = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

export function isEmail(email: string): boolean {
  return emailRe.test(email);
}

// 密码长度校验（>= 6 位）
// 注意: 参数可能为 undefined (body 缺失字段), 用 ?? "" 兜底避免 .length 抛 TypeError,
// 使缺字段走与 Go/PHP/Rust 端一致的"具体 422 文案", 而非依赖全局 onError 的隐式兜底
export function isPasswordTooShort(password: string | undefined): boolean {
  return (password ?? "").length < 6;
}

// 密码长度校验（> 64 位）
// 上限 64: bcrypt 只取前 72 字节, 上限保证任何合法密码都不会被静默截断产生碰撞
export function isPasswordTooLong(password: string | undefined): boolean {
  return (password ?? "").length > 64;
}

// 密码字符正则：仅允许字母、数字、_、@、$，长度 6~64 位
const passwordRe = /^[a-zA-Z0-9_@$]{6,64}$/;

// 说明：下面密码校验函数在"长度 < 6"这一规则上重叠
// （passwordRe 自带 {6,} 长度限制），分开定义只是为了：
//   - isPasswordTooShort：单独命中"过短"这一种情况，给出专属提示
//   - isPasswordInvalid：命中字符集非法（含过短），提示覆盖完整规则
// 调用方先查前者、再查后者，即可返回精确的错误提示。

// 密码字符校验（仅允许字母、数字、_、@、$，且长度 6~64 位）
export function isPasswordInvalid(password: string | undefined): boolean {
  return !passwordRe.test(password ?? "");
}

// 昵称过短校验（< 2 字）
// 注意: 用 [...nickname] 按 Unicode 码点计数 (与 Go RuneCountInString / Rust chars().count()
// / PHP mb_strlen 一致), 而非 string.length 的 UTF-16 单元数 — 否则 emoji 等代理对字符
// (如 😀 占 2 个 UTF-16 单元) 会被误判为 2 位
export function isNicknameTooShort(nickname: string | undefined): boolean {
  return [...(nickname ?? "")].length < 2;
}

// 昵称过长校验（> 20 字）
export function isNicknameTooLong(nickname: string | undefined): boolean {
  return [...(nickname ?? "")].length > 20;
}

// body 字段类型校验：给定字段（若存在）必须是字符串。
// 用途：防御类型断言（`body as {...}`）掩盖的错误类型输入——非字符串字段（如数字/布尔/数组）
// 会导致后续 Prisma/逻辑层报 500，而 Go/Rust 绑定层对此类输入统一返回 400。
// 注意：null 视为"无值"而非类型错误（与 Go/PHP 端一致），不在此拦截，由业务校验走 422。
// @returns 第一个类型非法的字段名；全部合法（或字段缺失/null）时返回 null
export function nonStringField(
  body: Record<string, unknown>,
  fields: string[]
): string | null {
  for (const field of fields) {
    const value = body[field];
    if (value !== undefined && value !== null && typeof value !== "string") {
      return field;
    }
  }
  return null;
}

// 将 body 中为 null 的字段归一化为空字符串（与 Go/PHP 端一致: null 视为"无值"）。
// 归一化后走既有业务校验，返回具体的 422 文案（如"昵称长度不能小于2位"）而非 400。
export function normalizeNulls(
  body: Record<string, unknown>,
  fields: string[]
): void {
  for (const field of fields) {
    if (body[field] === null) {
      body[field] = "";
    }
  }
}

// body 类型校验 + null 归一化的组合入口（供各 handler 复用, 消除样板重复）。
// 返回自描述联合类型 (而非魔数字符串), 让调用方的分支判断一目了然:
// - { ok: false, reason: "missing" }: 请求体整体缺失 (undefined), 调用方返回 422 "参数错误";
// - { ok: false, reason: "top-level" }: 顶层 body 非纯对象 (数组/标量/null), 调用方返回 400;
// - { ok: false, reason: 字段名 }: 该字段类型非字符串, 调用方返回 400;
// - { ok: true }: 校验通过 (body 已就地为 null 字段归一化为空串)。
// 注意:
// 1. body 存在但缺字段 (如 {}) 时, 校验函数已对 undefined 做空串兜底 (见 isNicknameTooShort 等),
//    会走具体 422 文案 (与 Go/PHP/Rust 端一致);
// 2. fields 白名单必须与 handler 解构/使用的字段保持一致, 新增字段忘记加入会导致
//    该字段绕过类型校验, 错误类型输入可能漏成 500。
// 3. "missing" 由调用方显式返回 422, 不再依赖全局 onError 对 TypeError 的宽泛兜底 (见 src/index.ts)。
export type BodyErrorResult =
  | { ok: true }
  | { ok: false; reason: "missing" | "top-level" | string };

export function bodyError(body: unknown, fields: string[]): BodyErrorResult {
  if (body === undefined) return { ok: false, reason: "missing" };
  if (body === null || typeof body !== "object" || Array.isArray(body)) {
    return { ok: false, reason: "top-level" };
  }
  const record = body as Record<string, unknown>;
  const bad = nonStringField(record, fields);
  if (bad !== null) return { ok: false, reason: bad };
  normalizeNulls(record, fields);
  return { ok: true };
}
