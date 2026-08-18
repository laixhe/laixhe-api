// 应用入口
// 初始化 Elysia HTTP 服务，注册路由、中间件，配置优雅关闭
// 路由前缀：/api/v1/auth/*（鉴权），/api/v1/user/*（用户）

import { Elysia } from "elysia";
import { cors } from "@elysiajs/cors";
import { openapi } from "@elysia/openapi";
import { stringify } from "yaml";
import { randomUUID } from "crypto";
import config from "./config";
import { prisma } from "./lib/prisma";
import { authRoutes } from "./routes/auth";
import { userRoutes } from "./routes/user";
import { healthRoutes } from "./routes/health";
import { UnauthorizedError } from "./middleware/authGuard";
import { globalRateLimit } from "./middleware/rateLimit";
import { schemas } from "./docs/schemas";

// OpenAPI 文档由 @elysia/openapi 插件 (官方推荐的 OpenAPI 插件, 替代已弃用的 @elysiajs/swagger)
// 从路由 detail 注解运行时动态生成 (替代原 swagger-jsdoc 静态生成)
// 端点: /api/v1/swagger (UI), /api/v1/swagger.json (JSON spec), /api/v1/swagger.yaml (YAML, 下方转换)
const app = new Elysia()
  .use(cors())
  // 每个请求生成 X-Request-Id 响应头 (与 Go/PHP/Rust 端 requestId 中间件对齐, 便于日志串联排查;
  // 若客户端已携带则透传, 否则生成 UUID)
  .onRequest(({ request, set }) => {
    const requestId = request.headers.get("x-request-id") || randomUUID();
    set.headers["x-request-id"] = requestId;
  })
  // 全局错误处理：生产环境不暴露内部错误详情
  // 注意：onError 需注册在 .use() 之前，否则插件内抛出的错误不会进入此处理
  // （Elysia 1.4.x 的已知行为）
  .onError(({ code, error, set }) => {
    // 鉴权插件抛出的统一 401 异常
    if (error instanceof UnauthorizedError) {
      set.status = 401;
      return { code: 401, message: error.message };
    }
    // 注意: "缺 body → 422" 已在各 handler 中通过 bodyError 显式处理 (见 util/validate.ts 与 routes),
    // 因此这里不再宽泛地把 TypeError 归类为 422, 避免掩盖业务代码中真实的 TypeError bug (应返回 500 暴露)
    // 请求体 JSON 解析失败 (如带 content-type 但空 body / 非法 JSON): 属请求格式错误
    if (code === "PARSE") {
      set.status = 400;
      return { code: 400, message: "Bad Request" };
    }
    if (code === "NOT_FOUND") {
      return { code: 404, message: "Not Found" };
    }
    console.error(`[Unhandled Error] ${code}:`, error);
    return { code: 500, message: "服务器内部错误" };
  })
  // 全局限流中间件 (对齐 Go/PHP/Rust 端): 单个 IP 在 60s 窗口内最多 1000 次请求, 超过返回 429 统一 JSON;
  // 健康检查路径豁免限流 (负载均衡/容器编排探活不应被 429 拦截)
  .onBeforeHandle(globalRateLimit)
  .use(
    openapi({
      path: "/api/v1/swagger",
      specPath: "/api/v1/swagger.json",
      provider: "swagger-ui",
      exclude: {
        // 隐藏含点号的静态文档路径 (如 /api/v1/swagger.json|yaml), 不进入文档
        staticFile: true,
      },
      // 与其余语言端 UI 一致, 使用 swagger-ui-dist@5 (unpkg CDN);
      // 插件默认生成的 UI spec URL 为相对路径 (specPath 去掉前导 /), 页面在 /api/v1/swagger 时会被浏览器
      // 解析成 /api/v1/api/v1/swagger.json (404), 此处用 url 覆盖为绝对路径修复。
      // 注意: 插件类型把 url 排除在 swagger 配置外, 故用 as never 断言 (运行时插件在 url 之后展开 swagger, 覆盖生效)
      swagger: {
        version: "5",
        url: "/api/v1/swagger.json",
      } as never,
      documentation: {
        info: {
          title: "API接口",
          description: "用户认证与用户管理 API 服务",
          version: "1.0",
        },
        tags: [
          { name: "Auth", description: "鉴权" },
          { name: "User", description: "用户" },
          { name: "Health", description: "健康检查" },
        ],
        security: [{ BearerAuth: [] }],
        components: {
          securitySchemes: {
            BearerAuth: {
              type: "http",
              scheme: "bearer",
              bearerFormat: "JWT",
              description: "在请求头携带 Authorization: Bearer <token>",
            },
          },
          schemas,
        },
      },
    })
  )
  .use(authRoutes)
  .use(userRoutes)
  .use(healthRoutes)
  // 根路径健康提示 (不进入 OpenAPI 文档)
  .get("/", () => "laixhe-api is running", { detail: { hide: true } });

// /api/v1/swagger.yaml: 插件仅提供 JSON spec (specPath), 此处请求 JSON 后动态转换为 YAML, 保持端点与 Go/PHP/Rust 版一致。
// 需在 app 声明完成后注册 (handler 内引用 app 自身, 拆开避免 TS 自引用类型推断失败)
const finalApp = app.get(
  "/api/v1/swagger.yaml",
  async () => {
    const json = await app
      .handle(new Request("http://internal/api/v1/swagger.json"))
      .then((res) => res.text());
    return new Response(stringify(JSON.parse(json)), {
      headers: {
        "Content-Type": "application/x-yaml",
        "Cache-Control": "public, max-age=300",
      },
    });
  },
  { detail: { hide: true } }
);

// 仅在作为入口直接运行时启动服务器
// （被测试或其它模块 import 时不监听端口，避免端口占用/进程挂起）
if (import.meta.main) {
  finalApp.listen({
    hostname: config.http.ip,
    port: config.http.port,
  });

  console.log(
    `Server is running at http://${finalApp.server?.hostname}:${finalApp.server?.port}`
  );

  // 优雅关闭：收到 SIGTERM/SIGINT 时断开 Prisma 连接
  async function shutdown() {
    console.log("Shutting down...");
    await prisma.$disconnect();
    process.exit(0);
  }
  process.on("SIGTERM", shutdown);
  process.on("SIGINT", shutdown);
}

export default finalApp;
