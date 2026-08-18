package com.laixhe.webapi.dto;

import io.swagger.v3.oas.annotations.media.Schema;

/**
 * 健康检查响应体 (对应 swagger controllers.HealthResponse)
 */
@Schema(description = "健康检查响应体")
public record HealthResponse(
        @Schema(description = "服务状态 (固定 ok)") String status,
        @Schema(description = "数据库状态 (正常时为 up; 数据库不可用时直接返回 503 错误体)") String database,
        @Schema(description = "服务版本") String version,
        @Schema(name = "started_at", description = "服务启动时间 (服务器本地时区)") String startedAt,
        @Schema(description = "当前时间 (服务器本地时区)") String now
) {
}
