package com.laixhe.webapi;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.ConfigurationPropertiesScan;

/**
 * 应用启动入口 (基于 swagger.yaml 生成, 与 Go/PHP/TS/Rust 版 API 对齐)
 */
@SpringBootApplication
@ConfigurationPropertiesScan
public class WebApiApplication {

    public static void main(String[] args) {
        SpringApplication.run(WebApiApplication.class, args);
    }
}
