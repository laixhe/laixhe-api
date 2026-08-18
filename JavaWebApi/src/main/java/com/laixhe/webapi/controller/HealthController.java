package com.laixhe.webapi.controller;

import com.laixhe.webapi.common.ApiException;
import com.laixhe.webapi.common.Error;
import com.laixhe.webapi.config.AppProperties;
import com.laixhe.webapi.dto.HealthResponse;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

/**
 * 健康检查 (对应 Go 版 controllers/health.go)
 * 数据库不可用时返回 503 统一错误体, 便于负载均衡探活。
 */
@Tag(name = "Health", description = "健康检查")
@RestController
@RequestMapping("/api/v1")
@RequiredArgsConstructor
public class HealthController {

    private static final DateTimeFormatter TIME_FORMAT = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    /** 数据库探测结果缓存时长(毫秒), 避免频繁探活压垮数据库 (与 Go healthPingInterval 对齐) */
    private static final long PING_INTERVAL_MS = 5000;

    private final JdbcTemplate jdbcTemplate;
    private final AppProperties appProperties;

    /** 服务启动时间 (服务器本地时区) */
    private final String startedAt = LocalDateTime.now().format(TIME_FORMAT);

    private volatile long lastPingAt = 0;
    private volatile boolean lastHealthy = true;

    @Operation(summary = "健康检查")
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "OK", content = @Content(schema = @Schema(implementation = HealthResponse.class))),
            @ApiResponse(responseCode = "503", description = "Service Unavailable", content = @Content(schema = @Schema(implementation = Error.class))),
    })
    @GetMapping("/health")
    public HealthResponse health() {
        if (!dbHealthy()) {
            throw new ApiException(503, "database unavailable");
        }
        return new HealthResponse("ok", "up", appProperties.getVersion(),
                startedAt, LocalDateTime.now().format(TIME_FORMAT));
    }

    /** 探测数据库连接, 结果缓存 PING_INTERVAL_MS 时长 (缓存有效期内并发读互不阻塞) */
    private boolean dbHealthy() {
        long now = System.currentTimeMillis();
        if (now - lastPingAt < PING_INTERVAL_MS) {
            return lastHealthy;
        }
        boolean ok = true;
        try {
            jdbcTemplate.queryForObject("SELECT 1", Integer.class);
        } catch (Exception e) {
            ok = false;
        }
        lastPingAt = System.currentTimeMillis();
        lastHealthy = ok;
        return ok;
    }
}
