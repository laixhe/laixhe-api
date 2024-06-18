package com.laixhe.service;

import com.laixhe.request.auth.RegisterRequest;
import com.laixhe.response.auth.LoginResponse;
import com.laixhe.response.auth.RegisterResponse;

/**
 * @author laixhe
 */
public interface AuthService {
    /**
     * 注册
     * @param req 注册请求参数
     */
    RegisterResponse register(RegisterRequest req);
    /**
     * 登录
     * @param email 邮箱
     * @param password 密码
     * @throws com.laixhe.exception.BusinessException 异常
     */
    LoginResponse login(String email, String password);
}
