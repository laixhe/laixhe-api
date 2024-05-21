package com.laixhe.service;

import com.laixhe.request.auth.RegisterRequest;
import com.laixhe.response.auth.LoginResponse;
import com.laixhe.response.auth.RegisterResponse;

/**
 * @author laixhe
 */
public interface AuthService {
    LoginResponse login(String email, String password);
    RegisterResponse register(RegisterRequest req);
}
