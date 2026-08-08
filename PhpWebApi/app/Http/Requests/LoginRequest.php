<?php

namespace App\Http\Requests;

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Support\Facades\Validator;

/**
 * 登录请求参数 (与 Go 端 controllers/auth.go 校验对齐)
 */
class LoginRequest implements IRequest
{
    public string $email = '';
    public string $password = '';


    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            'email' => ['required', 'email'],
            // 密码只能包含字母 数字 _ @ $, 长度不能小于 6 位
            'password' => ['required', 'string', 'min:6', 'regex:/^[a-zA-Z0-9_@$]+$/'],
        ],
            [
                'email' => '邮箱格式错误',
                'password.min' => '密码长度不能小于6位',
                'password.regex' => '密码格式错误，只能包含字母 数字 _ @ $',
            ]);
        if ($validator->fails()) {
            return new Result(ResultCode::Param, $validator->errors()->first());
        }

        $this->param($params);

        return null;
    }

    public function param(array $params): void
    {
        $this->email = $params['email'] ?? '';
        $this->password = $params['password'] ?? '';
    }
}
