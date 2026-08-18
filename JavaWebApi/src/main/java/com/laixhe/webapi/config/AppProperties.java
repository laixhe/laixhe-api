package com.laixhe.webapi.config;

import lombok.Getter;
import lombok.Setter;
import org.springframework.boot.context.properties.ConfigurationProperties;

/**
 * 应用配置 (前缀 app, 对应 application.yaml 中 app.*)
 */
@Getter
@Setter
@ConfigurationProperties(prefix = "app")
public class AppProperties {

    /** 服务版本 (构建时可注入, 健康检查接口返回) */
    private String version = "1.0.0";

    /** 请求超时时间(秒), 超时返回 408 (与 Go 版 http.timeout 对齐) */
    private int httpTimeout = 30;

    private Jwt jwt = new Jwt();

    private Limit limit = new Limit();

    @Getter
    @Setter
    public static class Jwt {
        /** 签名密钥 (生产环境务必修改; 环境变量 JWT_SECRET_KEY 优先, 与 Go 版 config.yaml 对齐) */
        private String secretKey;
        /** 过期时长(秒) */
        private long expireSeconds = 2592000;
    }

    @Getter
    @Setter
    public static class Limit {
        /** 是否启用接口限流 */
        private boolean enable = true;
        /** 单个 IP 在窗口内允许的最大请求数 */
        private int max = 1000;
        /** 滑动窗口时长(秒) */
        private int windowSeconds = 60;
    }
}
