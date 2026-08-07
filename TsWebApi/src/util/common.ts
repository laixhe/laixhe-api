// 共享工具函数：时间格式化、用户模型转响应实体

import type { UserInfo } from "../entity/auth";

// 格式化为 "YYYY-MM-DD HH:mm:ss"，与 Go time.DateTime 一致
export function formatDateTime(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}

// 将 Prisma User 转为 UserInfo（去除密码等敏感字段）
// 复合词 JSON key 使用 snake_case 以与 Go API 保持一致
export function toUserInfo(user: {
  id: number;
  typeId: number;
  account: string;
  mobile: string | null;
  email: string | null;
  nickname: string;
  avatarUrl: string | null;
  sex: number;
  states: number;
  createdAt: Date;
}): UserInfo {
  return {
    uid: user.id,
    type_id: user.typeId,
    account: user.account,
    mobile: user.mobile,
    email: user.email,
    nickname: user.nickname,
    avatar_url: user.avatarUrl || "",
    sex: user.sex,
    states: user.states,
    created_at: formatDateTime(user.createdAt),
  };
}
