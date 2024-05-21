package com.laixhe.exception;

import org.springframework.validation.ObjectError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import lombok.extern.slf4j.Slf4j;

import java.util.List;

import com.laixhe.result.Result;
import com.laixhe.result.ResultCode;


/**
 * 全局异常处理
 * 来处理各种异常,包括自己定义的异常和内部异常
 *
 * @author laixhe
 */
@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    /**
     * 处理请求参数异常
     */
    @ExceptionHandler(value = MethodArgumentNotValidException.class)
    public <T> Result<T> handleMethodArgumentNotValidException(MethodArgumentNotValidException e) {

        List<ObjectError> errors = e.getBindingResult().getAllErrors();
        StringBuilder errorMsg = new StringBuilder();
        errors.forEach(error -> errorMsg.append(error.getDefaultMessage()).append("; "));

        return Result.otherErr(ResultCode.Param, errorMsg.toString());
    }

    /**
     * 处理自定义异常
     */
    @ExceptionHandler(value = BusinessException.class)
    @ResponseBody
    public <T> Result<T> bizExceptionHandler(BusinessException e) {
        return Result.businessErr(e);
    }

    /**
     * 处理其他异常
     */
    @ExceptionHandler(value = Exception.class)
    @ResponseBody
    public <T> Result<T> exceptionHandler(Exception e) {
        return Result.otherErr(ResultCode.Service, e.getMessage());
    }

}

