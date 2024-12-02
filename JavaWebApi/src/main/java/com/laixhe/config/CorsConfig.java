package com.laixhe.config;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;
import org.springframework.web.filter.CorsFilter;

///**
// * 全局跨域处理
// *
// * @author laixhe
// */
//@Slf4j
//@Configuration
//public class CorsConfig implements WebMvcConfigurer {
//    @Override
//    public void addCorsMappings(CorsRegistry registry) {
//
//        registry
//                // 允许所有请求路径跨域访问
//                .addMapping("/**")
//                // 是否携带Cookie，默认false
//                .allowCredentials(true)
//                // 允许的请求头类型
//                .allowedHeaders("*")
//                // 预检请求的缓存时间（单位：秒）
//                .maxAge(3600)
//                // 允许的请求方法类型
//                .allowedMethods("*")
//                // 允许哪些域名进行跨域访问
//                .allowedOrigins("*");
//    }
//}

/**
 * 跨域处理
 *
 * @author laixhe
 */
@Slf4j
@Configuration
public class CorsConfig {
    @Bean
    public CorsFilter corsFilter() {
        CorsConfiguration config = new CorsConfiguration();
        // 允许的请求头类型
        config.addAllowedHeader("*");
        // 允许的请求方法类型
        config.addAllowedMethod("*");
        // 允许哪些域名进行跨域访问
//        config.addAllowedOrigin("http://127.0.0.1:5500");
        config.addAllowedOrigin("*");
        // 是否携带Cookie，默认false
        config.setAllowCredentials(true);

        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();

        // 允许所有请求路径跨域访问
        source.registerCorsConfiguration("/**", config);

        return new CorsFilter(source);
    }
}