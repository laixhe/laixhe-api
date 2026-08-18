package com.laixhe.webapi.common;

import jakarta.validation.ConstraintViolationException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.MissingServletRequestParameterException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;
import org.springframework.web.servlet.resource.NoResourceFoundException;

/**
 * 全局异常处理, 统一返回 {"code":xx,"message":xx} 错误体
 * (与 Go 版 ErrorHandler 对齐: 业务/参数错误原样返回, 未知错误统一 500 文案)
 */
@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    /** 业务异常: 状态码 + 文案原样返回 */
    @ExceptionHandler(ApiException.class)
    public ResponseEntity<Error> handleApiException(ApiException e) {
        return ResponseEntity.status(e.getStatus()).body(new Error(e.getStatus(), e.getMessage()));
    }

    /** @RequestBody 校验失败 → 422, 取第一个字段错误信息 */
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Error> handleValidation(MethodArgumentNotValidException e) {
        String message = e.getBindingResult().getFieldErrors().stream()
                .findFirst()
                .map(FieldError::getDefaultMessage)
                .orElse("参数错误");
        return ResponseEntity.status(422).body(new Error(422, message));
    }

    @ExceptionHandler(ConstraintViolationException.class)
    public ResponseEntity<Error> handleConstraintViolation(ConstraintViolationException e) {
        return ResponseEntity.status(422).body(new Error(422, e.getMessage()));
    }

    /** 请求体不可解析 (JSON 格式错误等) → 400 */
    @ExceptionHandler(HttpMessageNotReadableException.class)
    public ResponseEntity<Error> handleNotReadable(HttpMessageNotReadableException e) {
        return ResponseEntity.badRequest().body(new Error(400, "Bad Request"));
    }

    /** 缺少必填 query 参数 → 400 */
    @ExceptionHandler(MissingServletRequestParameterException.class)
    public ResponseEntity<Error> handleMissingParam(MissingServletRequestParameterException e) {
        return ResponseEntity.badRequest().body(new Error(400, "Bad Request"));
    }

    /** query 参数类型错误 → 400 */
    @ExceptionHandler(MethodArgumentTypeMismatchException.class)
    public ResponseEntity<Error> handleTypeMismatch(MethodArgumentTypeMismatchException e) {
        return ResponseEntity.badRequest().body(new Error(400, "Bad Request"));
    }

    /** 未匹配路由 → 404 (与 Go/Rust/TS/PHP 端对齐) */
    @ExceptionHandler(NoResourceFoundException.class)
    public ResponseEntity<Error> handleNotFound(NoResourceFoundException e) {
        return ResponseEntity.status(404).body(new Error(404, "Not Found"));
    }

    /** 其余未知错误: 记录服务端日志后统一返回固定 500 文案, 避免泄露内部实现细节 */
    @ExceptionHandler(Throwable.class)
    public ResponseEntity<Error> handleUnknown(Throwable e) {
        log.error("unhandled error: {}", e.getMessage(), e);
        return ResponseEntity.internalServerError().body(new Error(500, "internal server error"));
    }
}
