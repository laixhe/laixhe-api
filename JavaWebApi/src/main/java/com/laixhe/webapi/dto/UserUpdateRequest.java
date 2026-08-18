package com.laixhe.webapi.dto;

import io.swagger.v3.oas.annotations.media.Schema;

/**
 * 更新用户信息请求 (对应 swagger entity.UserUpdateRequest)
 * Uid 由 JWT 提供, 不入参; 昵称/头像格式由控制器校验 (与 Go 版一致)
 */
@Schema(description = "更新用户信息请求")
public record UserUpdateRequest(
        @Schema(description = "昵称") String nickname,
        @Schema(name = "avatar_url", description = "头像地址") String avatarUrl
) {
}
