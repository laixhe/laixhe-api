package com.laixhe.webapi.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * OpenAPI 文档元信息 (springdoc 3.x 已移除 springdoc.info.* 配置, 改为注入 OpenAPI Bean)
 */
@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI openAPI() {
        return new OpenAPI().info(new Info()
                .title("API接口")
                .description("用户认证与用户管理 API 服务")
                .version("1.0"));
    }
}
