import { Elysia, t } from "elysia";

import { jwtConfig } from "@middleware/jwt";
import { AuthService } from "@service/auth/index";
import { LoginRequestSchema, LoginResponseSchema } from "@entity/auth/login";

export const authLogin = new Elysia().use(jwtConfig).post(
  "login",
  async ({ JwtConfig, body }) => {
    const resp = await AuthService.Login(body);
    const token = await JwtConfig.sign({
      uid: resp.user.uid,
    });
    resp.token = token;
    return resp;
  },
  {
    body: LoginRequestSchema,
    response: {
      200: LoginResponseSchema,
    },
    detail: {
      summary: "登录",
    },
  }
);
