package com.laixhe.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;
import lombok.Data;

/**
 * @author laixhe
 */
@Data
@Component
@ConfigurationProperties(prefix = "jwt")
public class JwtConfig {
    // 密钥
    private String secret;
    // 有效期（秒）
    private int expire;
}
