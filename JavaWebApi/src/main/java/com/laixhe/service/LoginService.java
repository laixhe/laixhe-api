package com.laixhe.service;

import com.laixhe.entity.User;

/**
 * @author laixhe
 */
public interface LoginService {
    User login(String email, String password);
}
