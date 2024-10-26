package com.laixhe.response.auth;

import lombok.Data;
import lombok.ToString;

import com.laixhe.response.user.UserResponse;

/**
 * 登录响应参数
 *
 * @author laixhe
 */
@Data
@ToString
public class LoginResponse {
    private String token;
    private UserResponse user;
}
