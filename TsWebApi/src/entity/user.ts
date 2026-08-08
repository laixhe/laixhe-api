// User 相关类型：用户状态枚举、用户信息响应体（保持与 Go API 一致）

// UserState 用户状态枚举
//   - 0: 禁用
//   - 1: 正常
export enum UserState {
  Disabled = 0,
  Normal = 1,
}

// Sex 性别枚举（DB 层为 Int，枚举值即入库数值）
//   - 0: 未填写
//   - 1: 男
//   - 2: 女
export enum Sex {
  Unknown = 0,
  Male = 1,
  Female = 2,
}

/**
 * UserInfo 用户信息（不含密码）
 * JSON key 与 Go 结构体 json tag 保持一致：
 *   复合词字段使用 snake_case（type_id, avatar_url, created_at）
 *   单词字段保持原样（uid, account, mobile, email, nickname, sex, states）
 */
export interface UserInfo {
  uid: number;         // 用户id
  type_id: number;     // 类型 1-普通用户
  account: string;     // 账号
  mobile: string | null;  // 手机号
  email: string | null;   // 邮箱
  nickname: string;    // 昵称
  avatar_url: string;  // 头像地址（null 时归一化为空字符串）
  sex: Sex;            // 性别，取值见 Sex 枚举（0-未填写 1-男 2-女）
  states: UserState;   // 状态，取值见 UserState 枚举（0-禁用 1-正常）
  created_at: string;  // 创建时间，格式 "YYYY-MM-DD HH:mm:ss"
}
