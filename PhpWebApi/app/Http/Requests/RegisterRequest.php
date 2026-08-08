<?php

namespace App\Http\Requests;

use App\Result\Result;
use App\Result\ResultCode;
use Illuminate\Support\Facades\Validator;

/**
 * 注册请求参数 (与 Go 端 controllers/auth.go 校验对齐)
 */
class RegisterRequest implements IRequest
{
    public string $nickname = '';
    public string $email = '';
    public string $password = '';

    public function validator(array $params): ?Result
    {
        $validator = Validator::make($params, [
            // 昵称按"字符"计数 2~20 位
            'nickname' => ['required', 'string', 'min:2', 'max:20'],
            'email' => ['required', 'email'],
            // 密码只能包含字母 数字 _ @ $, 长度不能小于 6 位
            'password' => ['required', 'string', 'min:6', 'regex:/^[a-zA-Z0-9_@$]+$/'],
//            'type_id' => ['required', Rule::in([1,2,3])]
        ],
            [
                'nickname' => '昵称长度不能小于2位',
                'nickname.min' => '昵称长度不能小于2位',
                'nickname.max' => '昵称长度不能超过20位',
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
        $this->nickname = $params['nickname'] ?? '';
        $this->email = $params['email'] ?? '';
        $this->password = $params['password'] ?? '';
    }
}
