package com.laixhe.response.auth;

import lombok.Data;
import lombok.ToString;

import com.laixhe.response.user.UserResponse;

/**
 * @author laixhe
 */
@Data
@ToString
public class LoginResponse {
    private String token;
    private UserResponse user;
}
