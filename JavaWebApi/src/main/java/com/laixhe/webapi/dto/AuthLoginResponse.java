package com.laixhe.webapi.dto;

import io.swagger.v3.oas.annotations.media.Schema;

/**
 * 登录响应 (对应 swagger entity.AuthLoginResponse)
 */
@Schema(description = "登录响应")
public record AuthLoginResponse(
        @Schema(description = "jwt token") String token,
        @Schema(description = "用户") UserResponse user
) {
}
