package com.laixhe.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import com.auth0.jwt.interfaces.DecodedJWT;
import cn.hutool.core.util.StrUtil;
import lombok.extern.slf4j.Slf4j;

import java.util.HashMap;
import java.util.Map;

import com.laixhe.exception.BusinessException;
import com.laixhe.utils.AssertUtils;
import com.laixhe.result.ResultCode;
import com.laixhe.result.Result;
import com.laixhe.utils.JWTUtils;
import com.laixhe.entity.User;
import com.laixhe.service.LoginService;
import com.laixhe.request.LoginRequest;
import com.laixhe.response.LoginResponse;

/**
 * @author laixhe
 */
@Slf4j
@RestController
@RequestMapping("/api/auth")
public class AuthLoginController {

    @Autowired
    private RedisTemplate<String, Object> redisTemplate;

    private final LoginService loginService;

    @Autowired
    public AuthLoginController(LoginService loginService) {
        this.loginService = loginService;
    }

    @PostMapping("/login")
    public Result<LoginResponse> login(@RequestBody @Validated LoginRequest req) {
        log.info("login LoginRequest={}", req.toString());
//        AssertUtils.isTrue(StrUtil.isNotBlank(req.getEmail()), ResultCode.ERROR_PARAM, "邮箱不能空!");
//        AssertUtils.isTrue(StrUtil.isNotBlank(req.getPassword()), ResultCode.ERROR_PARAM, "密码不能空!");

        try {
            User user = loginService.login(req.getEmail(), req.getPassword());

            // 生成 token 数据
            Map<String, Object> payload = new HashMap<>();
            payload.put("uid", user.getId());
            payload.put("email", user.getEmail());
            // 生成 token
            String token = JWTUtils.createToken(payload);

            LoginResponse resp = new LoginResponse();
            resp.setUid(user.getId());
            resp.setEmail(user.getEmail());
            resp.setUname(user.getUname());
            resp.setScore(user.getScore());
            resp.setCreatedAt(user.getCreatedAt());
            resp.setToken(token);

            return Result.success(resp);
        } catch (BusinessException e) {
            log.info("login c={} e={}", e.getCode(), e.getMsg());
            return Result.businessErr(e);
        } catch (Exception e) {
            log.info("login e={}", e.getMessage());
            return Result.error(e.getMessage());
        }
    }

    @GetMapping("/verify-token")
    public Result verifyToken(String token) {
        log.info("verifyToken={}", token);

        // 设置键值对
        redisTemplate.opsForValue().set("key", "value");
        // 获取值
        String value = (String) redisTemplate.opsForValue().get("key");
        log.info("redisTemplate value={}", value);


        // token 验证
        DecodedJWT decodedJWT = JWTUtils.verifyToken(token);
        Integer uid = decodedJWT.getClaim("uid").asInt();
        String email = decodedJWT.getClaim("email").asString();

        log.info("tokenData k={} v={}", uid, email);
        return Result.success();
    }
}
