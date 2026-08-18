// User 用户相关路由
// 处理用户信息查询、列表、信息更新
// info 和 list 为公开接口，update 需要 JWT 鉴权

import { Elysia } from "elysia";
import { prisma } from "../lib/prisma";
import { requireAuth } from "../middleware/authGuard";
import { isNicknameTooShort, isNicknameTooLong, bodyError } from "../util/validate";
import { info, debug, warn, error } from "../util/logger";
import { toUserInfo, userPublicSelect } from "../util/common";
import { fail } from "../util/response";
import { rateLimit } from "../middleware/rateLimit";

// 用户总数缓存 (5s TTL, 与 health 探活缓存同理):
// count(*) 在 InnoDB 下为全表扫描, 缓存可避免高频翻页重复全表 count;
// 代价是新增用户后最多延迟 5s 反映到列表 total
// 注: 进程内缓存, 多实例部署时各实例独立计数 (与限流器同样的单实例假设)
const TOTAL_TTL = 5_000;
let totalCache: { at: number; total: number } | null = null;
let totalInFlight: Promise<number> | null = null;

async function getTotalUserCount(): Promise<number> {
  const now = Date.now();
  if (totalCache && now - totalCache.at < TOTAL_TTL) {
    return totalCache.total;
  }
  // 复用进行中的 count 请求 (single-flight), 避免 TTL 过期瞬间并发翻页同时触发全表 count
  if (totalInFlight) {
    return totalInFlight;
  }
  totalInFlight = prisma.user.count().then((total) => {
    totalCache = { at: Date.now(), total };
    return total;
  });
  try {
    return await totalInFlight;
  } finally {
    totalInFlight = null;
  }
}

export const userRoutes = new Elysia({ prefix: "/api/v1/user" })
  // GET /api/v1/user/info?uid=xxx（公开接口）
  .get("/info", async ({ query, set }) => {
    // 非数字 uid 按请求格式错误返回 400 (与 Go/Rust 端绑定层行为一致);
    // 数字但非法 (<=0) 仍走下方 422 "无效的用户ID"
    // 用 isSafeInteger: 拒绝超出 JS 安全整数范围的超大值 (Number 会丢失精度, 其他端绑定层直接 400)
    let uid: number;
    if (query.uid === undefined || query.uid === "") {
      uid = 0;
    } else {
      uid = Number(query.uid);
      if (!Number.isSafeInteger(uid)) {
        warn("user", "查询用户信息-uid非数字", { uid: query.uid });
        return fail(set, 400, "Bad Request");
      }
    }
    info("user", "GET /info 请求", { uid });
    if (!uid || uid <= 0) {
      warn("user", "查询用户信息-无效uid", { uid: query.uid });
      return fail(set, 422, "无效的用户ID");
    }

    debug("user", "查询用户信息-开始查询", { uid });
    const user = await prisma.user.findUnique({
      where: { id: uid },
      select: userPublicSelect,
    });
    if (!user) {
      warn("user", "查询用户信息-用户不存在", { uid });
      return fail(set, 422, "用户不存在");
    }

    info("user", "查询用户信息成功", { uid });
    return toUserInfo(user);
  }, {
    // 公开接口防刷：宽松限流（60 次/分钟）
    beforeHandle: rateLimit(60, 60_000, "user:info"),
    detail: {
      tags: ["User"],
      summary: "获取用户信息",
      parameters: [
        {
          name: "uid",
          in: "query",
          required: true,
          description: "用户id",
          schema: { type: "integer" },
        },
      ],
      responses: {
        "200": {
          description: "OK",
          content: {
            "application/json": { schema: { $ref: "#/components/schemas/User" } },
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
  })
  // GET /api/v1/user/list?page=1&page_size=12（公开接口）
  .get("/list", async ({ query, set }) => {
    // 分页参数钳制：page 最小 1、最大 MAX_PAGE；page_size 非正数回落默认 12、超过 MAX_PAGE_SIZE 钳制为 100，
    // 防止恶意大分页触发深 OFFSET 扫描 (其他端仅钳最小值, 此处额外设上限为教学增强)
    // 上限取值说明: 最坏 offset = (MAX_PAGE-1)*MAX_PAGE_SIZE = 9900 行, 保证偏移有界;
    // InnoDB 的 OFFSET 需扫描并丢弃前 N 行, 深翻页成本随偏移线性增长, 若要支持任意深度,
    // 生产建议改用 keyset 分页 (where: { id: { lt: 上页最小 id } }), 见 Go 端 ListUser 注释
    const MAX_PAGE_SIZE = 100;
    const MAX_PAGE = 100;
    // 非数字分页参数按请求格式错误返回 400 (与 Go/Rust 端绑定层行为一致);
    // 仅缺省/空字符串回落到默认值
    // 用 isSafeInteger: 拒绝超出 JS 安全整数范围的超大值 (Number 会丢失精度, 其他端绑定层直接 400)
    let page: number;
    if (query.page === undefined || query.page === "") {
      page = 1;
    } else {
      page = Number(query.page);
      if (!Number.isSafeInteger(page)) {
        warn("user", "查询用户列表-page非数字", { page: query.page });
        return fail(set, 400, "Bad Request");
      }
    }
    let pageSize: number;
    if (query.page_size === undefined || query.page_size === "") {
      pageSize = 12;
    } else {
      pageSize = Number(query.page_size);
      if (!Number.isSafeInteger(pageSize)) {
        warn("user", "查询用户列表-page_size非数字", { page_size: query.page_size });
        return fail(set, 400, "Bad Request");
      }
    }
    // 与 Go/Rust/PHP 端钳制语义一致: page_size <= 0 时回落默认 12,
    // 超过 MAX_PAGE_SIZE 钳制为 100; page 同时钳制上下限防深 OFFSET
    page = Math.min(Math.max(page, 1), MAX_PAGE);
    pageSize = pageSize > 0 ? Math.min(pageSize, MAX_PAGE_SIZE) : 12;
    const offset = (page - 1) * pageSize;
    info("user", "GET /list 请求", { page, pageSize });

    // SELECT count(*) FROM `user` + SELECT <userPublicSelect 列> FROM `user` ORDER BY `id` DESC LIMIT ? OFFSET ?
    // total 加 5s 短缓存 (与 health 探活缓存同理): count(*) 为全表扫描, 高频翻页时可显著降低数据库压力;
    // 代价是列表总人数最多延迟 5s 反映 (对齐 PHP 端 paginate 行为取舍)
    debug("user", "查询用户列表-开始查询", { page, pageSize, offset });
    const [users, total] = await Promise.all([
      prisma.user.findMany({
        skip: offset,
        take: pageSize,
        orderBy: { id: "desc" },
        select: userPublicSelect,
      }),
      getTotalUserCount(),
    ]);

    info("user", "查询用户列表", { page, pageSize, total, count: users.length });
    return {
      total,
      page,
      page_size: pageSize,
      list: users.map(toUserInfo),
    };
  }, {
    // 公开接口防刷：宽松限流（60 次/分钟）
    beforeHandle: rateLimit(60, 60_000, "user:list"),
    detail: {
      tags: ["User"],
      summary: "获取用户列表",
      parameters: [
        {
          name: "page",
          in: "query",
          required: false,
          description: "分页-当前页(默认 1)",
          schema: { type: "integer" },
        },
        {
          name: "page_size",
          in: "query",
          required: false,
          description: "分页-每页数量(默认 12)",
          schema: { type: "integer" },
        },
      ],
      responses: {
        "200": {
          description: "OK",
          content: {
            "application/json": {
              schema: { $ref: "#/components/schemas/UserListResponse" },
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
  })
  // 以下路由需要 JWT 鉴权（由 requireAuth 插件注入 user，无需重复查询用户）
  .use(requireAuth)
  // POST /api/v1/user/update
  .post(
    "/update",
    async ({ body, set, user }) => {
      // avatar_url 缺省为空字符串 (与 Go/PHP/Rust 端一致: 空串不更新), 此处仅做类型声明
      // 运行时类型校验 (missing → 422; 顶层/字段类型 → 400): 见 bodyError
      const bodyErr = bodyError(body, ["nickname", "avatar_url"]);
      if (!bodyErr.ok) {
        if (bodyErr.reason === "missing") {
          // 请求体整体缺失 → 422 (与 Go/PHP/Rust 端一致), 显式处理而非依赖全局 TypeError 兜底
          return fail(set, 422, "参数错误");
        }
        warn("user", "更新用户信息-body类型错误", { uid: user.id, reason: bodyErr.reason });
        return fail(set, 400, "Bad Request");
      }
      const { nickname, avatar_url = "" } = body as { nickname: string; avatar_url?: string };
      info("user", "POST /update 请求", { uid: user.id, nickname });

      // 验证昵称格式 (缺失字段同样命中以下规则, 与 Go/PHP/Rust 端 422 文案一致)
      if (isNicknameTooShort(nickname)) {
        warn("user", "更新用户信息-昵称过短", { uid: user.id, nickname });
        return fail(set, 422, "昵称长度不能小于2位");
      }
      if (isNicknameTooLong(nickname)) {
        warn("user", "更新用户信息-昵称过长", { uid: user.id, nickname });
        return fail(set, 422, "昵称长度不能超过20位");
      }
      // 验证头像地址格式
      if (avatar_url.length > 255) {
        warn("user", "更新用户信息-头像地址过长", { uid: user.id });
        return fail(set, 422, "头像地址长度不能超过255位");
      }
      // 必须精确以 http:// 或 https:// 开头 (不用 startsWith("http"), 否则 httpxxx:// 也能通过)
      if (
        avatar_url.length > 0 &&
        !avatar_url.startsWith("http://") &&
        !avatar_url.startsWith("https://")
      ) {
        warn("user", "更新用户信息-头像地址格式错误", { uid: user.id, avatar_url });
        return fail(set, 422, "头像地址必须以http或https开头");
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
          // 只取响应所需字段, 避免整行拉取含 password hash (与其它查询保持一致)
          select: userPublicSelect,
        });
        info("user", "更新用户信息成功", { uid: user.id, nickname });
        return toUserInfo(updatedUser);
      } catch (err) {
        error("user", "更新用户信息-数据库更新失败", { uid: user.id, error: String(err) });
        return fail(set, 500, "更新失败，请稍后再试");
      }
    }, {
      detail: {
        tags: ["User"],
        summary: "更新用户信息",
        security: [{ BearerAuth: [] }],
        requestBody: {
          required: true,
          content: {
            "application/json": {
              schema: { $ref: "#/components/schemas/UserUpdateRequest" },
            },
          },
        },
        responses: {
          "200": {
            description: "OK",
            content: {
              "application/json": { schema: { $ref: "#/components/schemas/User" } },
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
  );

export default userRoutes;
