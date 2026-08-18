package com.laixhe.webapi.common;

import lombok.Getter;

/**
 * 业务异常: 状态码即响应 HTTP 状态码, message 即响应错误文案
 */
@Getter
public class ApiException extends RuntimeException {

    private final int status;

    public ApiException(int status, String message) {
        super(message);
        this.status = status;
    }

    /** 参数错误 → 422 */
    public static ApiException paramError(String message) {
        return new ApiException(422, message);
    }

    /** 未授权 → 401 */
    public static ApiException unauthorized() {
        return new ApiException(401, "Unauthorized");
    }
}
