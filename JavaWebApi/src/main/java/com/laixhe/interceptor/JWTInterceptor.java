package com.laixhe.interceptor;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.web.servlet.HandlerInterceptor;
import lombok.extern.slf4j.Slf4j;

/**
 * 拦截器
 * @author laixhe
 */
@Slf4j
public class JWTInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler){

        // 获取请求头中的令牌
        String token = request.getHeader("token");
        log.info("当前token为：{}", token);

        return true;
    }
}
