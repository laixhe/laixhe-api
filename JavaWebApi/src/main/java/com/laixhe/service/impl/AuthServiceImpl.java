package com.laixhe.service.impl;

import cn.hutool.crypto.digest.BCrypt;
import com.mybatisflex.core.query.QueryWrapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import lombok.extern.slf4j.Slf4j;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.Map;

import com.laixhe.exception.BusinessException;
import com.laixhe.utils.JWTUtils;
import com.laixhe.service.AuthService;
import com.laixhe.entity.User;
import com.laixhe.mapper.UserMapper;
import com.laixhe.result.ResultCode;
import com.laixhe.response.auth.LoginResponse;
import com.laixhe.response.user.UserResponse;
import com.laixhe.request.auth.RegisterRequest;
import com.laixhe.response.auth.RegisterResponse;

/**
 * @author laihxe
 */
@Slf4j
@Service
public class AuthServiceImpl implements AuthService {

    private final UserMapper userMapper;

    @Autowired
    public AuthServiceImpl(UserMapper userMapper) {
        this.userMapper = userMapper;
    }

    @Override
    public RegisterResponse register(RegisterRequest req){

        LocalDateTime now = LocalDateTime.now();

        User user = new User();
        user.setPassword(BCrypt.hashpw(req.getPassword(), BCrypt.gensalt()));
        user.setEmail(req.getEmail());
        user.setUname(req.getUname());
        user.setAge(req.getAge());
        user.setScore(0.0);
        user.setLoginAt(now);
        user.setCreatedAt(now);
        user.setUpdatedAt(now);

        // INSERT INTO `user`(`password`, `email`, `uname`, `age`, `score`, `login_at`, `created_at`, `updated_at`, `deleted_at`) VALUES (?, ?, ?, ?, ?, ?, now(), now(), ?)
        userMapper.insert(user);

        // 生成 token 数据
        Map<String, Object> payload = new HashMap<>();
        payload.put("uid", user.getId());
        String token = JWTUtils.createToken(payload);

        UserResponse userResponse = new UserResponse();
        userResponse.setUid(user.getId());
        userResponse.setEmail(user.getEmail());
        userResponse.setUname(user.getUname());
        userResponse.setCreatedAt(user.getCreatedAt());

        RegisterResponse resp = new RegisterResponse();
        resp.setToken(token);
        resp.setUser(userResponse);
        return resp;
    }

    @Override
    public LoginResponse login(String email, String password) {

        QueryWrapper queryWrapper = QueryWrapper.create()
                .eq("email", email)
                .isNull("deleted_at");
        // SELECT * FROM `user` WHERE email = ? AND deleted_at IS NULL LIMIT 1
        User user = userMapper.selectOneByQuery(queryWrapper);
        if (user == null) {
            throw new BusinessException(ResultCode.AuthUserError, "");
        }
        // 判断密码是否匹配
        if (!BCrypt.checkpw(password, user.getPassword())) {
            throw new BusinessException(ResultCode.AuthUserError, "");
        }

        // 生成 token 数据
        Map<String, Object> payload = new HashMap<>();
        payload.put("uid", user.getId());
        String token = JWTUtils.createToken(payload);

        UserResponse userResponse = new UserResponse();
        userResponse.setUid(user.getId());
        userResponse.setEmail(user.getEmail());
        userResponse.setUname(user.getUname());
        userResponse.setCreatedAt(user.getCreatedAt());

        LoginResponse resp = new LoginResponse();
        resp.setToken(token);
        resp.setUser(userResponse);
        return resp;
    }

}
