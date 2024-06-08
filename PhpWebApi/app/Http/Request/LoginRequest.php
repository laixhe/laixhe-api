<?php

namespace App\Http\Request;

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Support\Facades\Validator;

/**
 * 登录请求参数
 */
class LoginRequest implements IRequest
{
    public string $email = '';
    public string $password = '';


    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            'email' => ['required', 'email'],
            'password' => ['required', 'string', 'min:6', 'max:20'],
        ],
            [
                'email' => '邮箱格式不正确！',
                'password' => '密码长度应该在6到20位！',
            ]);
        if ($validator->fails()) {
            return new Result(ResultCode::Param, $validator->errors()->first(), []);
        }

        $this->param($params);

        return null;
    }

    public function param(array $params): void
    {
        $this->email = $params['email'];
        $this->password = $params['password'];
    }
}
