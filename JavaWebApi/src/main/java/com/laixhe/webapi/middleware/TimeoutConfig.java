package com.laixhe.webapi.middleware;

import com.laixhe.webapi.config.AppProperties;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import tools.jackson.databind.ObjectMapper;

import java.util.Arrays;

/**
 * 注册请求超时过滤器。
 * 顺序 -300: 早于限流过滤器 (-200), 与 Go 版中间件顺序 (超时 → 限流 → 业务路由) 对齐。
 */
@Configuration
public class TimeoutConfig {

    @Bean
    public FilterRegistrationBean<TimeoutFilter> timeoutFilterRegistration(AppProperties appProperties,
                                                                           ObjectMapper objectMapper,
                                                                           Environment environment) {
        FilterRegistrationBean<TimeoutFilter> registration =
                new FilterRegistrationBean<>(new TimeoutFilter(appProperties.getHttpTimeout(), objectMapper));
        registration.setOrder(-300);
        // 测试 profile (h2) 关闭: 超时过滤器将请求转交独立线程执行,
        // 与 @Transactional 测试的线程绑定事务不兼容, 关闭以保持事务语义
        registration.setEnabled(!Arrays.asList(environment.getActiveProfiles()).contains("h2"));
        return registration;
    }
}
