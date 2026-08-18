package com.laixhe.webapi.security;

/**
 * JWT 载荷, 存储用户 UID (与 Go 版 middlewares.JwtClaims 对齐)
 *
 * @param uid 用户id (从 1 起, 0 视为无效, 防御伪造 {"uid":0} 的 token)
 */
public record JwtClaims(int uid) {
}
