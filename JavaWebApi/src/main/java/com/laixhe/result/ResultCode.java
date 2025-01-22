package com.laixhe.result;

import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * 响应码
 *
 * @author laixhe
 */
@AllArgsConstructor
@NoArgsConstructor
public enum ResultCode implements IResult, Serializable {
    // 通用 0 - 99
    Success(0, "成功"),
    Unknown(1, "未知错误"),
    Service(2, "服务错误"),
    Param(3, "参数错误"),
    // 用户  100 - 199
    AuthNotLogin(100, "未授权登录"),
    AuthExpire(101, "授权过期"),
    AuthInvalid(102, "授权无效"),
    AuthUserError(103, "用户或密码错误"),
    UserExist(104, "用户已存在"),
    UserNotExist(105, "用户不存在"),
    EmailExist(106, "邮箱已存在"),
    EmailNotExist(107, "邮箱不存在"),
    ;

    private int code;
    private String msg;

    @Override
    public int getCode() {
        return code;
    }

    @Override
    public String getMsg() {
        return msg;
    }

}
