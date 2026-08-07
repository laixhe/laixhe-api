// 用户请求/响应类型（保持与 Go API 一致）

/**
 * UserUpdateRequest 更新用户请求（昵称、头像）
 */
export interface UserUpdateRequest {
  nickname: string;
  avatar_url: string;
}

/**
 * UserInfoRequest 用户信息请求（通过 uid 查询）
 */
export interface UserInfoRequest {
  uid: number;
}

/**
 * UserListRequest 用户列表请求（分页）
 * Page 默认 1, PageSize 默认 12
 */
export interface UserListRequest {
  page?: number;
  pageSize?: number;
}

/**
 * UserListResponse 用户列表响应
 */
export interface UserListResponse {
  total: number;
  page: number;
  page_size: number;
  list: import("./auth").UserInfo[];
}

// UserType 用户类型枚举
//   - 0: 未知
//   - 1: 普通用户
export enum UserType {
  Normal = 1,
}

// UserSex 用户性别枚举
//   - 0: 未填写
//   - 1: 男
//   - 2: 女
export enum UserSex {
  Unknown = 0,
  Male = 1,
  Female = 2,
}

// UserState 用户状态枚举
//   - 0: 禁用
//   - 1: 正常
export enum UserState {
  Disabled = 0,
  Normal = 1,
}
