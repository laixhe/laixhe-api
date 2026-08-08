// 共享工具函数：时间格式化、用户模型转响应实体

import type { UserInfo, UserState, Sex } from "../entity/user";
import type { User } from "../generated/prisma/client";

// 格式化为 "YYYY-MM-DD HH:mm:ss"，与 Go time.DateTime 一致
function formatDateTime(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}

// 响应映射所需的字段子集：基于 Prisma 生成的 User 类型推导（Pick），
// schema 变更时自动同步，避免手写类型与模型漂移
type UserInfoSource = Pick<
  User,
  | "id"
  | "typeId"
  | "account"
  | "mobile"
  | "email"
  | "nickname"
  | "avatarUrl"
  | "sex"
  | "states"
  | "createdAt"
>;

// 将 Prisma User 转为 UserInfo（不含 password 等敏感字段）
// 复合词 JSON key 使用 snake_case 以与 Go API 保持一致
export function toUserInfo(user: UserInfoSource): UserInfo {
  return {
    uid: user.id,
    type_id: user.typeId,
    account: user.account,
    mobile: user.mobile,
    email: user.email,
    nickname: user.nickname,
    avatar_url: user.avatarUrl || "",
    sex: user.sex as Sex,
    states: user.states as UserState,
    created_at: formatDateTime(user.createdAt),
  };
}
