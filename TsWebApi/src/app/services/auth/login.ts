import dayjs from "dayjs";

import { ParamError } from "@core/error";
import { UserModel } from "@model/user";
import { LoginRequest, LoginResponse } from "@entity/auth/login";

export default async function Login(req: LoginRequest): Promise<LoginResponse> {
  const user = await UserModel.GetByEmail(req.email);
  if (!user) {
    throw new ParamError("邮箱或密码不正确");
  }
  return {
    token: "",
    user: {
      uid: user.id,
      nickname: user.nickname,
      email: user.email,
      created_at: dayjs(user.created_at).utc().format("YYYY-MM-DD HH:mm:ss"),
    },
  };
}
