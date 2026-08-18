// 统一 JWT 鉴权插件
// 将"校验 Token → 查询用户 → 检查状态"抽为可复用插件，
// 需要鉴权的路由只需 .use(requireAuth)，即可在 handler 中直接使用 user

import { Elysia } from "elysia";
import { prisma } from "../lib/prisma";
import { getJwtClaimsFromHeaders } from "./jwt";
import { warn } from "../util/logger";
import { userPublicSelect } from "../util/common";
import { UserState } from "../entity/user";

// 鉴权失败统一异常，由全局 onError 转换为 401 响应 (文案与 Go/PHP/Rust 版一致为 "Unauthorized")
export class UnauthorizedError extends Error {
  constructor(message = "Unauthorized") {
    super(message);
    this.name = "UnauthorizedError";
  }
}

/**
 * requireAuth 鉴权插件
 * 校验 JWT → 查询用户（不含密码）→ 检查状态，通过后向 context 注入 user
 * 任一步失败抛出 UnauthorizedError，由全局 onError 统一返回 401
 * 使用 as: "scoped" 作用域：注入的 user 仅对挂载后的路由生效，
 * 且类型可经 .use() 正确传播（默认作用域在此 Elysia 版本下无法跨插件传播类型）
 */
export const requireAuth = new Elysia({ name: "requireAuth" }).derive(
  { as: "scoped" },
  async ({ headers }) => {
    // 1. 解析并校验 JWT
    // 校验失败统一 throw UnauthorizedError，由全局 onError 设置 401 并返回响应
    const claims = await getJwtClaimsFromHeaders(headers);
    if (!claims) {
      warn("auth", "鉴权-JWT无效或缺失");
      throw new UnauthorizedError();
    }

    // 2. 查询用户 (仅取响应所需字段, 见 userPublicSelect)
    const user = await prisma.user.findUnique({
      where: { id: claims.uid },
      select: userPublicSelect,
    });
    // 用户不存在 / 账号被禁用统一返回 "Unauthorized" (与 Go 端一致, 不暴露账号状态)
    if (!user) {
      warn("auth", "鉴权-用户不存在", { uid: claims.uid });
      throw new UnauthorizedError();
    }

    // 3. 检查用户状态（禁用用户返回 401）
    if (user.states !== UserState.Normal) {
      warn("auth", "鉴权-账号已被禁用", { uid: claims.uid });
      throw new UnauthorizedError();
    }

    return { user };
  }
);
