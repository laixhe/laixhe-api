// 鉴权响应类型（保持与 Go API snake_case JSON 键名一致）

import type { UserInfo } from "./user";

/**
 * AuthTokenResponse 鉴权成功后的响应体
 *
 * 示例：
 * {
 *   "token": "eyJhbGciOiJIUzI1NiIs...",
 *   "user": {
 *     "uid": 1,
 *     "type_id": 1,
 *     "account": "3f9c...",
 *     "mobile": null,
 *     "email": "user@example.com",
 *     "nickname": "张三",
 *     "avatar_url": "",
 *     "sex": 1,
 *     "states": 1,
 *     "created_at": "2026-08-08 12:00:00"
 *   }
 * }
 */
export interface AuthTokenResponse {
  token: string;
  user: UserInfo;
}
