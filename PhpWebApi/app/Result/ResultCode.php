<?php

namespace App\Result;

/**
 * 响应码
 */
enum ResultCode: int
{
    case Success = 0;              // 成功
    case BadRequest = 400;         // 请求格式错误 (如参数非数字, 与 Go/Rust 端绑定层 400 对齐)
    case AuthInvalid = 401;        // 授权无效
    case NotFound = 404;           // 路由不存在 (统一 JSON 格式, 与 Go/Rust/TS 端对齐)
    case Param = 422;              // 参数错误
    case TooManyRequests = 429;    // 请求过于频繁
    case Service = 500;            // 服务错误
    case ServiceUnavailable = 503; // 服务不可用

    public function text(): string
    {
        return match ($this) {
            self::Success => '成功',
            self::BadRequest => 'Bad Request',
            self::AuthInvalid => 'Unauthorized',
            self::NotFound => 'Not Found',
            self::Param => '参数错误',
            self::TooManyRequests => '请求过于频繁，请稍后再试',
            self::Service => '服务错误',
            self::ServiceUnavailable => '服务不可用',
        };
    }

    public static function intToEnum(int $code) : ResultCode {
        return match ($code) {
            400 => self::BadRequest,
            401 => self::AuthInvalid,
            404 => self::NotFound,
            422 => self::Param,
            429 => self::TooManyRequests,
            503 => self::ServiceUnavailable,
            default => self::Service,
        };
    }

}
