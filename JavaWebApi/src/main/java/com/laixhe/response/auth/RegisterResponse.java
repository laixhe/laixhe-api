package com.laixhe.response.auth;

import lombok.Data;
import lombok.ToString;

import com.laixhe.response.user.UserResponse;

/**
 * 刷新Jwt响应参数
 *
 * @author laixhe
 */
@Data
@ToString
public class RegisterResponse {
    private String token;
    private UserResponse user;
}
