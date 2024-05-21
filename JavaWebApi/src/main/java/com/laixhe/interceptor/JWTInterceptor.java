package com.laixhe.interceptor;

import cn.hutool.core.util.StrUtil;
import com.auth0.jwt.interfaces.DecodedJWT;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.web.servlet.HandlerInterceptor;
import lombok.extern.slf4j.Slf4j;

import com.laixhe.exception.BusinessException;
import com.laixhe.utils.JWTUtils;
import com.laixhe.result.ResultCode;

/**
 * 拦截器
 *
 * @author laixhe
 */
@Slf4j
public class JWTInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {

        // 获取请求头中的令牌
        String authorization = request.getHeader("Authorization");
        if (StrUtil.isBlank(authorization) || !authorization.startsWith("Bearer ")) {
            throw new BusinessException(ResultCode.AuthInvalid, "");
        }
        String token = authorization.substring(7);
        // token 验证
        DecodedJWT decodedJWT = JWTUtils.verifyToken(token);
        Integer uid = decodedJWT.getClaim("uid").asInt();
        if (uid <= 0) {
            throw new BusinessException(ResultCode.AuthInvalid, "");
        }

        request.setAttribute("uid", uid);
        return true;
    }
}
