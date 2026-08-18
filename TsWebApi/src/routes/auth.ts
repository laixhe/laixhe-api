// Auth 鉴权相关路由
// 处理用户注册、登录、Token 刷新

import { Elysia } from "elysia";
import { randomUUID } from "crypto";
import { prisma } from "../lib/prisma";
import { Prisma } from "../generated/prisma/client";
import { generateToken } from "../middleware/jwt";
import { requireAuth } from "../middleware/authGuard";
import { isEmail, isPasswordTooShort, isPasswordTooLong, isPasswordInvalid, isNicknameTooShort, isNicknameTooLong, bodyError } from "../util/validate";
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
      // 参数校验为手动完成 (与 Go/PHP/Rust 端一致), 此处仅做类型声明
      // 运行时类型校验 (missing → 422; 顶层/字段类型 → 400): 见 bodyError
      const bodyErr = bodyError(body, ["nickname", "email", "password"]);
      if (!bodyErr.ok) {
        if (bodyErr.reason === "missing") {
          // 请求体整体缺失 → 422 (与 Go/PHP/Rust 端一致), 显式处理而非依赖全局 TypeError 兜底
          return fail(set, 422, "参数错误");
        }
        warn("auth", "注册-body类型错误", { reason: bodyErr.reason });
        return fail(set, 400, "Bad Request");
      }
      const { nickname, email, password } = body as { nickname: string; email: string; password: string };
      info("auth", "POST /register 请求", { email, nickname });

      // 参数校验 (缺失字段同样命中以下规则, 与 Go/PHP/Rust 端 422 文案一致)
      if (isNicknameTooShort(nickname)) {
        warn("auth", "注册-昵称过短", { nickname });
        return fail(set, 422, "昵称长度不能小于2位");
      }
      if (isNicknameTooLong(nickname)) {
        warn("auth", "注册-昵称过长", { nickname });
        return fail(set, 422, "昵称长度不能超过20位");
      }
      if (!isEmail(email)) {
        warn("auth", "注册-邮箱格式错误", { email });
        return fail(set, 422, "邮箱格式错误");
      }
      if (isPasswordTooShort(password)) {
        warn("auth", "注册-密码过短");
        return fail(set, 422, "密码长度不能小于6位");
      }
      if (isPasswordTooLong(password)) {
        warn("auth", "注册-密码过长");
        return fail(set, 422, "密码长度不能超过64位");
      }
      if (isPasswordInvalid(password)) {
        warn("auth", "注册-密码字符非法");
        return fail(set, 422, "密码格式错误，只能包含字母 数字 _ @ $");
      }

      // 先检查邮箱是否已注册，避免无效的 bcrypt 计算
      debug("auth", "注册-开始检查邮箱是否已存在", { email });
      const existUser = await prisma.user.findFirst({
        where: { email },
        select: { id: true },
      });
      if (existUser) {
        warn("auth", "注册-邮箱已存在", { email });
        return fail(set, 422, "邮箱已存在");
      }

      // bcrypt 加密密码
      // 用 Bun.password (原生实现, 在独立线程计算, 不阻塞 JS 主线程) 而非 bcryptjs:
      // bcryptjs 的异步 API 只是主线程上的协作式分片, CPU 计算仍占用主线程, 并发下会成为热点
      debug("auth", "注册-开始加密密码");
      const hashedPassword = await Bun.password.hash(password, {
        algorithm: "bcrypt",
        cost: 10,
      });
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
        // email 为唯一索引, 该兜底对 account/email/关联表 uid 的冲突均生效
        // (并发注册同邮箱等极端情况才触发, 正常重复邮箱已被上面的先查后插拦截)
        if (err instanceof Prisma.PrismaClientKnownRequestError && err.code === "P2002") {
          warn("auth", "注册-唯一键冲突", { email });
          return fail(set, 422, "注册失败，请稍后再试");
        }
        error("auth", "注册-创建用户失败", { email, error: String(err) });
        return fail(set, 500, "注册失败，请稍后再试");
      }
    },
    {
      // 速率限制：注册会执行 bcrypt（CPU 密集）并写 3 张表，
      // 防止脚本刷注册造成 CPU/存储滥用（与 login 同规格）
      beforeHandle: rateLimit(5, 60_000, "auth:register"),
      detail: {
        tags: ["Auth"],
        summary: "注册",
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: { $ref: "#/components/schemas/AuthRegisterRequest" },
            },
          },
        },
        responses: {
          "200": {
            description: "OK",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/AuthTokenResponse" },
              },
            },
          },
          "400": {
            description: "请求格式错误",
            content: {
              "application/json": { schema: { $ref: "#/components/schemas/Error" } },
            },
          },
          "422": {
            description: "参数错误",
            content: {
              "application/json": { schema: { $ref: "#/components/schemas/Error" } },
            },
          },
          "500": {
            description: "Internal Server Error",
            content: {
              "application/json": { schema: { $ref: "#/components/schemas/Error" } },
            },
          },
        },
      },
    }
  )
  // POST /api/v1/auth/login
  .post(
    "/login",
    async ({ body, set }) => {
      // 运行时类型校验 (missing → 422; 顶层/字段类型 → 400): 见 bodyError
      const bodyErr = bodyError(body, ["email", "password"]);
      if (!bodyErr.ok) {
        if (bodyErr.reason === "missing") {
          // 请求体整体缺失 → 422 (与 Go/PHP/Rust 端一致), 显式处理而非依赖全局 TypeError 兜底
          return fail(set, 422, "参数错误");
        }
        warn("auth", "登录-body类型错误", { reason: bodyErr.reason });
        return fail(set, 400, "Bad Request");
      }
      const { email, password } = body as { email: string; password: string };
      info("auth", "POST /login 请求", { email });

      // 参数校验 (缺失字段同样命中以下规则, 与 Go/PHP/Rust 端 422 文案一致)
      if (!isEmail(email)) {
        warn("auth", "登录-邮箱格式错误", { email });
        return fail(set, 422, "邮箱格式错误");
      }
      if (isPasswordTooShort(password)) {
        warn("auth", "登录-密码过短");
        return fail(set, 422, "密码长度不能小于6位");
      }
      if (isPasswordTooLong(password)) {
        warn("auth", "登录-密码过长");
        return fail(set, 422, "密码长度不能超过64位");
      }
      if (isPasswordInvalid(password)) {
        warn("auth", "登录-密码字符非法");
        return fail(set, 422, "密码格式错误，只能包含字母 数字 _ @ $");
      }

      // 一次查询取回全部字段（含密码），验证通过后直接映射响应，避免二次查询
      debug("auth", "登录-开始查询用户", { email });
      const user = await prisma.user.findFirst({
        where: { email },
      });
      if (!user) {
        warn("auth", "登录-用户不存在", { email });
        return fail(set, 422, "邮箱或密码错误");
      }

      // 封禁账号与密码错误返回同一提示, 避免暴露账号状态 (可被探测) (与 Go/PHP/Rust 端一致)
      if (user.states !== UserState.Normal) {
        warn("auth", "登录-账号已被禁用", { uid: user.id, email });
        return fail(set, 422, "邮箱或密码错误");
      }

      // 验证密码 (Bun.password 原生实现, 不阻塞主线程, 理由同注册处注释)
      debug("auth", "登录-开始验证密码", { uid: user.id });
      const isPasswordValid = await Bun.password.verify(password, user.password);
      if (!isPasswordValid) {
        warn("auth", "登录-密码错误", { uid: user.id, email });
        return fail(set, 422, "邮箱或密码错误");
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
      // 速率限制：每分钟最多 5 次尝试，防暴力破解
      beforeHandle: rateLimit(5, 60_000, "auth:login"),
      detail: {
        tags: ["Auth"],
        summary: "登录",
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: { $ref: "#/components/schemas/AuthLoginRequest" },
            },
          },
        },
        responses: {
          "200": {
            description: "OK",
            content: {
              "application/json": {
                schema: { $ref: "#/components/schemas/AuthTokenResponse" },
              },
            },
          },
          "400": {
            description: "请求格式错误",
            content: {
              "application/json": { schema: { $ref: "#/components/schemas/Error" } },
            },
          },
          "422": {
            description: "参数错误",
            content: {
              "application/json": { schema: { $ref: "#/components/schemas/Error" } },
            },
          },
          "500": {
            description: "Internal Server Error",
            content: {
              "application/json": { schema: { $ref: "#/components/schemas/Error" } },
            },
          },
        },
      },
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
  }, {
    detail: {
      tags: ["Auth"],
      summary: "刷新Jwt",
      security: [{ BearerAuth: [] }],
      responses: {
        "200": {
          description: "OK",
          content: {
            "application/json": {
              schema: { $ref: "#/components/schemas/AuthTokenResponse" },
            },
          },
        },
        "400": {
          description: "请求格式错误",
          content: {
            "application/json": { schema: { $ref: "#/components/schemas/Error" } },
          },
        },
        "401": {
          description: "未授权",
          content: {
            "application/json": { schema: { $ref: "#/components/schemas/Error" } },
          },
        },
        "500": {
          description: "Internal Server Error",
          content: {
            "application/json": { schema: { $ref: "#/components/schemas/Error" } },
          },
        },
      },
    },
  });

export default authRoutes;
