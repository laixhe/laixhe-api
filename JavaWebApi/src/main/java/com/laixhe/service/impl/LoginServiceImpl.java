package com.laixhe.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import lombok.extern.slf4j.Slf4j;

import java.util.HashMap;
import java.util.Map;

import com.laixhe.exception.BusinessException;
import com.laixhe.service.LoginService;
import com.laixhe.entity.User;
import com.laixhe.mapper.UserMapper;
import com.laixhe.result.ResultCode;

/**
 * @author laihxe
 */
@Slf4j
@Service
public class LoginServiceImpl implements LoginService {

    private final UserMapper userMapper;

    @Autowired
    public LoginServiceImpl(UserMapper userMapper) {
        this.userMapper = userMapper;
    }

    @Override
    public User login(String email, String password) {

        Map<String, Object> where = new HashMap<>();
        where.put("email", email);

        User user = userMapper.selectOneByMap(where);
        if (user == null) {
            throw new BusinessException(ResultCode.ERROR_PROMPT, "用户不存在!");
        }
        if (!user.getPassword().equals(password)) {
            throw new BusinessException(ResultCode.ERROR_PROMPT, "密码错误!");
        }
        return user;
    }

}
