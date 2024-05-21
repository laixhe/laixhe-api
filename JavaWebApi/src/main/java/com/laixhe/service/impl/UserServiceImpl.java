package com.laixhe.service.impl;

import com.mybatisflex.core.query.QueryWrapper;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

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
                .eq("id", uid)
                .isNull("deleted_at");
        // SELECT * FROM `user` WHERE id = ? AND deleted_at IS NULL LIMIT 1
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
    public ListResponse list() {

        QueryWrapper queryWrapper = QueryWrapper.create()
                .isNull("deleted_at")
                .orderBy("id DESC");
        // SELECT * FROM `user` WHERE deleted_at IS NULL ORDER BY id DESC
        List<User> users = userMapper.selectListByQuery(queryWrapper);

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
        resp.setUsers(userResponses);
        return resp;
    }

}
