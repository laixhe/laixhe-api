// 统一响应工具
// 消除各处重复的 `set.status = xxx; return { code, message }` 样板代码

// set 的最小结构（Elysia 的 SetObject 兼容此结构：status 为 number 或状态码字符串，可选）
interface StatusSetter {
  status?: number | string;
}

/**
 * 失败响应：设置 HTTP 状态码并返回统一错误结构 { code, message }
 * @param set Elysia 的 set 对象
 * @param httpCode HTTP 状态码（同时作为业务 code）
 * @param message 错误提示
 */
export function fail(set: StatusSetter, httpCode: number, message: string) {
  set.status = httpCode;
  return { code: httpCode, message };
}
