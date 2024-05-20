package com.laixhe.exception;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import com.laixhe.result.ResultCode;

/**
 * 自定义业务异常类
 * @author laixhe
 */
@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class BusinessException extends RuntimeException {
    private ResultCode code;
    private String msg;
}

