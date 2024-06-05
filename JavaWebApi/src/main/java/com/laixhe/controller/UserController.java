package com.laixhe.controller;

import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import com.laixhe.result.Result;
import com.laixhe.service.UserService;
import com.laixhe.response.user.InfoResponse;
import com.laixhe.response.user.ListResponse;

import java.util.ArrayList;
import java.util.List;

/**
 * @author laixhe
 */
@Slf4j
@RestController
@RequestMapping("/api/user")
public class UserController {

    private final UserService userService;

    @Autowired
    public UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping("/info")
    public Result<InfoResponse> info(HttpServletRequest request){
        int uid = (int)request.getAttribute("uid");
        log.info("info uid={}", uid);

        InfoResponse resp = userService.info(uid);
        return Result.success(resp);
    }

    @GetMapping("/list")
    public Result<ListResponse> list(){
        log.info("list");
        ListResponse resp = userService.list();
        return Result.success(resp);
    }

    @GetMapping("/test1")
    public Result<List<String>> test1(@RequestParam int uid, @RequestParam String uname){
        // http://webapi.laixhe.com/api/user/test1?uid=111&uname=laixhe
        List<String> list = new ArrayList<>();
        list.add("test1");
        list.add(String.valueOf(uid));
        list.add(uname);
        log.info("test1 uid={}", uid);
        return Result.success(list);
    }

    @GetMapping("/test2/{uid}")
    public Result<List<String>> test2(@PathVariable int uid, @RequestHeader("Authorization") String authorization){
        // http://webapi.laixhe.com/api/user/test2/111
        List<String> list = new ArrayList<>();
        list.add("test2");
        list.add(String.valueOf(uid));
        list.add(authorization);
        return Result.success(list);
    }

}
