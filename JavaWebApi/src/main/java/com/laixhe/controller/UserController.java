package com.laixhe.controller;

import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.laixhe.result.Result;
import com.laixhe.service.UserService;
import com.laixhe.response.user.InfoResponse;
import com.laixhe.response.user.ListResponse;

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
}
