// 鉴权请求/响应类型（保持与 Go API snake_case JSON 键名一致）

/**
 * AuthRegisterRequest 注册请求
 * nickname 昵称（2-20字）, email 邮箱, password 密码（>=6位 字母数字_@$）
 */
export interface AuthRegisterRequest {
  nickname: string;
  email: string;
  password: string;
}

/**
 * AuthLoginRequest 登录请求
 */
export interface AuthLoginRequest {
  email: string;
  password: string;
}

/**
 * AuthRefreshRequest 刷新Token请求
 */
export interface AuthRefreshRequest {
  uid: number;
}

/**
 * AuthTokenResponse 鉴权 Token 响应
 */
export interface AuthTokenResponse {
  token: string;
  user: UserInfo;
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
  avatar_url: string;  // 头像地址
  sex: number;         // 性别 0-未填写 1-男 2-女
  states: number;      // 状态 0-禁用 1-正常
  created_at: string;  // 创建时间，格式 "YYYY-MM-DD HH:mm:ss"
}
