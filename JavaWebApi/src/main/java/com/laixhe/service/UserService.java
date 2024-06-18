package com.laixhe.service;

import com.laixhe.response.user.InfoResponse;
import com.laixhe.response.user.ListResponse;

import java.time.LocalDateTime;

/**
 * @author laixhe
 */
public interface UserService {
    /**
     * 查询用户信息
     *
     * @param uid 用户ID
     * @throws com.laixhe.exception.BusinessException 异常
     */
    InfoResponse info(int uid);

    /**
     * 查询用户列表
     *
     * @param size 每页页数(数量)
     * @param page 分页当前页数
     */
    ListResponse list(int size, int page);

    /**
     * 修改用户信息
     *
     * @param uid     用户ID
     * @param uname   用户名
     * @param loginAt 登录时间
     * @throws com.laixhe.exception.BusinessException 异常
     */
    void update(int uid, String uname, LocalDateTime loginAt);
}
