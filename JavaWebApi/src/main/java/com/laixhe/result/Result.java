package com.laixhe.result;

import lombok.Data;

import java.io.Serializable;

import com.laixhe.exception.BusinessException;

/**
 * @author laixhe
 */
@Data
public class Result<T> implements Serializable {

    private int code;
    private String msg;
    private T data;

    public static <T> Result<T> success() {
        return success(null);
    }

    public static <T> Result<T> success(T data) {
        Result<T> result = new Result<>();
        result.setCode(ResultCode.Success.getCode());
        result.setMsg(ResultCode.Success.getMsg());
        result.setData(data);
        return result;
    }

    public static <T> Result<T> error(String msg) {
        Result<T> result = new Result<>();
        result.setCode(ResultCode.Service.getCode());
        result.setMsg(msg);
        return result;
    }

    // 自定义异常返回的结果
    public static <T> Result<T> businessErr(BusinessException e) {
        Result<T> result = new Result<>();
        result.setCode(e.getCode().getCode());
        result.setMsg(e.getMsg().isBlank() ? e.getCode().getMsg() : e.getMsg());
        result.setData(null);
        return result;
    }

    // 其他异常处理方法返回的结果
    public static <T> Result<T> otherErr(ResultCode resultCode, String msg) {
        Result<T> result = new Result<>();
        result.setCode(resultCode.getCode());
        result.setMsg(msg.isBlank() ? resultCode.getMsg() : msg);
        result.setData(null);
        return result;
    }


}


