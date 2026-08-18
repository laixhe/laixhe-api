package com.laixhe.webapi.dto;

import io.swagger.v3.oas.annotations.media.Schema;

/**
 * 刷新Jwt响应 (对应 swagger entity.AuthRefreshResponse)
 */
@Schema(description = "刷新Jwt响应")
public record AuthRefreshResponse(
        @Schema(description = "jwt token") String token,
        @Schema(description = "用户") UserResponse user
) {
}
