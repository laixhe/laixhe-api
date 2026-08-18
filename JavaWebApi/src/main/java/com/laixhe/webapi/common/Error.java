package com.laixhe.webapi.common;

import io.swagger.v3.oas.annotations.media.Schema;

/**
 * 统一错误响应体 (对应 swagger 中 core.Error)
 *
 * @param code    HTTP 状态码
 * @param message 错误信息
 */
@Schema(description = "统一错误响应体")
public record Error(
        @Schema(description = "HTTP 状态码") int code,
        @Schema(description = "错误信息") String message
) {
}
