// Auth 鉴权相关路由
// 处理用户注册、登录、Token 刷新

import { Elysia, t } from "elysia";
import bcrypt from "bcryptjs";
import { randomUUID } from "crypto";
import { prisma } from "../lib/prisma";
import { Prisma } from "../generated/prisma/client";
import { generateToken } from "../middleware/jwt";
import { requireAuth } from "../middleware/authGuard";
import { isEmail, isPasswordTooShort, isPasswordInvalid, isNicknameTooShort, isNicknameTooLong } from "../util/validate";
import { info, debug, warn, error } from "../util/logger";
import { toUserInfo } from "../util/common";
import { fail } from "../util/response";
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
        return fail(set, 400, "昵称长度不能小于2位");
      }
      if (isNicknameTooLong(nickname)) {
        warn("auth", "注册-昵称过长", { nickname });
        return fail(set, 400, "昵称长度不能超过20位");
      }
      if (!isEmail(email)) {
        warn("auth", "注册-邮箱格式错误", { email });
        return fail(set, 400, "邮箱格式错误");
      }
      if (isPasswordTooShort(password)) {
        warn("auth", "注册-密码过短");
        return fail(set, 400, "密码长度不能小于6位");
      }
      if (isPasswordInvalid(password)) {
        warn("auth", "注册-密码字符非法");
        return fail(set, 400, "密码格式错误，只能包含字母 数字 _ @ $");
      }

      // 先检查邮箱是否已注册，避免无效的 bcrypt 计算
      debug("auth", "注册-开始检查邮箱是否已存在", { email });
      const existUser = await prisma.user.findFirst({
        where: { email },
        select: { id: true },
      });
      if (existUser) {
        warn("auth", "注册-邮箱已存在", { email });
        return fail(set, 400, "邮箱已存在");
      }

      // bcrypt 加密密码
      debug("auth", "注册-开始加密密码");
      const hashedPassword = await bcrypt.hash(password, 10);
      debug("auth", "注册-密码加密完成");

      // 事务创建用户（同一事务中创建用户、扩展信息、第三方关联），任何错误自动回滚
      // 成功分支直接在 try 内 return，避免 let + 类型收窄技巧
      debug("auth", "注册-开始事务创建用户");
      try {
        const user = await prisma.$transaction(async (tx) => {
          const newUser = await tx.user.create({
            data: {
              typeId: 1,
              // 账号使用随机 UUID 保证全局唯一（当前为内部标识，后续可按需改为可读的展示账号）
              account: randomUUID(),
              email,
              password: hashedPassword,
              nickname,
              states: UserState.Normal,
            },
            omit: { password: true },
          });
          debug("auth", "注册-事务用户记录已创建", { uid: newUser.id });

          await tx.userExtend.create({
            data: { uid: newUser.id },
          });
          debug("auth", "注册-事务用户扩展已创建", { uid: newUser.id });

          await tx.userThirdParty.create({
            data: { uid: newUser.id },
          });
          debug("auth", "注册-事务第三方记录已创建", { uid: newUser.id });

          return newUser;
        });

        debug("auth", "注册-开始生成Token", { uid: user.id });
        const token = await generateToken(user.id);
        debug("auth", "注册-Token生成成功");
        const userInfo = toUserInfo(user);

        info("auth", "注册成功", { uid: user.id, email });
        const response: AuthTokenResponse = { token, user: userInfo };
        return response;
      } catch (err) {
        // 并发注册时唯一约束兜底（email/account 为 @unique，冲突码 P2002）
        if (err instanceof Prisma.PrismaClientKnownRequestError && err.code === "P2002") {
          warn("auth", "注册-邮箱已存在(唯一约束冲突)", { email });
          return fail(set, 400, "邮箱已存在");
        }
        error("auth", "注册-创建用户失败", { email, error: String(err) });
        return fail(set, 500, "注册失败，请稍后再试");
      }
    },
    {
      // Elysia schema 校验：请求参数不合法时自动返回 400（见 index.ts 全局 VALIDATION 处理）
      body: t.Object({
        nickname: t.String(),
        email: t.String(),
        password: t.String(),
      }),
      // 速率限制：注册会执行 bcrypt（CPU 密集）并写 3 张表，
      // 防止脚本刷注册造成 CPU/存储滥用（与 login 同规格）
      beforeHandle: rateLimit(5, 60_000),
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
        return fail(set, 400, "邮箱格式错误");
      }
      if (isPasswordTooShort(password)) {
        warn("auth", "登录-密码过短");
        return fail(set, 400, "密码长度不能小于6位");
      }
      if (isPasswordInvalid(password)) {
        warn("auth", "登录-密码字符非法");
        return fail(set, 400, "密码格式错误，只能包含字母 数字 _ @ $");
      }

      // 一次查询取回全部字段（含密码），验证通过后直接映射响应，避免二次查询
      debug("auth", "登录-开始查询用户", { email });
      const user = await prisma.user.findFirst({
        where: { email },
      });
      if (!user) {
        warn("auth", "登录-用户不存在", { email });
        return fail(set, 400, "邮箱或密码错误");
      }

      // 检查用户状态（禁用用户返回 401）
      if (user.states !== UserState.Normal) {
        warn("auth", "登录-账号已被禁用", { uid: user.id, email });
        return fail(set, 401, "登录失败，账号已被禁用");
      }

      // 验证密码
      debug("auth", "登录-开始验证密码", { uid: user.id });
      const isPasswordValid = await bcrypt.compare(password, user.password);
      if (!isPasswordValid) {
        warn("auth", "登录-密码错误", { uid: user.id, email });
        return fail(set, 400, "邮箱或密码错误");
      }
      debug("auth", "登录-密码验证通过", { uid: user.id });

      const token = await generateToken(user.id);
      debug("auth", "登录-Token生成成功", { uid: user.id });
      const userInfo = toUserInfo(user);

      info("auth", "登录成功", { uid: user.id, email });
      const response: AuthTokenResponse = { token, user: userInfo };
      return response;
    },
    {
      // Elysia schema 校验：请求参数不合法时自动返回 400（见 index.ts 全局 VALIDATION 处理）
      body: t.Object({
        email: t.String(),
        password: t.String(),
      }),
      // 速率限制：每分钟最多 5 次尝试，防暴力破解
      beforeHandle: rateLimit(5, 60_000),
    }
  )
  // 以下路由需要 JWT 鉴权（由 requireAuth 插件注入 user）
  .use(requireAuth)
  // POST /api/v1/auth/refresh
  .post("/refresh", async ({ user }) => {
    info("auth", "POST /refresh 请求", { uid: user.id });

    // JWT 校验、用户查询、状态检查已由 requireAuth 插件完成
    debug("auth", "刷新Token-开始生成新Token", { uid: user.id });
    const token = await generateToken(user.id);
    const userInfo = toUserInfo(user);

    info("auth", "Token刷新成功", { uid: user.id });
    const response: AuthTokenResponse = { token, user: userInfo };
    return response;
  });

export default authRoutes;
