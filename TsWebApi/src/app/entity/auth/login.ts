import { t, type Static } from "elysia";

export const LoginRequestSchema = t.Object({
  email: t.String({ format: "email", description: "用户邮箱" }),
  password: t.String({
    minLength: 6,
    maxLength: 20,
    description: "用户密码(6~20个字符)",
  }),
});

export type LoginRequest = Static<typeof LoginRequestSchema>;

export const LoginResponseSchema = t.Object({
  token: t.String({ description: "Token: Bearer xxx" }),
  user: t.Object(
    {
      uid: t.Number({ description: "用户ID" }),
      nickname: t.String({ description: "昵称" }),
      email: t.String({ description: "邮箱" }),
      created_at: t.String({ description: "创建时间" }),
    },
    { title: "用户信息" }
  ),
});

export type LoginResponse = Static<typeof LoginResponseSchema>;
