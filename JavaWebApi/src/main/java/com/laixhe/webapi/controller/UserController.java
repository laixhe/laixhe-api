package com.laixhe.webapi.controller;

import com.laixhe.webapi.common.Error;
import com.laixhe.webapi.common.Validators;
import com.laixhe.webapi.dto.UserListResponse;
import com.laixhe.webapi.dto.UserResponse;
import com.laixhe.webapi.dto.UserUpdateRequest;
import com.laixhe.webapi.security.ClaimsHolder;
import com.laixhe.webapi.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.enums.ParameterIn;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * 用户接口 (对应 Go 版 controllers/user.go)
 */
@Tag(name = "User", description = "用户")
@RestController
@RequestMapping("/api/v1/user")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    /** 更新用户信息 (需 Bearer 令牌, Uid 由 JWT 提供) */
    @Operation(summary = "更新用户信息",
            parameters = @Parameter(in = ParameterIn.HEADER, name = "Authorization", required = true, description = "Bearer XXX令牌"))
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "OK", content = @Content(schema = @Schema(implementation = UserResponse.class))),
            @ApiResponse(responseCode = "400", description = "Bad Request", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "422", description = "Unprocessable Entity", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "500", description = "Internal Server Error", content = @Content(schema = @Schema(implementation = Error.class))),
    })
    @PostMapping("/update")
    public UserResponse update(@RequestBody UserUpdateRequest req) {
        Validators.validateNickname(req.nickname());
        Validators.validateAvatarUrl(req.avatarUrl());
        return userService.update(ClaimsHolder.uid(), req);
    }

    /** 获取用户信息 */
    @Operation(summary = "获取用户信息")
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "OK", content = @Content(schema = @Schema(implementation = UserResponse.class))),
            @ApiResponse(responseCode = "400", description = "Bad Request", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "422", description = "Unprocessable Entity", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "500", description = "Internal Server Error", content = @Content(schema = @Schema(implementation = Error.class))),
    })
    @GetMapping("/info")
    public UserResponse info(@Parameter(description = "用户id", required = true)
                             @RequestParam("uid") int uid) {
        return userService.info(uid);
    }

    /** 获取用户列表 (page/page_size 缺省时归一化为 1/12, 与 Go 版一致) */
    @Operation(summary = "获取用户列表")
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "OK", content = @Content(schema = @Schema(implementation = UserListResponse.class))),
            @ApiResponse(responseCode = "400", description = "Bad Request", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "422", description = "Unprocessable Entity", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "500", description = "Internal Server Error", content = @Content(schema = @Schema(implementation = Error.class))),
    })
    @GetMapping("/list")
    public UserListResponse list(@Parameter(description = "分页-当前页(默认 1)", schema = @Schema(defaultValue = "1"))
                                 @RequestParam(value = "page", defaultValue = "0") int page,
                                 @Parameter(description = "分页-每页数量(默认 12)", schema = @Schema(defaultValue = "12"))
                                 @RequestParam(value = "page_size", defaultValue = "0") int pageSize) {
        return userService.list(page, pageSize);
    }
}
