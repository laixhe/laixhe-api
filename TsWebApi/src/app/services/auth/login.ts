import { Elysia } from "elysia";

import dayjs from "dayjs";

import { log } from "@core/log";
import { jwtConfig } from "@middleware/jwt";

import { ParamError } from "@core/error";
import { UserModel } from "@model/user/index";
import {
  LoginRequestSchema,
  LoginResponseSchema,
  LoginRequest,
  LoginResponse,
} from "@entity/auth/login";

const login = new Elysia()
  .use(log)
  .use(jwtConfig)
  .post(
    "login",
    async ({ JwtConfig, body }) => {
      const user = await UserModel.getByEmail(body.email);
      if (!user) {
        throw new ParamError("邮箱或密码不正确");
      }
      const token = await JwtConfig.sign({
        uid: user.id,
      });
      let resp: LoginResponse = {
        token: token,
        user: {
          uid: user.id,
          type_id: user.type_id,
          nickname: user.nickname,
          email: user.email,
          avatar_url: user.avatar_url,
          states: user.states,
          created_at: dayjs(user.created_at).format("YYYY-MM-DD HH:mm:ss"),
        },
      };
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

export default login;
