package com.laixhe.result;

import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * @author laixhe
 */
@AllArgsConstructor
@NoArgsConstructor
public enum ResultCode implements IResult, Serializable {
    SUCCESS(0,"成功"),
    ERROR_TOKEN_INVALID(1, "token无效"),
    ERROR_PARAM(2, "参数错误"),
    ERROR_PROMPT(3, "提示错误"),
    ERROR_SERVER(4,"系统异常"),
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
