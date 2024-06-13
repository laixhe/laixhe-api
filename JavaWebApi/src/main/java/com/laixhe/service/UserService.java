package com.laixhe.service;

import com.laixhe.response.user.InfoResponse;
import com.laixhe.response.user.ListResponse;

/**
 * @author laixhe
 */
public interface UserService {
    InfoResponse info(int uid);
    ListResponse list(int size, int page);
    void update(int uid, String uname);
}
