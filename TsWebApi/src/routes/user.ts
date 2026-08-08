// User 用户相关路由
// 处理用户信息查询、列表、信息更新
// info 和 list 为公开接口，update 需要 JWT 鉴权

import { Elysia, t } from "elysia";
import { prisma } from "../lib/prisma";
import { requireAuth } from "../middleware/authGuard";
import { isNicknameTooShort, isNicknameTooLong } from "../util/validate";
import { info, debug, warn, error } from "../util/logger";
import { toUserInfo } from "../util/common";
import { fail } from "../util/response";
import { rateLimit } from "../middleware/rateLimit";

export const userRoutes = new Elysia({ prefix: "/api/v1/user" })
  // GET /api/v1/user/info?uid=xxx（公开接口）
  .get("/info", async ({ query, set }) => {
    const uid = parseInt(query.uid || "0", 10);
    info("user", "GET /info 请求", { uid });
    if (!uid || uid <= 0) {
      warn("user", "查询用户信息-无效uid", { uid: query.uid });
      return fail(set, 400, "无效的用户ID");
    }

    debug("user", "查询用户信息-开始查询", { uid });
    const user = await prisma.user.findUnique({
      where: { id: uid },
      omit: { password: true },
    });
    if (!user) {
      warn("user", "查询用户信息-用户不存在", { uid });
      return fail(set, 400, "用户不存在");
    }

    info("user", "查询用户信息成功", { uid });
    return toUserInfo(user);
  }, {
    query: t.Object({ uid: t.String() }),
    // 公开接口防刷：宽松限流（60 次/分钟）
    beforeHandle: rateLimit(60, 60_000),
  })
  // GET /api/v1/user/list?page=1&page_size=12（公开接口）
  .get("/list", async ({ query }) => {
    // 分页参数钳制：page 最小 1；page_size 限制在 [1, MAX_PAGE_SIZE]，
    // 防止恶意大分页触发全表级 LIMIT/OFFSET 查询
    const MAX_PAGE_SIZE = 100;
    let page = parseInt(query.page || "1", 10);
    page = Number.isFinite(page) ? Math.max(page, 1) : 1;
    let pageSize = parseInt(query.page_size || "12", 10);
    pageSize = Number.isFinite(pageSize) ? Math.min(Math.max(pageSize, 1), MAX_PAGE_SIZE) : 12;
    const offset = (page - 1) * pageSize;
    info("user", "GET /list 请求", { page, pageSize });

    // SELECT count(*) FROM `user` + SELECT * FROM `user` ORDER BY `id` DESC LIMIT ? OFFSET ?
    debug("user", "查询用户列表-开始查询", { page, pageSize, offset });
    const [users, total] = await Promise.all([
      prisma.user.findMany({
        skip: offset,
        take: pageSize,
        orderBy: { id: "desc" },
        omit: { password: true },
      }),
      prisma.user.count(),
    ]);

    info("user", "查询用户列表", { page, pageSize, total, count: users.length });
    return {
      total,
      page,
      page_size: pageSize,
      list: users.map(toUserInfo),
    };
  }, {
    query: t.Object({
      page: t.Optional(t.String()),
      page_size: t.Optional(t.String()),
    }),
    // 公开接口防刷：宽松限流（60 次/分钟）
    beforeHandle: rateLimit(60, 60_000),
  })
  // 以下路由需要 JWT 鉴权（由 requireAuth 插件注入 user，无需重复查询用户）
  .use(requireAuth)
  // POST /api/v1/user/update
  .post(
    "/update",
    async ({ body, set, user }) => {
      const { nickname, avatar_url } = body;
      info("user", "POST /update 请求", { uid: user.id, nickname });

      // 验证昵称格式
      if (isNicknameTooShort(nickname)) {
        warn("user", "更新用户信息-昵称过短", { uid: user.id, nickname });
        return fail(set, 400, "昵称长度不能小于2位");
      }
      if (isNicknameTooLong(nickname)) {
        warn("user", "更新用户信息-昵称过长", { uid: user.id, nickname });
        return fail(set, 400, "昵称长度不能超过20位");
      }
      // 验证头像地址格式
      if (avatar_url.length > 255) {
        warn("user", "更新用户信息-头像地址过长", { uid: user.id });
        return fail(set, 400, "头像地址长度不能超过255位");
      }
      if (avatar_url.length > 0 && !avatar_url.startsWith("http")) {
        warn("user", "更新用户信息-头像地址格式错误", { uid: user.id, avatar_url });
        return fail(set, 400, "头像地址必须以http或https开头");
      }

      // 更新用户（仅非空字段更新，avatar_url 为空时保留旧值，与 Go 行为一致）
      // 成功分支直接在 try 内 return，避免 let + 类型收窄技巧
      try {
        const updateData: { nickname: string; avatarUrl?: string } = { nickname };
        if (avatar_url !== "") {
          updateData.avatarUrl = avatar_url;
        }
        debug("user", "更新用户信息-开始更新数据库", { uid: user.id });
        const updatedUser = await prisma.user.update({
          where: { id: user.id },
          data: updateData,
        });
        info("user", "更新用户信息成功", { uid: user.id, nickname });
        return toUserInfo(updatedUser);
      } catch (err) {
        error("user", "更新用户信息-数据库更新失败", { uid: user.id, error: String(err) });
        return fail(set, 500, "更新失败，请稍后再试");
      }
    },
    {
      // Elysia schema 校验：请求参数不合法时自动返回 400（见 index.ts 全局 VALIDATION 处理）
      body: t.Object({
        nickname: t.String(),
        avatar_url: t.String(),
      }),
    }
  );

export default userRoutes;
