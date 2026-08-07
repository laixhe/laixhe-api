// User 用户相关路由
// 处理用户信息查询、列表、信息更新
// info 和 list 为公开接口，update 需要 JWT 鉴权

import { Elysia, t } from "elysia";
import { prisma } from "../lib/prisma";
import { getJwtClaimsFromHeaders } from "../middleware/jwt";
import { isNicknameTooShort, isNicknameTooLong } from "../util/validate";
import { info, warn, error } from "../util/logger";
import { toUserInfo } from "../util/common";
import { UserState } from "../entity/user";

export const userRoutes = new Elysia({ prefix: "/api/v1/user" })
  // GET /api/v1/user/info?uid=xxx（公开接口）
  .get("/info", async ({ query, set }) => {
    const uid = parseInt(query.uid || "0", 10);
    info("user", "GET /info 请求", { uid });
    if (!uid || uid <= 0) {
      warn("user", "查询用户信息-无效uid", { uid: query.uid });
      set.status = 400;
      return { code: 400, message: "无效的用户ID" };
    }

    info("user", "查询用户信息-开始查询", { uid });
    const user = await prisma.user.findUnique({
      where: { id: uid },
      omit: { password: true },
    });
    if (!user) {
      warn("user", "查询用户信息-用户不存在", { uid });
      set.status = 400;
      return { code: 400, message: "用户不存在" };
    }

    info("user", "查询用户信息成功", { uid });
    return toUserInfo(user);
  }, {
    query: t.Object({ uid: t.String() }),
  })
  // GET /api/v1/user/list?page=1&page_size=12（公开接口）
  .get("/list", async ({ query, set }) => {
    const page = Math.max(parseInt(query.page || "1", 10), 1);
    const pageSize = Math.max(parseInt(query.page_size || "12", 10), 1);
    const offset = (page - 1) * pageSize;
    info("user", "GET /list 请求", { page, pageSize });

    // SELECT count(*) FROM `user` + SELECT * FROM `user` ORDER BY `id` DESC LIMIT ? OFFSET ?
    info("user", "查询用户列表-开始查询", { page, pageSize, offset });
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
  })
  // POST /api/v1/user/update
  .post(
    "/update",
    async ({ body, headers, set }) => {
      const claims = await getJwtClaimsFromHeaders(headers);
      if (!claims) {
        warn("user", "更新用户信息-JWT无效或缺失");
        set.status = 401;
        return { code: 401, message: "未授权" };
      }

      const { nickname, avatar_url } = body;
      info("user", "POST /update 请求", { uid: claims.uid, nickname });

      // 验证昵称格式
      if (isNicknameTooShort(nickname)) {
        warn("user", "更新用户信息-昵称过短", { uid: claims.uid, nickname });
        set.status = 400;
        return { code: 400, message: "昵称长度不能小于2位" };
      }
      if (isNicknameTooLong(nickname)) {
        warn("user", "更新用户信息-昵称过长", { uid: claims.uid, nickname });
        set.status = 400;
        return { code: 400, message: "昵称长度不能超过20位" };
      }
      // 验证头像地址格式
      if (avatar_url.length > 255) {
        warn("user", "更新用户信息-头像地址过长", { uid: claims.uid });
        set.status = 400;
        return { code: 400, message: "头像地址长度不能超过255位" };
      }
      if (avatar_url.length > 0 && !avatar_url.startsWith("http")) {
        warn("user", "更新用户信息-头像地址格式错误", { uid: claims.uid, avatar_url });
        set.status = 400;
        return { code: 400, message: "头像地址必须以http或https开头" };
      }

      // 查询用户
      info("user", "更新用户信息-开始查询用户", { uid: claims.uid });
      const user = await prisma.user.findUnique({
        where: { id: claims.uid },
        omit: { password: true },
      });
      if (!user) {
        warn("user", "更新用户信息-用户不存在", { uid: claims.uid });
        set.status = 400;
        return { code: 400, message: "用户不存在" };
      }
      info("user", "更新用户信息-用户查询成功", { uid: user.id });

      // 检查用户状态（禁用用户返回 401）
      if (user.states !== UserState.Normal) {
        warn("user", "更新用户信息-账号已被禁用", { uid: user.id });
        set.status = 401;
        return { code: 401, message: "登录失败，账号已被禁用" };
      }

      // 更新用户（仅非空字段更新，avatar_url 为空时保留旧值，与 Go 行为一致）
      let updatedUser;
      try {
        const updateData: { nickname: string; avatarUrl?: string } = { nickname };
        if (avatar_url !== "") {
          updateData.avatarUrl = avatar_url;
        }
        info("user", "更新用户信息-开始更新数据库", { uid: claims.uid });
        updatedUser = await prisma.user.update({
          where: { id: claims.uid },
          data: updateData,
        });
      } catch (err) {
        error("user", "更新用户信息-数据库更新失败", { uid: claims.uid, error: String(err) });
        set.status = 500;
        return { code: 500, message: "更新失败，请稍后再试" };
      }

      info("user", "更新用户信息成功", { uid: claims.uid, nickname });
      return toUserInfo(updatedUser);
    },
    {
      body: t.Object({
        nickname: t.String(),
        avatar_url: t.String(),
      }),
    }
  );

export default userRoutes;
