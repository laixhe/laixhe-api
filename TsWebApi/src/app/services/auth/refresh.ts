import { Elysia, t } from "elysia";

import { log } from "@core/log";
import { jwtConfig, jwtAuth } from "@middleware/jwt";

const refresh = new Elysia()
  .use(log)
  .use(jwtConfig)
  .use(jwtAuth)
  .post(
    "refresh",
    async (context) => {
      const token = await context.JwtConfig.sign({
        uid: context.uid,
      });
      let resp: {
        token: string;
        user: {
          uid: number;
          uname: string;
          email: string;
          created_at: string;
        };
      } = {
        token: token,
        user: {
          uid: context.uid,
          uname: "refresh.password",
          email: "refresh.email",
          created_at: "2025-05-02",
        },
      };
      return resp;
    },
    {
      response: {
        200: t.Object({
          token: t.String({ description: "Token: Bearer xxx" }),
          user: t.Object({
            uid: t.Number({ description: "用户ID" }),
            uname: t.String({ description: "用户名" }),
            email: t.String({ description: "用户邮箱" }),
            created_at: t.String({ description: "创建时间" }),
          }),
        }),
      },
      detail: {
        summary: "刷新Jwt",
      },
    }
  );

export default refresh;
