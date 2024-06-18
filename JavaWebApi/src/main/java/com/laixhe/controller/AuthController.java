package com.laixhe.controller;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import lombok.extern.slf4j.Slf4j;

import com.laixhe.utils.JWTUtils;
import com.laixhe.result.Result;
import com.laixhe.service.AuthService;
import com.laixhe.request.auth.LoginRequest;
import com.laixhe.response.auth.LoginResponse;
import com.laixhe.request.auth.RegisterRequest;
import com.laixhe.response.auth.RegisterResponse;
import com.laixhe.response.auth.RefreshResponse;

import java.util.HashMap;
import java.util.Map;

/**
 * @author laixhe
 */
@Slf4j
@RestController
@RequestMapping("/api/auth")
public class AuthController {

    private final AuthService authService;

    @Autowired
    public AuthController(AuthService authService) {
        this.authService = authService;
    }

    /**
     * 注册
     */
    @PostMapping("/register")
    public Result<RegisterResponse> register(@RequestBody @Validated RegisterRequest req) {
        log.info("Register {}", req.toString());

        RegisterResponse resp = authService.register(req);
        return Result.success(resp);
    }

    /**
     * 登录
     */
    @PostMapping("/login")
    public Result<LoginResponse> login(@RequestBody @Validated LoginRequest req) {
        log.info("login {}", req.toString());

        LoginResponse resp = authService.login(req.getEmail(), req.getPassword());
        return Result.success(resp);
    }

    /**
     * 刷新Jwt
     */
    @PostMapping("/refresh")
    public Result<RefreshResponse> refresh(HttpServletRequest request) {
        int uid = (int) request.getAttribute("uid");
        log.info("refresh uid={}", uid);

        // 生成 token 数据
        Map<String, Object> payload = new HashMap<>();
        payload.put("uid", uid);
        String token = JWTUtils.createToken(payload);

        RefreshResponse resp = new RefreshResponse();
        resp.setToken(token);
        return Result.success(resp);
    }

}
