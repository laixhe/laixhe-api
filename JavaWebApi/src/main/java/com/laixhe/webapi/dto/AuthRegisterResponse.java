package com.laixhe.webapi.dto;

import io.swagger.v3.oas.annotations.media.Schema;

/**
 * 注册响应 (对应 swagger entity.AuthRegisterResponse)
 */
@Schema(description = "注册响应")
public record AuthRegisterResponse(
        @Schema(description = "jwt token") String token,
        @Schema(description = "用户") UserResponse user
) {
}
