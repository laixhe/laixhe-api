package com.laixhe.webapi.middleware;

import com.laixhe.webapi.common.Error;
import com.laixhe.webapi.config.AppProperties;
import io.github.bucket4j.Bandwidth;
import io.github.bucket4j.Bucket;
import io.github.bucket4j.Refill;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import tools.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

/**
 * 基于 Bucket4j 的 IP 限流过滤器 (与 Go 版 core/middlewares/rate_limit.go 对齐):
 * 每个 IP 一个令牌桶, 窗口内最多 {@code max} 次请求, 超过阈值返回 429 统一 JSON。
 * 健康检查路径豁免限流, 便于负载均衡/容器编排探活。
 */
@Slf4j
public class RateLimitFilter implements Filter {

    /** 健康检查路径 (限流豁免, 避免负载均衡探活被误伤) */
    private static final String HEALTH_PATH = "/api/v1/health";

    /** 内存保护: IP 桶空闲超过 2 个窗口时长即被后台清理, 防止伪造 IP 导致内存无限增长 */
    private static final long IDLE_WINDOWS = 2;

    private final AppProperties appProperties;
    private final ObjectMapper objectMapper;
    private final int max;
    private final Duration window;
    private final Map<String, IpBucket> buckets = new ConcurrentHashMap<>();

    private ScheduledExecutorService janitor;

    public RateLimitFilter(AppProperties appProperties, ObjectMapper objectMapper) {
        this.appProperties = appProperties;
        this.objectMapper = objectMapper;
        // max < 1 时按 1 处理 (与 Go 版 NewRateLimiter 对齐)
        this.max = Math.max(appProperties.getLimit().getMax(), 1);
        this.window = Duration.ofSeconds(Math.max(appProperties.getLimit().getWindowSeconds(), 1));
    }

    /** 启动后台清理任务: 周期性扫描空闲 IP 桶, 与请求路径解耦 (对应 Go 版 janitor) */
    @PostConstruct
    void startJanitor() {
        janitor = Executors.newSingleThreadScheduledExecutor(r -> {
            Thread t = new Thread(r, "rate-limit-janitor");
            t.setDaemon(true);
            return t;
        });
        janitor.scheduleWithFixedDelay(this::cleanup, window.toSeconds(), window.toSeconds(), TimeUnit.SECONDS);
    }

    @PreDestroy
    void stopJanitor() {
        if (janitor != null) {
            janitor.shutdownNow();
        }
    }

    private void cleanup() {
        long now = System.nanoTime();
        long idleNanos = window.toNanos() * IDLE_WINDOWS;
        buckets.values().removeIf(entry -> now - entry.lastAccessNanos > idleNanos);
    }

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {
        // 配置关闭限流时直接放行
        if (!appProperties.getLimit().isEnable()) {
            chain.doFilter(request, response);
            return;
        }
        if (!(request instanceof HttpServletRequest httpRequest) || !(response instanceof HttpServletResponse httpResponse)) {
            chain.doFilter(request, response);
            return;
        }
        // 健康检查路径豁免限流
        if (HEALTH_PATH.equals(httpRequest.getRequestURI())) {
            chain.doFilter(request, response);
            return;
        }
        String ip = resolveClientIp(httpRequest);
        if (!allow(ip)) {
            log.warn("接口限流触发 ip={}", ip);
            httpResponse.setStatus(429);
            httpResponse.setContentType("application/json");
            httpResponse.setCharacterEncoding(StandardCharsets.UTF_8.name());
            httpResponse.getWriter().write(objectMapper.writeValueAsString(new Error(429, "请求过于频繁，请稍后再试")));
            return;
        }
        chain.doFilter(request, response);
    }

    /** 尝试为当前 IP 消费 1 个令牌: 成功放行, 失败即超限 */
    private boolean allow(String ip) {
        long now = System.nanoTime();
        IpBucket entry = buckets.compute(ip, (key, existing) -> {
            if (existing == null) {
                // 令牌桶: 容量 = max, 按窗口时长匀速补充 max 个令牌 (滑动窗口效果)
                Bucket bucket = Bucket.builder()
                        .addLimit(Bandwidth.classic(max, Refill.greedy(max, window)))
                        .build();
                return new IpBucket(bucket, now);
            }
            existing.lastAccessNanos = now;
            return existing;
        });
        return entry.bucket.tryConsume(1);
    }

    /**
     * 解析客户端 IP: 优先 X-Forwarded-For 首个地址, 其次真实连接地址, 兜底 "unknown"。
     * 注意: 直接信任 X-Forwarded-For 时客户端可伪造该头绕过限流, 生产环境建议仅在有可信反向代理时使用。
     */
    private String resolveClientIp(HttpServletRequest request) {
        String xff = request.getHeader("X-Forwarded-For");
        if (xff != null) {
            int comma = xff.indexOf(',');
            String first = (comma >= 0 ? xff.substring(0, comma) : xff).trim();
            if (!first.isEmpty()) {
                return first;
            }
        }
        String remote = request.getRemoteAddr();
        return remote != null && !remote.isEmpty() ? remote : "unknown";
    }

    /** IP 对应的令牌桶及最后访问时间 */
    private static final class IpBucket {
        final Bucket bucket;
        volatile long lastAccessNanos;

        IpBucket(Bucket bucket, long lastAccessNanos) {
            this.bucket = bucket;
            this.lastAccessNanos = lastAccessNanos;
        }
    }
}
