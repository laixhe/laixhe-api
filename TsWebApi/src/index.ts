import { Elysia } from "elysia";
import { swagger } from "@elysiajs/swagger";
import { cors } from "@elysiajs/cors";
import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";
import customParseFormat from "dayjs/plugin/customParseFormat";

import { log } from "@core/log";
import { routeV1 } from "@route/route";

dayjs.extend(utc);
dayjs.extend(customParseFormat);

const app = new Elysia()
  .use(log)
  .use(cors())
  .use(
    swagger({
      documentation: {
        info: {
          title: "API",
          description: "API接口文档",
          version: "0.0.1",
        },
        components: {
          securitySchemes: {
            bearerAuth: {
              type: "http",
              scheme: "bearer",
              bearerFormat: "JWT",
            },
          },
        },
        tags: [
          { name: "Auth", description: "鉴权相关" },
          { name: "User", description: "用户相关" },
        ],
      },
    })
  )
  .use(routeV1)
  .listen(process.env.PORT || 6600);

console.log(
  `🦊 Elysia is running at ${app.server?.hostname}:${app.server?.port}`
);
