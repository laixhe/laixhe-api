import { Elysia } from "elysia";

import log from "@core/log";
import { jwtConfig } from "@middleware/jwt";
import { AuthService } from "@service/auth/index";
import { LoginRequestSchema, LoginResponseSchema } from "@entity/auth/login";

const Login = new Elysia()
  .use(log)
  .use(jwtConfig)
  .post(
    "login",
    async (context) => {
      const resp = await AuthService.Login(context.body);
      const token = await context.JwtConfig.sign({
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

export default Login;
