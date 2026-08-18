package com.laixhe.webapi.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;

/**
 * 登录请求 (对应 swagger entity.AuthLoginRequest)
 */
@Schema(description = "登录请求")
public record AuthLoginRequest(
        @NotBlank(message = "邮箱格式错误")
        @Email(message = "邮箱格式错误")
        @Schema(description = "邮箱")
        String email,
        @NotBlank(message = "密码格式错误，需 6~64 位，只能包含字母 数字 _ @ $")
        @Pattern(regexp = "^[a-zA-Z0-9_@$]{6,64}$", message = "密码格式错误，需 6~64 位，只能包含字母 数字 _ @ $")
        @Schema(description = "密码")
        String password
) {
}
