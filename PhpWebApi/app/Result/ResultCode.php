<?php

namespace App\Result;

/**
 * 响应码
 */
enum ResultCode: int
{
    case Success = 0;              // 成功
    case AuthInvalid = 401;        // 授权无效
    case Param = 422;              // 参数错误
    case TooManyRequests = 429;    // 请求过于频繁
    case Service = 500;            // 服务错误
    case ServiceUnavailable = 503; // 服务不可用

    public function text(): string
    {
        return match ($this) {
            self::Success => '成功',
            self::Service => '服务错误',
            self::Param => '参数错误',
            self::AuthInvalid => '授权无效',
            self::TooManyRequests => '请求过于频繁，请稍后再试',
            self::ServiceUnavailable => '服务不可用',
        };
    }

    public static function intToEnum(int $code) : ResultCode {
        return match ($code) {
            401 => self::AuthInvalid,
            422 => self::Param,
            429 => self::TooManyRequests,
            503 => self::ServiceUnavailable,
            default => self::Service,
        };
    }

}
