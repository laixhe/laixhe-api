package com.laixhe.webapi.controller;

import com.laixhe.webapi.common.Error;
import io.swagger.v3.oas.annotations.Operation;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

/**
 * 未匹配路由统一返回 {"code":404,"message":"Not Found"} (与 Go/Rust/TS/PHP 端对齐)
 */
@RestController
public class NotFoundController {

    /** 隐藏: 兜底路由不进入 OpenAPI 文档 */
    @Operation(hidden = true)
    @ResponseStatus(HttpStatus.NOT_FOUND)
    @RequestMapping("/**")
    public Error notFound() {
        return new Error(HttpStatus.NOT_FOUND.value(), "Not Found");
    }
}
