package com.laixhe.controller;

import lombok.extern.slf4j.Slf4j;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import cn.hutool.core.date.LocalDateTimeUtil;
import cn.hutool.core.date.DatePattern;

import com.laixhe.result.ResultCode;
import com.laixhe.result.Result;
import com.laixhe.service.UserService;
import com.laixhe.response.user.InfoResponse;
import com.laixhe.response.user.ListResponse;
import com.laixhe.request.user.UpdateRequest;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

/**
 * @author laixhe
 */
@Slf4j
@RestController
@RequestMapping("api/user")
public class UserController {

    private final UserService userService;

    @Autowired
    public UserController(UserService userService) {
        this.userService = userService;
    }

    /**
     * 查询用户信息
     */
    @GetMapping("info")
    public Result<InfoResponse> info(HttpServletRequest request) {
        int uid = (int) request.getAttribute("uid");
        log.info("info uid={}", uid);

        InfoResponse resp = userService.info(uid);
        return Result.success(resp);
    }

    /**
     * 查询用户列表
     */
    @GetMapping("list")
    public Result<ListResponse> list(@RequestParam int size, @RequestParam int page) {
        // page 分页当前页数
        // size 每页页数(数量)
        if (size <= 0) {
            size = 20;
        }
        if (page <= 0) {
            page = 1;
        }
        log.info("list size={} page={}", size, page);
        ListResponse resp = userService.list(size, page);
        return Result.success(resp);
    }

    /**
     * 修改用户信息
     */
    @PostMapping("update")
    public Result<List<String>> update(HttpServletRequest request, @RequestBody @Validated UpdateRequest req) {
        int uid = (int) request.getAttribute("uid");
        log.info("update {}", req.toString());

        LocalDateTime loginAt = null;
        try {
            loginAt = LocalDateTimeUtil.parse(req.getLoginAt(), DatePattern.NORM_DATETIME_PATTERN);
        } catch (Exception e) {
            return Result.otherErr(ResultCode.Param, "登录时间格式不对！" + e.getMessage());
        }

        userService.update(uid, req.getUname(), loginAt);

        return Result.success();
    }

    @GetMapping("test1")
    public Result<List<String>> test1(@RequestParam int uid, @RequestParam String uname) {
        // http://webapi.laixhe.com/api/user/test1?uid=111&uname=laixhe
        List<String> list = new ArrayList<>();
        list.add("test1");
        list.add(String.valueOf(uid));
        list.add(uname);
        log.info("test1 uid={}", uid);
        return Result.success(list);
    }

    @GetMapping("test2/{uid}")
    public Result<List<String>> test2(@PathVariable int uid, @RequestHeader("Authorization") String authorization) {
        // http://webapi.laixhe.com/api/user/test2/111
        List<String> list = new ArrayList<>();
        list.add("test2");
        list.add(String.valueOf(uid));
        list.add(authorization);
        return Result.success(list);
    }

}
