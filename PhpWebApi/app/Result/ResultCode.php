<?php

namespace App\Result;

/**
 * 响应码
 */
enum ResultCode: int
{
    // 通用 0 - 99
    case Success = 0; // 成功
    case Unknown = 1; // 未知错误
    case Service = 2; // 服务错误
    case Param = 3;   // 参数错误
    // 用户  100 - 199
    case AuthNotLogin = 100; // 未授权登录
    case AuthExpire = 101;  // 授权过期
    case AuthInvalid = 102;  // 授权无效
    case AuthUserError = 103; // 用户或密码错误
    case UserExist = 104;   // 用户已存在
    case UserNotExist = 105; // 用户不存在
    case EmailExist = 106;   // 邮箱已存在
    case EmailNotExist = 107; // 邮箱不存在


    public function text(): string
    {
        return match ($this) {
            self::Success => '成功',
            self::Unknown => '未知错误',
            self::Service => '服务错误',
            self::Param => '参数错误',
            self::AuthNotLogin => '未授权登录',
            self::AuthExpire => '授权过期',
            self::AuthInvalid => '授权无效',
            self::AuthUserError => '用户或密码错误',
            self::UserExist => '用户已存在',
            self::UserNotExist => '用户不存在',
            self::EmailExist => '邮箱已存在',
            self::EmailNotExist => '邮箱不存在',
        };
    }

    public static function intToEnum(int $code) : ResultCode {
        return match ($code) {
            2 => self::Service,
            3 => self::Param,
            100 => self::AuthNotLogin,
            101 => self::AuthExpire,
            102 => self::AuthInvalid,
            103 => self::AuthUserError,
            104 => self::UserExist,
            105 => self::UserNotExist,
            106 => self::EmailExist,
            107 => self::EmailNotExist,
            default => self::Unknown,
        };
    }

}
