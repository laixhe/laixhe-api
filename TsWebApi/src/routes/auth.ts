// Auth 鉴权相关路由
// 处理用户注册、登录、Token 刷新

import { Elysia, t } from "elysia";
import bcrypt from "bcryptjs";
import { randomUUID } from "crypto";
import { prisma } from "../lib/prisma";
import { generateToken, getJwtClaimsFromHeaders } from "../middleware/jwt";
import { isEmail, isPasswordTooShort, isPasswordInvalid, isNicknameTooShort, isNicknameTooLong } from "../util/validate";
import { info, warn, error } from "../util/logger";
import { toUserInfo } from "../util/common";
import { rateLimit } from "../middleware/rateLimit";
import type {
  AuthTokenResponse,
} from "../entity/auth";
import { UserState } from "../entity/user";

export const authRoutes = new Elysia({ prefix: "/api/v1/auth" })
  // POST /api/v1/auth/register
  .post(
    "/register",
    async ({ body, set }) => {
      const { nickname, email, password } = body;
      info("auth", "POST /register 请求", { email, nickname });

      // 参数校验
      if (isNicknameTooShort(nickname)) {
        warn("auth", "注册-昵称过短", { nickname });
        set.status = 400;
        return { code: 400, message: "昵称长度不能小于2位" };
      }
      if (isNicknameTooLong(nickname)) {
        warn("auth", "注册-昵称过长", { nickname });
        set.status = 400;
        return { code: 400, message: "昵称长度不能超过20位" };
      }
      if (!isEmail(email)) {
        warn("auth", "注册-邮箱格式错误", { email });
        set.status = 400;
        return { code: 400, message: "邮箱格式错误" };
      }
      if (isPasswordTooShort(password)) {
        warn("auth", "注册-密码过短");
        set.status = 400;
        return { code: 400, message: "密码长度不能小于6位" };
      }
      if (isPasswordInvalid(password)) {
        warn("auth", "注册-密码字符非法");
        set.status = 400;
        return { code: 400, message: "密码格式错误，只能包含字母 数字 _ @ $" };
      }

      // 先检查邮箱是否已注册，避免无效的 bcrypt 计算
      info("auth", "注册-开始检查邮箱是否已存在", { email });
      const existUser = await prisma.user.findFirst({
        where: { email },
        select: { id: true },
      });
      if (existUser) {
        warn("auth", "注册-邮箱已存在", { email });
        set.status = 400;
        return { code: 400, message: "邮箱已存在" };
      }

      // bcrypt 加密密码
      info("auth", "注册-开始加密密码");
      const hashedPassword = await bcrypt.hash(password, 10);
      info("auth", "注册-密码加密完成");

      // 事务创建用户（在同一事务中创建用户、扩展信息、第三方关联）
      // 事务中任何错误都会自动回滚
      info("auth", "注册-开始事务创建用户");
      let user;
      try {
        user = await prisma.$transaction(async (tx) => {
          // INSERT INTO `user` (...)
          const newUser = await tx.user.create({
            data: {
              typeId: 1,
              account: randomUUID(), // 生成全局唯一账号
              email,
              password: hashedPassword,
              nickname,
              states: 1,
            },
            omit: { password: true },
          });
          info("auth", "注册-事务用户记录已创建", { uid: newUser.id });

          // INSERT INTO `user_extend` (uid) VALUES (?)
          await tx.userExtend.create({
            data: { uid: newUser.id },
          });
          info("auth", "注册-事务用户扩展已创建", { uid: newUser.id });

          // INSERT INTO `user_third_party` (uid) VALUES (?)
          await tx.userThirdParty.create({
            data: { uid: newUser.id },
          });
          info("auth", "注册-事务第三方记录已创建", { uid: newUser.id });

          return newUser;
        });
      } catch (err) {
        error("auth", "注册-事务创建用户失败", { email, error: String(err) });
        set.status = 500;
        return { code: 500, message: "注册失败，请稍后再试" };
      }

      info("auth", "注册-开始生成Token", { uid: user.id });
      const token = await generateToken(user.id);
      info("auth", "注册-Token生成成功");
      const userInfo = toUserInfo(user);

      info("auth", "注册成功", { uid: user.id, email });
      const response: AuthTokenResponse = { token, user: userInfo };
      return response;
    },
    {
      body: t.Object({
        nickname: t.String(),
        email: t.String(),
        password: t.String(),
      }),
    }
  )
  // POST /api/v1/auth/login
  .post(
    "/login",
    async ({ body, set }) => {
      const { email, password } = body;
      info("auth", "POST /login 请求", { email });

      // 参数校验
      if (!isEmail(email)) {
        warn("auth", "登录-邮箱格式错误", { email });
        set.status = 400;
        return { code: 400, message: "邮箱格式错误" };
      }
      if (isPasswordTooShort(password)) {
        warn("auth", "登录-密码过短");
        set.status = 400;
        return { code: 400, message: "密码长度不能小于6位" };
      }
      if (isPasswordInvalid(password)) {
        warn("auth", "登录-密码字符非法");
        set.status = 400;
        return { code: 400, message: "密码格式错误，只能包含字母 数字 _ @ $" };
      }

      // 查询用户
      info("auth", "登录-开始查询用户", { email });
      const user = await prisma.user.findFirst({
        where: { email },
      });
      if (!user) {
        warn("auth", "登录-用户不存在", { email });
        set.status = 400;
        return { code: 400, message: "邮箱或密码错误" };
      }
      info("auth", "登录-用户查询成功", { uid: user.id });

      // 检查用户状态（禁用用户返回 401）
      if (user.states !== UserState.Normal) {
        warn("auth", "登录-账号已被禁用", { uid: user.id, email });
        set.status = 401;
        return { code: 401, message: "登录失败，账号已被禁用" };
      }

      // 验证密码
      info("auth", "登录-开始验证密码", { uid: user.id });
      const isPasswordValid = await bcrypt.compare(password, user.password);
      if (!isPasswordValid) {
        warn("auth", "登录-密码错误", { uid: user.id, email });
        set.status = 400;
        return { code: 400, message: "邮箱或密码错误" };
      }
      info("auth", "登录-密码验证通过", { uid: user.id });

      const token = await generateToken(user.id);
      info("auth", "登录-Token生成成功", { uid: user.id });
      const userInfo = toUserInfo(user);

      info("auth", "登录成功", { uid: user.id, email });
      const response: AuthTokenResponse = { token, user: userInfo };
      return response;
    },
    {
      body: t.Object({
        email: t.String(),
        password: t.String(),
      }),
      // 速率限制：每分钟最多 5 次尝试，防暴力破解
      beforeHandle: rateLimit(5, 60_000),
    }
  )
  // POST /api/v1/auth/refresh
  .post("/refresh", async ({ headers, set }) => {
    info("auth", "POST /refresh 请求");
    const claims = await getJwtClaimsFromHeaders(headers);
    if (!claims) {
      warn("auth", "刷新Token-JWT无效或缺失");
      set.status = 401;
      return { code: 401, message: "未授权" };
    }

    // 查询用户
    info("auth", "刷新Token-开始查询用户", { uid: claims.uid });
    const user = await prisma.user.findUnique({
      where: { id: claims.uid },
      omit: { password: true },
    });
    if (!user) {
      warn("auth", "刷新Token-用户不存在", { uid: claims.uid });
      set.status = 401;
      return { code: 401, message: "用户不存在" };
    }
    info("auth", "刷新Token-用户查询成功", { uid: user.id });

    // 检查用户状态（禁用用户返回 401）
    if (user.states !== UserState.Normal) {
      warn("auth", "刷新Token-账号已被禁用", { uid: user.id });
      set.status = 401;
      return { code: 401, message: "登录失败，账号已被禁用" };
    }

    info("auth", "刷新Token-开始生成新Token", { uid: user.id });
    const token = await generateToken(user.id);
    const userInfo = toUserInfo(user);

    info("auth", "Token刷新成功", { uid: user.id });
    const response: AuthTokenResponse = { token, user: userInfo };
    return response;
  });

export default authRoutes;
