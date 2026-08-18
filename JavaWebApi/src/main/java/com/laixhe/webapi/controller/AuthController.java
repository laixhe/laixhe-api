package com.laixhe.webapi.controller;

import com.laixhe.webapi.common.Error;
import com.laixhe.webapi.common.Validators;
import com.laixhe.webapi.dto.AuthLoginRequest;
import com.laixhe.webapi.dto.AuthLoginResponse;
import com.laixhe.webapi.dto.AuthRefreshResponse;
import com.laixhe.webapi.dto.AuthRegisterRequest;
import com.laixhe.webapi.dto.AuthRegisterResponse;
import com.laixhe.webapi.security.ClaimsHolder;
import com.laixhe.webapi.service.AuthService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.enums.ParameterIn;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * 鉴权接口 (对应 Go 版 controllers/auth.go)
 */
@Tag(name = "Auth", description = "鉴权")
@RestController
@RequestMapping("/api/v1/auth")
@RequiredArgsConstructor
public class AuthController {

    private final AuthService authService;

    /** 注册 */
    @Operation(summary = "注册")
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "OK", content = @Content(schema = @Schema(implementation = AuthRegisterResponse.class))),
            @ApiResponse(responseCode = "400", description = "Bad Request", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "422", description = "Unprocessable Entity", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "500", description = "Internal Server Error", content = @Content(schema = @Schema(implementation = Error.class))),
    })
    @PostMapping("/register")
    public AuthRegisterResponse register(@Valid @RequestBody AuthRegisterRequest req) {
        Validators.validateNickname(req.nickname());
        return authService.register(req);
    }

    /** 登录 */
    @Operation(summary = "登录")
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "OK", content = @Content(schema = @Schema(implementation = AuthLoginResponse.class))),
            @ApiResponse(responseCode = "400", description = "Bad Request", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "422", description = "Unprocessable Entity", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "500", description = "Internal Server Error", content = @Content(schema = @Schema(implementation = Error.class))),
    })
    @PostMapping("/login")
    public AuthLoginResponse login(@Valid @RequestBody AuthLoginRequest req) {
        return authService.login(req);
    }

    /** 刷新Jwt (需 Bearer 令牌, Uid 由 JWT 提供) */
    @Operation(summary = "刷新Jwt",
            parameters = @Parameter(in = ParameterIn.HEADER, name = "Authorization", required = true, description = "Bearer XXX令牌"))
    @ApiResponses({
            @ApiResponse(responseCode = "200", description = "OK", content = @Content(schema = @Schema(implementation = AuthRefreshResponse.class))),
            @ApiResponse(responseCode = "400", description = "Bad Request", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "401", description = "Unauthorized", content = @Content(schema = @Schema(implementation = Error.class))),
            @ApiResponse(responseCode = "500", description = "Internal Server Error", content = @Content(schema = @Schema(implementation = Error.class))),
    })
    @PostMapping("/refresh")
    public AuthRefreshResponse refresh() {
        return authService.refresh(ClaimsHolder.uid());
    }
}
