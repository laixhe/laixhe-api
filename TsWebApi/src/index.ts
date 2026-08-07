// 应用入口
// 初始化 Elysia HTTP 服务，注册路由、中间件，配置优雅关闭
// 路由前缀：/api/v1/auth/*（鉴权），/api/v1/user/*（用户）

import { Elysia } from "elysia";
import { cors } from "@elysiajs/cors";
import config from "./config";
import { prisma } from "./lib/prisma";
import { authRoutes } from "./routes/auth";
import { userRoutes } from "./routes/user";

const app = new Elysia()
  .use(cors())
  .use(authRoutes)
  .use(userRoutes)
  .get("/", () => "laixhe-api is running")
  // 全局错误处理：生产环境不暴露内部错误详情
  .onError(({ code, error }) => {
    if (code === "NOT_FOUND") {
      return { code: 404, message: "Not Found" };
    }
    if (code === "VALIDATION") {
      return { code: 400, message: "请求参数错误" };
    }
    console.error(`[Unhandled Error] ${code}:`, error);
    return { code: 500, message: "服务器内部错误" };
  })
  .listen({
    hostname: config.http.ip,
    port: config.http.port,
  });

console.log(
  `Server is running at http://${app.server?.hostname}:${app.server?.port}`
);

// 优雅关闭：收到 SIGTERM/SIGINT 时断开 Prisma 连接
async function shutdown() {
  console.log("Shutting down...");
  await prisma.$disconnect();
  process.exit(0);
}
process.on("SIGTERM", shutdown);
process.on("SIGINT", shutdown);

export default app;
