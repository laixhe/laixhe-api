package com.laixhe.webapi.middleware;

import com.laixhe.webapi.config.AppProperties;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import tools.jackson.databind.ObjectMapper;

/**
 * 注册 IP 限流过滤器。
 * 顺序 -200: 早于 Spring Security 过滤器链 (-100), 保证限流先于 JWT 校验执行,
 * 与 Go 版中间件顺序 (限流 → 业务路由/JWT) 对齐。
 */
@Configuration
public class RateLimitConfig {

    @Bean
    public FilterRegistrationBean<RateLimitFilter> rateLimitFilterRegistration(AppProperties appProperties,
                                                                               ObjectMapper objectMapper) {
        FilterRegistrationBean<RateLimitFilter> registration =
                new FilterRegistrationBean<>(new RateLimitFilter(appProperties, objectMapper));
        registration.setOrder(-200);
        return registration;
    }
}
