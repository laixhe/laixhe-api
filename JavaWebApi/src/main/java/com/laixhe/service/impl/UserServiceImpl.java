package com.laixhe.service.impl;

import lombok.extern.slf4j.Slf4j;
import com.mybatisflex.core.paginate.Page;
import com.mybatisflex.core.query.QueryWrapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

import com.laixhe.entity.User;
import com.laixhe.exception.BusinessException;
import com.laixhe.mapper.UserMapper;
import com.laixhe.response.user.InfoResponse;
import com.laixhe.response.user.ListResponse;
import com.laixhe.response.user.UserResponse;
import com.laixhe.result.ResultCode;
import com.laixhe.service.UserService;

/**
 * @author laihxe
 */
@Slf4j
@Service
public class UserServiceImpl implements UserService {

    private final UserMapper userMapper;

    @Autowired
    public UserServiceImpl(UserMapper userMapper) {
        this.userMapper = userMapper;
    }

    @Override
    public InfoResponse info(int uid) {

        QueryWrapper queryWrapper = QueryWrapper.create()
                .select("id", "email", "uname", "created_at")
                .eq("id", uid)
                .isNull("deleted_at");
        // SELECT id, email, uname, created_at FROM `user` WHERE id = ? AND deleted_at IS NULL  LIMIT 1
        User user = userMapper.selectOneByQuery(queryWrapper);
        if (user == null) {
            throw new BusinessException(ResultCode.AuthNotLogin, "");
        }

        UserResponse userResponse = new UserResponse();
        userResponse.setUid(user.getId());
        userResponse.setEmail(user.getEmail());
        userResponse.setUname(user.getUname());
        userResponse.setCreatedAt(user.getCreatedAt());

        InfoResponse resp = new InfoResponse();
        resp.setUser(userResponse);
        return resp;
    }

    @Override
    public ListResponse list(int size, int page) {

        //QueryWrapper queryWrapper = QueryWrapper.create()
        //        .select("id", "email", "uname", "created_at")
        //        .isNull("deleted_at")
        //        .orderBy("id DESC");
        // SELECT id, email, uname, created_at FROM `user` WHERE deleted_at IS NULL ORDER BY id DESC
        //List<User> users = userMapper.selectListByQuery(queryWrapper);

        // SELECT COUNT(*) AS `total` FROM `user` WHERE deleted_at IS NULL
        // SELECT id, email, uname, created_at FROM `user` WHERE deleted_at IS NULL ORDER BY id DESC LIMIT 2, 20
        QueryWrapper queryWrapper = QueryWrapper.create()
                .select("id", "email", "uname", "created_at")
                .isNull("deleted_at")
                .orderBy("id DESC");

        Page<User> userPage = userMapper.paginate(page, size, queryWrapper);
        List<User> users = userPage.getRecords();

        List<UserResponse> userResponses = new ArrayList<>();
        for (User user : users) {
            UserResponse userResponse = new UserResponse();
            userResponse.setUid(user.getId());
            userResponse.setEmail(user.getEmail());
            userResponse.setUname(user.getUname());
            userResponse.setCreatedAt(user.getCreatedAt());
            userResponses.add(userResponse);
        }

        ListResponse resp = new ListResponse();
        resp.setList(userResponses);
        resp.setPage((int) userPage.getPageNumber());
        resp.setSize((int) userPage.getPageSize());
        resp.setTotal((int) userPage.getTotalRow());
        return resp;
    }

    @Override
    public void update(int uid, String uname, LocalDateTime loginAt) {

        // SELECT id, email FROM `user` WHERE uname = ? AND deleted_at IS NULL LIMIT 1
        QueryWrapper queryWrapper = QueryWrapper.create()
                .select("id", "email")
                .eq("uname", uname)
                .isNull("deleted_at");
        User user = userMapper.selectOneByQuery(queryWrapper);
        if (user != null) {
            if (user.getId().equals(uid)) {
                return;
            }
            throw new BusinessException(ResultCode.Param, "用户名已存在！");
        }

        // UPDATE `user` SET `uname` = ? , `login_at` = ? , `updated_at` = now() WHERE id = ? LIMIT 1
        user = new User();
        user.setUname(uname);
        user.setLoginAt(loginAt);
        queryWrapper = QueryWrapper.create().eq("id", uid).limit(1);
        userMapper.updateByQuery(user, queryWrapper);
    }

}
